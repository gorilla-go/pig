package pig

import (
	"context"
	"github.com/samber/do"
	"net/http"
)

type Kernel struct {
	context      context.Context
	injector     *do.Injector
	middleware   []IMiddleware
	router       IRouter
	errorHandler IHttpErrorHandler
}

func NewKernel(r IRouter) *Kernel {
	if r == nil {
		r = NewRouter()
	}

	return &Kernel{
		injector: do.New(),
		context:  context.Background(),
		router:   r,
	}
}

func (k *Kernel) Through(middleware []IMiddleware) *Kernel {
	k.middleware = middleware
	return k
}

func (k *Kernel) ErrorHandler(h IHttpErrorHandler) *Kernel {
	k.errorHandler = h
	return k
}

func (k *Kernel) Handle(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if errno := recover(); errno != nil {
			errorHandler, err := do.Invoke[IHttpErrorHandler](k.injector)
			if err != nil {
				if k.errorHandler != nil {
					k.errorHandler.Handle(errno.(error), k.injector)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			errorHandler.Handle(errno.(error), k.injector)
		}
	}()

	controllerAction, routerParams := k.router.Route(req.URL.Path)
	if controllerAction == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if routerParams != nil {
		do.ProvideValue(k.injector, routerParams)
	}

	k.Inject(w, req)

	pipeline := NewPipeline().Send(k.injector)
	for _, middleware := range k.middleware {
		pipeline.Through(middleware.Handle)
	}
	pipeline.Then(controllerAction)
}

func (k *Kernel) Inject(w http.ResponseWriter, req *http.Request) {
	do.Provide(k.injector, func(*do.Injector) (http.ResponseWriter, error) {
		return w, nil
	})
	do.Provide(k.injector, func(*do.Injector) (*http.Request, error) {
		return req, nil
	})
	do.Provide(k.injector, func(*do.Injector) (context.Context, error) {
		return k.context, nil
	})
	do.Provide(k.injector, func(*do.Injector) (IRouter, error) {
		return k.router, nil
	})
}
