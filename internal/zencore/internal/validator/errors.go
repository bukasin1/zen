package validator

import (
	"strings"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	var errStr []string
	for _, err := range v {
		errStr = append(errStr, err.Field+": "+err.Message)
	}
	return "Validation failed: " + strings.Join(errStr, "; ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}
