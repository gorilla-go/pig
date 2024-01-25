package main

import (
	"github.com/gorilla-go/pig"
)

type CorsMiddleware struct {
}

func (c CorsMiddleware) Handle(context *pig.Context, f func(*pig.Context)) {
	if context.Request().IsOption() {
		context.Response().Raw().Header().Set("Access-Control-Allow-Origin", "*")
		context.Response().Raw().Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
		context.Response().Raw().Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
		context.Response().Raw().WriteHeader(204)
		return
	}
	f(context)
}

func main() {
	r := pig.NewRouter()
	r.ANY("/", func(context *pig.Context) {
		context.Response().Text("hello world")
	}, &CorsMiddleware{})
	pig.New().Router(r).Run(8081)
}
