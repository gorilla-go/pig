package pig

import (
	"github.com/samber/do"
)

type ILogger interface {
	Info(message string)
	Debug(message string)
	Warning(message string)
	Fatal(message string)
}

type IRouter interface {
	Route(string) (func(*do.Injector), RouterParams)
}

type IMiddleware interface {
	Handle(*do.Injector, func(*do.Injector))
}

type IHttpErrorHandler interface {
	Handle(error, *do.Injector)
}
