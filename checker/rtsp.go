package checker

import (
	"net"
	"strings"
	"time"
)

type RTSPChecker struct {
	Host    string
	Timeout int
}

func (c *RTSPChecker) Check() Result {
	addr := c.Host
	if !strings.Contains(addr, ":") {
		addr = addr + ":554"
	}
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Duration(c.Timeout)*time.Second)
	if err != nil {
		return Result{Status: "offline", Error: err.Error()}
	}
	defer conn.Close()
	ms := time.Since(start).Milliseconds()
	return Result{Status: "online", ResponseMs: &ms}
}
