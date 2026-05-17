package framework

func cloneRouteMetadata(metadata RouteMetadata) RouteMetadata {
	if metadata == nil {
		return nil
	}

	cloned := make(RouteMetadata, len(metadata))

	for key, value := range metadata {
		cloned[key] = value
	}

	return cloned
}

func cloneRouteDefinition(route RouteDefinition) RouteDefinition {
	cloned := route

	if route.Middlewares != nil {
		cloned.Middlewares = append([]string(nil), route.Middlewares...)
	}

	if route.Metadata != nil {
		cloned.Metadata = make(RouteMetadata, len(route.Metadata))

		for key, value := range route.Metadata {
			cloned.Metadata[key] = value
		}
	}

	return cloned
}
