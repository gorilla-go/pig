package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.GET("/", func(c *pig.Context) {
		postId := c.ParamVar()["post_id"].Int()
		_ = c.FileVar()["ok"]
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	router.GET("/:id", func(context *pig.Context) {
		postId := context.ParamVar()["id"].Int()
		context.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
