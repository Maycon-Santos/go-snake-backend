package validator

import (
	"fmt"
	"reflect"
	"strings"
)

type validatorConfig struct {
	field      string
	required   *bool
	minLen     *uint
	maxLen     *uint
	noContains []string
}

type errType = string

type validator = func(config validatorConfig, field string, value interface{}) (errType errType, err error)

var validators []validator

const (
	Required   errType = "required"
	MinLen     errType = "minLen"
	MaxLen     errType = "maxLen"
	NoContains errType = "noContains"
)

func init() {
	validators = append(
		validators,
		required,
		minLen,
		maxLen,
		noContains,
	)
}

func required(config validatorConfig, field string, value interface{}) (errType, error) {
	if config.required == nil {
		return "", nil
	}

	reflectValue := reflect.ValueOf(value)

	if reflectValue.Interface() == reflect.Zero(reflectValue.Type()).Interface() {
		return Required, fmt.Errorf("missing %s field", field)
	}

	return "", nil
}

func minLen(config validatorConfig, field string, value interface{}) (errType, error) {
	if config.minLen == nil {
		return "", nil
	}

	reflectValue := reflect.ValueOf(value)
	valueKind := reflectValue.Kind()

	if valueKind != reflect.Array &&
		valueKind != reflect.Slice &&
		valueKind != reflect.String &&
		valueKind != reflect.Map {
		panic("unsupported type. Must be string, slice, array or map")
	}

	if uint(reflectValue.Len()) < *config.minLen {
		return MinLen, fmt.Errorf("the %s field must be greater than %d", field, *config.minLen)
	}

	return "", nil
}

func maxLen(config validatorConfig, field string, value interface{}) (errType, error) {
	if config.maxLen == nil {
		return "", nil
	}

	reflectValue := reflect.ValueOf(value)
	valueKind := reflectValue.Kind()

	if valueKind != reflect.Array &&
		valueKind != reflect.Slice &&
		valueKind != reflect.String &&
		valueKind != reflect.Map {
		panic("unsupported type. Must be string, slice, array or map")
	}

	if uint(reflectValue.Len()) > *config.maxLen {
		return MaxLen, fmt.Errorf("the %s field must be less than %d", field, *config.maxLen)
	}

	return "", nil
}

func noContains(config validatorConfig, field string, value interface{}) (errType, error) {
	if config.noContains == nil {
		return "", nil
	}

	valueString := fmt.Sprint(value)

	for _, substr := range config.noContains {
		if strings.Contains(valueString, substr) {
			return NoContains, fmt.Errorf("%s field cannot have the following characters: %s", field, strings.Join(config.noContains, ", "))
		}
	}

	return "", nil
}
