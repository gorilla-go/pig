package pig

import "github.com/samber/do"

type Router struct {
	regRouteMap map[string]func(*do.Injector)
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Bind(regMap map[string]func(*do.Injector)) {
	r.regRouteMap = regMap
}

func (r *Router) Route(path string) func(*do.Injector) {
	return r.regRouteMap[path]
}
