package zencore

import "net/http"

type Group struct {
	prefix      string
	app         *App
	middlewares []Middleware

	middlewareDefinitions []MiddlewareDefinition
}

func (app *App) Group(prefix string) *Group {
	return &Group{
		prefix:      prefix,
		app:         app,
		middlewares: []Middleware{},
	}
}

// Allow for nested route groups
func (g *Group) Group(prefix string) *Group {
	return &Group{
		prefix: g.prefix + prefix,
		app:    g.app,
		// pass on the parent group middlewares
		middlewares:           append([]Middleware{}, g.middlewares...),
		middlewareDefinitions: append([]MiddlewareDefinition{}, g.middlewareDefinitions...),
	}
}

// Use adds middleware to the group.
// Middlewares added to a group will be applied to all handlers in that group.
// Call this before any route definitions (Get, Post, etc).
func (g *Group) Use(m ...Middleware) {
	g.middlewares = append(g.middlewares, m...)
}

// UseNamed adds named middleware to the group.
// Middlewares added to a group will be applied to all handlers in that group.
//
// Call this before any route definitions (Get, Post, etc).
//
// Example:
//
//	loggerMiddleware := zencore.NamedMiddleware("logger", zencore.Logger())
//	api := app.Group("/api")
//	api.UseNamed(loggerMiddleware)
func (g *Group) UseNamed(mds ...MiddlewareDefinition) {
	for _, md := range mds {
		g.middlewares = append(g.middlewares, md.Func)

		g.middlewareDefinitions = append(g.middlewareDefinitions, md)
	}
}

// Route returns a new RouteBuilder for the given path.
// It is used to define a new route under the current group.
// Example:
//
//	api := app.Group("/api")
//	userRoutes := api.Group("/users")
//	userRoutes.Route("/{id}").Get(getUser)
func (g *Group) Route(path string) *RouteBuilder {
	return &RouteBuilder{
		app:   g.app,
		path:  g.prefix + path,
		group: g,
		// pass on the group middlewares
		middlewares:           append([]Middleware{}, g.middlewares...),
		middlewareDefinitions: append([]MiddlewareDefinition{}, g.middlewareDefinitions...),
		metadata:              make(RouteMetadata),
	}
}

func (g *Group) Get(path string, handler HandlerFunc) {
	fullPath, wrapped := g.wrapPathMiddleware(path, handler)
	g.app.router.Handle(http.MethodGet, fullPath, wrapped)
}

func (g *Group) Post(path string, handler HandlerFunc) {
	fullPath, wrapped := g.wrapPathMiddleware(path, handler)
	g.app.router.Handle(http.MethodPost, fullPath, wrapped)
}

func (g *Group) wrapPathMiddleware(path string, handler HandlerFunc) (string, HandlerFunc) {
	return g.prefix + path, chainMiddlewares(handler, g.middlewares)
}
