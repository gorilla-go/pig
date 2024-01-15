package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/:id", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("id", "")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/test/<id:[a-z]+>", func(c *pig.Context) {
		c.Echo("ok")
	})

	pig.New().Router(r).Run()
}
