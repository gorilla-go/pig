package pig

import (
	"github.com/samber/do"
	"net/http"
)

type Kernel struct {
	injector   *do.Injector
	middleware []IMiddleware
	router     IRouter
}

func NewKernel(r IRouter) *Kernel {
	return &Kernel{
		injector: do.New(),
		router:   r,
	}
}

func (k *Kernel) Through(middleware []IMiddleware) *Kernel {
	k.middleware = middleware
	return k
}

func (k *Kernel) Handle(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if errno := recover(); errno != nil {
			errorHandler, err := do.Invoke[IHttpErrorHandler](k.injector)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			errorHandler.Handle(errno.(error), k.injector)
			return
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
	do.Provide(k.injector, func(*do.Injector) (IRouter, error) {
		return k.router, nil
	})
}
