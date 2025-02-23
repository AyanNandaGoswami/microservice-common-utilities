package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
)

func ReturnErrorMessage(w http.ResponseWriter, errMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIResponse{Message: errMessage, ExtraData: nil})
}
