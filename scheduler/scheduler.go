package scheduler

import (
	"sync"
	"time"

	"github.com/matelabdev/watcher-agent/checker"
	"github.com/matelabdev/watcher-agent/config"
	"github.com/matelabdev/watcher-agent/reporter"
)

type Scheduler struct {
	cfg         *config.Config
	rep         *reporter.Reporter
	defaultTick time.Duration
	stopChs     []chan struct{}
	wg          sync.WaitGroup
}

func New(cfg *config.Config, rep *reporter.Reporter, defaultTick time.Duration) *Scheduler {
	return &Scheduler{cfg: cfg, rep: rep, defaultTick: defaultTick}
}

func (s *Scheduler) Start() {
	for _, m := range s.cfg.Monitors {
		stop := make(chan struct{})
		s.stopChs = append(s.stopChs, stop)
		s.wg.Add(1)
		go s.runMonitor(m, stop)
	}
}

func (s *Scheduler) Stop() {
	for _, ch := range s.stopChs {
		close(ch)
	}
	s.wg.Wait()
}

func (s *Scheduler) runMonitor(m config.MonitorConfig, stop <-chan struct{}) {
	defer s.wg.Done()

	tick := s.defaultTick
	if m.Interval > 0 {
		tick = time.Duration(m.Interval) * time.Second
	}

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			s.check(m)
		}
	}
}

func (s *Scheduler) check(m config.MonitorConfig) {
	c := checker.New(m)
	result := c.Check()

	metadata := map[string]interface{}{
		"host": m.Host,
		"port": m.Port,
	}
	if result.ResponseMs != nil {
		metadata["response_ms"] = *result.ResponseMs
	}
	if result.Error != "" {
		metadata["error"] = result.Error
	}

	_ = s.rep.SendHeartbeat(m.Key, result.Status, metadata)
}
