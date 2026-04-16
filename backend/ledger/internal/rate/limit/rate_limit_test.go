package limit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rafaeldepontes/ledger/internal/rate/limit"
)

// TestRateLimit checks that the rate limiter enforces limits correctly
func TestRateLimit(t *testing.T) {
	mw := limit.NewMiddleware()

	// Simple handler that returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	limitedHandler := mw.RateLimit(handler)

	// First two requests should succeed (burst of 5 allowed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		limitedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d on request %d", rec.Code, i+1)
		}
	}

	// Sixth request should be rejected
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	limitedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d on request 6", rec.Code)
	}
}
