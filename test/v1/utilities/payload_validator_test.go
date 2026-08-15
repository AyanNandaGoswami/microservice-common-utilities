package utilities_test

import (
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

type TestUserPayload struct {
	Email    string `json:"email_address" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=18"`
	Username string `json:"user_name" validate:"required,min=3"`
}

func TestValidatePayload(t *testing.T) {
	t.Run("Valid Payload", func(t *testing.T) {
		validPayload := TestUserPayload{
			Email:    "test@example.com",
			Age:      25,
			Username: "john_doe",
		}

		resp := utilities.ValidatePayload(validPayload)
		if resp != nil {
			t.Errorf("expected nil error response for valid payload, got: %#v", resp)
		}
	})

	t.Run("Nil Payload", func(t *testing.T) {
		resp := utilities.ValidatePayload(nil)
		if resp == nil {
			t.Fatal("expected error response for nil payload, got nil")
		}
		if resp.Message != "Payload cannot be nil" {
			t.Errorf("expected 'Payload cannot be nil', got '%s'", resp.Message)
		}
	})

	t.Run("Non Struct Payload", func(t *testing.T) {
		resp := utilities.ValidatePayload("invalid_string_payload")
		if resp == nil {
			t.Fatal("expected error response for non-struct payload, got nil")
		}
		if resp.Message != "Payload must be a struct" {
			t.Errorf("expected 'Payload must be a struct', got '%s'", resp.Message)
		}
	})

	t.Run("Invalid Fields Payload", func(t *testing.T) {
		invalidPayload := TestUserPayload{
			Email:    "invalid-email",
			Age:      15,
			Username: "ab",
		}

		resp := utilities.ValidatePayload(invalidPayload)
		if resp == nil {
			t.Fatal("expected validation error response, got nil")
		}

		if resp.Message != "Invalid request payload" {
			t.Errorf("expected 'Invalid request payload', got '%s'", resp.Message)
		}

		fieldErrors, ok := resp.ExtraData.([]models.FieldValidationErrorResponse)
		if !ok {
			t.Fatalf("expected ExtraData to be []FieldValidationErrorResponse, got %#v", resp.ExtraData)
		}

		if len(fieldErrors) != 3 {
			t.Errorf("expected 3 field validation errors, got %d", len(fieldErrors))
		}
	})
}
