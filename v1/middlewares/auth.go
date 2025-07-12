package middlewares

import (
	"context"
	"net/http"
	"strings"

	auth "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

type contextKey string

const (
	UserIdKey          contextKey = "userId"
	PrimitiveUserIdKey contextKey = "primitiveUserId"
	TokenKey           contextKey = "token"
)

// AuthValidateMiddleware is a middleware to validate Headers and JWT token
func AuthValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Authorization header is present
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			ReturnErrorMessage(w, "Authorization header is missing", 401)
			return
		}

		// Split the Authorization header into "Bearer" and the token
		splitedInfo := strings.Split(authorization, " ")
		if len(splitedInfo) != 2 || splitedInfo[0] != "Bearer" {
			ReturnErrorMessage(w, "Invalid Authorization header format", 401)
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
			ReturnErrorMessage(w, errorMessageParts[len(errorMessageParts)-1], 401)
			return
		}

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
