package models

// FieldValidationErrorResponse represents a validation error message associated with a specific JSON payload field name.
type FieldValidationErrorResponse struct {
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

// APIResponse represents the standard JSON response format used across HTTP endpoints.
type APIResponse struct {
	Message   string      `json:"message"`
	ExtraData interface{} `json:"extra_data"`
}
