package models

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nats-io/nats.go"
)

type NatStatusType string

const (
	NATSuccess NatStatusType = "success"
	NATFailed  NatStatusType = "failed"
)

type NATSResponse struct {
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"`
	Status  NatStatusType `json:"status"`
}

// ParseNatMsgToStruct parses the NATS message into the targeted struct.
// targetedStruct must be a pointer to a struct.
func ParseNatMsgToStruct(natMsg *nats.Msg, targetedStruct interface{}) error {
	// targetedStruct must be a pointer to a struct
	t := reflect.TypeOf(targetedStruct)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("targetedStruct must be a pointer to a struct")
	}

	// Unmarshal the raw NATS message into NATSResponse
	var natResp NATSResponse
	if err := json.Unmarshal(natMsg.Data, &natResp); err != nil {
		return fmt.Errorf("failed to unmarshal NATS message: %w", err)
	}

	fmt.Printf("Received from NATS: %v", natResp)

	// check if status is failed
	if natResp.Status == NATFailed {
		return fmt.Errorf("NATS response indicates failure: %s", natResp.Message)
	}

	// Marshal the Data field to JSON
	dataBytes, err := json.Marshal(natResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal Data field: %w", err)
	}

	// Unmarshal JSON into the targeted struct
	if err := json.Unmarshal(dataBytes, targetedStruct); err != nil {
		return fmt.Errorf("failed to unmarshal into targeted struct: %w", err)
	}

	return nil
}

// PrepareNATSResponse prepares a NATSResponse struct with given parameters.
func PrepareNATSResponse(message string, data interface{}, status NatStatusType) NATSResponse {
	return NATSResponse{
		Message: message,
		Data:    data,
		Status:  status,
	}
}
