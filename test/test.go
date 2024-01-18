package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Response().Echo("Hello World!")
	})

	pig.New().Pprof().Router(r).Run(8848)
}
