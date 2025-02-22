package main

import (
	"fmt"
	"log"

	auth "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/utilities"
)

func main() {
	// Generate a new token for a user
	userId := "123456"
	primitiveId := "67b0e9c917ddc6790e248882"
	tokenString, err := auth.GenerateNewJWToken(userId, primitiveId)
	if err != nil {
		log.Fatalf("Error generating JWT token: %v", err)
	}
	fmt.Printf("Generated JWT Token: %s\n", tokenString)

	// Retrieve the user ID from the JWT token
	details, err := auth.RetrieveDetilsFromJWT(tokenString)
	if err != nil {
		log.Fatalf("Error retrieving user ID from JWT token: %v", err)
	}

	fmt.Println(details.UserId, details.PrimitiveUserId)

	// Display the retrieved user ID
	// fmt.Printf("User ID extracted from JWT: %s\n %s", userId, primitiveUserId)
}
