package framework

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, bool) {
	methodRoutes, ok := r.routes[method]
	if !ok {
		return nil, false
	}

	pathHandler, ok := methodRoutes[path]
	if !ok {
		return nil, false
	}

	return pathHandler, true
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.FindRoute(req.Method, req.URL.Path); ok {
		ctx := NewContext(w, req)
		handler(ctx)
		return
	}

	http.NotFound(w, req)
}
