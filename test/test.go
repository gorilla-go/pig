package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Student struct {
	age int `query:"id" validate:"max=10,min=1,oneOf=2|3" msg:"无效的年龄参数"`
}

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
