package framework

func (r RouteDefinition) Tags() []string {
	value, ok := r.Metadata[RouteMetadataKeyTags]

	if !ok {
		return nil
	}

	tags, ok := value.([]string)

	if !ok {
		return nil
	}

	cloned := make([]string, len(tags))
	copy(cloned, tags)

	return cloned
}

func (r RouteDefinition) Summary() string {
	value, ok := r.Metadata[RouteMetadataKeySummary]

	if !ok {
		return ""
	}

	summary, ok := value.(string)

	if !ok {
		return ""
	}

	return summary
}

func (r RouteDefinition) Description() string {
	value, ok := r.Metadata[RouteMetadataKeyDescription]

	if !ok {
		return ""
	}

	description, ok := value.(string)

	if !ok {
		return ""
	}

	return description
}

func (r RouteDefinition) IsDeprecated() bool {
	value, ok := r.Metadata[RouteMetadataKeyDeprecated]

	if !ok {
		return false
	}

	deprecated, ok := value.(bool)

	if !ok {
		return false
	}

	return deprecated
}

func (r RouteDefinition) IsInternal() bool {
	value, ok := r.Metadata[RouteMetadataKeyInternal]

	if !ok {
		return false
	}

	internal, ok := value.(bool)

	if !ok {
		return false
	}

	return internal
}

func (r RouteDefinition) Version() string {
	value, ok := r.Metadata[RouteMetadataKeyVersion]

	if !ok {
		return ""
	}

	version, ok := value.(string)

	if !ok {
		return ""
	}

	return version
}

func (r RouteDefinition) OperationID() string {
	value, ok := r.Metadata[RouteMetadataKeyOperationID]

	if !ok {
		return ""
	}

	id, ok := value.(string)

	if !ok {
		return ""
	}

	return id
}
