package zencore

type RouteDoc struct {
	Method string `json:"method"`
	Path   string `json:"path"`

	Name        string   `json:"name,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Version     string   `json:"version,omitempty"`
	OperationID string   `json:"operationId,omitempty"`
	Deprecated  bool     `json:"deprecated,omitempty"`
	Internal    bool     `json:"internal,omitempty"`
	Middlewares []string `json:"middlewares,omitempty"`
}

func buildRouteDoc(
	route RouteDefinition,
) RouteDoc {
	return RouteDoc{
		Method: route.Method,
		Path:   route.Path,

		Name: route.Name,

		Summary: route.Summary(),

		Description: route.Description(),

		Tags: route.Tags(),

		Version: route.Version(),

		OperationID: route.OperationID(),

		Deprecated: route.IsDeprecated(),

		Internal: route.IsInternal(),

		Middlewares: append(
			[]string(nil),
			route.Middlewares...,
		),
	}
}

type RouteDocOptions struct {
	IncludeInternal bool
}

// RouteDocs returns lightweight route documentation.
func (a *App) RouteDocs(
	options ...RouteDocOptions,
) []RouteDoc {
	var opts RouteDocOptions

	if len(options) > 0 {
		opts = options[0]
	}

	docs := make([]RouteDoc, 0)

	for _, route := range a.routeRegistry {
		if route.IsInternal() &&
			!opts.IncludeInternal {
			continue
		}

		docs = append(
			docs,
			buildRouteDoc(route),
		)
	}

	return docs
}
