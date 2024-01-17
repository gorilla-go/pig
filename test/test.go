package main

import (
	"github.com/gorilla-go/pig"
)

type Middleware struct {
}

func (m *Middleware) Handle(context *pig.Context, f func(*pig.Context)) {
	context.Response().Echo("error")
}

func main() {
	r := pig.NewRouter()
	r.GET("/:id", func(context *pig.Context) {
		context.Response().Echo("hello world")
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8848)
}
