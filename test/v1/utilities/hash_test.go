package utilities_test

import (
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

func TestPasswordHashing(t *testing.T) {
	password := "mySecretPassword123!"

	hashed, err := utilities.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hashed == "" || hashed == password {
		t.Fatalf("invalid hash returned: %s", hashed)
	}

	// Verify matching password
	if !utilities.CheckPassword(password, hashed) {
		t.Errorf("expected CheckPassword to return true for matching password")
	}

	// Verify incorrect password
	if utilities.CheckPassword("wrongPassword", hashed) {
		t.Errorf("expected CheckPassword to return false for incorrect password")
	}
}
