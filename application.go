package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"net"
	"net/http"
)

type Application struct {
	middleware []IMiddleware
	router     IRouter
}

func New() *Application {
	return &Application{
		middleware: []IMiddleware{},
	}
}

func (a *Application) Use(m ...IMiddleware) *Application {
	a.middleware = append(a.middleware, m...)
	return a
}

func (a *Application) Router(router IRouter) *Application {
	a.router = router
	return a
}

func (a *Application) Run(port ...int) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		NewKernel(a.router).Through(a.middleware).Handle(w, req)
	})

	err := http.ListenAndServe(
		fmt.Sprintf(
			"%s:%d",
			net.IPv4(0, 0, 0, 0).String(),
			foundation.DefaultParam(port, 8080),
		),
		nil,
	)
	panic(err)
}
