package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Student struct {
	age int `query:"id"`
}

func main() {
	r := pig.NewRouter()
	r.GET("/favicon.ico", func(context *pig.Context) {

	})

	r.GET("/<id:\\d+>", func(context *pig.Context) {
		s := &Student{}
		context.Request().Bind(s)
		fmt.Println(s.age)
	})

	pig.New().Router(r).Run()
}
