package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservice-common-utilities/v1/utilities"
)

func TestAuthValidateMiddleware(t *testing.T) {
	userId := "user-123"
	primitiveId := "prim-456"
	validToken, err := utilities.GenerateNewJWToken(userId, primitiveId)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserId, _ := r.Context().Value(middlewares.UserIdKey).(string)
		ctxPrimId, _ := r.Context().Value(middlewares.PrimitiveUserIdKey).(string)
		ctxToken, _ := r.Context().Value(middlewares.TokenKey).(string)

		if ctxUserId != userId {
			t.Errorf("expected context UserId '%s', got '%s'", userId, ctxUserId)
		}

		if ctxPrimId != primitiveId {
			t.Errorf("expected context PrimitiveUserId '%s', got '%s'", primitiveId, ctxPrimId)
		}

		if ctxToken != validToken {
			t.Errorf("expected context Token '%s', got '%s'", validToken, ctxToken)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := middlewares.AuthValidateMiddleware(nextHandler)

	t.Run("Valid Bearer Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}
	})

	t.Run("Valid Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "authToken",
			Value: validToken,
		})
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}
	})

	t.Run("Valid Query Param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected?token="+validToken, nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("Invalid Header Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat "+validToken)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("Invalid Token String", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.jwt.token")
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})
}
