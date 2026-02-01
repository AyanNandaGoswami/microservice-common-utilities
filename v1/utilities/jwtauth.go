package utilities

import (
	"errors"
	"log"
	"time"

	auth "github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId          string `json:"user_id"`
	PrimitiveUserId string `json:"primitive_user_id"`
	jwt.RegisteredClaims
}

const (
	expireAfter = 60 // in minutes
	jwtKey      = "my_secret_key"
)

// Generate a new JWT token for given user ID
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
