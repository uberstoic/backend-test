package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	// Register Success
	t.Run("Register Success", func(t *testing.T) {
		user := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "testuser", response["login"])
	})

	// Register Conflict (login is taken by kto-to)
	t.Run("Register Conflict", func(t *testing.T) {
		user := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Login Success
	t.Run("Login Success", func(t *testing.T) {
		user := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
	})

	// Login Wrong Password
	t.Run("Login Wrong Password", func(t *testing.T) {
		user := map[string]string{
			"login":    "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
