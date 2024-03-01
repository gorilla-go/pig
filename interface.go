package pig

import "github.com/gorilla-go/pig/foundation/constant"

type IRouter interface {
	Route(path string, method constant.RequestMethod) (func(*Context), RouterParams, []IMiddleware)
	Url(routerName string, params map[string]any) string
}

type IMiddleware interface {
	Handle(*Context, func(*Context))
}

type IHttpErrorHandler interface {
	Handle(any, *Context)
}

type ICache[T any] interface {
	Get(key string) T
	Set(key string, value T, expire int64) bool
	Has(key string) bool
	Remove(key string) bool
}
