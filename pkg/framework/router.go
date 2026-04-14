package framework

import "net/http"

type HandlerFunc func(*Context)

type Router struct {
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	key := method + ":" + path
	r.routes[key] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + ":" + req.URL.Path

	if handler, ok := r.routes[key]; ok {
		ctx := NewContext(w, req)
		handler(ctx)
		return
	}

	http.NotFound(w, req)
}
