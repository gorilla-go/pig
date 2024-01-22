package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/validate"
)

type Form struct {
	Username string `json:"username" query:"username" validate:"required,email" msg:"username is required"`
}

func main() {
	v := validate.New(map[string]validate.Checker{
		"required": validate.Required,
		"email":    validate.Email,
	})
	router := pig.NewRouter()
	router.GET("/", func(ctx *pig.Context) {
		form := &Form{}
		ctx.Request().Bind(form)
		err := v.Validate(form)
		if err != nil {
			ctx.Response().Text(err.Error())
			return
		}
		ctx.Response().Text(ctx.Request().ParamVar().TrimString("username"))
	})

	pig.New().Router(router).Run(8081)
}
