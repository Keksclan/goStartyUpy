package checks

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HTTPGetCheck performs an HTTP GET request and verifies the response status
// code falls within an expected range.
type HTTPGetCheck struct {
	// URL is the full URL to probe (e.g. "http://localhost:8080/healthz").
	URL string
	// Label is the human-readable name for this check.
	Label string
	// ExpectedStatusMin is the lower bound (inclusive) of acceptable status
	// codes. Defaults to 200 when zero.
	ExpectedStatusMin int
	// ExpectedStatusMax is the upper bound (inclusive) of acceptable status
	// codes. Defaults to 399 when zero.
	ExpectedStatusMax int
}

// Name returns the check label.
func (c HTTPGetCheck) Name() string { return c.Label }

// Run performs the GET request using the provided context for cancellation
// and timeout.
func (c HTTPGetCheck) Run(ctx context.Context) Result {
	start := time.Now()

	minStatus := c.ExpectedStatusMin
	if minStatus == 0 {
		minStatus = 200
	}
	maxStatus := c.ExpectedStatusMax
	if maxStatus == 0 {
		maxStatus = 399
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}

	client := &http.Client{
		// No custom timeout here – we rely on the context deadline set by the
		// Runner so behaviour is consistent across check types.
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}
	resp.Body.Close()

	if resp.StatusCode < minStatus || resp.StatusCode > maxStatus {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("unexpected status %d (expected %d–%d)", resp.StatusCode, minStatus, maxStatus),
		}
	}

	return Result{
		Name:     c.Label,
		OK:       true,
		Duration: time.Since(start),
	}
}
