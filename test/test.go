package main

import (
	"github.com/gorilla-go/pig"
	"github.com/samber/do"
	"net/http"
)

func main() {
	router := pig.NewRouter()
	router.ANY("/", func(ctx *pig.Context) {
		ctx.ResponseWriter().Write([]byte("Hello World!"))
	})

	router.Miss(func(ctx *pig.Context) {
		response := do.MustInvoke[http.ResponseWriter](ctx.Injector())
		response.WriteHeader(404)
		response.Write([]byte("Not Found"))
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
