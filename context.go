package pig

import (
	"github.com/gorilla-go/pig/di"
)

type Context struct {
	container *di.Container
	config    IConfig
}

func NewContext() *Context {
	return &Context{
		container: di.New(),
	}
}

func (c *Context) Container() *di.Container {
	return c.container
}

func (c *Context) Request() *Request {
	return di.MustInvoke[*Request](c.container)
}

func (c *Context) Response() *Response {
	return di.MustInvoke[*Response](c.container)
}

func (c *Context) Logger() ILogger {
	return di.MustInvoke[ILogger](c.container)
}

func (c *Context) Config(s string) any {
	config, err := di.MustInvoke[IConfig](c.container).Get(s)
	if err != nil {
		panic(err)
	}
	return config
}
