package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	r := pig.NewRouter()
	r.POST("/", func(context *pig.Context) {
		u := &User{}
		context.Request().JsonBind(u)
		fmt.Println(u.Name, u.Age)
	})

	pig.New().Router(r).Run(8081)
}
