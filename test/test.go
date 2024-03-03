package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"os"
	"path/filepath"
)

func main() {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}

	r := pig.NewRouter()
	r.Static("/static", filepath.Clean(getwd+"/di"))
	r.POST("/upload", func(ctx *pig.Context) {
		file := ctx.Request().FileVar().MustFile("file")
		fmt.Println(file.FilePath)
	})
	pig.New().Router(r).Run(8081)
}
