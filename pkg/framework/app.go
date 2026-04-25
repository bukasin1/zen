package framework

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework/share/logger"
)

type App struct {
	router            *Router
	middlewares       []Middleware
	systemMiddlewares []Middleware

	logger logger.Logger

	// ✅ NEW
	server *http.Server

	onStartHooks    []func(ctx context.Context) error
	onShutdownHooks []func(ctx context.Context) error

	shutdownTimeout time.Duration
}

// TODO: add new app configs (Probable future updates)
func New() *App {
	app := &App{
		router:      NewRouter(),
		middlewares: []Middleware{},
		// auto install system middlewares
		systemMiddlewares: []Middleware{Logger(), Recovery()},

		logger: logger.NewConsoleLogger(true),

		shutdownTimeout: 10 * time.Second,
	}

	return app
}

func (a *App) SetLogger(l logger.Logger) {
	a.logger = l
}

func (a *App) Use(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

func (a *App) UseSystem(m Middleware) {
	a.systemMiddlewares = append(a.systemMiddlewares, m)
}

func (a *App) Static(path, dir string) {
	fs := http.FileServer(http.Dir(dir))
	prefix := "/" + strings.Trim(path, "/*")

	// Strip the prefix from the request path
	// This is done so that the file server can find the files in the directory
	// For example, if the prefix is "/static" and the request path is "/static/file.txt",
	// the file server will look for "file.txt" in the directory.
	// This only needs to be done on the file server handler not the router
	fs = http.StripPrefix(prefix, fs)

	a.router.Handle(http.MethodGet, path, HandlerFunc(func(ctx *Context) {
		fs.ServeHTTP(ctx.Writer, ctx.Request)
		// run context extended hooks AFTER static write attempt
		ctx.runAfterResponseHooks()
	}))
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

func (a *App) OnStart(fn func(ctx context.Context) error) {
	a.onStartHooks = append(a.onStartHooks, fn)
}

func (a *App) OnShutdown(fn func(ctx context.Context) error) {
	a.onShutdownHooks = append(a.onShutdownHooks, fn)
}

func (a *App) runStartHooks(ctx context.Context) error {
	for _, hook := range a.onStartHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) runShutdownHooks(ctx context.Context) {
	for _, hook := range a.onShutdownHooks {
		if err := hook(ctx); err != nil {
			a.logger.Error("shutdown hook failed", logger.Fields{
				"error": err.Error(),
			})
		}
	}
}

func (a *App) buildAppHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, a.logger)

		handler := func(c *Context) {
			a.router.ServeHTTP(c)
		}

		handler = a.applyMiddlewares(handler)
		handler(ctx)
	})
}

func (a *App) Run(addr string) error {
	handler := a.buildAppHandler()
	// http.ListenAndServe(addr, handler)

	// create server instance
	a.server = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	rootCtx := context.Background()

	// 1. Run startup hooks
	if err := a.runStartHooks(rootCtx); err != nil {
		a.logger.Error("startup hook failed", logger.Fields{
			"error": err.Error(),
		})
		return err
	}

	// 2. Start server in goroutine
	go func() {
		a.logger.Info(fmt.Sprintf("server starting on http://localhost%v/\n", addr), logger.Fields{
			"addr": addr,
		})

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("server error", logger.Fields{
				"error": err.Error(),
			})
		}
	}()

	// 3. Listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	a.logger.Info("shutdown signal received", logger.Fields{
		"signal": sig.String(),
	})

	// 4. Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	// 5. Graceful shutdown (drains active requests)
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown failed", logger.Fields{
			"error": err.Error(),
		})
	} else {
		a.logger.Info("server shutdown complete", nil)
	}

	// 6. Run shutdown hooks AFTER draining
	a.runShutdownHooks(ctx)

	a.logger.Info("application shutdown complete", nil)

	return nil
}

// ------------------------- Deprecated Functions ------------------

// Deprecated: use Listen instead
// TODO: remove this function
func (a *App) ListenOld(addr string) error {
	handler := a.buildAppHandlerOld()
	return http.ListenAndServe(addr, handler)
}

// Deprecated: use buildAppHandler instead
// TODO: remove this function
func (a *App) buildAppHandlerOld() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, a.logger)

		handler := func(c *Context) {
			a.router.ServeHTTPOld(c)
		}

		handler = a.applyMiddlewares(handler)
		handler(ctx)
	})
}

// Deprecated: use Static instead
// TODO: remove this function
func (a *App) StaticOld(prefix, dir string) {
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
