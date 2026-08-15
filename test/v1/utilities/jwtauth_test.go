package utilities_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuth(t *testing.T) {
	userId := "user-uuid-1234"
	primitiveId := "prim-5678"

	tokenString, err := utilities.GenerateNewJWToken(userId, primitiveId)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("generated JWT token string is empty")
	}

	// Retrieve details from generated JWT
	claims, err := utilities.RetrieveDetilsFromJWT(tokenString)
	if err != nil {
		t.Fatalf("failed to retrieve details from valid JWT: %v", err)
	}

	if claims.UserId != userId {
		t.Errorf("expected UserId '%s', got '%s'", userId, claims.UserId)
	}

	if claims.PrimitiveUserId != primitiveId {
		t.Errorf("expected PrimitiveUserId '%s', got '%s'", primitiveId, claims.PrimitiveUserId)
	}

	// Test invalid token string
	_, err = utilities.RetrieveDetilsFromJWT("invalid.jwt.token")
	if err == nil {
		t.Errorf("expected error for invalid token string, got nil")
	}
}

func TestVerifyRS256Token(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key pair: %v", err)
	}

	now := time.Now()
	claims := struct {
		UserId          string `json:"user_id"`
		PrimitiveUserId string `json:"primitive_user_id"`
		Email           string `json:"email"`
		Firstname       string `json:"firstname"`
		Lastname        string `json:"lastname"`
		Country         string `json:"country"`
		ContactNumber   string `json:"contact_number"`
		IsSystemAdmin   bool   `json:"is_system_admin"`
		jwt.RegisteredClaims
	}{
		UserId:          "user-123",
		PrimitiveUserId: "prim-456",
		Email:           "test@example.com",
		Firstname:       "John",
		Lastname:        "Doe",
		Country:         "US",
		ContactNumber:   "+1234567890",
		IsSystemAdmin:   true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign RS256 token: %v", err)
	}

	decoded, err := utilities.VerifyRS256Token(tokenString, &privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to verify RS256 token: %v", err)
	}

	if decoded.UserId != "user-123" || decoded.Email != "test@example.com" || !decoded.IsSystemAdmin {
		t.Errorf("decoded claims mismatch: %+v", decoded)
	}
}
