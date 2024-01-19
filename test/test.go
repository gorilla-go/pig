package main

import (
	"fmt"
	"github.com/gorilla-go/pig/di"
)

type Form struct {
	id   int    `query:"id" form:"id"`
	name string `query:"name" form:"name"`
}

type IncludeForm struct {
	form *Form `di:""`
}

func main() {
	container := di.New()
	di.ProvideLazy(container, func(*di.Container) (*Form, error) {
		return &Form{
			id:   1,
			name: "test",
		}, nil
	})

	inject := di.Inject(container, &IncludeForm{})
	fmt.Println(inject.form.id)
}
