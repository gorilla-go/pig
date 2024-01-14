package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("post_id")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/:id", func(context *pig.Context) {
		context.Json(map[string]interface{}{
			"id": context.ParamVar().TrimString("id"),
		})
	})

	err := pig.New().Router(r).Start()
	if err != nil {
		panic(err)
	}
}
