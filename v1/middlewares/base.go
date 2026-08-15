package middlewares

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

// Deprecated: ReturnErrorMessage is deprecated. Use utilities.HandleError instead.
func ReturnErrorMessage(w http.ResponseWriter, errMessage string, statusCode int) {
	utilities.HandleError(w, &utilities.AppError{
		Message: errMessage,
		Code:    statusCode,
	})
}
