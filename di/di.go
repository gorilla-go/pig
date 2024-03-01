package di

import (
	"errors"
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"github.com/gorilla-go/pig/foundation/injector"
	"reflect"
	"strings"
	"unsafe"
)

type Provider func(c *Container) (any, error)

func ProvideLazy[T any](c *Container, provider Provider, name ...string) {
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

func ProvideNew[T any](c *Container, provider Provider, name ...string) {
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
		a, err := v.(Provider)(c)
		return a.(T), err
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	if v, ok := c.eager[typeStr]; ok {
		return v.(T), nil
	}

	if v, ok := c.lazy[typeStr]; ok {
		t, err := v.(Provider)(c)
		if err != nil {
			return t.(T), err
		}

		if _, ok := c.eager[typeStr]; ok {
			panic(fmt.Sprintf("DI: type %s already exists", typeStr))
		}
		c.eager[typeStr] = t
		return t.(T), nil
	}

	return *new(T), errors.New(fmt.Sprintf("DI: type %s not found", typeStr))
}

func InvokeNamed[T any](c *Container, name string) (T, error) {
	if v, ok := c.rebuildNamed[name]; ok {
		a, err := c.rebuild[v].(Provider)(c)
		if err != nil {
			return a.(T), err
		}
		return a.(T), nil
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	if v, ok := c.eagerNamed[name]; ok {
		return c.eager[v].(T), nil
	}

	if v, ok := c.lazyNamed[name]; ok {
		t, err := c.lazy[v].(Provider)(c)
		if err != nil {
			return t.(T), err
		}

		if _, ok := c.eagerNamed[name]; ok {
			panic(fmt.Sprintf("DI: name %s already exists", name))
		}

		if _, ok := c.eager[v]; ok {
			panic(fmt.Sprintf("DI: type %s already exists", v))
		}

		c.eager[v] = t
		c.eagerNamed[name] = v
		return t.(T), nil
	}

	return *new(T), errors.New(fmt.Sprintf("DI: name %s not found", name))
}

func MustInvoke[T any](c *Container, t ...T) T {
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

func Autowire[T any](c *Container, s T) T {
	v := reflect.ValueOf(s)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("inject target must be struct or struct pointer.")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag, ok := field.Tag.Lookup("di")
		if !ok {
			continue
		}

		tag = strings.TrimSpace(tag)
		if tag != "" {
			injector.ServiceInjector(
				field.Type,
				MustInvokeNamed[any](c, tag),
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
			continue
		}

		s, err := invokeTypeName(c, typeName(field.Type.String()))
		if err != nil {
			panic(err)
		}
		injector.ServiceInjector(
			field.Type,
			s,
			unsafe.Pointer(v.Field(i).UnsafeAddr()),
		)
	}
	return s
}

func invokeTypeName(c *Container, typeStr typeName) (any, error) {
	if v, ok := c.rebuild[typeStr]; ok {
		return v.(Provider)(c)
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	if v, ok := c.eager[typeStr]; ok {
		return v, nil
	}

	if v, ok := c.lazy[typeStr]; ok {
		t, err := v.(Provider)(c)
		if err != nil {
			return t, err
		}

		if _, ok := c.eager[typeStr]; ok {
			panic(fmt.Sprintf("DI: type %s already exists", typeStr))
		}
		c.eager[typeStr] = t
		return t, nil
	}

	panic(fmt.Sprintf("DI: type %s not found", typeStr))
}

func typeToString[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
