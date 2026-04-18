package reporter_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matelabdev/watcher-agent/config"
	"github.com/matelabdev/watcher-agent/reporter"
)

func TestReporter_SendHeartbeat(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/heartbeat" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer prj_test" {
			t.Errorf("unexpected auth: %s", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		MasterURL:    srv.URL,
		ProjectToken: "prj_test",
	}
	r := reporter.New(cfg)
	err := r.SendHeartbeat("camera-giris", "online", map[string]interface{}{"host": "192.168.1.100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["monitor_key"] != "camera-giris" {
		t.Errorf("got monitor_key=%v, want camera-giris", received["monitor_key"])
	}
	if received["status"] != "online" {
		t.Errorf("got status=%v, want online", received["status"])
	}
}

func TestReporter_SyncMonitors(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/monitors/sync" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		MasterURL:    srv.URL,
		ProjectToken: "prj_test",
		Monitors: []config.MonitorConfig{
			{Key: "camera-giris", Name: "Giriş Kamerası", Type: "tcp", Host: "192.168.1.100", Port: 554, Timeout: 5, Interval: 60},
		},
	}
	r := reporter.New(cfg)
	err := r.SyncMonitors()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	monitors, ok := received["monitors"].([]interface{})
	if !ok || len(monitors) != 1 {
		t.Fatalf("expected 1 monitor in payload, got %v", received["monitors"])
	}
}