package models

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
