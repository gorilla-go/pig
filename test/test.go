package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/:name/<id:\\d+>", func(context *pig.Context) {
		fmt.Println(
			context.Request().ParamVar().TrimString("name"),
			context.Request().ParamVar().Int("id"),
		)
	})
	pig.New().Router(r).Run(8081)
}
