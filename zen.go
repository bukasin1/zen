package zen

import (
	"net/http/httptest"

	"github.com/bukasin1/zen/internal/zencore"
	"github.com/bukasin1/zen/internal/zencore/errors"
	"github.com/bukasin1/zen/internal/zencore/logger"
)

type (
	App                  = zencore.App
	Context              = zencore.Context
	AppError             = errors.AppError
	Middleware           = zencore.Middleware
	Config               = zencore.Config
	Logger               = logger.Logger
	HTTPConfig           = zencore.HTTPConfig
	LogConfig            = zencore.LogConfig
	MiddlewareDefinition = zencore.MiddlewareDefinition

	// HandlerFunc is a function that handles an HTTP request.
	// It takes a *Context as an argument and returns nothing.
	HandlerFunc     = zencore.HandlerFunc
	RouteDocOptions = zencore.RouteDocOptions
)

// Error code aliases
const (
	ErrRequestTimeout      = errors.ErrRequestTimeout
	ErrBadRequest          = errors.ErrBadRequest
	ErrUnauthorized        = errors.ErrUnauthorized
	ErrForbidden           = errors.ErrForbidden
	ErrNotFound            = errors.ErrNotFound
	ErrConflict            = errors.ErrConflict
	ErrInternal            = errors.ErrInternal
	ErrValidation          = errors.ErrValidation
	ErrRequestBodyTooLarge = errors.ErrRequestBodyTooLarge
)

