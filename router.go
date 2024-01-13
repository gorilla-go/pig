package pig

import (
	"strings"
)

type Router struct {
	regRouteMap map[string]func(*Context)
	missRoute   func(*Context)
}

type RouterParams ReqParams

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Map(regMap map[string]func(*Context)) {
	for s, f := range regMap {
		if len(s) == 0 {
			panic("route path can't be empty")
		}

		if f == nil {
			panic("route action can't be nil")
		}
	}
	r.regRouteMap = regMap
}

func (r *Router) Miss(f func(*Context)) *Router {
	r.missRoute = f
	return r
}

func (r *Router) Route(path string) (func(*Context), RouterParams) {
	for regexp, fn := range r.regRouteMap {
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
