package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

type contextKey string

const (
	// UserIdKey is the request context key used to store the user ID extracted from the JWT.
	UserIdKey contextKey = "userId"

	// PrimitiveUserIdKey is the request context key used to store the primitive user ID extracted from the JWT.
	PrimitiveUserIdKey contextKey = "primitiveUserId"

	// TokenKey is the request context key used to store the raw JWT authorization token.
	TokenKey contextKey = "token"
)

// AuthValidateMiddleware is an HTTP middleware that validates the incoming authorization token
// from Authorization headers (Bearer token), cookies (authToken), or query parameters (token).
// If valid, it attaches UserIdKey, PrimitiveUserIdKey, and TokenKey to the request context.
func AuthValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Check if the Authorization header is present
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			// Split the Authorization header into "Bearer" and the token
			splitedInfo := strings.Split(authorization, " ")
			if len(splitedInfo) != 2 || splitedInfo[0] != "Bearer" {
				utilities.HandleError(w, utilities.Unauthorized("invalid Authorization header format"))
				return
			}

			// Extract the token from the Authorization header
			token = splitedInfo[1]
		}

		// If token is empty, try retrieving token from cookies
		if token == "" {
			cookie, _ := r.Cookie("authToken")
			if cookie != nil {
				token = cookie.Value
			}
		}

		// If token is empty, try retrieving token from query parameters
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			utilities.HandleError(w, utilities.Unauthorized("Authorization token is missing"))
			return
		}

		// Validate token is blacklisted or not
		// if models.IsTokenBlacklisted(token) {
		// 	returnErrorMessage(w, "Token is not alive. Please login again.")
		// 	return
		// }

		// Retrieve user ID from JWT token
		info, err := utilities.RetrieveDetilsFromJWT(token)
		if err != nil {
			// Split the error message by ":"
			errorMessageParts := strings.Split(err.Error(), ":")

			// Send the error message without ":"
			utilities.HandleError(w, utilities.Unauthorized(errorMessageParts[len(errorMessageParts)-1]))
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
