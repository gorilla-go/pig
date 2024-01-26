package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/di"
	"log"
	"net/http"
	"runtime/debug"
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
			errorHandler, err := di.Invoke[IHttpErrorHandler](k.context.container)
			if err != nil {
				log.Println(fmt.Sprintf("%s\n\r%s", errno, string(debug.Stack())))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			errorHandler.Handle(errno, k.context)
			return
		}
	}()

	if k.router == nil {
		panic("router unset.")
	}
	controllerAction, routerParams, cusMiddleware := k.router.Route(req.URL.Path, req.Method)
	if controllerAction == nil {
		controllerAction = func(context *Context) {
			context.Response().Raw().WriteHeader(http.StatusNotFound)
		}
	}

	if cusMiddleware != nil && len(cusMiddleware) > 0 {
		k.middleware = cusMiddleware
	}

	di.ProvideValue[*Context](k.context.container, k.context)
	di.ProvideValue[IRouter](k.context.container, k.router)
	di.ProvideValue[*Request](k.context.container, NewRequest(req, routerParams))
	di.ProvideValue[*Response](k.context.container, NewResponse(w, req))

	pipeline := NewPipeline[*Context]().Send(k.context)
	for _, middleware := range k.middleware {
		pipeline.Through(middleware.Handle)
	}
	pipeline.Then(controllerAction)
}
