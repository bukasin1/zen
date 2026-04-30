package framework

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework/logger"
)

type serviceEntry struct {
	once sync.Once
	init func() any
	inst any
}

type App struct {
	router            *Router
	middlewares       []Middleware
	systemMiddlewares []Middleware

	logger logger.Logger

	server *http.Server

	onStartHooks    []func(ctx context.Context) error
	onShutdownHooks []func(ctx context.Context) error

	shutdownTimeout time.Duration

	config     Config
	services   map[string]*serviceEntry
	servicesMu sync.RWMutex
}

func New() *App {
	cfg := DefaultConfig()

	app := &App{
		router:      NewRouter(),
		middlewares: []Middleware{},
		// auto install system middlewares
		systemMiddlewares: []Middleware{RequestLogger(), Recovery()},

		services: make(map[string]*serviceEntry),
	}

	app.SetAppConfig(cfg)

	return app
}

// Sets the application configuration after instantiation.
//
// Note:
// HTTP config ShutdownTimeout default is 10s,
// HTTP config Addr default is :8080,
// Log config Pretty default is false,
// Log config EnableJSON default is false
func (a *App) SetAppConfig(cfg Config) {
	if cfg.AppName != "" {
		a.config.AppName = cfg.AppName
	}
	if cfg.Env != "" {
		a.config.Env = cfg.Env
	}
	a.SetHTTPConfig(cfg.HTTP)

	a.SetLoggerConfig(cfg.Log)
}

// Sets App's HTTP configuration after instantiation.
//
// Note:
// HTTP config Addr default is :8080,
// HTTP config ShutdownTimeout default is 10s
func (a *App) SetHTTPConfig(h HTTPConfig) {
	if h.Addr != "" {
		a.config.HTTP.Addr = h.Addr
	}
	if h.ShutdownTimeout != 0 {
		a.config.HTTP.ShutdownTimeout = h.ShutdownTimeout
	}

	a.shutdownTimeout = a.config.HTTP.ShutdownTimeout
}

// Sets App's Logger configuration after instantiation.
//
// Note:
// Log config Pretty default is false,
// Log config EnableJSON default is false
func (a *App) SetLoggerConfig(l LogConfig) {
	if l.Level != "" {
		a.config.Log.Level = l.Level
	}
	a.config.Log.Pretty = l.Pretty
	a.config.Log.EnableJSON = l.EnableJSON

	// set default console logger for app if logger is nil
	if a.logger == nil {
		a.SetLogger(logger.NewConsoleLogger(a.config.Log.Pretty))
	} else {
		if logger, ok := a.logger.(*logger.ConsoleLogger); ok {
			logger.Pretty = a.config.Log.Pretty
		}
	}
}

func (a *App) SetLogger(l logger.Logger) {
	a.logger = l
}

// Use adds a middleware to the application.
//
// The middleware will be executed in the order they are added.
//
// Note:
// System middlewares are executed after regular middlewares.
func (a *App) Use(m Middleware) {
	a.middlewares = append(a.middlewares, m)
}

// UseSystem adds a system middleware to the application.
//
// Note:
// System middlewares are executed after regular middlewares.
// This function should be called by the framework itself or by extensions, not by the user.
func (a *App) UseSystem(m Middleware) {
	a.systemMiddlewares = append(a.systemMiddlewares, m)
}

// RegisterService registers a service with the application.
//
// Note: Service init functions must be idempotent and side-effect safe.
func (a *App) RegisterService(name string, init func() any) {
	a.servicesMu.Lock()
	defer a.servicesMu.Unlock()

	if _, exists := a.services[name]; exists {
		panic("service already registered: " + name)
	}

	if init == nil {
		panic("service init function cannot be nil: " + name)
	}

	a.services[name] = &serviceEntry{
		init: init,
	}
}

// Service returns the service with the given name.
func (a *App) Service(name string) any {
	a.servicesMu.RLock()
	entry, ok := a.services[name]
	a.servicesMu.RUnlock()

	if !ok {
		panic("service not found: " + name)
	}

	entry.once.Do(func() {
		entry.inst = entry.init()
	})

	return entry.inst
}

// GetService returns the service with the given name.
// It is a type-safe wrapper around the App Service function.
//
// Note:
// If the service is not found, it will panic.
// If the service type assertion fails, it will panic.
func GetService[T any](a *App, name string) T {
	svc := a.Service(name)
	return svc.(T)
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
			a.LogError("shutdown hook failed", logger.Fields{
				"error": err.Error(),
			})
		}
	}
}

func (a *App) LogInfo(msg string, fields logger.Fields) {
	if a.logger == nil {
		return
	}
	a.logger.Info(msg, fields)
}

func (a *App) LogError(msg string, fields logger.Fields) {
	if a.logger == nil {
		return
	}
	a.logger.Error(msg, fields)
}

func (a *App) LogWarn(msg string, fields logger.Fields) {
	if a.logger == nil {
		return
	}
	a.logger.Warn(msg, fields)
}

func (a *App) LogDebug(msg string, fields logger.Fields) {
	if a.logger == nil {
		return
	}
	a.logger.Debug(msg, fields)
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

func (a *App) Run(_ string) error {
	addr := a.config.HTTP.Addr
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
		a.LogError("startup hook failed", logger.Fields{
			"error": err.Error(),
		})
		return err
	}

	// 2. Start server in goroutine
	go func() {
		a.LogInfo("server starting", logger.Fields{
			"addr": addr,
		})

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.LogError("server error", logger.Fields{
				"error": err.Error(),
			})
		}
	}()

	// 3. Listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	signal.Stop(quit)

	a.LogInfo("shutdown signal received", logger.Fields{
		"signal": sig.String(),
	})

	// 4. Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	// 5. Graceful shutdown (drains active requests)
	if err := a.server.Shutdown(ctx); err != nil {
		a.LogError("server shutdown failed", logger.Fields{
			"error": err.Error(),
		})
	} else {
		a.LogInfo("server shutdown complete", nil)
	}

	// 6. Run shutdown hooks AFTER draining
	a.runShutdownHooks(ctx)

	a.LogInfo("application shutdown complete", nil)

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
