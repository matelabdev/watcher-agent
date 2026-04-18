package checker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matelabdev/watcher-agent/checker"
	"github.com/matelabdev/watcher-agent/config"
)

func TestHTTPChecker_Online(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	mon := config.MonitorConfig{
		Key:     "test-http",
		Type:    "http",
		URL:     srv.URL,
		Timeout: 5,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "online" {
		t.Errorf("got Status=%q, want online", result.Status)
	}
}

func TestHTTPChecker_Degraded_On5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	mon := config.MonitorConfig{
		Key:     "test-http",
		Type:    "http",
		URL:     srv.URL,
		Timeout: 5,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "degraded" {
		t.Errorf("got Status=%q, want degraded", result.Status)
	}
}

func TestHTTPChecker_Offline_OnConnectionError(t *testing.T) {
	mon := config.MonitorConfig{
		Key:     "test-http",
		Type:    "http",
		URL:     "http://127.0.0.1:19998",
		Timeout: 1,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "offline" {
		t.Errorf("got Status=%q, want offline", result.Status)
	}
}
