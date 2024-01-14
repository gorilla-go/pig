package pig

type ILogger interface {
	Info(message string)
	Debug(message string)
	Warning(message string)
	Fatal(message string)
}

type IRouter interface {
	Route(path string, method string) (func(*Context), RouterParams)
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
