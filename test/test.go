package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello World!")
	})

	r.POST("/upload", func(context *pig.Context) {
		fmt.Println("ok")
		path := context.FileVar()["file"].ArchiveMove("~/Downloads").FilePath
		context.Echo(path)
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
	r.Debug()
	pig.New().Router(r).Run()
}
