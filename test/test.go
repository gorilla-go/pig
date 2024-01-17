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
	r.GET("/", func(context *pig.Context) {
		context.Response().Echo("hello world")
	}, &Middleware{})

	pig.New().Router(r).Run(8848)
}
