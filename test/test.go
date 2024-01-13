package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"io/ioutil"
	"net/http"
)

func main() {
	router := pig.NewRouter()
	router.POST("/", func(ctx *pig.Context) {
		fmt.Println("--- Request Headers ---")
		for key, values := range ctx.Request().Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		bodyBytes, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			return
		}

		fmt.Println("--- Request Body ---")
		fmt.Println(string(bodyBytes))

		// 处理请求...

		ctx.ResponseWriter().WriteHeader(http.StatusOK)
		ctx.ResponseWriter().Write([]byte("Hello, world!"))
	})

	router.POST("/ok", func(context *pig.Context) {
		fmt.Println(context.PostVar())
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
