package checker_test

import (
	"net"
	"testing"

	"github.com/matelabdev/watcher-agent/checker"
	"github.com/matelabdev/watcher-agent/config"
)

func TestTCPChecker_Online(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	host := "127.0.0.1"
	port := ln.Addr().(*net.TCPAddr).Port

	mon := config.MonitorConfig{
		Key:     "test",
		Type:    "tcp",
		Host:    host,
		Port:    port,
		Timeout: 2,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "online" {
		t.Errorf("got Status=%q, want online (error: %s)", result.Status, result.Error)
	}
	if result.ResponseMs == nil {
		t.Error("expected ResponseMs to be set")
	}
}

func TestTCPChecker_Offline(t *testing.T) {
	mon := config.MonitorConfig{
		Key:     "test",
		Type:    "tcp",
		Host:    "127.0.0.1",
		Port:    19999,
		Timeout: 1,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "offline" {
		t.Errorf("got Status=%q, want offline", result.Status)
	}
	if result.Error == "" {
		t.Error("expected Error to be set")
	}
}
