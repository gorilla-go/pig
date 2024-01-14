package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	router := pig.NewRouter()
	router.GET("/", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("post_id", "0")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	err := pig.New().Router(router).Start()
	if err != nil {
		panic(err)
	}
}
