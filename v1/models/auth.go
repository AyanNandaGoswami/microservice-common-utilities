package models

// DecodedJwtClaims holds the decoded user identity parameters extracted from a validated JWT token.
type DecodedJwtClaims struct {
	UserId          string `json:"user_id"`
	PrimitiveUserId string `json:"primitive_user_id"`
	Email           string `json:"email,omitempty"`
	Firstname       string `json:"firstname,omitempty"`
	Lastname        string `json:"lastname,omitempty"`
	Country         string `json:"country,omitempty"`
	ContactNumber   string `json:"contact_number,omitempty"`
	IsSystemAdmin   bool   `json:"is_system_admin,omitempty"`
}
