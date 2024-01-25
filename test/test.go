package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/di"
)

type HttpErrorHandler struct {
}

func (h *HttpErrorHandler) Handle(a any, context *pig.Context) {
	fmt.Println(a)
}

type ErrorCollection struct {
}

func (m *ErrorCollection) Handle(context *pig.Context, f func(*pig.Context)) {
	di.ProvideValue[pig.IHttpErrorHandler](context.Container(), &HttpErrorHandler{})
	f(context)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		panic("target error")
	})

	pig.New().Use(&ErrorCollection{}).Router(r).Run(8081)
}
