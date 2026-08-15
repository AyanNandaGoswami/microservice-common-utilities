package utilities_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

func TestAppError_Error(t *testing.T) {
	err := utilities.BadRequest("custom error message")
	if err.Error() != "custom error message" {
		t.Errorf("expected 'custom error message', got '%s'", err.Error())
	}
}

func TestErrorConstructors(t *testing.T) {
	t.Run("BadRequest", func(t *testing.T) {
		err := utilities.BadRequest("bad request")
		if err.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", err.Code)
		}
		if err.Message != "bad request" {
			t.Errorf("expected 'bad request', got '%s'", err.Message)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		err := utilities.Unauthorized("unauthorized access")
		if err.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", err.Code)
		}
		if err.Message != "unauthorized access" {
			t.Errorf("expected 'unauthorized access', got '%s'", err.Message)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		err := utilities.NotFound("resource not found")
		if err.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", err.Code)
		}
		if err.Message != "resource not found" {
			t.Errorf("expected 'resource not found', got '%s'", err.Message)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		underlying := errors.New("db error")
		err := utilities.InternalError("internal error", underlying)
		if err.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", err.Code)
		}
		if err.Message != "internal error" {
			t.Errorf("expected 'internal error', got '%s'", err.Message)
		}
		if err.Err != underlying {
			t.Errorf("expected underlying error %v, got %v", underlying, err.Err)
		}
	})
}

func TestHandleError(t *testing.T) {
	t.Run("AppError handling", func(t *testing.T) {
		rec := httptest.NewRecorder()
		appErr := utilities.BadRequest("invalid body")

		utilities.HandleError(rec, appErr)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var resp models.APIResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if resp.Message != "Invalid body" {
			t.Errorf("expected 'Invalid body', got '%s'", resp.Message)
		}
	})

	t.Run("Generic error handling fallback", func(t *testing.T) {
		rec := httptest.NewRecorder()
		genericErr := errors.New("unknown internal error")

		utilities.HandleError(rec, genericErr)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}

		var resp models.APIResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if resp.Message != "Internal server error" {
			t.Errorf("expected 'Internal server error', got '%s'", resp.Message)
		}
	})
}
