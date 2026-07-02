package zencore

import "strings"

// cloneRouteMetadata creates a shallow copy of route metadata.
// Nested reference values are not deeply cloned.
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

// cloneRouteDefinition creates a safe copy of a route definition.
// Metadata values themselves are NOT deeply cloned.
// Nested reference types inside metadata should be treated as immutable.
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

func normalizeHTTPMethod(method string) string {
	return strings.ToUpper(strings.TrimSpace(method))
}
