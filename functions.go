package utils

import (
	"encoding/json"
	"go/parser"
	"reflect"

	validate "github.com/go-playground/validator/v10"
	"github.com/palantir/stacktrace"
)

// ExtractValidationRules extracts validation rules from the given input struct.
// It iterates over the fields of the struct and checks for the presence of "json" and "validate" tags.
// If a "json" tag is found, it maps the tag value to the corresponding "validate" tag value in the result map.
// If a field is of type struct, it recursively calls ExtractValidationRules to extract validation rules from the nested struct fields.
// The extracted validation rules are returned as a map[string]string, where the key represents the "json" tag value and the value represents the "validate" tag value.
func ExtractValidationRules(input interface{}) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(input)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")
		if jsonTag != "" {
			result[jsonTag] = validateTag
		}
		if field.Type.Kind() == reflect.Struct {
			nestedResult := ExtractValidationRules(reflect.New(field.Type).Elem().Interface())
			for k, v := range nestedResult {
				result[k] = v
			}
		}
	}
	return result
}

// NewValidator creates a new instance of the validator.Validate struct.
// It registers a custom validation functions for different types.
func NewValidator() *validate.Validate {
	v := validate.New()

	// Register custom validation for javascript code
	v.RegisterValidation("javascript", func(fl validate.FieldLevel) bool {
		_, err := parser.ParseFile(nil, "input", fl.Field().String(), 0)
		return err == nil
	})
	return v
}

// ConvertToStruct copies the input map to the output struct.
func ConvertToStruct(input interface{}, output interface{}) error {
	out, err := json.Marshal(input)
	if err != nil {
		return stacktrace.Propagate(err, "Error marshalling input")
	}

	err = json.Unmarshal(out, output)
	if err != nil {
		return stacktrace.Propagate(err, "Error unmarshalling input")
	}
	return nil
}
