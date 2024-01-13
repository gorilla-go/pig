package pig

type ILogger interface {
	Info(message string)
	Debug(message string)
	Warning(message string)
	Fatal(message string)
}

type IRouter interface {
	Route(string) (func(*Context), RouterParams)
}

type IMiddleware interface {
	Handle(*Context, func(*Context))
}

type IHttpErrorHandler interface {
	Handle(error, *Context)
}

type ISession[T any] interface {
	Get(string) T
	Set(string, T)
}
