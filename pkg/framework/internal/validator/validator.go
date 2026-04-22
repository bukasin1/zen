package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func ValidateStruct(s any) ValidationErrors {
	var errs ValidationErrors

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errs
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldName := fieldType.Name
		rules := strings.Split(tag, ",")

		for _, rule := range rules {
			if err := applyRule(fieldVal, rule); err != nil {
				errs = append(errs, FieldError{
					Field:   fieldName,
					Message: err.Error(),
				})
			}
		}
	}

	return errs
}

func applyRule(v reflect.Value, rule string) error {
	switch {
	case rule == "required":
		if isZero(v) {
			return fmt.Errorf("is required")
		}

	case strings.HasPrefix(rule, "min="):
		return validateMin(v, rule)

	case strings.HasPrefix(rule, "max="):
		return validateMax(v, rule)

	case rule == "email":
		return validateEmail(v)
	}

	return nil
}

func isZero(v reflect.Value) bool {
	return v.IsZero()
}

func validateMin(v reflect.Value, rule string) error {
	// example: min=3
	// implement for string length / int
	return nil
}

func validateMax(v reflect.Value, rule string) error {
	return nil
}

func validateEmail(v reflect.Value) error {
	// simple check (upgrade later)
	str, ok := v.Interface().(string)
	if !ok || !strings.Contains(str, "@") {
		return fmt.Errorf("invalid email")
	}
	return nil
}
