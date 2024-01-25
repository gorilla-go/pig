package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"regexp"
	"strings"
)

type Router struct {
	group           string
	groupMiddleware []IMiddleware
	regRouteMap     *foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]
	missRoute       func(*Context)
	middlewareMap   map[string][]IMiddleware
}

type RouterParams foundation.ReqParams

func NewRouter() *Router {
	return &Router{
		group:           "",
		groupMiddleware: make([]IMiddleware, 0),
		regRouteMap: &foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]{
			K: make([]string, 0),
			M: make(map[string]*foundation.LinkedHashMap[string, func(*Context)]),
		},
		middlewareMap: make(map[string][]IMiddleware),
	}
}

func (r *Router) addRoute(t string, path string, f func(*Context), middleware []IMiddleware) {
	if len(path) == 0 {
		panic("invalid router path.")
	}

	t = strings.ToUpper(t)
	if r.group != "" {
		if path[0] != '/' {
			path = fmt.Sprintf("/%s", path)
		}

		if r.group[len(r.group)-1] == '/' {
			r.group = r.group[:len(r.group)-1]
		}

		path = fmt.Sprintf("%s%s", r.group, path)
	}

	if r.regRouteMap.ContainsKey(path) == false {
		r.regRouteMap.Put(path, &foundation.LinkedHashMap[string, func(*Context)]{
			K: make([]string, 0),
			M: make(map[string]func(*Context)),
		})
	}
	r.regRouteMap.Get(path).Put(t, f)

	if len(middleware) > 0 {
		r.middlewareMap[r.ReqUniPath(path, t)] = middleware
		return
	}

	if len(r.groupMiddleware) > 0 {
		r.middlewareMap[r.ReqUniPath(path, t)] = r.groupMiddleware
		return
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

func (r *Router) Group(path string, f func(r *Router), middleware ...IMiddleware) {
	if len(path) == 0 {
		panic("invalid group router path.")
	}
	router := &Router{
		group:           path,
		groupMiddleware: middleware,
		regRouteMap: &foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]{
			K: make([]string, 0),
			M: make(map[string]*foundation.LinkedHashMap[string, func(*Context)]),
		},
		middlewareMap: make(map[string][]IMiddleware),
	}
	f(router)

	router.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		methodMap.ForEach(func(method string, fn func(*Context)) bool {
			r.addRoute(method, uri, fn, router.middlewareMap[r.ReqUniPath(uri, method)])
			return true
		})
		return true
	})
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

	r.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		originUri := uri
		uri = strings.Trim(uri, "/")
		path = strings.Trim(path, "/")

		patternMode := strings.Contains(uri, ":") ||
			strings.Contains(uri, "<") ||
			strings.Contains(uri, ">")
		if !patternMode && uri == path {
			fn, middlewares = r.findFn(methodMap, originUri, requestMethod)
			return false
		}

		if patternMode {
			uriParts := strings.Split(uri, "/")
			pathParts := strings.Split(path, "/")
			if len(uriParts) != len(pathParts) {
				return true
			}

			for i, part := range uriParts {
				if len(pathParts[i]) == 0 {
					return true
				}

				if len(part) > 4 && part[0] == '<' && part[len(part)-1] == '>' {
					if strings.Contains(pathParts[i], ".") {
						pathParts[i] = (strings.Split(pathParts[i], "."))[0]
					}

					pathFormatArr := strings.SplitN(part[1:len(part)-1], ":", 2)
					if len(pathFormatArr) < 2 {
						return true
					}

					regexpStr := strings.TrimSpace(pathFormatArr[1])
					match, err := regexp.Match(regexpStr, []byte(pathParts[i]))
					if err != nil || !match {
						return true
					}

					routerParams[pathFormatArr[0]] = foundation.NewReqParamV([]string{pathParts[i]})
					continue
				}

				if len(part) > 1 && part[0] == ':' {
					if strings.Contains(pathParts[i], ".") {
						pathParts[i] = (strings.Split(pathParts[i], "."))[0]
					}
					routerParams[part[1:]] = foundation.NewReqParamV([]string{pathParts[i]})
					continue
				}

				if part != pathParts[i] {
					return true
				}
			}

			fn, middlewares = r.findFn(methodMap, originUri, requestMethod)
			return false
		}
		return true
	})

	if r.missRoute != nil && fn == nil {
		return r.missRoute, nil, middlewares
	}
	return fn, routerParams, middlewares
}

func (r *Router) findFn(
	methodMap *foundation.LinkedHashMap[string, func(*Context)],
	originUri string,
	requestMethod string,
) (fn func(*Context), middlewares []IMiddleware) {
	if m, ok := r.middlewareMap[r.ReqUniPath(originUri, requestMethod)]; ok {
		fn = methodMap.Get(requestMethod)
		middlewares = m
	}

	if fn == nil {
		if m, ok := r.middlewareMap[r.ReqUniPath(originUri, "ANY")]; ok {
			fn = methodMap.Get("ANY")
			middlewares = m
		}
	}
	return
}
