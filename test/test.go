package main

import (
	"fmt"
	"github.com/gorilla-go/pig/di"
)

type Form struct {
	id   int    `query:"id" form:"id"`
	name string `query:"name" form:"name"`
}

func (f *Form) Name() {
	fmt.Println(f.name)
}

type IForm interface {
	Name()
}

type IncludeForm struct {
	form IForm `di:""`
}

func main() {
	container := di.New()
	di.ProvideLazy[IForm](container, func(c *di.Container) (any, error) {
		return &Form{
			id:   1,
			name: "test",
		}, nil
	})

	inject := di.Inject(container, &IncludeForm{})
	inject.form.Name()
}