var (
	New                 = zencore.New
	NewContext          = zencore.NewContext
	NewConsoleLogger    = logger.NewConsoleLogger
	NewDevConsoleLogger = logger.NewDevConsoleLogger

	// load app config from environment variables
	// it uses the [DefaultConfig] as a fallback
	// Environment variables are case sensitive and should be uppercase
	//
	//	APP_NAME="Application name" (default: "Zen")
	//	APP_ENV="Application environment" (default: "development")
	//	HTTP_ADDR="HTTP server address" (default: ":8080")
	//	HTTP_SHUTDOWN_TIMEOUT="HTTP server shutdown timeout" (default: 2s)
	//	LOG_LEVEL="Log level" (default: "debug")
	//	LOG_PRETTY="Enable pretty printing" (default: true)
	//	LOG_ENABLE_JSON="Enable JSON logging (default: true)
	LoadConfigFromEnv = zencore.LoadConfigFromEnv
	// DefaultConfig returns a default configuration for the application.
	DefaultConfig = zencore.DefaultConfig

	// --------------- Middlewares start ---------------

	// NamedMiddleware returns a MiddlewareDefinition for the given name and middleware.
	//
	// Example:
	//
	//	loggerMiddleware := zencore.NamedMiddleware("logger", zencore.Logger())
	//	app.UseNamed(loggerMiddleware)
	NamedMiddleware = zencore.NamedMiddleware

	// CORS returns a middleware that applies CORS headers to the response.
	//
	// Note:
	//   - [DefaultCORSConfig] should be used when you want to allow all origins, methods, and headers.
	//   - [DefaultCORSConfig] is used if no config is provided.
	//   - When custom config is provided, it merges with default values if the field is empty/zero value
	//     otherwise it overrides the default config fields.
	//   - If credentials are allowed, wildcard origin cannot be used.
	//   - Wildcard origin ("*") must be the first and only entry in AllowOrigins if it's included.
	//     Mixing wildcard and specific origins is not allowed.
	//
	// Example:
	//
	//	// use default cors config
	//	app.Use(zen.CORS(zen.DefaultCORSConfig()))
	//	// This works as well - with default config
	//	app.Use(zen.CORS(zen.CORSConfig{}))
	//
	//	// use custom cors config
	//	app.Use(zen.CORS(zen.CORSConfig{
	//		AllowOrigins:     []string{"http://localhost:3000"},
	//		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	//		AllowHeaders:     []string{"Content-Type", "Authorization"},
	//		ExposeHeaders:    []string{"Content-Length"},
	//		AllowCredentials: true,
	//		MaxAge:           86400,
	//	}))
	CORS = zencore.CORS

	// DefaultCORSConfig returns the default CORS configuration.
	// It allows all origins, methods, and headers.
	// It also sets the MaxAge to 300 seconds (5 minutes).
	// Use [DefaultCORSConfig] when you want to allow all origins, methods, and headers.
	DefaultCORSConfig = zencore.DefaultCORSConfig

	// Recovery is a middleware that recovers from panics and returns a 500 error.
	Recover = zencore.Recovery
	// RequestLogger is a middleware that logs requests to the console.
	//
	// It logs the following information:
	// - Request ID
	// - Timestamp
	// - HTTP method
	// - HTTP status code
	// - Request path
	// - Request duration
	// - Request size
	// - Client IP address
	RequestLogger = zencore.RequestLogger
	// MaxBodySize is a middleware that limits the size of the request body.
	//
	// It reads the request body using http.MaxBytesReader.
	//
	// Parameters:
	//   - limit: The maximum size of the request body in bytes. (limit = 0 is same as not calling middleware)
	//
	// Example:
	//
	//	app.Use(zencore.MaxBodySize(10 * 1024 * 1024)) // 10MB limit
	//
	// Important: This middleware should be placed before any middleware that reads the request body.
	MaxBodySize = zencore.MaxBodySize
	// RateLimit is a middleware that rate limits requests to a handler.
	//
	// It uses a custom key function to determine the key for each request.
	//
	// This can allow ratelimiting of a specific route or endpoint request
	// Use [RateLimitIP] for IP-based rate limiting.
	RateLimit = zencore.RateLimit
	// NewRateLimiter creates a new RateLimiter.
	//
	// limit is the maximum number of requests allowed in a window.
	// window is the time duration for the rate limit.
	//
	// Safe defaults are used for cleanupPeriod.
	NewRateLimiter = zencore.NewRateLimiter
	AuthMiddleware = zencore.AuthMiddleware
	RequireAuth    = zencore.RequireAuth

	// --------------- Middlewares End ----------

	// Timeout returns a middleware that injects a timeout-aware context
	// into the request lifecycle.
	//
	// IMPORTANT:
	//
	// This middleware DOES NOT forcibly terminate handler execution.
	// Go does not support safely killing goroutines.
	//
	// Instead, this middleware:
	//
	//   - creates a request-scoped timeout context
	//   - propagates cancellation via c.Request.Context()
	//   - allows downstream services (DBs, HTTP clients, etc.) to stop work
	//   - preserves framework response lifecycle integrity
	//   - preserves panic recovery behavior
	//
	// Handlers and services are expected to respect:
	//
	//	select {
	//	case <-ctx.Done():
	//	    return
	//	default:
	//	}
	//
	// Example:
	//
	//	app.Use(zen.Timeout(5 * time.Second))
	//
	//	db.QueryContext(c.StdContext(), ...)
	//
	//	http.NewRequestWithContext(c.StdContext(), ...)
	Timeout = zencore.Timeout

	// GenerateETag generates an ETag header for the given body.
	//
	// It uses the SHA256 hash of the body to generate the ETag.
	GenerateETag = zencore.GenerateETag

	// --------------- Testing helpers ----------

	// PerformTestRequest performs a test request to the application.
	// It is a helper function for testing the application.
	//
	// Parameters:
	//   - app: The application to test.
	//   - method: The HTTP method to use.
	//   - path: The path to test.
	//   - body: The request body to use.
	//   - headers: The headers to use.
	//
	// Example:
	//
	//	rec := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
	PerformTestRequest = zencore.PerformTestRequest

	// NewTestContext creates a new test context.
	NewTestContext = zencore.NewTestContext

	HasStatus = zencore.HasStatus

	// --------------- Testing helpers End ----------
)

// GetService returns the service with the given name.
// It is a type-safe wrapper around the App Service function.
//
// Note:
// If the service is not found, it will panic.
// If the service type assertion fails, it will panic.
func GetService[T any](a *App, name string) T {
	return zencore.GetService[T](a, name)
}

// ------------------------- Testing Helpers -------------------------

// DecodeJSONResponse decodes the JSON response from the response recorder.
// It is a helper function for testing the application.
//
// Parameters:
//   - rec: The response recorder to decode the JSON response from.
//   - target: The target to decode the JSON response into.
//
// Example:
//
//	res := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
//	var body map[string]any
//	if err := zen.DecodeJSONResponse(res, &body); err != nil {
//		t.Fatal(err)
//	}
func DecodeJSONResponse(rec *httptest.ResponseRecorder, target any) error {
	return zencore.DecodeJSONResponse(rec, target)
}

// DecodeJSONResponseAs decodes the JSON response from the response recorder.
// It is a helper function for testing the application.
//
// Parameters:
//   - rec: The response recorder to decode the JSON response from.
//
// Example:
//
//	res := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
//	body, err := zen.DecodeJSONResponseAs[map[string]any](res)
//	if err != nil {
//		t.Fatal(err)
//	}
func DecodeJSONResponseAs[T any](rec *httptest.ResponseRecorder) (T, error) {
	return zencore.DecodeJSONResponseAs[T](rec)
}
