// Package zen is the public API facade for the Zen backend framework.
//
// Zen is a lightweight, production-oriented Go backend framework built on
// top of the Go standard library. This package re-exports the stable,
// public surface of the framework's internal implementation
// (internal/zencore) as types, constants, constructors, configuration
// helpers, middleware, utilities, and testing helpers.
//
// Consumers of the framework should only ever import this package; the
// internal/zencore package is not part of the public API and may change
// without notice.
package zen

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bukasin1/zen/internal/zencore"
	"github.com/bukasin1/zen/internal/zencore/errors"
	"github.com/bukasin1/zen/internal/zencore/logger"
)

// ------------------------------------------------------------------------
// Type aliases
// ------------------------------------------------------------------------

type (
	App                  = zencore.App
	Context              = zencore.Context
	AppError             = errors.AppError
	Middleware           = zencore.Middleware
	Config               = zencore.Config
	Logger               = logger.Logger
	HTTPConfig           = zencore.HTTPConfig
	LogConfig            = zencore.LogConfig
	CORSConfig           = zencore.CORSConfig
	MiddlewareDefinition = zencore.MiddlewareDefinition

	// HandlerFunc is a function that handles an HTTP request.
	// It takes a *Context as an argument and returns nothing.
	HandlerFunc     = zencore.HandlerFunc
	RouteDocOptions = zencore.RouteDocOptions
)

// ------------------------------------------------------------------------
// Error constants
// ------------------------------------------------------------------------

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

// ------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------

// New creates a new Zen application instance.
func New() *App {
	return zencore.New()
}

// NewConsoleLogger creates a new console logger.
//
// When pretty is true, log output is formatted for human readability;
// otherwise, structured JSON output is used.
func NewConsoleLogger(pretty bool) *logger.ConsoleLogger {
	return logger.NewConsoleLogger(pretty)
}

// NewDevConsoleLogger creates a new development-oriented console logger.
func NewDevConsoleLogger() *logger.DevConsoleLogger {
	return logger.NewDevConsoleLogger()
}

// GetService returns the service with the given name.
// It is a type-safe wrapper around the App Service function.
//
// Note:
// If the service is not found, it will panic.
// If the service type assertion fails, it will panic.
func GetService[T any](a *App, name string) T {
	return zencore.GetService[T](a, name)
}

// ------------------------------------------------------------------------
// Configuration
// ------------------------------------------------------------------------

// LoadConfigFromEnv loads app config from environment variables.
// It uses [DefaultConfig] as a fallback.
// Environment variables are case sensitive and should be uppercase:
//
//	APP_NAME="Application name" (default: "Zen")
//	APP_ENV="Application environment" (default: "development")
//	HTTP_ADDR="HTTP server address" (default: ":8080")
//	HTTP_SHUTDOWN_TIMEOUT="HTTP server shutdown timeout" (default: 2s)
//	LOG_LEVEL="Log level" (default: "debug")
//	LOG_PRETTY="Enable pretty printing" (default: true)
//	LOG_ENABLE_JSON="Enable JSON logging (default: true)
func LoadConfigFromEnv() Config {
	return zencore.LoadConfigFromEnv()
}

// DefaultConfig returns a default configuration for the application.
func DefaultConfig() Config {
	return zencore.DefaultConfig()
}

// DefaultCORSConfig returns the default CORS configuration.
// It allows all origins, methods, and headers.
// It also sets the MaxAge to 300 seconds (5 minutes).
// Use [DefaultCORSConfig] when you want to allow all origins, methods, and headers.
func DefaultCORSConfig() CORSConfig {
	return zencore.DefaultCORSConfig()
}

// ------------------------------------------------------------------------
// Environment helpers
// ------------------------------------------------------------------------

// GetEnv returns the value of the environment variable named by key, or
// fallback if the variable is not set.
func GetEnv(key string, fallback string) string {
	return zencore.GetEnv(key, fallback)
}

// GetEnvInt returns the integer value of the environment variable named by
// key, or fallback if the variable is not set or cannot be parsed.
func GetEnvInt(key string, fallback int) int {
	return zencore.GetEnvInt(key, fallback)
}

// GetEnvBool returns the boolean value of the environment variable named by
// key, or fallback if the variable is not set or cannot be parsed.
func GetEnvBool(key string, fallback bool) bool {
	return zencore.GetEnvBool(key, fallback)
}

// GetEnvDuration returns the duration value of the environment variable
// named by key, or fallback if the variable is not set or cannot be parsed.
func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	return zencore.GetEnvDuration(key, fallback)
}

// MustGetEnv returns the value of the environment variable named by key.
// It panics if the variable is not set.
func MustGetEnv(key string) string {
	return zencore.MustGetEnv(key)
}

// MustGetEnvInt returns the integer value of the environment variable named
// by key. It panics if the variable is not set or cannot be parsed.
func MustGetEnvInt(key string) int {
	return zencore.MustGetEnvInt(key)
}

// MustGetEnvBool returns the boolean value of the environment variable named
// by key. It panics if the variable is not set or cannot be parsed.
func MustGetEnvBool(key string) bool {
	return zencore.MustGetEnvBool(key)
}

// MustGetEnvDuration returns the duration value of the environment variable
// named by key. It panics if the variable is not set or cannot be parsed.
func MustGetEnvDuration(key string) time.Duration {
	return zencore.MustGetEnvDuration(key)
}

// ------------------------------------------------------------------------
// Middleware
// ------------------------------------------------------------------------

