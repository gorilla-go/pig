package main

import (
	"github.com/gorilla-go/pig"
	"time"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		time.Sleep(time.Second * 10)
	})

	pig.New().Router(r).Run(8848)
}
