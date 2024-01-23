package pig

import (
	"github.com/gorilla-go/pig/di"
)

type Context struct {
	container *di.Container
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
