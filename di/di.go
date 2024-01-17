package di

import (
	"errors"
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"reflect"
)

type Provider[T any] func(c *Container) (T, error)

func ProvideLazy[T any](c *Container, provider Provider[T], name ...string) {
	registerName := foundation.DefaultParam(name, "")
	typeStr := typeName(typeToString[T]())

	typeExistStr := fmt.Sprintf("DI: type %s already exists", typeStr)
	if _, ok := c.eager[typeStr]; ok {
		panic(typeExistStr)
	}

	if _, ok := c.lazy[typeStr]; ok {
		panic(typeExistStr)
	}
	if registerName == "" {
		c.lazy[typeStr] = provider
		return
	}

	nameExistStr := fmt.Sprintf("DI: name %s already exists", registerName)
	if _, ok := c.eagerNamed[registerName]; ok {
		panic(nameExistStr)
	}

	if _, ok := c.lazyNamed[registerName]; ok {
		panic(nameExistStr)
	}
	c.lazyNamed[registerName] = typeStr
	c.lazy[typeStr] = provider
}

func ProvideValue[T any](c *Container, value T, name ...string) {
	registerName := foundation.DefaultParam(name, "")
	typeStr := typeName(typeToString[T]())

	existStr := fmt.Sprintf("DI: type %s already exists", typeStr)
	if _, ok := c.eager[typeStr]; ok {
		panic(existStr)
	}

	if _, ok := c.lazy[typeStr]; ok {
		panic(existStr)
	}

	if registerName == "" {
		c.eager[typeStr] = value
		return
	}

	nameExistStr := fmt.Sprintf("DI: name %s already exists", registerName)
	if _, ok := c.eagerNamed[registerName]; ok {
		panic(nameExistStr)
	}

	if _, ok := c.lazyNamed[registerName]; ok {
		panic(nameExistStr)
	}
	c.eagerNamed[registerName] = typeStr
	c.eager[typeStr] = value
}

func ProvideNew[T any](c *Container, provider Provider[T], name ...string) {
	registerName := foundation.DefaultParam(name, "")
	typeStr := typeName(typeToString[T]())

	if _, ok := c.rebuild[typeStr]; ok {
		panic(fmt.Sprintf("DI: type %s already exists", typeStr))
	}

	if registerName == "" {
		c.rebuild[typeStr] = provider
		return
	}

	if _, ok := c.rebuildNamed[registerName]; ok {
		panic(fmt.Sprintf("DI: name %s already exists", registerName))
	}
	c.rebuildNamed[registerName] = typeStr
	c.rebuild[typeStr] = provider
}

func Invoke[T any](c *Container) (T, error) {
	typeStr := typeName(typeToString[T]())
	if v, ok := c.rebuild[typeStr]; ok {
		return v.(Provider[T])(c)
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	if v, ok := c.eager[typeStr]; ok {
		return v.(T), nil
	}

	if v, ok := c.lazy[typeStr]; ok {
		t, err := v.(Provider[T])(c)
		if err != nil {
			return t, err
		}

		if _, ok := c.eager[typeStr]; ok {
			panic(fmt.Sprintf("DI: type %s already exists", typeStr))
		}
		c.eager[typeStr] = t
		return t, nil
	}

	return *new(T), errors.New(fmt.Sprintf("DI: type %s not found", typeStr))
}

func InvokeNamed[T any](c *Container, name string) (T, error) {
	if v, ok := c.rebuildNamed[name]; ok {
		return c.rebuild[v].(Provider[T])(c)
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	if v, ok := c.eagerNamed[name]; ok {
		return c.eager[v].(T), nil
	}

	if v, ok := c.lazyNamed[name]; ok {
		t, err := c.lazy[v].(Provider[T])(c)
		if err != nil {
			return t, err
		}

		if _, ok := c.eagerNamed[name]; ok {
			panic(fmt.Sprintf("DI: name %s already exists", name))
		}

		if _, ok := c.eager[v]; ok {
			panic(fmt.Sprintf("DI: type %s already exists", v))
		}

		c.eager[v] = t
		c.eagerNamed[name] = v
		return t, nil
	}

	return *new(T), errors.New(fmt.Sprintf("DI: name %s not found", name))
}

func MustInvoke[T any](c *Container) T {
	v, err := Invoke[T](c)
	if err != nil {
		panic(err)
	}
	return v
}

func MustInvokeNamed[T any](c *Container, name string) T {
	v, err := InvokeNamed[T](c, name)
	if err != nil {
		panic(err)
	}
	return v
}

func typeToString[T any]() string {
	t := fmt.Sprintf("%T", *new(T))
	if t == "<nil>" {
		return reflect.TypeOf((*T)(nil)).Elem().String()
	}
	return t
}
