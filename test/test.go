package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.POST("/upload", func(context *pig.Context) {
		filePath := context.FileVar()["file"].FilePath
		context.Echo(filePath)
	})

	r.GET("/download", func(context *pig.Context) {
		context.Download(
			pig.NewFile("/your/file/path.jpg"),
			"filename.jpg",
		)
	})

	pig.New().Router(r).Run(8088)
}
