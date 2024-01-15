package pig

type ILogger interface {
	Info(message string, c *Context)
	Debug(message string, c *Context)
	Warning(message string, c *Context)
	Fatal(message string, c *Context)
}

type IRouter interface {
	Route(path string, method string) (func(*Context), RouterParams, []IMiddleware)
}

type IMiddleware interface {
	Handle(*Context, func(*Context))
}

type IHttpErrorHandler interface {
	Handle(any, *Context)
}

type ISession[T any] interface {
	Get(string) T
	Set(string, T)
}

type IConfig interface {
	Get(string) (any, error)
}
