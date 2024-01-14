package pig

import (
	"github.com/gorilla-go/pig/foundation"
	"strings"
)

type Router struct {
	regRouteMap *foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]
	missRoute   func(*Context)
}

type RouterParams ReqParams

func NewRouter() *Router {
	return &Router{
		regRouteMap: &foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]{
			K: make([]string, 0),
			M: make(map[string]*foundation.LinkedHashMap[string, func(*Context)]),
		},
	}
}

func (r *Router) addRoute(t string, path string, f func(*Context)) {
	t = strings.ToUpper(t)
	if r.regRouteMap.ContainsKey(path) == false {
		r.regRouteMap.Put(path, &foundation.LinkedHashMap[string, func(*Context)]{
			K: make([]string, 0),
			M: make(map[string]func(*Context)),
		})
	}
	r.regRouteMap.Get(path).Put(t, f)
}

func (r *Router) GET(path string, f func(*Context)) {
	r.addRoute("GET", path, f)
}

func (r *Router) POST(path string, f func(*Context)) {
	r.addRoute("POST", path, f)
}

func (r *Router) PUT(path string, f func(*Context)) {
	r.addRoute("PUT", path, f)
}

func (r *Router) DELETE(path string, f func(*Context)) {
	r.addRoute("DELETE", path, f)
}

func (r *Router) PATCH(path string, f func(*Context)) {
	r.addRoute("PATCH", path, f)
}

func (r *Router) OPTIONS(path string, f func(*Context)) {
	r.addRoute("OPTIONS", path, f)
}

func (r *Router) HEAD(path string, f func(*Context)) {
	r.addRoute("HEAD", path, f)
}

func (r *Router) CONNECT(path string, f func(*Context)) {
	r.addRoute("CONNECT", path, f)
}

func (r *Router) TRACE(path string, f func(*Context)) {
	r.addRoute("TRACE", path, f)
}

func (r *Router) ANY(path string, f func(*Context)) {
	r.addRoute("ANY", path, f)
}

func (r *Router) Miss(f func(*Context)) *Router {
	r.missRoute = f
	return r
}

func (r *Router) Route(path string, requestMethod string) (func(*Context), RouterParams) {
	requestMethod = strings.ToUpper(requestMethod)
	var fn func(*Context) = nil
	routerParams := make(RouterParams)

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
			return false
		}

		if patternMode {
			regexpParts := strings.Split(regexpTrim, "/")
			pathParts := strings.Split(path, "/")
			if len(regexpParts) != len(pathParts) {
				return false
			}

			for i, part := range regexpParts {
				if part[0] == ':' && len(pathParts[i]) > 0 {
					routerParams[part[1:]] = NewReqParamV([]string{pathParts[i]})
					continue
				}
			}
			fn = methodMap.Get(requestMethod)
		}
		return true
	})

	if r.missRoute != nil {
		return r.missRoute, nil
	}
	return fn, routerParams
}
