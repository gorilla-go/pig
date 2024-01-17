package main

import (
	"github.com/gorilla-go/pig/di"
)

type User struct {
	Name string
}

func main() {
	container := di.New()
	di.ProvideValue(container, &User{Name: "pig"})
	di.MustInvoke[*User](container)
}
