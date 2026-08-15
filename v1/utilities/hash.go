package utilities

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password using bcrypt with a cost factor of 14.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compares a plain-text password with a bcrypt hashed password, returning true if they match.
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
