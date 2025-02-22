package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	auth "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/utilities"
)

type contextKey string

const (
	UserIdKey          contextKey = "userId"
	PrimitiveUserIdKey contextKey = "primitiveUserId"
	TokenKey           contextKey = "token"
)

func returnErrorMessage(w http.ResponseWriter, errMessage string) {
	// Convert the error message to JSON format
	errorMessage := map[string]string{"message": errMessage}
	jsonErrorMessage, _ := json.Marshal(errorMessage)

	// Return the error message as JSON with status code 401 (Unauthorized)
	http.Error(w, string(jsonErrorMessage), http.StatusUnauthorized)
}

// AuthValidateMiddleware is a middleware to validate Headers and JWT token
func AuthValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Authorization header is present
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			returnErrorMessage(w, "Authorization header is missing")
			return
		}

		// Split the Authorization header into "Bearer" and the token
		splitedInfo := strings.Split(authorization, " ")
		if len(splitedInfo) != 2 || splitedInfo[0] != "Bearer" {
			returnErrorMessage(w, "Invalid Authorization header format")
			return
		}

		// Extract the token from the Authorization header
		token := splitedInfo[1]

		// Validate token is blacklisted or not
		// if models.IsTokenBlacklisted(token) {
		// 	returnErrorMessage(w, "Token is not alive. Please login again.")
		// 	return
		// }

		// Retrieve user ID from JWT token
		info, err := auth.RetrieveDetilsFromJWT(token)
		if err != nil {
			// Split the error message by ":"
			errorMessageParts := strings.Split(err.Error(), ":")

			// Send the error message without ":"
			returnErrorMessage(w, errorMessageParts[len(errorMessageParts)-1])
			return
		}

		// Add user id to request context if needed
		// r = r.WithContext(context.WithValue(r.Context(), "userId", info.UserId))
		// r = r.WithContext(context.WithValue(r.Context(), "primitiveUserId", info.PrimitiveUserId))
		// r = r.WithContext(context.WithValue(r.Context(), "token", token))

		// Add user id to request context if needed
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIdKey, info.UserId)
		ctx = context.WithValue(ctx, PrimitiveUserIdKey, info.PrimitiveUserId)
		ctx = context.WithValue(ctx, TokenKey, token)

		// Pass the new context along with the request to the next handler
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
