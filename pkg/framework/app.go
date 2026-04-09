package framework

import (
	"net/http"
)

type App struct {
	router *Router
}

func New() *App {
	return &App{
		router: NewRouter(),
	}
}

func (a *App) Get(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodGet, path, handler)
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.router.Handle(http.MethodPost, path, handler)
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
