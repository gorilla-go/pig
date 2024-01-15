package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/test/ss", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("post_id", "")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/:id", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("id", "")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/test/<id:[a-z]+>", func(context *pig.Context) {
		context.Json(map[string]interface{}{
			"test": context.ParamVar().TrimString("id", ""),
		})
	})

	err := pig.New().Router(r).Run()
	if err != nil {
		panic(err)
	}
}
