package framework

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkRequest performs a benchmark request
// against the application's real ServeHTTP lifecycle.
func BenchmarkRequest(
	b *testing.B,
	app *App,
	method string,
	path string,
	body []byte,
) {
	b.Helper()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			method,
			path,
			bytes.NewReader(body),
		)

		recorder := httptest.NewRecorder()

		app.ServeHTTP(
			recorder,
			req,
		)
	}
}

// BenchmarkJSONRequest performs a JSON request benchmark.
func BenchmarkJSONRequest(
	b *testing.B,
	app *App,
	method string,
	path string,
	body []byte,
) {
	b.Helper()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			method,
			path,
			bytes.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

		recorder := httptest.NewRecorder()

		app.ServeHTTP(
			recorder,
			req,
		)
	}
}

// BenchmarkMiddlewareChain benchmarks
// middleware pipeline execution overhead.
func BenchmarkMiddlewareChain(
	b *testing.B,
	middlewares ...Middleware,
) {
	b.Helper()

	handler := HandlerFunc(func(ctx *Context) {})

	handler = chainMiddlewares(
		handler,
		middlewares,
	)

	app := New()

	app.Route("/benchmark").
		Get(handler)

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/benchmark",
		nil,
	)
}
