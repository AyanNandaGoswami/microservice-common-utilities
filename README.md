# Microservice Common Utilities (`v1`)

A standardized Go utility library providing common helpers, middleware, error handling, payload validation, JWT authentication, password hashing, NATS communication helpers, and shared data models for microservice architectures.

---

## Table of Contents
- [Installation](#installation)
- [Packages Overview](#packages-overview)
- [Utilities (`v1/utilities`)](#utilities-v1utilities)
  - [Payload Validation](#payload-validation)
  - [Error Handling](#error-handling)
  - [JWT Authentication](#jwt-authentication)
  - [Password Hashing](#password-hashing)
  - [NATS Messaging](#nats-messaging)
- [Middlewares (`v1/middlewares`)](#middlewares-v1middlewares)
  - [Auth Validation Middleware](#auth-validation-middleware)
- [Models (`v1/models`)](#models-v1models)
- [Deprecations](#deprecations)

---

## Installation

Add the package to your Go module:

```bash
go get github.com/AyanNandaGoswami/microservice-common-utilities/v1
```

---

## Packages Overview

| Package | Import Path | Description |
| :--- | :--- | :--- |
| **`utilities`** | `.../v1/utilities` | Core helper utilities for error handling, payload validation, JWT, hashing, and NATS. |
| **`middlewares`** | `.../v1/middlewares` | HTTP middleware wrappers for authentication and context population. |
| **`models`** | `.../v1/models` | Shared struct definitions for API responses, user models, and NATS messages. |

---

## Utilities (`v1/utilities`)

### Payload Validation

`ValidatePayload` validates struct instances against `validator` tags (`validate:"required,email,min=3"` etc.) and formats error responses using field JSON tags.

```go
import "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)

    if apiResp := utilities.ValidatePayload(req); apiResp != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(apiResp)
        return
    }
}
```

### Error Handling

Construct structured `*AppError` types and write clean HTTP JSON error responses:

```go
import "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"

// Helper constructors
err400 := utilities.BadRequest("Invalid request parameters")
err401 := utilities.Unauthorized("Invalid or expired session token")
err404 := utilities.NotFound("User profile not found")
err500 := utilities.InternalError("Failed to persist record", err)

// Send structured JSON error response
utilities.HandleError(w, err404)
```

`HandleError` converts the error into a standardized `models.APIResponse` with sentence-cased error messages.

### JWT Authentication

Generate and decode JWT tokens containing `UserId` and `PrimitiveUserId`:

```go
import "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"

// Generate a new token (valid for 60 minutes)
tokenString, err := utilities.GenerateNewJWToken("user-123", "prim-456")
if err != nil {
    log.Fatal(err)
}

// Retrieve claims from token string
claims, err := utilities.RetrieveDetilsFromJWT(tokenString)
if err != nil {
    log.Println("Invalid token:", err)
} else {
    fmt.Println("User ID:", claims.UserId)
}
```

### Password Hashing

Securely hash and verify passwords using bcrypt (cost 14):

```go
import "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"

// Hash password
hashed, err := utilities.HashPassword("supersecret")

// Compare password
isValid := utilities.CheckPassword("supersecret", hashed) // returns true
```

### NATS Messaging

Manage NATS connection lifecycle and handle request-reply pattern payloads:

```go
import "github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"

// Initialize connection (uses NATS_URI environment variable or default URL)
if err := utilities.InitializeNATS(); err != nil {
    log.Fatalf("NATS connection failed: %v", err)
}
defer utilities.CloseNATS()

// Request and parse structured response
var userResponse models.UserDetail
err := utilities.RequestAndParse("user.get_details", payload, &userResponse)

// Send structured reply in message handler
respPayload := utilities.PrepareNATSResponse("User found", userDetail, models.NATSuccess)
utilities.Reply(respPayload, msg)
```

---

## Middlewares (`v1/middlewares`)

### Auth Validation Middleware

`AuthValidateMiddleware` inspects incoming HTTP requests for JWT tokens in:
1. `Authorization` Header (`Bearer <token>`)
2. HTTP Cookie (`authToken`)
3. URL Query Parameter (`token`)

On success, it adds `UserIdKey`, `PrimitiveUserIdKey`, and `TokenKey` to the request context.

```go
import (
    "net/http"
    "github.com/AyanNandaGoswami/microservice-common-utilities/v1/middlewares"
)

mux := http.NewServeMux()
protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    userId := r.Context().Value(middlewares.UserIdKey).(string)
    w.Write([]byte("Hello User " + userId))
})

// Wrap handler with AuthValidateMiddleware
http.Handle("/api/profile", middlewares.AuthValidateMiddleware(protectedHandler))
```

---

## Models (`v1/models`)

- **`models.APIResponse`**: Standard envelope for API responses (`Message`, `ExtraData`).
- **`models.FieldValidationErrorResponse`**: Contains `FieldName` and validation `Message`.
- **`models.DecodedJwtClaims`**: Contains `UserId` and `PrimitiveUserId`.
- **`models.NATSResponse`**: Struct wrapper for NATS message responses (`Message`, `Data`, `Status`).
- **`models.User` / `models.UserDetail`**: Standardized user models with `json` and `bson` struct tags.

---

## Deprecations

- **`middlewares.ReturnErrorMessage(w, msg, statusCode)`**:
  - **Deprecated**: Use `utilities.HandleError` instead (e.g. `utilities.HandleError(w, utilities.BadRequest(msg))`).
