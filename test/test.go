package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/samber/do/v2"
)

type HttpErrorHandler struct {
}

func (h *HttpErrorHandler) Handle(a any, context *pig.Context) {
	fmt.Println("error targeted")
	context.Echo("500", 500)
}

type Middleware struct {
}

func (*Middleware) Handle(c *pig.Context, next func(*pig.Context)) {
	do.ProvideValue[pig.IHttpErrorHandler](c.Injector(), &HttpErrorHandler{})
	next(c)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		panic("error")
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8088)
}
