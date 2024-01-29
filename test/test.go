package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"os"
)

type m struct {
}

func (m *m) Handle(context *pig.Context, f func(*pig.Context)) {
	fmt.Println("ok2")
	f(context)
}

type m2 struct {
}

func (m *m2) Handle(context *pig.Context, f func(*pig.Context)) {
	fmt.Println("ok3")
	f(context)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		panic("target error")
	}, &m2{})

	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	r.Static("/static/", getwd+"/test/")

	pig.New().Use(&m{}).Router(r).Run(8081)
}
