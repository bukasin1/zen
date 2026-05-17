package framework

type RouteMetadata map[string]any

type RouteDefinition struct {
	Method string
	Path   string

	Name string

	HandlerName string

	Middlewares []string

	Metadata RouteMetadata
}

func (a *App) registerRoute(def RouteDefinition) {
	a.routeRegistry = append(a.routeRegistry, def)
}

func (a *App) Routes() []RouteDefinition {
	routes := make([]RouteDefinition, len(a.routeRegistry))

	for i, route := range a.routeRegistry {
		routes[i] = cloneRouteDefinition(route)
	}

	return routes
}
