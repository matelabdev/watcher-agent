package checker

import (
	"net/http"
	"time"
)

type HTTPChecker struct {
	URL     string
	Timeout int
}

func (c *HTTPChecker) Check() Result {
	client := &http.Client{Timeout: time.Duration(c.Timeout) * time.Second}
	start := time.Now()
	resp, err := client.Get(c.URL)
	if err != nil {
		return Result{Status: "offline", Error: err.Error()}
	}
	defer resp.Body.Close()
	ms := time.Since(start).Milliseconds()
	if resp.StatusCode >= 500 {
		return Result{Status: "degraded", ResponseMs: &ms}
	}
	return Result{Status: "online", ResponseMs: &ms}
}