// NamedMiddleware returns a MiddlewareDefinition for the given name and middleware.
//
// Example:
//
//	loggerMiddleware := zencore.NamedMiddleware("logger", zencore.Logger())
//	app.UseNamed(loggerMiddleware)
func NamedMiddleware(name string, middleware Middleware) MiddlewareDefinition {
	return zencore.NamedMiddleware(name, middleware)
}

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
func CORS(config CORSConfig) Middleware {
	return zencore.CORS(config)
}

// Recovery is a middleware that recovers from panics and returns a 500 error.
func Recovery() Middleware {
	return zencore.Recovery()
}

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
func RequestLogger() Middleware {
	return zencore.RequestLogger()
}

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
func MaxBodySize(limit int64) Middleware {
	return zencore.MaxBodySize(limit)
}

// GzipCompression is a middleware that compresses response bodies using gzip.
func GzipCompression() Middleware {
	return zencore.GzipCompression()
}

// RateLimit is a middleware that rate limits requests to a handler.
//
// It uses a custom key function to determine the key for each request.
//
// This can allow ratelimiting of a specific route or endpoint request
// Use [RateLimitIP] for IP-based rate limiting.
func RateLimit(rl *zencore.RateLimiter, keyFn func(*Context) string) Middleware {
	return zencore.RateLimit(rl, keyFn)
}

// RateLimitIP creates a middleware that rate limits requests based on IP address.
func RateLimitIP(rl *zencore.RateLimiter) Middleware {
	return zencore.RateLimitIP(rl)
}

// NewRateLimiter creates a new RateLimiter.
//
// limit is the maximum number of requests allowed in a window.
// window is the time duration for the rate limit.
//
// Safe defaults are used for cleanupPeriod.
func NewRateLimiter(limit int, window time.Duration) *zencore.RateLimiter {
	return zencore.NewRateLimiter(limit, window)
}

// AuthMiddleware returns a middleware that authenticates requests using the
// given token validator.
func AuthMiddleware(validator zencore.TokenValidator) Middleware {
	return zencore.AuthMiddleware(validator)
}

// RequireAuth returns a middleware that requires a request to be authenticated.
func RequireAuth() Middleware {
	return zencore.RequireAuth()
}

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
func Timeout(duration time.Duration) Middleware {
	return zencore.Timeout(duration)
}

// ------------------------------------------------------------------------
// Utilities
// ------------------------------------------------------------------------

// GenerateETag generates an ETag header for the given body.
//
// It uses the SHA256 hash of the body to generate the ETag.
func GenerateETag(data []byte) string {
	return zencore.GenerateETag(data)
}

// ------------------------------------------------------------------------
// Testing helpers
// ------------------------------------------------------------------------

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
func PerformTestRequest(
	handler http.Handler,
	method string,
	path string,
	body []byte,
	headers map[string]string,
) *httptest.ResponseRecorder {
	return zencore.PerformTestRequest(handler, method, path, body, headers)
}

// PerformTestRequestFromRequest performs a test request to the application from a request.
// It is a helper function for testing the application.
//
// Parameters:
//   - app: The application to test.
//   - req: The request to test.
//
// Example:
//
//	rec := zen.PerformTestRequestFromRequest(app, req)
func PerformTestRequestFromRequest(
	handler http.Handler,
	req *http.Request,
) *httptest.ResponseRecorder {
	return zencore.PerformTestRequestFromRequest(handler, req)
}

// PerformTestJSONRequest performs a test request to the application.
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
//	rec := zen.PerformTestJSONRequest(app, "GET", "/health", nil, nil)
func PerformTestJSONRequest(
	handler http.Handler,
	method string,
	path string,
	body []byte,
	headers map[string]string,
) *httptest.ResponseRecorder {
	return zencore.PerformTestJSONRequest(handler, method, path, body, headers)
}

// NewTestingContext creates a new test context.
func NewTestingContext(
	method string,
	path string,
	body []byte,
) (*Context, *httptest.ResponseRecorder) {
	return zencore.NewTestingContext(method, path, body)
}

// HasHeader reports whether the response recorder's header matches the
// given key and expected value.
//
// Parameters:
//   - rec: The response recorder to check.
//   - key: The header key to check.
//   - expected: The expected header value.
//
// Example:
//
//	res := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
//	if !zen.HasHeader(res, "Content-Type", "application/json") {
//		t.Fatalf("expected header Content-Type: application/json, got %s", res.Header().Get("Content-Type"))
//	}
func HasHeader(
	rec *httptest.ResponseRecorder,
	key string,
	expected string,
) bool {
	return zencore.HasHeader(rec, key, expected)
}

// HasStatus reports whether the response recorder's status code matches the
// given status.
//
// Parameters:
//   - rec: The response recorder to check.
//   - status: The status code to check.
//
// Example:
//
//	res := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
//	if !zen.HasStatus(res, http.StatusOK) {
//		t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
//	}
func HasStatus(
	rec *httptest.ResponseRecorder,
	status int,
) bool {
	return zencore.HasStatus(rec, status)
}

// ResponseBody returns the response body from the response recorder.
// It is a helper function for testing the application.
//
// Parameters:
//   - rec: The response recorder to get the response body from.
//
// Example:
//
//	rec := zen.PerformTestRequest(app, "GET", "/health", nil, nil)
//	body := zen.ResponseBody(rec)
func ResponseBody(rec *httptest.ResponseRecorder) string {
	return zencore.ResponseBody(rec)
}

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
