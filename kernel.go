package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/di"
	"github.com/gorilla-go/pig/foundation/constant"
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
			errorHandler, err := di.Invoke[IHttpErrorHandler](k.context.Container())
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
	controllerAction, routerParams, cusMiddleware := k.router.Route(req.URL.Path, constant.RequestMethod(req.Method))
	if controllerAction == nil {
		controllerAction = func(context *Context) {
			context.Response().Code(http.StatusNotFound)
		}
	}

	if cusMiddleware != nil {
		k.middleware = cusMiddleware
	}

	container := k.context.Container()
	di.ProvideValue[*Context](container, k.context)
	di.ProvideValue[IRouter](container, k.router)
	di.ProvideValue[*Request](container, NewRequest(req, routerParams))
	di.ProvideValue[*Response](container, NewResponse(w, req))

	pipeline := NewPipeline[*Context]().Send(k.context)
	for _, middleware := range k.middleware {
		pipeline.Through(middleware.Handle)
	}
	pipeline.Then(controllerAction)
}
