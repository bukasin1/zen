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

	middlewareDefinitions []MiddlewareDefinition
	name                  string
	metadata              RouteMetadata
}

// Route returns a new RouteBuilder for the given path.
// It is used to define a new route under the current app.
//
// Example:
//
//	app.Route("/users").Get(func(ctx *Context) { ... })
func (a *App) Route(path string) *RouteBuilder {
	return &RouteBuilder{
		app:  a,
		path: path,
		// name:                  "",
		middlewareDefinitions: append([]MiddlewareDefinition{}, a.middlewareDefinitions...),
		metadata:              make(RouteMetadata),
	}
}

// Use adds middleware to the route.
// Call this before any route definitions (Get, Post, etc).
func (rb *RouteBuilder) Use(m ...Middleware) *RouteBuilder {
	rb.middlewares = append(rb.middlewares, m...)
	return rb
}

// UseNamed adds named middleware to the route.
// Call this before any route definitions (Get, Post, etc).
//
// Example:
//
//	loggerMiddleware := framework.NamedMiddleware("logger", framework.Logger())
//	app.Route("/users").
//		UseNamed(loggerMiddleware).
//		Get(func(ctx *Context) { ... })
func (rb *RouteBuilder) UseNamed(mds ...MiddlewareDefinition) *RouteBuilder {
	for _, md := range mds {
		rb.middlewares = append(rb.middlewares, md.Func)

		rb.middlewareDefinitions = append(rb.middlewareDefinitions, md)
	}

	return rb
}

// Name sets the name of the route.
//
// Example:
//
//	app.Route("/users").Name("users").Get(func(ctx *Context) { ... })
func (rb *RouteBuilder) Name(name string) *RouteBuilder {
	rb.name = name
	return rb
}

// Meta adds metadata to the route.
//
// Example:
//
//	app.Route("/users").Meta("key", "value").Get(func(ctx *Context) { ... })
func (rb *RouteBuilder) Meta(key string, value any) *RouteBuilder {
	rb.metadata[key] = value
	return rb
}

func (rb *RouteBuilder) applyMiddlewares(h HandlerFunc) HandlerFunc {
	h = chainMiddlewares(h, rb.middlewares)
	return h
}

func (rb *RouteBuilder) middlewareNames() []string {
	names := make([]string, 0, len(rb.middlewareDefinitions))

	for _, md := range rb.middlewareDefinitions {
		names = append(names, md.Name)
	}

	return names
}

func (rb *RouteBuilder) registerRoute(method, handlerName string) {
	rb.app.registerRoute(RouteDefinition{
		Method: method,
		Path:   rb.path,
		Name:   rb.name,

		HandlerName: handlerName,

		Middlewares: rb.middlewareNames(),

		Metadata: cloneRouteMetadata(rb.metadata),
	})
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

	rb.registerRoute(http.MethodGet, "StaticFileServer")
}

func (rb *RouteBuilder) Get(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodGet, rb.path, handler)

	rb.registerRoute(http.MethodGet, "GET_HANDLER")
}

func (rb *RouteBuilder) Post(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPost, rb.path, handler)

	rb.registerRoute(http.MethodPost, "POST_HANDLER")
}

func (rb *RouteBuilder) Put(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPut, rb.path, handler)

	rb.registerRoute(http.MethodPut, "PUT_HANDLER")
}

func (rb *RouteBuilder) Delete(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodDelete, rb.path, handler)

	rb.registerRoute(http.MethodDelete, "DELETE_HANDLER")
}

func (rb *RouteBuilder) Patch(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodPatch, rb.path, handler)

	rb.registerRoute(http.MethodPatch, "PATCH_HANDLER")
}

func (rb *RouteBuilder) Head(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodHead, rb.path, handler)

	rb.registerRoute(http.MethodHead, "HEAD_HANDLER")
}

func (rb *RouteBuilder) Options(handler HandlerFunc) {
	handler = rb.applyMiddlewares(handler)
	rb.app.router.Handle(http.MethodOptions, rb.path, handler)

	rb.registerRoute(http.MethodOptions, "OPTIONS_HANDLER")
}

func (rb *RouteBuilder) Any(handler HandlerFunc) {
	rb.Get(handler)
	rb.Post(handler)
	rb.Put(handler)
	rb.Delete(handler)
	rb.Patch(handler)
	rb.Head(handler)
	rb.Options(handler)
}
