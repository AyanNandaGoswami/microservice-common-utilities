package utilities_test

import (
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

func TestPrepareNATSResponse(t *testing.T) {
	msg := "Operation completed successfully"
	data := map[string]string{"key": "value"}

	resp := utilities.PrepareNATSResponse(msg, data, models.NATSuccess)

	if resp.Message != msg {
		t.Errorf("expected message '%s', got '%s'", msg, resp.Message)
	}

	if resp.Status != models.NATSuccess {
		t.Errorf("expected status '%s', got '%s'", models.NATSuccess, resp.Status)
	}

	dataMap, ok := resp.Data.(map[string]string)
	if !ok || dataMap["key"] != "value" {
		t.Errorf("unexpected data payload: %#v", resp.Data)
	}
}
