package pig

type IRouter interface {
	Route(path string, method string) (func(*Context), RouterParams, []IMiddleware)
}

type IMiddleware interface {
	Handle(*Context, func(*Context))
}

type IHttpErrorHandler interface {
	Handle(any, *Context)
}
