package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Form struct {
	id   int    `query:"id" form:"id"`
	name string `query:"name" form:"name"`
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		form := &Form{}
		c.Request().Bind(form)
		fmt.Println(form)
		c.Response().Echo("Hello, World!")
	})

	pig.New().Router(r).Run(8081)
}
