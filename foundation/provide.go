package foundation

import "github.com/samber/do"

func Provide[T any](i *do.Injector, t T) {
	do.Provide[T](i, func(injector *do.Injector) (T, error) {
		return t, nil
	})
}

func ProvideValue[T any](i *do.Injector, t T) {
	do.ProvideValue[T](i, t)
}

func Invoke[T any](i *do.Injector) (T, error) {
	return do.Invoke[T](i)
}
