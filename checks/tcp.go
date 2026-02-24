package checks

import (
	"context"
	"net"
	"time"
)

// TCPDialCheck verifies that a TCP endpoint is reachable by dialing it.
type TCPDialCheck struct {
	// Address is the host:port to dial (e.g. "localhost:5432").
	Address string
	// Label is the human-readable name for this check.
	Label string
}

// Name returns the check label.
func (c TCPDialCheck) Name() string { return c.Label }

// Run dials the TCP address using the context deadline/timeout and closes the
// connection on success.
func (c TCPDialCheck) Run(ctx context.Context) Result {
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
	conn.Close()

	return Result{
		Name:     c.Label,
		OK:       true,
		Duration: time.Since(start),
	}
}
