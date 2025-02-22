package models

type DecodedJwtClaims struct {
	UserId          string `json:"user_id"`
	PrimitiveUserId string `json:"primitive_user_id"`
}
