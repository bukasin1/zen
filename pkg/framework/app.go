package framework

import (
	"net/http"
)

type App struct {
	mux *http.ServeMux
}

func New() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.mux.HandleFunc(path, handler)
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.mux.HandleFunc(path, handler)
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a.mux)
}
