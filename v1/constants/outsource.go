package constants

import "fmt"

const AUTHORIZATION_SERVICE_BASE_URL = "http://127.0.0.1:4002"

var VALIDATE_AUTHORIZATION_ENDPOINTS = fmt.Sprintf("%s/authorization/v1/validate-permisison/", AUTHORIZATION_SERVICE_BASE_URL)
