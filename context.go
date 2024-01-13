package pig

import (
	"github.com/samber/do"
	"net/http"
)

type Context struct {
	injector *do.Injector
}

func NewContext() *Context {
	return &Context{
		injector: do.New(),
	}
}

func (c *Context) Injector() *do.Injector {
	return c.injector
}

func (c *Context) RouterParams() RouterParams {
	routerParams, err := do.Invoke[RouterParams](c.injector)
	if err != nil {
		return nil
	}

	return routerParams
}

func (c *Context) Params() map[string]*ReqParamV {
	reqParamMap := make(map[string]*ReqParamV)

	request, err := do.Invoke[*http.Request](c.Injector())
	if err == nil {
		for n, v := range request.URL.Query() {
			reqParamMap[n] = NewReqParamV(v)
		}
	}

	routerParams := c.RouterParams()
	if routerParams != nil {
		for n, v := range routerParams {
			reqParamMap[n] = v
		}
	}
	return reqParamMap
}
