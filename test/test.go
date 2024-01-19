package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Middleware1 struct {
}

func (m *Middleware1) Handle(context *pig.Context, f func(*pig.Context)) {
	//TODO implement me
	fmt.Println("middleware1")
	f(context)
}

type Middleware2 struct {
}

func (m *Middleware2) Handle(context *pig.Context, f func(*pig.Context)) {
	//TODO implement me
	fmt.Println("middleware2")
	f(context)
}

type Middleware3 struct {
}

func (m *Middleware3) Handle(context *pig.Context, f func(*pig.Context)) {
	//TODO implement me
	fmt.Println("middleware3")
	f(context)
}

func main() {
	router := pig.NewRouter()
	router.Group("/", func(r *pig.Router) {
		r.GET("test", func(c *pig.Context) {
			fmt.Println("ok")
		})
	})

	pig.New().Use(&Middleware3{}).Router(router).Run(8081)
}
