package checker

import (
	"fmt"
	"net"
	"time"
)

type TCPChecker struct {
	Host    string
	Port    int
	Timeout int
}

func (c *TCPChecker) Check() Result {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Duration(c.Timeout)*time.Second)
	if err != nil {
		return Result{Status: "offline", Error: err.Error()}
	}
	conn.Close()
	ms := time.Since(start).Milliseconds()
	return Result{Status: "online", ResponseMs: &ms}
}
