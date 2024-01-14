package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"strings"
)

type Router struct {
	regRouteMap   *foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]
	missRoute     func(*Context)
	middlewareMap map[string][]IMiddleware
}

type RouterParams ReqParams

func NewRouter() *Router {
	return &Router{
		regRouteMap: &foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]{
			K: make([]string, 0),
			M: make(map[string]*foundation.LinkedHashMap[string, func(*Context)]),
		},
		middlewareMap: make(map[string][]IMiddleware),
	}
}

func (r *Router) addRoute(t string, path string, f func(*Context), middleware []IMiddleware) {
	t = strings.ToUpper(t)
	if r.regRouteMap.ContainsKey(path) == false {
		r.regRouteMap.Put(path, &foundation.LinkedHashMap[string, func(*Context)]{
			K: make([]string, 0),
			M: make(map[string]func(*Context)),
		})
	}
	r.regRouteMap.Get(path).Put(t, f)

	if len(middleware) > 0 {
		r.middlewareMap[r.ReqUniPath(path, t)] = middleware
	}
}

func (r *Router) ReqUniPath(path, method string) string {
	return fmt.Sprintf("%s://%s", method, path)
}

func (r *Router) GET(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("GET", path, f, middleware)
}

func (r *Router) POST(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("POST", path, f, middleware)
}

func (r *Router) PUT(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("PUT", path, f, middleware)
}

func (r *Router) DELETE(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("DELETE", path, f, middleware)
}

func (r *Router) PATCH(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("PATCH", path, f, middleware)
}

func (r *Router) OPTIONS(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("OPTIONS", path, f, middleware)
}

func (r *Router) HEAD(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("HEAD", path, f, middleware)
}

func (r *Router) CONNECT(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("CONNECT", path, f, middleware)
}

func (r *Router) TRACE(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("TRACE", path, f, middleware)
}

func (r *Router) ANY(path string, f func(*Context), middleware ...IMiddleware) {
	r.addRoute("ANY", path, f, middleware)
}

func (r *Router) Miss(f func(*Context)) *Router {
	r.missRoute = f
	return r
}

func (r *Router) Route(path string, requestMethod string) (func(*Context), RouterParams, []IMiddleware) {
	requestMethod = strings.ToUpper(requestMethod)
	var fn func(*Context) = nil
	routerParams := make(RouterParams)
	middlewares := make([]IMiddleware, 0)

	r.regRouteMap.ForEach(func(regexp string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		ok := methodMap.ContainsKey(requestMethod)
		if !ok {
			ok = methodMap.ContainsKey("ANY")
			if ok {
				fn = methodMap.Get("ANY")
			} else {
				return false
			}
		}

		regexpTrim := strings.Trim(regexp, "/")
		path = strings.Trim(path, "/")

		patternMode := strings.Contains(regexpTrim, ":")
		if !patternMode && regexpTrim == path {
			fn = methodMap.Get(requestMethod)
			if m, ok := r.middlewareMap[r.ReqUniPath(regexp, requestMethod)]; ok {
				middlewares = m
			}
			return false
		}

		if patternMode {
			regexpParts := strings.Split(regexpTrim, "/")
			pathParts := strings.Split(path, "/")
			if len(regexpParts) != len(pathParts) {
				return true
			}

			for i, part := range regexpParts {
				if len(path) > 1 && part[0] == ':' && len(pathParts[i]) > 0 {
					if strings.Contains(pathParts[i], ".") {
						pathParts[i] = (strings.Split(pathParts[i], "."))[0]
					}
					routerParams[part[1:]] = NewReqParamV([]string{pathParts[i]})
					continue
				}

				if part != pathParts[i] {
					return true
				}
			}
			fn = methodMap.Get(requestMethod)
			if m, ok := r.middlewareMap[r.ReqUniPath(regexp, requestMethod)]; ok {
				middlewares = m
			}
			return false
		}
		return true
	})

	if r.missRoute != nil {
		return r.missRoute, nil, middlewares
	}
	return fn, routerParams, middlewares
}
