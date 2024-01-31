package pig

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Router struct {
	group           string
	groupMiddleware []IMiddleware
	regRouteMap     *foundation.LinkedHashMap[string, *foundation.LinkedHashMap[string, func(*Context)]]
	routerAliasMap  map[string]*RouterAlias
	missRoute       func(*Context)
	middlewareMap   map[string][]IMiddleware
	static          map[string]string
}

type RouterAlias struct {
	name string
}

func (a *RouterAlias) Name(name string) {
	a.name = name
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
		routerAliasMap: make(map[string]*RouterAlias),
		middlewareMap:  make(map[string][]IMiddleware),
		static:         make(map[string]string),
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

	_, ok := r.routerAliasMap[path]
	if !ok {
		r.routerAliasMap[path] = &RouterAlias{}
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

func (r *Router) Static(path string, realPath string) {
	if r.group != "" {
		panic("group router nonsupport static files.")
	}

	r.static[path] = realPath
}

func (r *Router) GET(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("GET", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) POST(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("POST", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) PUT(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("PUT", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) DELETE(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("DELETE", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) PATCH(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("PATCH", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) OPTIONS(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("OPTIONS", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) HEAD(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("HEAD", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) CONNECT(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("CONNECT", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) TRACE(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("TRACE", path, f, middleware)
	return r.routerAliasMap[path]
}

func (r *Router) ANY(path string, f func(*Context), middleware ...IMiddleware) *RouterAlias {
	r.addRoute("ANY", path, f, middleware)
	return r.routerAliasMap[path]
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
		routerAliasMap: make(map[string]*RouterAlias),
		middlewareMap:  make(map[string][]IMiddleware),
	}
	f(router)

	router.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		methodMap.ForEach(func(method string, fn func(*Context)) bool {
			r.addRoute(method, uri, fn, router.middlewareMap[r.ReqUniPath(uri, method)])
			return true
		})
		return true
	})

	for s, alias := range router.routerAliasMap {
		if _, ok := r.routerAliasMap[s]; ok {
			panic(fmt.Sprintf("router alias %s already exists.", s))
		}
		r.routerAliasMap[s] = alias
	}
}

func (r *Router) Miss(f func(*Context)) *Router {
	r.missRoute = f
	return r
}

func (r *Router) Url(routerName string, params map[string]any) string {
alias:
	for s, alias := range r.routerAliasMap {
		if alias.name == routerName {
			if (params == nil || len(params) == 0) && !r.isPatternMode(s) {
				return s
			}

			paramsCopy := make(map[string]any)
			for k, v := range params {
				paramsCopy[k] = v
			}

			uriParts := strings.Split(s, "/")
			for i, part := range uriParts {
				if len(part) > 4 && part[0] == '<' && part[len(part)-1] == '>' {
					pathFormatArr := strings.SplitN(part[1:len(part)-1], ":", 2)
					if len(pathFormatArr) < 2 {
						continue alias
					}

					paramName := strings.TrimSpace(pathFormatArr[0])
					if v, ok := paramsCopy[paramName]; ok {
						match, err := regexp.Match(strings.TrimSpace(pathFormatArr[1]), []byte(fmt.Sprintf("%v", v)))
						if err != nil || !match {
							continue alias
						}
						uriParts[i] = fmt.Sprintf("%v", v)
						delete(paramsCopy, paramName)
						continue
					}
					continue alias
				}

				if len(part) > 1 && part[0] == ':' {
					paramName := strings.TrimSpace(part[1:])
					if v, ok := paramsCopy[paramName]; ok {
						uriParts[i] = fmt.Sprintf("%v", v)
						delete(paramsCopy, paramName)
						continue
					}
					continue alias
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
			return strings.Join(uriParts, "/") + query
		}
	}

	panic(fmt.Sprintf("router %s not exists.", routerName))
}

func (r *Router) Route(path string, requestMethod string) (func(*Context), RouterParams, []IMiddleware) {
	requestMethod = strings.ToUpper(requestMethod)
	var fn func(*Context) = nil
	routerParams := make(RouterParams)
	middlewares := make([]IMiddleware, 0)

	if requestMethod == "GET" && len(r.static) > 0 {
		fn = r.fetchStatic(path)
		if fn != nil {
			return fn, routerParams, middlewares
		}
	}

	r.regRouteMap.ForEach(func(uri string, methodMap *foundation.LinkedHashMap[string, func(*Context)]) bool {
		originUri := uri
		uri = strings.Trim(uri, "/")
		path = strings.Trim(path, "/")

		patternMode := r.isPatternMode(uri)
		if !patternMode && uri == path {
			if methodMap.ContainsKey(requestMethod) {
				fn = methodMap.Get(requestMethod)
				if m, ok := r.middlewareMap[r.ReqUniPath(originUri, requestMethod)]; ok {
					middlewares = m
				}
				return false
			}

			fn = methodMap.Get("ANY")
			if m, ok := r.middlewareMap[r.ReqUniPath(originUri, "ANY")]; ok {
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
				pathParts[i] = strings.TrimSpace(pathParts[i])
				part = strings.TrimSpace(part)

				if len(pathParts[i]) == 0 {
					return true
				}

				if len(part) > 4 && part[0] == '<' && part[len(part)-1] == '>' {
					pathFormatArr := strings.SplitN(part[1:len(part)-1], ":", 2)
					if len(pathFormatArr) < 2 {
						return true
					}

					match, err := regexp.Match(strings.TrimSpace(pathFormatArr[1]), []byte(pathParts[i]))
					if err != nil || !match {
						return true
					}

					routerParams[strings.TrimSpace(pathFormatArr[0])] = foundation.NewReqParamV([]string{pathParts[i]})
					continue
				}

				if len(part) > 1 && part[0] == ':' {
					routerParams[strings.TrimSpace(part[1:])] = foundation.NewReqParamV([]string{pathParts[i]})
					continue
				}

				if part != pathParts[i] {
					return true
				}
			}

			if methodMap.ContainsKey(requestMethod) {
				fn = methodMap.Get(requestMethod)
				if m, ok := r.middlewareMap[r.ReqUniPath(originUri, requestMethod)]; ok {
					middlewares = m
				}
				return false
			}

			fn = methodMap.Get("ANY")
			if m, ok := r.middlewareMap[r.ReqUniPath(originUri, "ANY")]; ok {
				middlewares = m
			}
			return false
		}
		return true
	})

	if r.missRoute != nil && fn == nil {
		return r.missRoute, nil, nil
	}

	if len(middlewares) > 0 {
		return fn, routerParams, middlewares
	}

	return fn, routerParams, nil
}

func (r *Router) isPatternMode(uri string) bool {
	return strings.Contains(uri, ":") ||
		strings.Contains(uri, "<") ||
		strings.Contains(uri, ">")
}

func (r *Router) fetchStatic(path string) func(*Context) {
	for uriPrefix, realPath := range r.static {
		if strings.HasPrefix(path, uriPrefix) {
			file := fmt.Sprintf(
				"%s/%s",
				filepath.Dir(realPath),
				strings.TrimPrefix(path, uriPrefix),
			)

			st, err := os.Stat(file)
			if err == nil && !st.IsDir() {
				return func(context *Context) {
					fi := NewFile(file)
					f, err := os.Open(fi.FilePath)
					if err != nil {
						panic(err)
					}
					defer func() {
						err := f.Close()
						if err != nil {
							panic(err)
						}
					}()
					context.Response().Raw().Header().Set("Content-Type", fi.ContentType)
					_, err = io.Copy(context.Response().Raw(), f)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
	return nil
}
