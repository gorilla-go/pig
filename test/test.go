package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/samber/do"
	"net/http"
)

type SessionMiddle struct {
}

func (s *SessionMiddle) Handle(injector *do.Injector, next func(*do.Injector)) {
	fmt.Println("session middle")
	next(injector)
}

type FilterMiddle struct {
}

func (f *FilterMiddle) Handle(injector *do.Injector, next func(*do.Injector)) {
	next(injector)
	fmt.Println("filter middle")
}

func main() {
	router := pig.NewRouter()
	router.Map(map[string]func(*do.Injector){
		"/": func(injector *do.Injector) {
			fmt.Println("controller action")
			response := do.MustInvoke[http.ResponseWriter](injector)
			response.WriteHeader(200)
			response.Write([]byte("Hello, World!"))
		},
		"/:id": func(injector *do.Injector) {
			fmt.Println("controller action with params")
			response := do.MustInvoke[http.ResponseWriter](injector)
			response.WriteHeader(200)
			response.Write([]byte("Hello, World!"))
		},
	})

	err := pig.New().Use(&SessionMiddle{}, &FilterMiddle{}).Router(router).Start()
	if err != nil {
		panic(err)
	}
}
