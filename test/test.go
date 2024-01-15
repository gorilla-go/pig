package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello World!")
	})

	r.GET("/:id", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("id", "")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/test/<id:[a-z]+>", func(c *pig.Context) {
		c.Echo("ok")
	})

	r.GET("/upload", func(context *pig.Context) {
		context.Download(pig.NewFile("/Users/yehua/Downloads/202401/15/1746723873245892608.zip"), "test.zip")
	})

	pig.New().Router(r).Run()
}
