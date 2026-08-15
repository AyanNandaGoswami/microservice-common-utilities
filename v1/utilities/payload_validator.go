package utilities

import (
	"reflect"
	"strings"

	common_models "github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {

	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {

		name := fld.Tag.Get("json")

		if idx := strings.Index(name, ","); idx >= 0 {
			name = name[:idx]
		}

		if name == "-" {
			return ""
		}

		return name
	})
}

// ValidatePayload validates a struct or pointer-to-struct payload against validator struct tags (mapping field names using json tags).
// It returns nil if the payload is valid, or an *APIResponse containing detailed validation errors if validation fails.
func ValidatePayload(payload interface{}) *common_models.APIResponse {

	if payload == nil {
		return &common_models.APIResponse{
			Message: "Payload cannot be nil",
		}
	}

	rootType := reflect.TypeOf(payload)

	if rootType.Kind() == reflect.Ptr {
		rootType = rootType.Elem()
	}

	if rootType.Kind() != reflect.Struct {
		return &common_models.APIResponse{
			Message: "Payload must be a struct",
		}
	}

	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &common_models.APIResponse{
			Message: err.Error(),
		}
	}

	fieldErrors := make(
		[]common_models.FieldValidationErrorResponse,
		0,
		len(validationErrors),
	)

	for _, validationErr := range validationErrors {

		fieldErrors = append(
			fieldErrors,
			common_models.FieldValidationErrorResponse{
				FieldName: getJSONFieldPath(
					rootType,
					validationErr,
				),
				Message: getValidationMessage(validationErr),
			},
		)
	}

	return &common_models.APIResponse{
		Message:   "Invalid request payload",
		ExtraData: fieldErrors,
	}
}

func getJSONFieldPath(
	rootType reflect.Type,
	validationErr validator.FieldError,
) string {

	namespace := validationErr.StructNamespace()

	parts := strings.Split(namespace, ".")

	if len(parts) > 0 {
		parts = parts[1:]
	}

	currentType := rootType
	jsonParts := make([]string, 0, len(parts))

	for _, part := range parts {

		fieldName := part

		if idx := strings.Index(fieldName, "["); idx >= 0 {
			fieldName = fieldName[:idx]
		}

		if currentType.Kind() == reflect.Ptr {
			currentType = currentType.Elem()
		}

		if currentType.Kind() != reflect.Struct {
			jsonParts = append(jsonParts, part)
			continue
		}

		field, ok := currentType.FieldByName(fieldName)
		if !ok {
			jsonParts = append(jsonParts, part)
			continue
		}

		jsonTag := field.Tag.Get("json")

		jsonName := jsonTag

		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			jsonName = jsonTag[:idx]
		}

		if jsonName == "" || jsonName == "-" {
			jsonName = fieldName
		}

		if idx := strings.Index(part, "["); idx >= 0 {
			jsonName += part[idx:]
		}

		jsonParts = append(jsonParts, jsonName)

		currentType = field.Type

		if currentType.Kind() == reflect.Ptr {
			currentType = currentType.Elem()
		}

		if currentType.Kind() == reflect.Slice {
			currentType = currentType.Elem()
		}
	}

	return strings.Join(jsonParts, ".")
}

func getValidationMessage(err validator.FieldError) string {

	if err.Param() != "" &&
		(err.Tag() == "required" || err.Tag() == "email") {
		return err.Param()
	}

	switch err.Tag() {

	case "required":
		return "This field is required."

	case "email":
		return "Invalid email address."

	case "min":
		return buildMinMessage(err)

	case "max":
		return buildMaxMessage(err)

	case "len":
		return "must be exactly " + err.Param() + " characters"

	case "gt":
		return "must be greater than " + err.Param()

	case "gte":
		return "must be greater than or equal to " + err.Param()

	case "lt":
		return "must be less than " + err.Param()

	case "lte":
		return "must be less than or equal to " + err.Param()

	case "oneof":
		return "invalid choice"

	default:
		return "invalid value"
	}
}

func buildMinMessage(err validator.FieldError) string {

	switch err.Kind() {

	case reflect.Slice,
		reflect.Array:

		return "minimum " + err.Param() + " item(s) is required"

	case reflect.String:

		return "minimum " + err.Param() + " characters required"

	default:

		return "minimum value allowed is " + err.Param()
	}
}

func buildMaxMessage(err validator.FieldError) string {

	switch err.Kind() {

	case reflect.String:

		return "maximum " + err.Param() + " characters allowed"

	case reflect.Slice,
		reflect.Array:

		return "maximum " + err.Param() + " item(s) allowed"

	default:

		return "maximum value allowed is " + err.Param()
	}
}
