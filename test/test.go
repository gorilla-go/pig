package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/di"
)

func main() {
	r := pig.NewRouter()
	r.GET("/<id:\\d+>", func(ctx *pig.Context) {
		router := di.MustInvoke[pig.IRouter](ctx.Container())
		ctx.Request().ParamVar().Raw().ToString()
		ctx.Response().Text(router.Url("index", map[string]any{
			"id":   ctx.Request().ParamVar().MustInt("id"),
			"name": "pig",
		}))
	}).Name("index")
	pig.New().Router(r).Run(8081)
}
