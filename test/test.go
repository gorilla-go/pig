package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/favicon.ico", func(context *pig.Context) {

	})

	r.GET("/test*111/<exp: \\d+>", func(context *pig.Context) {
		i := context.Request().ParamVar().MustInt("exp")
		fmt.Println(i)
		fmt.Println("ok")
	})

	pig.New().Router(r).Run()
}
