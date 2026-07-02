package zencore

import (
	"net/http"
	"strconv"
	"testing"
)

func BenchmarkRouterStaticRoutes_1000(
	b *testing.B,
) {
	app := New()

	for i := 0; i < 1000; i++ {
		path := "/static/" +
			strconv.Itoa(i)

		app.Route(path).
			Get(func(ctx *Context) {})
	}

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/static/500",
		nil,
	)
}

func BenchmarkRouterDynamicRoutes_1000(
	b *testing.B,
) {
	app := New()

	for i := 0; i < 1000; i++ {
		path := "/users/" +
			strconv.Itoa(i) +
			"/posts/{id}"

		app.Route(path).
			Get(func(ctx *Context) {})
	}

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/users/500/posts/10",
		nil,
	)
}

func BenchmarkRouterDeepPaths(
	b *testing.B,
) {
	app := New()

	app.Route(
		"/api/v1/users/{id}/posts/{postId}/comments/{commentId}",
	).Get(func(ctx *Context) {})

	BenchmarkRequest(
		b,
		app,
		http.MethodGet,
		"/api/v1/users/1/posts/2/comments/3",
		nil,
	)
}
