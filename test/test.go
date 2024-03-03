package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/di"
)

func main() {
	r := pig.NewRouter()
	r.GET("/upload", func(ctx *pig.Context) {
		a := &struct {
			Name   string      `json:"name"`
			router pig.IRouter `di:""`
		}{}
		di.Autowire(ctx.Container(), a)
		fmt.Println(a.router)
	})
	pig.New().Router(r).Run(8081)
}
