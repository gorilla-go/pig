package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/samber/do"
	"net/http"
)

func main() {
	router := pig.NewRouter()
	router.Map(map[string]func(*pig.Context){
		"/": func(ctx *pig.Context) {
			ctx.ResponseWriter().Write([]byte("Hello World!"))
		},
		"/:id": func(ctx *pig.Context) {
			fmt.Println(ctx.ParamVar()["id"].TrimString())
			fmt.Println(ctx.ParamVar()["q"].Bool())
		},
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
