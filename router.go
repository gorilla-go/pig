package pig

import (
	"github.com/samber/do"
	"strings"
)

type Router struct {
	regRouteMap map[string]func(*do.Injector)
}

type RouterParams map[string]string

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Map(regMap map[string]func(*do.Injector)) {
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

func (r *Router) Route(path string) (func(*do.Injector), RouterParams) {
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

			routerParams := make(map[string]string)
			for i, part := range regexpParts {
				if part[0] == ':' {
					routerParams[part[1:]] = pathParts[i]
					continue
				}
			}
			return fn, routerParams
		}
	}
	return nil, nil
}
