package models

type PermissionValidadtionRequest struct {
	ValidateBy      string `json:"validate_by"`
	PrimitiveUserId string `json:"primitiveUserId"`
	RequestedUrl    string `json:"requestedUrl"`
	RequestedMethod string `json:"requestedMethod"`
}
