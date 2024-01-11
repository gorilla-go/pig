package p_i_g

import (
	"github.com/samber/do"
)

type Logger interface {
	Info(message string)
	Debug(message string)
	Warning(message string)
	Fatal(message string)
}

type Router interface {
	Route(string) func(*do.Injector)
}

type Middleware interface {
	Handle(*do.Injector, func(*do.Injector))
}
