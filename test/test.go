package main

import (
	p_i_g "github.com/gorilla-go/p.i.g"
	"github.com/samber/do"
	"net/http"
)

func main() {
	router := p_i_g.NewRouter()
	router.Bind(map[string]func(*do.Injector){
		"/": func(injector *do.Injector) {
			response := do.MustInvoke[http.ResponseWriter](injector)
			response.WriteHeader(200)
			response.Write([]byte("Hello, World!"))
		},
	})

	err := p_i_g.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
