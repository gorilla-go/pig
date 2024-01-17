package main

import (
	"github.com/gorilla-go/pig"
)

type User struct {
	Name string
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Response().Echo("Hello World!")
	})

	pig.New().Router(r).Run(8848)
}
