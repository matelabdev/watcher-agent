package config_test

import (
	"testing"

	"github.com/matelabdev/watcher-agent/config"
)

func TestLoad_ValidFile(t *testing.T) {
	cfg, err := config.Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.MasterURL != "http://watcher.test" {
		t.Errorf("got MasterURL=%q, want %q", cfg.MasterURL, "http://watcher.test")
	}
	if cfg.ProjectToken != "prj_testtoken" {
		t.Errorf("got ProjectToken=%q, want %q", cfg.ProjectToken, "prj_testtoken")
	}
	if cfg.ReportInterval != 30 {
		t.Errorf("got ReportInterval=%d, want 30", cfg.ReportInterval)
	}
	if len(cfg.Monitors) != 2 {
		t.Fatalf("got %d monitors, want 2", len(cfg.Monitors))
	}

	tcp := cfg.Monitors[0]
	if tcp.Key != "camera-giris" {
		t.Errorf("got Key=%q, want camera-giris", tcp.Key)
	}
	if tcp.Type != "tcp" {
		t.Errorf("got Type=%q, want tcp", tcp.Type)
	}
	if tcp.Host != "192.168.1.100" {
		t.Errorf("got Host=%q, want 192.168.1.100", tcp.Host)
	}
	if tcp.Port != 554 {
		t.Errorf("got Port=%d, want 554", tcp.Port)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("testdata/nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}