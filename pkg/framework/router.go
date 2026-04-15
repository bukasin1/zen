package framework

import (
	"net/http"
)

type HandlerFunc func(*Context)

type route struct {
	pattern string
	handler HandlerFunc
}

type Router struct {
	routes map[string][]route
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
		pattern: path[1:],
		handler: handler,
	})
}

func matchRoute(pattern, path string) bool {
	return pattern == path[1:]
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, bool) {
	methodRoutes, ok := r.routes[method]
	if !ok {
		return nil, false
	}

	for _, routes := range methodRoutes {
		matched := matchRoute(routes.pattern, path)
		if matched {
			return routes.handler, true
		}
	}

	return nil, false
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.FindRoute(req.Method, req.URL.Path); ok {
		ctx := NewContext(w, req)
		handler(ctx)
		return
	}

	http.NotFound(w, req)
}
