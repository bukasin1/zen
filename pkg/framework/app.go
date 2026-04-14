package framework

import (
	"net/http"
)

type App struct {
	router      *Router
	middlewares []Middleware
}

func New() *App {
	return &App{
		router:      NewRouter(),
		middlewares: []Middleware{},
	}
}

func (a *App) Use(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

func (a *App) Get(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodGet, path, chainMiddlewares(handler, a.middlewares))
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodPost, path, chainMiddlewares(handler, a.middlewares))
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
