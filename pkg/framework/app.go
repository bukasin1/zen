package framework

import (
	"net/http"
	"strings"
)

type App struct {
	router            *Router
	middlewares       []Middleware
	systemMiddlewares []Middleware
}

// TODO: add new app configs (Probable future updates)
func New() *App {
	app := &App{
		router:      NewRouter(),
		middlewares: []Middleware{},
		// auto install system middlewares
		systemMiddlewares: []Middleware{Recovery()},
	}

	return app
}

func (a *App) Use(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

func (a *App) UseSystem(m Middleware) {
	a.systemMiddlewares = append(a.systemMiddlewares, m)
}

func (a *App) Static(prefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	prefix = "/" + strings.Trim(prefix, "/")

	// Strip the prefix from the request path
	// This is done so that the file server can find the files in the directory
	// For example, if the prefix is "/static" and the request path is "/static/file.txt",
	// the file server will look for "file.txt" in the directory.
	// This only needs to be done on the file server handler not the router
	fs = http.StripPrefix(prefix, fs)

	a.router.HandleStatic(prefix, fs)
}

func (a *App) Get(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodGet, path, handler)
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodPost, path, handler)
}

func (a *App) applyMiddlewares(h HandlerFunc) HandlerFunc {
	h = chainMiddlewares(h, a.middlewares)
	h = chainMiddlewares(h, a.systemMiddlewares)
	return h
}

func (a *App) buildAppHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)

		handler := func(c *Context) {
			a.router.ServeHTTP(c.Writer, c.Request)
		}

		handler = a.applyMiddlewares(handler)
		handler(ctx)
	})
}

func (a *App) Listen(addr string) error {
	handler := a.buildAppHandler()
	return http.ListenAndServe(addr, handler)
}
