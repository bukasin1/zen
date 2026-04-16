package framework

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type route struct {
	pattern string
	handler HandlerFunc
}

type staticRoute struct {
	prefix  string
	handler http.Handler
}

type Router struct {
	routes       map[string][]route
	staticRoutes []staticRoute
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string][]route),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = []route{}
	}
	r.routes[method] = append(r.routes[method], route{
		pattern: path,
		handler: handler,
	})
}

func (r *Router) HandleStatic(prefix string, handler http.Handler) {
	r.staticRoutes = append(r.staticRoutes, staticRoute{
		prefix:  prefix,
		handler: handler,
	})
}

func matchRoute(pattern, path string) (bool, map[string]string) {
	// check if match route path matches pattern and return a map of params
	pattern = strings.Trim(pattern, "/")
	path = strings.Trim(path, "/")

	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)

	for i := range patternParts {
		part := patternParts[i]
		value := pathParts[i]

		// param part
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = value
			continue
		}

		// static part
		if part != value {
			return false, nil
		}
	}

	return true, params
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, map[string]string, bool) {
	methodRoutes, ok := r.routes[method]
	if !ok {
		return nil, nil, false
	}

	for _, route := range methodRoutes {
		matched, params := matchRoute(route.pattern, path)
		if matched {
			return route.handler, params, true
		}
	}

	return nil, nil, false
}

func matchStaticPrefix(path, prefix string) bool {
	path = "/" + strings.Trim(path, "/")
	prefix = "/" + strings.Trim(prefix, "/")

	// root fallback matches everything
	if prefix == "/" {
		return true
	}

	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, params, ok := r.FindRoute(req.Method, req.URL.Path); ok {
		ctx := NewContext(w, req)
		ctx.params = params
		handler(ctx)
		return
	}

	for _, static := range r.staticRoutes {
		if matchStaticPrefix(req.URL.Path, static.prefix) {
			static.handler.ServeHTTP(w, req)
			return
		}
	}

	ctx := NewContext(w, req)
	ctx.Error(http.StatusNotFound, "404 page not found!")
}
