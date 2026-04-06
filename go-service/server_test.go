package main

import (
	"bytes"
	"campus-connect/go-service/handlers"
	"campus-connect/go-service/middleware"
	"campus-connect/go-service/webhook"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Helper to setup mock server router
func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/analytics", handlers.AnalyticsHandler)
	protectedMux.HandleFunc("/notifications", handlers.NotificationsHandler)

	protectedHandler := middleware.RateLimitMiddleware(middleware.APIKeyAuthMiddleware(protectedMux))

	mux.Handle("/analytics", protectedHandler)
	mux.Handle("/notifications", protectedHandler)
	mux.HandleFunc("/webhook", webhook.ProcessWebhook)
	return mux
}

// 1. Test Webhook POST processes successfully
func TestProcessWebhookSuccess(t *testing.T) {
	router := setupRouter()
	body := []byte(`{"eventId": 1, "title": "Test", "action": "CREATE", "timestamp": "2026-04-06"}`)
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("expected 202, got %v", status)
	}
}

// 2. Test Webhook Invalid Method
func TestProcessWebhookInvalidMethod(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %v", rr.Code)
	}
}

// 3. Test Webhook Bad JSON
func TestProcessWebhookBadJSON(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte("{bad json")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %v", rr.Code)
	}
}

// 4. Test API Key Missing
func TestAPIKeyMissing(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/analytics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %v", rr.Code)
	}
}

// 5. Test API Key Valid
func TestAPIKeyValid(t *testing.T) {
	os.Setenv("API_KEY", "testing123")
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/analytics", nil)
	req.Header.Set("X-API-Key", "testing123")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", rr.Code)
	}
}

// 6. Test Rate Limiter Exceeded Simulation
func TestRateLimiterExceeded(t *testing.T) {
	os.Setenv("API_KEY", "testkey")
	router := setupRouter()

	for i := 0; i < 6; i++ {
		req, _ := http.NewRequest("GET", "/analytics", nil)
		req.Header.Set("X-API-Key", "testkey")
		req.RemoteAddr = "127.0.0.1:12345"
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		
		if i == 5 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected 429 on 6th request, got %v", rr.Code)
		}
	}
}

// 7. Test Analytics Data Fetch
func TestAnalyticsData(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/analytics", nil)
	req.Header.Set("X-API-Key", os.Getenv("API_KEY"))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var stats map[string]int
	json.NewDecoder(rr.Body).Decode(&stats)
	if _, ok := stats["total_events"]; !ok {
		t.Errorf("expected total_events key in response")
	}
}

// 8. Test Notifications Data Fetch
func TestNotificationsData(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/notifications", nil)
	req.Header.Set("X-API-Key", os.Getenv("API_KEY"))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var notes []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&notes)
	if err != nil {
		t.Errorf("failed to decode notifications: %v", err)
	}
}
