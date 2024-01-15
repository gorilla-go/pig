package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()

	// 上传
	r.POST("/upload", func(context *pig.Context) {
		filePath := context.FileVar()["file"].FilePath
		context.Echo(filePath)
	})

	// 归档存储
	r.POST("/upload/archive", func(context *pig.Context) {
		file := context.FileVar()["file"]
		file = file.ArchiveMove("/your/dest/dir")
		context.Echo(file.FilePath)
	})

	// 移动文件
	r.POST("/upload/rename", func(context *pig.Context) {
		file := context.FileVar()["file"]
		file = file.Move("/your/dest/file.jpg")
		context.Echo(file.FilePath)
	})

	pig.New().Router(r).Run(8088)
}
