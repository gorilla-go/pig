package pig

import (
	"github.com/samber/do"
)

type Middleware struct {
}

func NewSysMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(context *Context, f func(*Context)) {
	// system middleware
	// error handler module
	do.Provide[ILogger](context.injector, func(injector *do.Injector) (ILogger, error) {
		return NewLogger(), nil
	})

	// error handler module
	do.Provide[IHttpErrorHandler](context.injector, func(injector *do.Injector) (IHttpErrorHandler, error) {
		return NewHttpErrorHandler(), nil
	})

	f(context)
}
