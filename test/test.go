package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()

	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello, World!")
	})

	r.GET("/:test", func(context *pig.Context) {
		context.Echo("Hello, World!2")
	})

	pig.New().Router(r).Run(8848)
}
