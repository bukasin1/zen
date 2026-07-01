package framework

import "testing"

func TestAppRoutes(
	t *testing.T,
) {
	app := New()

	authMiddleware := NamedMiddleware(
		"auth",
		func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) {
				next(ctx)
			}
		},
	)

	app.Route("/users").
		Name("users.list").
		Meta("tag", "Users").
		UseNamed(authMiddleware).
		Get(func(ctx *Context) {})

	routes := app.Routes()

	if len(routes) != 1 {
		t.Fatalf(
			"expected 1 route, got %d",
			len(routes),
		)
	}

	route := routes[0]

	if route.Method != "GET" {
		t.Fatalf(
			"expected GET method, got %s",
			route.Method,
		)
	}

	if route.Path != "/users" {
		t.Fatalf(
			"expected /users path, got %s",
			route.Path,
		)
	}

	if route.Name != "users.list" {
		t.Fatalf(
			"expected route name users.list, got %s",
			route.Name,
		)
	}

	if len(route.Middlewares) != 1 {
		t.Fatalf(
			"expected 1 middleware, got %d",
			len(route.Middlewares),
		)
	}

	if route.Middlewares[0] != "auth" {
		t.Fatalf(
			"expected auth middleware, got %s",
			route.Middlewares[0],
		)
	}

	tag, ok := route.Metadata["tag"]

	if !ok {
		t.Fatal("expected tag metadata")
	}

	if tag != "Users" {
		t.Fatalf(
			"expected Users tag, got %v",
			tag,
		)
	}
}

func TestRouteByName(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Name("users.list").
		Get(func(ctx *Context) {})

	route, found := app.RouteByName(
		"users.list",
	)

	if !found {
		t.Fatal("expected route to exist")
	}

	if route.Path != "/users" {
		t.Fatalf(
			"expected /users path, got %s",
			route.Path,
		)
	}
}

func TestRoutesByMethod(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Get(func(ctx *Context) {})

	app.Route("/posts").
		Get(func(ctx *Context) {})

	app.Route("/login").
		Post(func(ctx *Context) {})

	routes := app.RoutesByMethod("GET")

	if len(routes) != 2 {
		t.Fatalf(
			"expected 2 GET routes, got %d",
			len(routes),
		)
	}
}

func TestHasRoute(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Get(func(ctx *Context) {})

	if !app.HasRoute(
		"GET",
		"/users",
	) {
		t.Fatal(
			"expected route to exist",
		)
	}

	if app.HasRoute(
		"POST",
		"/users",
	) {
		t.Fatal(
			"expected route not to exist",
		)
	}
}

func TestDuplicateRouteRegistrationPanics(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Get(func(ctx *Context) {})

	defer func() {
		if recover() == nil {
			t.Fatal(
				"expected duplicate route panic",
			)
		}
	}()

	app.Route("/users").
		Get(func(ctx *Context) {})
}

func TestRouteMetadataHelpers(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Name("users.list").
		Tags("Users", "Admin").
		Summary("List users").
		Description("Returns all users").
		Version("v1").
		OperationID("listUsers").
		Deprecated().
		Internal().
		Get(func(ctx *Context) {})

	route, found := app.RouteByName(
		"users.list",
	)

	if !found {
		t.Fatal("expected route")
	}

	tags := route.Tags()

	if len(tags) != 2 {
		t.Fatalf(
			"expected 2 tags, got %d",
			len(tags),
		)
	}

	if route.Summary() != "List users" {
		t.Fatal(
			"unexpected summary",
		)
	}

	if route.Description() != "Returns all users" {
		t.Fatal(
			"unexpected description",
		)
	}

	if route.Version() != "v1" {
		t.Fatal(
			"unexpected version",
		)
	}

	if route.OperationID() != "listUsers" {
		t.Fatal(
			"unexpected operation id",
		)
	}

	if !route.IsDeprecated() {
		t.Fatal(
			"expected deprecated route",
		)
	}

	if !route.IsInternal() {
		t.Fatal(
			"expected internal route",
		)
	}
}

func TestRouteDocs(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Name("users.list").
		Tags("Users").
		Summary("List users").
		Version("v1").
		Get(func(ctx *Context) {})

	app.Route("/internal").
		Internal().
		Get(func(ctx *Context) {})

	docs := app.RouteDocs()

	if len(docs) != 1 {
		t.Fatalf(
			"expected 1 public route doc, got %d",
			len(docs),
		)
	}

	doc := docs[0]

	if doc.Path != "/users" {
		t.Fatalf(
			"expected /users path, got %s",
			doc.Path,
		)
	}

	if doc.Summary != "List users" {
		t.Fatal(
			"unexpected summary",
		)
	}

	if len(doc.Tags) != 1 {
		t.Fatal(
			"expected tags",
		)
	}
}

func TestRouteDocsIncludeInternal(
	t *testing.T,
) {
	app := New()

	app.Route("/users").
		Get(func(ctx *Context) {})

	app.Route("/internal").
		Internal().
		Get(func(ctx *Context) {})

	docs := app.RouteDocs(
		RouteDocOptions{
			IncludeInternal: true,
		},
	)

	if len(docs) != 2 {
		t.Fatalf(
			"expected 2 route docs, got %d",
			len(docs),
		)
	}
}
