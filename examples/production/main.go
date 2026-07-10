package main

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/bukasin1/zen"
)

var _, filename, _, _ = runtime.Caller(0)

// Gets the directory of the current source file
var currentDir = filepath.Dir(filename)

var staticDir = currentDir + "/public"

func newApp() *zen.App {
	// Load application configuration.
	cfg := loadConfig()

	// Create the application.
	app := zen.New()

	// Configure the HTTP server.
	app.SetAppConfig(zen.Config{
		HTTP: zen.HTTPConfig{
			Addr:         cfg.Server.Address,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
	})

	// Configure logging.
	if cfg.Logging.JSON {
		app.SetLogger(zen.NewConsoleLogger(cfg.Logging.JSON))
	} else {
		app.SetLogger(zen.NewDevConsoleLogger())
	}

	// Global middleware.
	app.Use(
		// zen.RequestIDMiddleware(),
		zen.RequestLogger(),
		zen.Recovery(),
		zen.Timeout(time.Second*3),
		// zen.GzipCompression(),
		// zen.ETag(),
		// zen.CacheControl(),
		zen.CORS(zen.DefaultCORSConfig()),
	)

	// Enable operational endpoints.
	// app.RegisterHealthRoutes()
	// app.RegisterMetricsRoute()
	// app.RegisterRuntimeRoutes()
	app.RegisterOperationalRoutes()

	// Static assets.
	app.Static("/static/*", staticDir)

	server := &Server{}

	// Application routes.
	app.Route("/").
		Summary("Home page").
		Get(server.Home)

	api := app.Group("/api")

	api.Route("/hello").
		Summary("Hello endpoint").
		Description("Returns a simple JSON response.").
		Tags("Example").
		Get(server.Hello)

	return app
}

func main() {
	app := newApp()

	// Start the application.
	if err := app.Run(""); err != nil {
		panic(err)
	}
}
