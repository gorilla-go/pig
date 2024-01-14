package main

import (
	"github.com/gorilla-go/pig"
	"github.com/gorilla-go/pig/foundation"
)

type Middleware struct {
}

func NewSysMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(context *pig.Context, f func(*pig.Context)) {
	foundation.Provide[pig.ILogger](context.Injector(), pig.NewLogger())
	foundation.Provide[pig.IHttpErrorHandler](context.Injector(), pig.NewHttpErrorHandler())
	f(context)
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

	err := pig.New().Use(NewSysMiddleware()).Router(r).Start()
	if err != nil {
		panic(err)
	}
}
