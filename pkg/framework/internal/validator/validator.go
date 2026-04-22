package validator

import (
	"fmt"
	"reflect"
	"strconv"
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

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldName := fieldType.Name

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			name := strings.Split(jsonTag, ",")[0]
			if name != "" && name != "-" {
				fieldName = name
			}
		}

		msgTag := fieldType.Tag.Get("msg")

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			if err := applyRule(fieldVal, rule); err != nil {
				message := err.Error()
				if msgTag != "" {
					message = msgTag
				}
				errs = append(errs, FieldError{
					Field:   fieldName,
					Message: message,
				})
				break
			}
		}
	}

	return errs
}

func applyRule(v reflect.Value, rule string) error {
	// allow optional fields to pass validation if the rule is not "required"
	// and the field is zero-valued
	if isZero(v) && rule != "required" {
		return nil
	}

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
	minStr := strings.TrimPrefix(rule, "min=")
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf(`invalid min rule: "%v". should be a number`, minStr)
	}

	switch v.Kind() {
	case reflect.String:
		if v.Len() < min {
			return fmt.Errorf("must be at least %d characters long. characters count is %d", min, v.Len())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < int64(min) {
			return fmt.Errorf("must be at least %d. value is %d", min, v.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() < uint64(min) {
			return fmt.Errorf("must be at least %d. value is %d", min, v.Uint())
		}
	default:
		return fmt.Errorf("min validation not supported for type %s", v.Kind())
	}

	return nil
}

func validateMax(v reflect.Value, rule string) error {
	maxStr := strings.TrimPrefix(rule, "max=")
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf(`invalid max rule: "%v". should be a number`, maxStr)
	}

	switch v.Kind() {
	case reflect.String:
		if v.Len() > max {
			return fmt.Errorf("must be at most %d characters long", max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > int64(max) {
			return fmt.Errorf("must be at most %d", max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() > uint64(max) {
			return fmt.Errorf("must be at most %d", max)
		}
	default:
		return fmt.Errorf("max validation not supported for type %s", v.Kind())
	}

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
