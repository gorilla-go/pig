package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.POST("/", func(ctx *pig.Context) {

	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
