package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matelabdev/watcher-agent/config"
	"github.com/matelabdev/watcher-agent/reporter"
	"github.com/matelabdev/watcher-agent/scheduler"
)

func main() {
	configPath := flag.String("config", "/etc/matelabdev-watcher/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if cfg.MasterURL == "" || cfg.ProjectToken == "" {
		fmt.Fprintln(os.Stderr, "error: master_url and project_token must be set in config")
		os.Exit(1)
	}

	rep := reporter.New(cfg)

	log.Printf("syncing %d monitor(s) with master...", len(cfg.Monitors))
	if err := rep.SyncMonitors(); err != nil {
		log.Printf("warn: monitor sync failed: %v", err)
	} else {
		log.Println("monitor sync OK")
	}

	tick := time.Duration(cfg.ReportInterval) * time.Second
	if tick == 0 {
		tick = 30 * time.Second
	}

	s := scheduler.New(cfg, rep, tick)
	s.Start()
	log.Printf("scheduler started with %d monitor(s), tick=%s", len(cfg.Monitors), tick)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	s.Stop()
	log.Println("done")
}
