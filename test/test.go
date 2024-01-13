package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.GET("/", func(ctx *pig.Context) {
		panic("test panic")
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
