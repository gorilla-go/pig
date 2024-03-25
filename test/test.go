package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/validate"
)

type Student struct {
	age int `query:"id" validate:"max=10,min=1,oneOf=2|3" msg:"无效的年龄参数"`
}

func main() {
	r := pig.NewRouter()
	r.GET("/favicon.ico", func(context *pig.Context) {

	})

	r.GET("/<id:\\d+>", func(context *pig.Context) {
		s := &Student{}
		context.Request().Bind(s)
		v := validate.New()
		err := v.CheckStruct(s)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s.age)
	})

	pig.New().Router(r).Run()
}
