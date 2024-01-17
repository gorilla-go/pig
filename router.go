package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"regexp"
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

func (r *Router) Debug() {
	fmt.Println("---------------------- router debug ----------------------")
	r.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		methodMap.ForEach(func(method string, fn func(*Context)) bool {
			fmt.Println(fmt.Sprintf("%s %s", method, uri))
			return true
		})
		return true
	})
}

func (r *Router) Route(path string, requestMethod string) (func(*Context), RouterParams, []IMiddleware) {
	requestMethod = strings.ToUpper(requestMethod)
	var fn func(*Context) = nil
	routerParams := make(RouterParams)
	middlewares := make([]IMiddleware, 0)

	r.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		ok := methodMap.ContainsKey(requestMethod)
		if !ok {
			ok = methodMap.ContainsKey("ANY")
			if ok {
				fn = methodMap.Get("ANY")
			} else {
				return true
			}
		}

		originUri := uri
		uri = strings.Trim(uri, "/")
		path = strings.Trim(path, "/")

		patternMode := strings.Contains(uri, ":") ||
			strings.Contains(uri, "<") ||
			strings.Contains(uri, ">")
		if !patternMode && uri == path {
			fn = methodMap.Get(requestMethod)
			if m, ok := r.middlewareMap[r.ReqUniPath(originUri, requestMethod)]; ok {
				middlewares = m
			}
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

					routerParams[pathFormatArr[0]] = NewReqParamV([]string{pathParts[i]})
					continue
				}

				if len(part) > 1 && part[0] == ':' {
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
			if m, ok := r.middlewareMap[r.ReqUniPath(originUri, requestMethod)]; ok {
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
