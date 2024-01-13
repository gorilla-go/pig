package pig

import (
	"strings"
)

type RequestMethodType string
type RouteIndex map[string]func(*Context)

type Router struct {
	regRouteMap map[string]RouteIndex
	missRoute   func(*Context)
}

type RouterParams ReqParams

func NewRouter() *Router {
	return &Router{
		regRouteMap: make(map[string]RouteIndex),
	}
}

func (r *Router) addRoute(t string, path string, f func(*Context)) {
	if _, ok := r.regRouteMap[path]; !ok {
		r.regRouteMap[path] = make(RouteIndex)
	}

	r.regRouteMap[path][t] = f
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

	for regexp, routeIndex := range r.regRouteMap {
		fn, ok := routeIndex[requestMethod]
		if !ok {
			if _, ok := routeIndex["ANY"]; ok {
				fn = routeIndex["ANY"]
			} else {
				continue
			}
		}

		regexpTrim := strings.Trim(regexp, "/")
		path = strings.Trim(path, "/")

		patternMode := strings.Contains(regexpTrim, ":")
		if !patternMode && regexpTrim == path {
			return fn, nil
		}

		if patternMode {
			regexpParts := strings.Split(regexpTrim, "/")
			pathParts := strings.Split(path, "/")
			if len(regexpParts) != len(pathParts) {
				continue
			}

			routerParams := make(RouterParams)
			for i, part := range regexpParts {
				if part[0] == ':' {
					routerParams[part[1:]] = NewReqParamV([]string{pathParts[i]})
					continue
				}
			}
			return fn, routerParams
		}
	}

	if r.missRoute != nil {
		return r.missRoute, nil
	}

	return nil, nil
}
