package zencore

import "strings"

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
	for _, existing := range a.routeRegistry {
		if existing.Method == def.Method && existing.Path == def.Path {
			panic(
				newFrameworkPanic("duplicate route registration: " + def.Method + " " + def.Path),
			)
		}
	}

	a.routeRegistry = append(a.routeRegistry, def)
}

func (a *App) Routes() []RouteDefinition {
	routes := make([]RouteDefinition, len(a.routeRegistry))

	for i, route := range a.routeRegistry {
		routes[i] = cloneRouteDefinition(route)
	}

	return routes
}

// RouteByName returns a route definition by its name.
func (a *App) RouteByName(name string) (RouteDefinition, bool) {
	name = strings.TrimSpace(name)

	if name == "" {
		return RouteDefinition{}, false
	}

	for _, route := range a.routeRegistry {
		if route.Name == name {
			return cloneRouteDefinition(route), true
		}
	}

	return RouteDefinition{}, false
}

// RoutesByMethod returns all routes registered for an HTTP method.
func (a *App) RoutesByMethod(method string) []RouteDefinition {
	method = normalizeHTTPMethod(method)

	if method == "" {
		return nil
	}

	routes := make([]RouteDefinition, 0)

	for _, route := range a.routeRegistry {
		if route.Method == method {
			routes = append(routes, cloneRouteDefinition(route))
		}
	}

	return routes
}

// RoutesByPath returns all routes matching a path.
func (a *App) RoutesByPath(path string) []RouteDefinition {
	path = normalizeRoutePath(path)

	routes := make([]RouteDefinition, 0)

	for _, route := range a.routeRegistry {
		if route.Path == path {
			routes = append(routes, cloneRouteDefinition(route))
		}
	}

	return routes
}

// HasRoute checks whether a route exists.
func (a *App) HasRoute(method string, path string) bool {
	method = normalizeHTTPMethod(method)
	path = normalizeRoutePath(path)

	for _, route := range a.routeRegistry {
		if route.Method == method && route.Path == path {
			return true
		}
	}

	return false
}
