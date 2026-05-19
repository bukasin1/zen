package framework

import "testing"

func noopMiddleware(
	next HandlerFunc,
) HandlerFunc {
	return func(ctx *Context) {
		next(ctx)
	}
}

func BenchmarkMiddlewareChain_Empty(
	b *testing.B,
) {
	BenchmarkMiddlewareChain(b)
}

func BenchmarkMiddlewareChain_1(
	b *testing.B,
) {
	BenchmarkMiddlewareChain(
		b,
		noopMiddleware,
	)
}

func BenchmarkMiddlewareChain_20(
	b *testing.B,
) {
	middlewares := make(
		[]Middleware,
		20,
	)

	for i := range middlewares {
		middlewares[i] = noopMiddleware
	}

	BenchmarkMiddlewareChain(
		b,
		middlewares...,
	)
}
