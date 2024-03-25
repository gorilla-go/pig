package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/validate"
)

func main() {
	r := pig.NewRouter()
	r.GET("/favicon.ico", func(context *pig.Context) {

	})

	r.GET("/<id:\\d+>", func(context *pig.Context) {
		v := validate.New()
		ok := v.CheckVar(
			context.Request().ParamVar().MustInt("id"),
			"min=1,max=10",
		)
		fmt.Println(ok)
	})

	pig.New().Router(r).Run()
}
