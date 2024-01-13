package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.PUT("/", func(ctx *pig.Context) {

		ctx.PostVar()["id"].TrimString()
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
