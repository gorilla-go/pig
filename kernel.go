package pig

import (
	"github.com/gorilla-go/pig/foundation"
	"github.com/samber/do"
	"net/http"
)

type Kernel struct {
	middleware []IMiddleware
	router     IRouter
	context    *Context
}

func NewKernel(r IRouter) *Kernel {
	return &Kernel{
		router:  r,
		context: NewContext(),
	}
}

func (k *Kernel) Through(middleware []IMiddleware) *Kernel {
	k.middleware = middleware
	return k
}

func (k *Kernel) Handle(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if errno := recover(); errno != nil {
			errorHandler, err := do.Invoke[IHttpErrorHandler](k.context.Injector())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			errorHandler.Handle(errno, k.context)
			return
		}
	}()

	k.Inject(w, req)

	controllerAction, routerParams, cusMiddleware := k.router.Route(req.URL.Path, req.Method)
	if controllerAction == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if cusMiddleware != nil && len(cusMiddleware) > 0 {
		k.middleware = cusMiddleware
	}

	if routerParams != nil && len(routerParams) > 0 {
		foundation.ProvideValue[RouterParams](k.context.Injector(), routerParams)
	}

	pipeline := NewPipeline[*Context]().Send(k.context)
	for _, middleware := range k.middleware {
		pipeline.Through(middleware.Handle)
	}
	pipeline.Then(controllerAction)
}

func (k *Kernel) Inject(w http.ResponseWriter, req *http.Request) {
	foundation.Provide(k.context.Injector(), req)
	foundation.Provide[http.ResponseWriter](k.context.Injector(), w)
	foundation.Provide(k.context.Injector(), k.context)
	foundation.Provide[IRouter](k.context.Injector(), k.router)
}
