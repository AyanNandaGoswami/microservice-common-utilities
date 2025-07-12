package outsource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/constants"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/models"
)

func HasPermission(requestBody *models.PermissionValidadtionRequest) error {
	// Define the URL for authorization validation
	url := constants.VALIDATE_AUTHORIZATION_ENDPOINTS

	// Marshal the request body to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create a new POST request with the marshaled JSON as the body
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set custom headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body to handle any error responses
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Define a variable to hold the parsed APIResponse struct
	var apiResponse models.APIResponse
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling response body: %w", err)
	}

	// Check the HTTP status code and process the response body
	switch resp.StatusCode {
	case http.StatusOK:
		// Successful response
		return nil
	case http.StatusForbidden:
		// Permission denied
		return fmt.Errorf("%s", apiResponse.Message)
	default:
		// Log or print the response body for error handling
		return fmt.Errorf("%s", apiResponse.Message)
	}
}
