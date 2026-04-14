package framework

import (
	"net/http"
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

func (a *App) Get(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodGet, path, a.applyMiddlewares(handler))
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodPost, path, a.applyMiddlewares(handler))
}

func (a *App) applyMiddlewares(h HandlerFunc) HandlerFunc {
	h = chainMiddlewares(h, a.middlewares)
	h = chainMiddlewares(h, a.systemMiddlewares)
	return h
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
