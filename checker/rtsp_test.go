package checker_test

import (
	"net"
	"testing"

	"github.com/matelabdev/watcher-agent/checker"
	"github.com/matelabdev/watcher-agent/config"
)

func TestRTSPChecker_Online(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	host := "127.0.0.1"
	addr := ln.Addr().(*net.TCPAddr)

	mon := config.MonitorConfig{
		Key:     "test-rtsp",
		Type:    "rtsp",
		Host:    host,
		Port:    addr.Port,
		Timeout: 2,
	}
	mon.Host = addr.String()

	c := checker.New(mon)
	result := c.Check()

	if result.Status != "online" {
		t.Errorf("got Status=%q, want online (error: %s)", result.Status, result.Error)
	}
}

func TestRTSPChecker_Offline(t *testing.T) {
	mon := config.MonitorConfig{
		Key:     "test-rtsp",
		Type:    "rtsp",
		Host:    "127.0.0.1:19997",
		Timeout: 1,
	}
	c := checker.New(mon)
	result := c.Check()

	if result.Status != "offline" {
		t.Errorf("got Status=%q, want offline", result.Status)
	}
}
