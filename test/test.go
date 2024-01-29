package main

import (
	"github.com/gorilla-go/pig"
	"os"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		panic("target error")
	})

	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	r.Static("/", getwd+"/test/")

	pig.New().Router(r).Run(8081)
}
