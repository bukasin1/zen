package framework

import "strings"

const (
	RouteMetadataKeyTags        = "tags"
	RouteMetadataKeySummary     = "summary"
	RouteMetadataKeyDescription = "description"
	RouteMetadataKeyDeprecated  = "deprecated"
	RouteMetadataKeyInternal    = "internal"
	RouteMetadataKeyVersion     = "version"
	RouteMetadataKeyOperationID = "operationId"
)

// Tags adds documentation or grouping tags to the route.
func (rb *RouteBuilder) Tags(
	tags ...string,
) *RouteBuilder {
	cleaned := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			continue
		}

		cleaned = append(cleaned, tag)
	}

	if len(cleaned) == 0 {
		return rb
	}

	rb.metadata[RouteMetadataKeyTags] = cleaned

	return rb
}

// Summary adds a short summary for the route.
func (rb *RouteBuilder) Summary(
	summary string,
) *RouteBuilder {
	summary = strings.TrimSpace(summary)

	if summary == "" {
		return rb
	}

	rb.metadata[RouteMetadataKeySummary] = summary

	return rb
}

// Description adds a detailed route description.
func (rb *RouteBuilder) Description(
	description string,
) *RouteBuilder {
	description = strings.TrimSpace(description)

	if description == "" {
		return rb
	}

	rb.metadata[RouteMetadataKeyDescription] = description

	return rb
}

// Deprecated marks the route as deprecated.
func (rb *RouteBuilder) Deprecated() *RouteBuilder {
	rb.metadata[RouteMetadataKeyDeprecated] = true

	return rb
}

// Internal marks the route as internal-only.
func (rb *RouteBuilder) Internal() *RouteBuilder {
	rb.metadata[RouteMetadataKeyInternal] = true

	return rb
}

// Version sets the route version metadata.
func (rb *RouteBuilder) Version(
	version string,
) *RouteBuilder {
	version = strings.TrimSpace(version)

	if version == "" {
		return rb
	}

	rb.metadata[RouteMetadataKeyVersion] = version

	return rb
}

// OperationID sets a unique operation identifier.
func (rb *RouteBuilder) OperationID(
	id string,
) *RouteBuilder {
	id = strings.TrimSpace(id)

	if id == "" {
		return rb
	}

	rb.metadata[RouteMetadataKeyOperationID] = id

	return rb
}
