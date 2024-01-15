package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"path/filepath"
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

	r.POST("/upload", func(context *pig.Context) {
		path := context.FileVar()["file"].FilePath
		fmt.Println(path)
		context.Download(context.FileVar()["file"], filepath.Base(path))
	})

	pig.New().Router(r).Run(8088)
}
