package middlewares

import (
	"context"
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

// Define the interface for getting user permissions.
type PermissionGetter interface {
	GetUserPermissionEndpoints(primitiveUserId string) (map[string]string, error)
}

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

// PermissionValidationMiddleware expects a concrete implementation of PermissionGetter
func PermissionValidationMiddleware(permissionGetter PermissionGetter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the primitiveUserId from the request context
			primitiveUserId := r.Context().Value(PrimitiveUserIdKey).(string)

			// Call the GetUserPermissionEndpoints method on the concrete PermissionGetter instance
			permissionEndpoints, err := permissionGetter.GetUserPermissionEndpoints(primitiveUserId)
			if err != nil {
				ReturnErrorMessage(w, err.Error(), 400)
				return
			}

			// Check if the user has permission to access the requested URL and method
			if !hasPermission(permissionEndpoints, r.URL.Path, r.Method) {
				ReturnErrorMessage(w, "You do not have permission to perform this action", 403)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasPermission(userPermissionEndpoints map[string]string, requestedUrl string, requesedMethod string) bool {
	method, exists := userPermissionEndpoints[requestedUrl]
	if exists {
		return method == requesedMethod
	}
	return false
}
