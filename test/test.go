package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/foundation"
)

type Middleware struct {
}

func (m *Middleware) Handle(c *pig.Context, f func(*pig.Context)) {
	foundation.Provide[pig.ILogger](c.Injector(), pig.NewLogger())
	foundation.Provide[pig.IHttpErrorHandler](c.Injector(), pig.NewHttpErrorHandler())
	f(c)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		postId := c.ParamVar().TrimString("post_id", "")
		c.Json(map[string]interface{}{
			"post_id": postId,
		})
	})

	r.GET("/:id", func(c *pig.Context) {
		c.Json(map[string]interface{}{
			"id": c.ParamVar().Lmt("id", []string{"1", "2", "3"}, "0").Int(),
		})
	})

	err := pig.New().Use(&Middleware{}).Router(r).Start()
	if err != nil {
		panic(err)
	}
}
