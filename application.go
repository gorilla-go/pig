package p_i_g

import (
	"fmt"
	"net"
	"net/http"
)

type Application struct {
	port       int
	address    net.IP
	middleware []IMiddleware
	router     IRouter
}

func New() *Application {
	return &Application{
		port:    8080,
		address: net.IPv4(0, 0, 0, 0),
	}
}

func (a *Application) Use(m ...IMiddleware) {
	a.middleware = append(a.middleware, m...)
}

func (a *Application) Router(router IRouter) *Application {
	a.router = router
	return a
}

func (a *Application) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		NewKernel(a.router).Through(a.middleware).Handle(w, req)
	})

	err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", a.address.String(), a.port),
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
