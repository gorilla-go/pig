package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation/constant"
	"github.com/gorilla-go/pig/foundation/mapping"
	"github.com/gorilla-go/pig/param"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

type Router struct {
	group           string
	groupMiddleware []IMiddleware
	regRouteMap     *mapping.LinkedHashMap[string, *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]]
	routerConfigMap map[string]*RouterConfig
	missRoute       func(*Context)
	middlewareMap   map[string][]IMiddleware
}

type RouterConfig struct {
	alias       string
	patternMode bool
}

func (a *RouterConfig) Name(name string) *RouterConfig {
	a.alias = name
	return a
}

func NewRouter() *Router {
	return &Router{
		group:           "",
		groupMiddleware: make([]IMiddleware, 0),
		regRouteMap:     mapping.NewLinkedHashMap[string, *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]](),
		routerConfigMap: make(map[string]*RouterConfig),
		middlewareMap:   make(map[string][]IMiddleware),
	}
}

func (r *Router) addRoute(
	requestMethod constant.RequestMethod,
	presetRequestPath string,
	function func(*Context),
	middleware []IMiddleware,
) *RouterConfig {
	presetRequestPath = strings.TrimSpace(presetRequestPath)
	if strings.HasSuffix(presetRequestPath, constant.WebSystemSeparator) {
		presetRequestPath += constant.DefaultResourcePath
	}
	presetRequestPath = strings.TrimPrefix(presetRequestPath, constant.WebSystemSeparator)
	for _, item := range strings.Split(presetRequestPath, constant.WebSystemSeparator) {
		if len(strings.TrimSpace(item)) == 0 {
			panic(fmt.Sprintf("invalid router %s", presetRequestPath))
		}
	}

	if r.group != "" {
		if strings.HasPrefix(presetRequestPath, constant.WebSystemSeparator) {
			presetRequestPath = presetRequestPath[1:]
		}
		presetRequestPath = r.group + presetRequestPath
	}

	if r.regRouteMap.ContainsKey(presetRequestPath) == false {
		r.regRouteMap.Put(presetRequestPath, mapping.NewLinkedHashMap[constant.RequestMethod, func(*Context)]())
	}

	requestPrefixMap := r.regRouteMap.MustGet(presetRequestPath)
	if requestPrefixMap.ContainsKey(requestMethod) {
		panic(fmt.Sprintf("router %s already exists.", presetRequestPath))
	}
	requestPrefixMap.Put(requestMethod, function)
	r.routerConfigMap[presetRequestPath] = &RouterConfig{
		patternMode: r.isPatternMode(presetRequestPath),
	}

	pid := r.RequestPathId(presetRequestPath, requestMethod)
	if len(middleware) > 0 {
		r.middlewareMap[pid] = middleware
		return r.routerConfigMap[presetRequestPath]
	}

	if len(r.groupMiddleware) > 0 {
		r.middlewareMap[pid] = r.groupMiddleware
		return r.routerConfigMap[presetRequestPath]
	}

	return r.routerConfigMap[presetRequestPath]
}

func (r *Router) RequestPathId(path string, method constant.RequestMethod) string {
	return fmt.Sprintf("%s://%s", method, path)
}

func (r *Router) Group(path string, fn func(r *Router), middleware ...IMiddleware) {
	path = strings.Trim(strings.TrimSpace(path), constant.WebSystemSeparator)
	path += constant.WebSystemSeparator

	router := &Router{
		group:           path,
		groupMiddleware: middleware,
		regRouteMap:     mapping.NewLinkedHashMap[string, *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]](),
		routerConfigMap: make(map[string]*RouterConfig),
		middlewareMap:   make(map[string][]IMiddleware),
	}
	fn(router)

	router.regRouteMap.ForEach(func(path string, methodMap *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]) bool {
		methodMap.ForEach(func(method constant.RequestMethod, fn func(*Context)) bool {
			r.addRoute(method, path, fn, router.middlewareMap[r.RequestPathId(path, method)])
			r.routerConfigMap[path] = router.routerConfigMap[path]
			return true
		})
		return true
	})
}

func (r *Router) Miss(f func(*Context)) *Router {
	r.missRoute = f
	return r
}

