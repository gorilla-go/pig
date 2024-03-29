package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation/constant"
	"github.com/gorilla-go/pig/foundation/mapping"
	"github.com/gorilla-go/pig/param"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	static          map[string]string
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
		static:          make(map[string]string),
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

func (r *Router) Static(path string, realPath string) {
	if r.group != "" {
		panic("group router nonsupport static files.")
	}
	realPath = filepath.Clean(realPath)
	if !strings.HasSuffix(realPath, constant.FileSystemSeparator) {
		realPath += constant.FileSystemSeparator
	}

	// is dir
	stat, err := os.Stat(realPath)
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		panic("static path must be a directory.")
	}

	path = strings.TrimPrefix(strings.TrimSpace(path), constant.WebSystemSeparator)
	if !strings.HasSuffix(path, constant.WebSystemSeparator) {
		path += constant.WebSystemSeparator
	}
	r.static[path] = realPath
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

func (r *Router) Url(routerName string, params map[string]any) string {
next:
	for presetPath, config := range r.routerConfigMap {
		if config.alias != routerName {
			continue
		}

		if !config.patternMode {
			return presetPath
		}

		paramsCopy := make(map[string]any)
		for k, v := range params {
			paramsCopy[k] = v
		}

		presetPaths := strings.Split(presetPath, constant.WebSystemSeparator)
		for i, presetPathItem := range presetPaths {
			presetPathItem = strings.TrimSpace(presetPathItem)

			if len(presetPathItem) > 4 &&
				strings.HasPrefix(presetPathItem, "<") &&
				strings.HasSuffix(presetPathItem, ">") {
				presetPathParamPair := strings.SplitN(presetPathItem[1:len(presetPathItem)-1], ":", 2)
				if len(presetPathParamPair) < 2 {
					continue next
				}

				key := strings.TrimSpace(presetPathParamPair[0])
				pregStr := strings.TrimSpace(presetPathParamPair[1])
				if v, ok := paramsCopy[key]; ok {
					match, err := regexp.Match(pregStr, []byte(fmt.Sprintf("%v", v)))
					if err != nil || !match {
						continue next
					}
					presetPaths[i] = fmt.Sprintf("%v", v)
					delete(paramsCopy, key)
					continue
				}
				continue next
			}

			if len(presetPathItem) > 1 && strings.HasPrefix(presetPathItem, ":") {
				key := strings.TrimSpace(presetPathItem[1:])
				if v, ok := paramsCopy[key]; ok {
					presetPaths[i] = fmt.Sprintf("%v", v)
					delete(paramsCopy, key)
					continue
				}
				continue next
			}
		}

		uv := url.Values{}
		for k, v := range paramsCopy {
			uv.Add(k, fmt.Sprintf("%v", v))
		}
		query := uv.Encode()
		if len(query) > 0 {
			query = "?" + query
		}
		return strings.Join(presetPaths, constant.WebSystemSeparator) + query
	}

	panic(fmt.Sprintf("router %s not exists.", routerName))
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

	// search for static file.
	if requestMethod == constant.GET && len(r.static) > 0 {
		fn = r.fetchStatic(path)
		if fn != nil {
			return fn, routerParams, middlewares
		}
	}

	r.regRouteMap.ForEach(func(presetPath string, methodMap *mapping.LinkedHashMap[constant.RequestMethod, func(*Context)]) bool {
		patternMode := r.routerConfigMap[presetPath].patternMode
		if !patternMode && presetPath == path {
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
			if len(presetPathParts) != len(pathParts) {
				return true
			}

			for i, presetPathItem := range presetPathParts {
				pathItem := strings.TrimSpace(pathParts[i])
				presetPathItem = strings.TrimSpace(presetPathItem)

				if len(pathItem) == 0 && len(presetPathItem) > 0 {
					return true
				}

				if len(presetPathItem) > 4 &&
					strings.HasPrefix(presetPathItem, "<") &&
					strings.HasSuffix(presetPathItem, ">") {
					presetPathParamPair := strings.SplitN(presetPathItem[1:len(presetPathItem)-1], ":", 2)
					if len(presetPathParamPair) < 2 {
						return true
					}
					pregStr := strings.TrimSpace(presetPathParamPair[1])
					if pregStr[0] != '^' {
						pregStr = "^" + pregStr
					}

					if pregStr[len(pregStr)-1] != '$' {
						pregStr = pregStr + "$"
					}
					match, err := regexp.Match(pregStr, []byte(pathItem))
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
	match, _ := regexp.MatchString("[:<>]", path)
	return match
}

func (r *Router) fetchStatic(path string) func(*Context) {
	for staticPathPrefix, realPath := range r.static {
		if strings.HasPrefix(path, staticPathPrefix) {
			file := filepath.Join(realPath, filepath.Clean(strings.TrimPrefix(path, staticPathPrefix)))

			st, err := os.Stat(file)
			if err == nil && !st.IsDir() && strings.HasPrefix(file, realPath) {
				return func(context *Context) {
					http.ServeFile(context.Response().Raw(), context.Request().Raw(), file)
				}
			}
		}
	}
	return nil
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
