package main

import (
	"github.com/gorilla-go/pig"
	"github.com/samber/do"
	"net/http"
)

func main() {
	router := pig.NewRouter()
	router.Map(map[string]func(*do.Injector){
		"/": func(injector *do.Injector) {
			response := do.MustInvoke[http.ResponseWriter](injector)
			response.WriteHeader(200)
			response.Write([]byte("Hello, World!"))
		},
	})
	router.Miss(func(injector *do.Injector) {
		response := do.MustInvoke[http.ResponseWriter](injector)
		response.WriteHeader(404)
		response.Write([]byte("Not Found"))
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