func (r *Router) Route(
	path string,
	requestMethod constant.RequestMethod,
) (func(*Context), *param.RequestParamPairs[*param.RequestParamItems[string]], []IMiddleware) {
	path = strings.TrimSpace(path)
	if strings.HasSuffix(path, constant.WebSystemSeparator) {
		path += constant.DefaultResourcePath
	}
	path = strings.TrimPrefix(path, constant.WebSystemSeparator)

	var fn func(*Context) = nil
	routerParams := param.NewRequestParamPairs[*param.RequestParamItems[string]]()
	middlewares := make([]IMiddleware, 0)

	var includedWildcard = lo.IndexOf(strings.Split(path, ""), "*") >= 0
	r.regRouteMap.ForEach(func(presetPath string, methodMap *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]) bool {
		patternMode := r.routerConfigMap[presetPath].patternMode
		if !patternMode && presetPath == path && !includedWildcard {
			if methodMap.ContainsKey(requestMethod) {
				fn = methodMap.MustGet(requestMethod)
				if m, ok := r.middlewareMap[r.RequestPathId(presetPath, requestMethod)]; ok {
					middlewares = m
				}
				return false
			}

			if methodMap.ContainsKey(constant.ANY) {
				fn = methodMap.MustGet(constant.ANY)
				if m, ok := r.middlewareMap[r.RequestPathId(presetPath, constant.ANY)]; ok {
					middlewares = m
				}
			}
			return false
		}

		if patternMode && (methodMap.ContainsKey(requestMethod) || methodMap.ContainsKey(constant.ANY)) {
			presetPathParts := strings.Split(presetPath, constant.WebSystemSeparator)
			pathParts := strings.Split(path, constant.WebSystemSeparator)
			if len(presetPathParts) != len(pathParts) && !strings.HasSuffix(presetPath, "*") {
				return true
			}

			for i, presetPathItem := range presetPathParts {
				if i > len(pathParts)-1 {
					return true
				}

				pathItem := strings.TrimSpace(pathParts[i])
				presetPathItem = strings.TrimSpace(presetPathItem)

				if len(pathItem) == 0 || lo.IndexOf(strings.Split(pathItem, ""), "*") >= 0 {
					return true
				}

				if len(presetPathItem) > 4 &&
					strings.HasPrefix(presetPathItem, "<") &&
					strings.HasSuffix(presetPathItem, ">") {
					presetPathParamPair := strings.SplitN(presetPathItem[1:len(presetPathItem)-1], ":", 2)
					if len(presetPathParamPair) < 2 {
						return true
					}
					regexpStr := strings.TrimSpace(presetPathParamPair[1])
					regexpStr = strings.TrimPrefix(regexpStr, "^")
					regexpStr = "^" + regexpStr
					regexpStr = strings.TrimSuffix(regexpStr, "$")
					regexpStr = regexpStr + "$"

					match, err := regexp.MatchString(regexpStr, pathItem)
					if err != nil || !match {
						return true
					}

					key := strings.TrimSpace(presetPathParamPair[0])
					routerParams.Raw().Set(
						key,
						param.NewRequestParamItems([]string{pathItem}),
					)
					continue
				}

				if len(presetPathItem) > 1 && strings.HasPrefix(presetPathItem, ":") {
					routerParams.Raw().Set(
						presetPathItem[1:],
						param.NewRequestParamItems([]string{pathItem}),
					)
					continue
				}

				if lo.IndexOf(strings.Split(presetPathItem, ""), "*") != -1 {
					presetPathArr := strings.Split(presetPathItem, "*")
					var presetPathArrTmp []string
					for _, s := range presetPathArr {
						s = strings.TrimSpace(s)
						if len(s) > 0 {
							presetPathArrTmp = append(
								presetPathArrTmp,
								regexp.QuoteMeta(s),
							)
						}
					}

					regexpStr := strings.Join(presetPathArrTmp, ".*")
					if len(presetPathArrTmp) == 0 {
						regexpStr = ".*"
					}

					matched, err := regexp.MatchString(regexpStr, pathItem)
					if err == nil && matched {
						if i == (len(presetPathParts)-1) && strings.HasSuffix(presetPathItem, "*") {
							break
						}
						continue
					}
				}

				if pathItem != presetPathItem {
					return true
				}
			}

			if methodMap.ContainsKey(requestMethod) {
				fn = methodMap.MustGet(requestMethod)
				if m, ok := r.middlewareMap[r.RequestPathId(presetPath, requestMethod)]; ok {
					middlewares = m
				}
				return false
			}

			fn = methodMap.MustGet(constant.ANY)
			if m, ok := r.middlewareMap[r.RequestPathId(presetPath, constant.ANY)]; ok {
				middlewares = m
			}
			return false
		}

		return true
	})

	if r.missRoute != nil && fn == nil {
		fn = r.missRoute
	}

	if len(middlewares) > 0 {
		return fn, routerParams, middlewares
	}
	return fn, routerParams, nil
}

func (r *Router) isPatternMode(path string) bool {
	match, _ := regexp.MatchString("[:<>*]", path)
	return match
}

func (r *Router) GET(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.GET, path, f, middleware)
}

func (r *Router) POST(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.POST, path, f, middleware)
}

func (r *Router) PUT(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.PUT, path, f, middleware)
}

func (r *Router) DELETE(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.DELETE, path, f, middleware)
}

func (r *Router) PATCH(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.PATCH, path, f, middleware)
}

func (r *Router) OPTION(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.OPTION, path, f, middleware)
}

func (r *Router) HEAD(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.HEAD, path, f, middleware)
}

func (r *Router) CONNECT(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.CONNECT, path, f, middleware)
}

func (r *Router) TRACE(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.TRACE, path, f, middleware)
}

func (r *Router) ANY(path string, f func(*Context), middleware ...IMiddleware) *RouterConfig {
	return r.addRoute(constant.ANY, path, f, middleware)
}
