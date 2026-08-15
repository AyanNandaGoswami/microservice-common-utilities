package utilities

import (
	"encoding/json"
	"net/http"
	"strings"

	common_models "github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
)

// AppError represents a custom application error containing an HTTP status code, a human-readable message, and an optional underlying error.
type AppError struct {
	Message string
	Code    int
	Err     error
}

// Error implements the standard Go error interface for AppError.
func (e *AppError) Error() string {
	return e.Message
}

// BadRequest constructs a new AppError representing an HTTP 400 Bad Request error.
func BadRequest(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

// Unauthorized constructs a new AppError representing an HTTP 401 Unauthorized error.
func Unauthorized(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

// NotFound constructs a new AppError representing an HTTP 404 Not Found error.
func NotFound(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

// InternalError constructs a new AppError representing an HTTP 500 Internal Server Error, wrapping the underlying error.
func InternalError(msg string, err error) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func toSentenceCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// HandleError formats and writes a structured APIResponse JSON error to the provided http.ResponseWriter.
// If err is a *AppError, its embedded status code and message are used. Otherwise, it defaults to a 500 Internal Server Error.
func HandleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(common_models.APIResponse{
			Message: toSentenceCase(appErr.Message),
		})
		return
	}

	// fallback
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(common_models.APIResponse{
		Message: toSentenceCase("internal server error"),
	})
}
