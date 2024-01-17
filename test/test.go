package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/:id", func(context *pig.Context) {
		context.Response().Json(map[string]interface{}{
			"code": 0,
		})
	})

	pig.New().Router(r).Run(8848)
}
