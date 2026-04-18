package scheduler_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matelabdev/watcher-agent/config"
	"github.com/matelabdev/watcher-agent/reporter"
	"github.com/matelabdev/watcher-agent/scheduler"
)

func TestScheduler_CallsCheckerAndReports(t *testing.T) {
	var callCount int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/heartbeat" {
			atomic.AddInt64(&callCount, 1)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		MasterURL:    srv.URL,
		ProjectToken: "prj_test",
		Monitors: []config.MonitorConfig{
			{
				Key:      "test-monitor",
				Type:     "tcp",
				Host:     "127.0.0.1",
				Port:     19996,
				Timeout:  1,
				Interval: 0,
			},
		},
		ReportInterval: 0,
	}

	rep := reporter.New(cfg)
	s := scheduler.New(cfg, rep, 50*time.Millisecond)
	s.Start()
	time.Sleep(200 * time.Millisecond)
	s.Stop()

	count := atomic.LoadInt64(&callCount)
	if count < 2 {
		t.Errorf("expected at least 2 heartbeat calls, got %d", count)
	}
}
