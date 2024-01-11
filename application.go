package p_i_g

import (
	"fmt"
	"net"
	"net/http"
)

type Application struct {
	port       int
	address    net.IP
	middleware []Middleware
}

func New() *Application {
	return &Application{
		port:    8080,
		address: net.IPv4(127, 0, 0, 1),
	}
}

func (a *Application) Use(m Middleware) {
	a.middleware = append(a.middleware, m)
}

func (a *Application) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		NewKernel().Through(a.middleware).Handle(w, req)
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
