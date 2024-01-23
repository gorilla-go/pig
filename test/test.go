package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/validate"
)

type Form struct {
	Username  string `query:"username" validate:"required" msg:"用户名不能为空"`
	Pass      string `query:"pass" validate:"required,minLen=6" msg:"密码不能为空"`
	PassAgain string `query:"pass_again" validate:"sameAs=Pass" msg:"两次密码不一致"`
	Email     string `query:"email" validate:"required,email" msg:"邮箱格式不正确"`
	Phone     string `query:"phone" validate:"cnPhone" msg:"手机号格式不正确"`
}

func main() {
	v := validate.New(map[string]validate.Checker{
		"required": validate.Required,
		"email":    validate.Email,
		"len":      validate.Len,
		"oneOf":    validate.OneOf,
		"sameAs":   validate.SameAs,
		"minLen":   validate.MinLen,
		"cnPhone":  validate.CnPhone,
	})

	router := pig.NewRouter()
	router.GET("/", func(ctx *pig.Context) {
		form := &Form{}
		ctx.Request().Bind(form)
		if err := v.Validate(form); err != nil {
			ctx.Response().Text(err.Error())
			return
		}
		ctx.Response().Text(ctx.Request().ParamVar().TrimString("username"))
	})

	pig.New().Router(router).Run(8081)
}
