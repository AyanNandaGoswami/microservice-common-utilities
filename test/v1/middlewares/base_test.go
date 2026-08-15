package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

func TestHandleErrorInMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	errMsg := "forbidden action"
	statusCode := http.StatusForbidden

	// Use utilities.HandleError instead of ReturnErrorMessage
	utilities.HandleError(rec, &utilities.AppError{
		Message: errMsg,
		Code:    statusCode,
	})

	if rec.Code != statusCode {
		t.Errorf("expected status code %d, got %d", statusCode, rec.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	if resp.Message != "Forbidden action" {
		t.Errorf("expected 'Forbidden action', got '%s'", resp.Message)
	}
}
