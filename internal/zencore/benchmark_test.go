package zencore

import (
	"net/http"
	"strconv"
	"testing"
)

func BenchmarkStaticRoute(
	b *testing.B,
) {
	app := New()

	app.Route("/health").
		Get(func(ctx *Context) {
			ctx.SuccessOK("ok")
		})

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/health",
		nil,
	)
}

func BenchmarkDynamicRoute(
	b *testing.B,
) {
	app := New()

	app.Route("/users/{id}").
		Get(func(ctx *Context) {
			ctx.SuccessOK(
				ctx.Param("id"),
			)
		})

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/users/123",
		nil,
	)
}

func BenchmarkJSONResponse(
	b *testing.B,
) {
	app := New()

	app.Route("/json").
		Get(func(ctx *Context) {
			ctx.JSON(http.StatusOK, map[string]any{
				"message": "hello",
				"success": true,
			})
		})

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/json",
		nil,
	)
}

func BenchmarkJSONBinding(
	b *testing.B,
) {
	type Payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	app := New()

	app.Route("/users").
		Post(func(ctx *Context) {
			var payload Payload

			_ = ctx.BindJSON(&payload)

			ctx.SuccessOK(payload)
		})

	body := []byte(`{
		"name":"john",
		"age":30
	}`)

	BenchmarkJSONRequest(
		b,
		app,
		http.MethodPost,
		"/users",
		body,
	)
}

func BenchmarkMiddlewarePipeline_1(
	b *testing.B,
) {
	middleware := func(
		next HandlerFunc,
	) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}

	BenchmarkMiddlewareChain(
		b,
		middleware,
	)
}

func BenchmarkMiddlewarePipeline_5(
	b *testing.B,
) {
	middleware := func(
		next HandlerFunc,
	) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}

	BenchmarkMiddlewareChain(
		b,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
	)
}

func BenchmarkMiddlewarePipeline_10(
	b *testing.B,
) {
	middleware := func(
		next HandlerFunc,
	) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}

	BenchmarkMiddlewareChain(
		b,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
		middleware,
	)
}

func BenchmarkRouteLookup_100(
	b *testing.B,
) {
	app := New()

	for i := 0; i < 100; i++ {
		path := "/route/" + strconv.Itoa(i)

		app.Route(path).
			Get(func(ctx *Context) {
				ctx.SuccessOK("ok")
			})
	}

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/route/5",
		nil,
	)
}

func BenchmarkRouteLookup_1000(
	b *testing.B,
) {
	app := New()

	for i := 0; i < 1000; i++ {
		path := "/route/" + strconv.Itoa(i)

		app.Route(path).
			Get(func(ctx *Context) {
				ctx.SuccessOK("ok")
			})
	}

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/route/5",
		nil,
	)
}
