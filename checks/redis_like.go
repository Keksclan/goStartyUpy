package checks

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// RedisPingCheck performs a minimal RESP PING against a Redis-compatible
// server using raw TCP – no external dependencies required.
type RedisPingCheck struct {
	// Address is the host:port of the Redis server (e.g. "localhost:6379").
	Address string
	// Label is the human-readable name for this check.
	Label string
}

// Name returns the check label.
func (c RedisPingCheck) Name() string { return c.Label }

// Run connects via TCP, sends a RESP-encoded PING command, and expects
// a "+PONG\r\n" response.
func (c RedisPingCheck) Run(ctx context.Context) Result {
	start := time.Now()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", c.Address)
	if err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}
	defer conn.Close()

	// Respect context deadline on the connection itself.
	if dl, ok := ctx.Deadline(); ok {
		conn.SetDeadline(dl)
	}

	// Send RESP inline array: *1\r\n$4\r\nPING\r\n
	ping := "*1\r\n$4\r\nPING\r\n"
	if _, err := conn.Write([]byte(ping)); err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("write PING: %v", err),
		}
	}

	// Read response – "+PONG\r\n" is only 7 bytes but we allow a generous
	// buffer in case the server sends extra data.
	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("read PONG: %v", err),
		}
	}

	reply := strings.TrimSpace(string(buf[:n]))
	if reply != "+PONG" {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("unexpected reply: %q", reply),
		}
	}

	return Result{
		Name:     c.Label,
		OK:       true,
		Duration: time.Since(start),
	}
}
