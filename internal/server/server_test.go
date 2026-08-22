package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWaveHandler(t *testing.T) {
	s := NewServer("example")
	body := `{"amplitude":1,"period":8,"depth":50,"layers":12,"rho":1025}`
	req := httptest.NewRequest(http.MethodPost, "/api/wave", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/wave status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("content-type = %q, want application/json", ct)
	}
	if !strings.Contains(rec.Body.String(), "wavelength") {
		t.Errorf("response missing wavelength field: %s", rec.Body.String())
	}
}

func TestWaveRejectsBadInput(t *testing.T) {
	s := NewServer("example")
	req := httptest.NewRequest(http.MethodPost, "/api/wave", strings.NewReader(`{"period":-1}`))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := NewServer("example")
	req := httptest.NewRequest(http.MethodGet, "/api/wave", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET should be 405, got %d", rec.Code)
	}
}
