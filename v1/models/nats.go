package models

// NatStatusType defines the status string type for NATS messaging responses.
type NatStatusType string

const (
	// NATSuccess indicates a successful NATS response operation status.
	NATSuccess NatStatusType = "success"

	// NATFailed indicates a failed NATS response operation status.
	NATFailed NatStatusType = "failed"
)

// NATSResponse represents the standard response structure exchanged over NATS request-reply topics.
type NATSResponse struct {
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"`
	Status  NatStatusType `json:"status"`
}
