package p_i_g

import (
	"context"
	"github.com/samber/do"
	"log"
	"net/http"
)

type Kernel struct {
	context    context.Context
	injector   *do.Injector
	middleware []Middleware
}

func NewKernel() *Kernel {
	return &Kernel{
		injector: do.New(),
		context:  context.Background(),
	}
}

func (k *Kernel) Through(middleware []Middleware) *Kernel {
	k.middleware = middleware
	return k
}

func (k *Kernel) Handle(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errorHandler, err := do.Invoke[HttpErrorHandler](k.injector)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			errorHandler.Handle(err.(error), k.injector)
		}
	}()

	k.Inject(w, req)

	pipeline := NewPipeline().Send(k.injector)
	for _, middleware := range k.middleware {
		pipeline.Through(middleware.Handle)
	}

	_, err := do.Invoke[Logger](k.injector)
	if err != nil {
		panic(err)
	}

	router, err := do.Invoke[Router](k.injector)
	if err != nil {
		panic(err)
	}

	controllerAction := router.Route(req.URL.Path)
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
}
