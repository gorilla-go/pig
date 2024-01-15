package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/redirect", func(context *pig.Context) {
		context.Redirect("/redirected", 302)
	})

	pig.New().Router(r).Run(8088)
}
