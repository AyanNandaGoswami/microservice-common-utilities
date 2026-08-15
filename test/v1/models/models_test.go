package models_test

import (
	"encoding/json"
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
)

func TestModelsJSONSerialization(t *testing.T) {
	t.Run("APIResponse", func(t *testing.T) {
		apiResp := models.APIResponse{
			Message:   "Success",
			ExtraData: map[string]string{"foo": "bar"},
		}

		data, err := json.Marshal(apiResp)
		if err != nil {
			t.Fatalf("failed to marshal APIResponse: %v", err)
		}

		var decoded models.APIResponse
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal APIResponse: %v", err)
		}

		if decoded.Message != apiResp.Message {
			t.Errorf("expected Message '%s', got '%s'", apiResp.Message, decoded.Message)
		}
	})

	t.Run("FieldValidationErrorResponse", func(t *testing.T) {
		fieldErr := models.FieldValidationErrorResponse{
			FieldName: "email",
			Message:   "Invalid email address",
		}

		data, err := json.Marshal(fieldErr)
		if err != nil {
			t.Fatalf("failed to marshal FieldValidationErrorResponse: %v", err)
		}

		var decoded models.FieldValidationErrorResponse
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal FieldValidationErrorResponse: %v", err)
		}

		if decoded.FieldName != "email" || decoded.Message != "Invalid email address" {
			t.Errorf("unexpected unmarshaled struct: %#v", decoded)
		}
	})

	t.Run("DecodedJwtClaims", func(t *testing.T) {
		claims := models.DecodedJwtClaims{
			UserId:          "usr-1",
			PrimitiveUserId: "prim-1",
		}

		data, err := json.Marshal(claims)
		if err != nil {
			t.Fatalf("failed to marshal DecodedJwtClaims: %v", err)
		}

		var decoded models.DecodedJwtClaims
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal DecodedJwtClaims: %v", err)
		}

		if decoded.UserId != claims.UserId || decoded.PrimitiveUserId != claims.PrimitiveUserId {
			t.Errorf("unexpected unmarshaled claims: %#v", decoded)
		}
	})

	t.Run("NATSResponse", func(t *testing.T) {
		natsResp := models.NATSResponse{
			Message: "OK",
			Status:  models.NATSuccess,
			Data:    "payload",
		}

		data, err := json.Marshal(natsResp)
		if err != nil {
			t.Fatalf("failed to marshal NATSResponse: %v", err)
		}

		var decoded models.NATSResponse
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal NATSResponse: %v", err)
		}

		if decoded.Status != models.NATSuccess || decoded.Message != "OK" {
			t.Errorf("unexpected unmarshaled NATSResponse: %#v", decoded)
		}
	})

	t.Run("UserDetail", func(t *testing.T) {
		user := models.UserDetail{
			User: models.User{
				UUID:       "uuid-123",
				Firstname:  "John",
				Fullname:   "John Doe",
				ProfileImg: "http://example.com/img.png",
			},
			Lastname:   "Doe",
			Middlename: "Edward",
			Email:      "john@example.com",
		}

		data, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("failed to marshal UserDetail: %v", err)
		}

		var decoded models.UserDetail
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal UserDetail: %v", err)
		}

		if decoded.UUID != user.UUID || decoded.Email != user.Email {
			t.Errorf("unexpected unmarshaled UserDetail: %#v", decoded)
		}
	})
}
