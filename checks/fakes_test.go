package checks

import (
	"context"
	"time"
)

// fakeCheck is a deterministic Check implementation used in tests.
type fakeCheck struct {
	name    string
	ok      bool
	err     string
	latency time.Duration // simulated latency
}

func (f fakeCheck) Name() string { return f.name }

func (f fakeCheck) Run(ctx context.Context) Result {
	// Simulate work that respects context cancellation.
	select {
	case <-time.After(f.latency):
		// completed normally
	case <-ctx.Done():
		return Result{
			Name:  f.name,
			OK:    false,
			Error: ctx.Err().Error(),
		}
	}
	return Result{
		Name:  f.name,
		OK:    f.ok,
		Error: f.err,
	}
}

// panicCheck always panics when Run is called.
type panicCheck struct {
	name string
}

func (p panicCheck) Name() string { return p.name }

func (p panicCheck) Run(_ context.Context) Result {
	panic("boom")
}
