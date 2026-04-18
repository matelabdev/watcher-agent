package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/matelabdev/watcher-agent/config"
)

type Reporter struct {
	cfg    *config.Config
	client *http.Client
}

func New(cfg *config.Config) *Reporter {
	return &Reporter{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *Reporter) SendHeartbeat(monitorKey, status string, metadata map[string]interface{}) error {
	payload := map[string]interface{}{
		"monitor_key": monitorKey,
		"status":      status,
		"metadata":    metadata,
	}
	return r.post("/api/heartbeat", payload)
}

func (r *Reporter) SyncMonitors() error {
	type monitorItem struct {
		Key    string      `json:"key"`
		Name   string      `json:"name"`
		Type   string      `json:"type"`
		Config interface{} `json:"config"`
	}

	items := make([]monitorItem, 0, len(r.cfg.Monitors))
	for _, m := range r.cfg.Monitors {
		cfg := map[string]interface{}{
			"host":     m.Host,
			"port":     m.Port,
			"url":      m.URL,
			"timeout":  m.Timeout,
			"interval": m.Interval,
		}
		items = append(items, monitorItem{Key: m.Key, Name: m.Name, Type: m.Type, Config: cfg})
	}
	return r.post("/api/monitors/sync", map[string]interface{}{"monitors": items})
}

func (r *Reporter) post(path string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := strings.TrimRight(r.cfg.MasterURL, "/") + path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.cfg.ProjectToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}