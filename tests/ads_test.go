package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAds(t *testing.T) {
	// Register and Login
	user := map[string]string{"login": "ad_tester", "password": "password"}
	body, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	// Create Ad Success
	t.Run("Create Ad Success", func(t *testing.T) {
		ad := map[string]interface{}{
			"title":     "Test Ad",
			"text":      "This is a test ad.",
			"image_url": "http://example.com/test.jpg",
			"price":     100,
		}
		body, _ := json.Marshal(ad)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/ads", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var adResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &adResponse)
		assert.Equal(t, "Test Ad", adResponse["title"])
	})

	// Create Ad Unauthorized
	t.Run("Create Ad Unauthorized", func(t *testing.T) {
		ad := map[string]interface{}{
			"title": "Unauthorized Ad",
			"text":  "This should fail.",
		}
		body, _ := json.Marshal(ad)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/ads", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Get Ads
	t.Run("Get Ads", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/ads", nil)
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var adsResponse []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &adsResponse)
		assert.NotEmpty(t, adsResponse)
		assert.Equal(t, "Test Ad", adsResponse[0]["title"])
	})
}
