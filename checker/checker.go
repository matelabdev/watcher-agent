package checker

import (
	"github.com/matelabdev/watcher-agent/config"
)

type Result struct {
	Status     string // "online" | "offline" | "degraded"
	ResponseMs *int64
	Error      string
}

type Checker interface {
	Check() Result
}

func New(m config.MonitorConfig) Checker {
	switch m.Type {
	case "http":
		return &HTTPChecker{URL: m.URL, Timeout: m.Timeout}
	case "rtsp":
		return &RTSPChecker{Host: m.Host, Timeout: m.Timeout}
	default: // "tcp"
		return &TCPChecker{Host: m.Host, Port: m.Port, Timeout: m.Timeout}
	}
}
