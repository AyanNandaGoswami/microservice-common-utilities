package utilities

import (
	"crypto/rsa"
	"errors"
	"log"
	"time"

	auth "github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT custom claims structure containing user ID and primitive user ID alongside standard registered claims.
type Claims struct {
	UserId          string `json:"user_id"`
	PrimitiveUserId string `json:"primitive_user_id"`
	Email           string `json:"email,omitempty"`
	Firstname       string `json:"firstname,omitempty"`
	Lastname        string `json:"lastname,omitempty"`
	Country         string `json:"country,omitempty"`
	ContactNumber   string `json:"contact_number,omitempty"`
	IsSystemAdmin   bool   `json:"is_system_admin,omitempty"`
	jwt.RegisteredClaims
}

const (
	expireAfter = 60 // in minutes
	jwtKey      = "my_secret_key"
)

// GenerateNewJWToken generates a signed JWT token string using legacy HS256 HMAC.
//
// Deprecated: HS256 is deprecated in favor of central OAuth 2.0 / OIDC RS256 asymmetric token issuance.
func GenerateNewJWToken(userId string, primitiveId string) (string, error) {
	// Set the JWT secret key
	jwtkeyBytes := []byte(jwtKey)

	// Set the token expiration duration
	expirationTime := time.Now().Add(time.Duration(expireAfter) * time.Minute)

	// Define the token claims
	claims := &Claims{
		UserId:          userId,
		PrimitiveUserId: primitiveId,
		// Registered token claims
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create a new JWT token with the specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the signed token string
	tokenString, err := token.SignedString(jwtkeyBytes)

	// Return the token string and error if any
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func checkJWTKey(token *jwt.Token) (interface{}, error) {
	// Check the signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	// Return the key
	return []byte(jwtKey), nil
}

// RetrieveDetilsFromJWT parses and validates a legacy HS256 JWT token string.
//
// Deprecated: Use VerifyRS256Token with the Central Auth Server's public RSA key instead.
func RetrieveDetilsFromJWT(tokenString string) (*auth.DecodedJwtClaims, error) {
	claims := &Claims{}

	// Parse the token string into the claim
	token, err := jwt.ParseWithClaims(tokenString, claims, checkJWTKey)
	var decodedClaim auth.DecodedJwtClaims

	if err != nil {
		return nil, err
	}

	// Verify token validity
	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	decodedClaim.UserId = claims.UserId
	decodedClaim.PrimitiveUserId = claims.PrimitiveUserId

	// Return user ID
	return &decodedClaim, nil
}

// VerifyRS256Token validates an RS256 signed Central SSO access token using an RSA Public Key.
func VerifyRS256Token(tokenString string, pubKey *rsa.PublicKey) (*auth.DecodedJwtClaims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method, expected RS256")
		}
		return pubKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired RS256 token")
	}

	userId := claims.UserId
	if userId == "" {
		userId = claims.Subject
	}

	return &auth.DecodedJwtClaims{
		UserId:          userId,
		PrimitiveUserId: claims.PrimitiveUserId,
		Email:           claims.Email,
		Firstname:       claims.Firstname,
		Lastname:        claims.Lastname,
		Country:         claims.Country,
		ContactNumber:   claims.ContactNumber,
		IsSystemAdmin:   claims.IsSystemAdmin,
	}, nil
}
