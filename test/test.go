package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/di"
)

func main() {
	r := pig.NewRouter()
	r.GET("/:id", func(ctx *pig.Context) {
		router := di.MustInvoke[pig.IRouter](ctx.Container())
		ctx.Response().Text(router.Url("index", map[string]any{
			"id":   1,
			"name": "pig",
		}))
	}).Name("index")
	pig.New().Router(r).Run(8081)
}
