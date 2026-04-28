package framework

import (
	"net/http"
	"strings"
)

type RouteBuilder struct {
	app         *App
	group       *Group
	path        string
	middlewares []Middleware
}

func (a *App) Route(path string) *RouteBuilder {
	return &RouteBuilder{
		app:  a,
		path: path,
	}
}

func (rb *RouteBuilder) Use(m ...Middleware) *RouteBuilder {
	rb.middlewares = append(rb.middlewares, m...)
	return rb
}

func (rb *RouteBuilder) applyMiddlewares(h HandlerFunc) HandlerFunc {
	h = chainMiddlewares(h, rb.middlewares)
	return h
}

func (rb *RouteBuilder) Static(dir string) {
	fs := http.FileServer(http.Dir(dir))

	prefix := "/" + strings.Trim(rb.path, "/*")

	// Strip the prefix from the request path
	// This is done so that the file server can find the files in the directory
	// For example, if the prefix is "/static" and the request path is "/static/file.txt",
	// the file server will look for "file.txt" in the directory.
	// This only needs to be done on the file server handler not the router
	fs = http.StripPrefix(prefix, fs)

	handler := HandlerFunc(func(ctx *Context) {
		fs.ServeHTTP(ctx.Writer, ctx.Request)
		// run context extended hooks AFTER static write attempt
		ctx.runAfterResponseHooks()
	})
	handler = rb.applyMiddlewares(handler)

	rb.app.router.Handle(http.MethodGet, rb.path, handler)
}

func (rb *RouteBuilder) Get(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodGet, rb.path, handler)
}

func (rb *RouteBuilder) Post(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPost, rb.path, handler)
}

func (rb *RouteBuilder) Put(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPut, rb.path, handler)
}

func (rb *RouteBuilder) Delete(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodDelete, rb.path, handler)
}

func (rb *RouteBuilder) Patch(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPatch, rb.path, handler)
}

func (rb *RouteBuilder) Head(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodHead, rb.path, handler)
}

func (rb *RouteBuilder) Options(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodOptions, rb.path, handler)
}

func (rb *RouteBuilder) Any(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodGet, rb.path, handler)
	rb.app.router.Handle(http.MethodPost, rb.path, handler)
	rb.app.router.Handle(http.MethodPut, rb.path, handler)
	rb.app.router.Handle(http.MethodDelete, rb.path, handler)
	rb.app.router.Handle(http.MethodPatch, rb.path, handler)
	rb.app.router.Handle(http.MethodHead, rb.path, handler)
	rb.app.router.Handle(http.MethodOptions, rb.path, handler)
}
