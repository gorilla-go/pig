package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.GET("/", func(ctx *pig.Context) {
		postId := ctx.ParamVar()["post_id"].Int()
		ctx.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
