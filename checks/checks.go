// Package checks provides a simple framework for running startup health
// checks (database pings, TCP dials, HTTP probes, etc.) with per-check
// timeouts and optional parallelism. Results are returned as structured
// values – the caller decides how to log or display them.
package checks

import (
	"context"
	"sync"
	"time"
)

// Check is the interface that each startup probe must implement.
type Check interface {
	// Name returns a short human-readable label for the check (e.g. "postgres").
	Name() string
	// Run executes the check using the provided context for cancellation and
	// deadlines. Implementations must never panic.
	Run(ctx context.Context) Result
}

// Result holds the outcome of a single check execution.
type Result struct {
	// Name is the label of the check that produced this result.
	Name string
	// OK is true when the check passed.
	OK bool
	// Duration is the wall-clock time the check took.
	Duration time.Duration
	// Error contains a human-readable error message; empty when OK is true.
	Error string
}

// Runner executes a set of Check values with a per-check timeout.
type Runner struct {
	// TimeoutPerCheck is applied to each individual check. Zero means no
	// extra timeout beyond the parent context.
	TimeoutPerCheck time.Duration
	// Parallel, when true, runs all checks concurrently. Results are still
	// returned in the same order as the input slice.
	Parallel bool
}

// Run executes the given checks and returns one Result per check, in input
// order. It never panics; individual check failures are captured in the
// corresponding Result.
func (r Runner) Run(ctx context.Context, chks ...Check) []Result {
	results := make([]Result, len(chks))

	if r.Parallel {
		var wg sync.WaitGroup
		for i, c := range chks {
			wg.Add(1)
			go func(idx int, chk Check) {
				defer wg.Done()
				results[idx] = r.runOne(ctx, chk)
			}(i, c)
		}
		wg.Wait()
	} else {
		for i, c := range chks {
			results[i] = r.runOne(ctx, c)
		}
	}

	return results
}

// runOne executes a single check with a per-check timeout derived from the
// parent context. Panics inside the check are recovered and reported as
// failures.
func (r Runner) runOne(ctx context.Context, c Check) (res Result) {
	// Recover from panics so the runner never crashes.
	defer func() {
		if rv := recover(); rv != nil {
			res.OK = false
			res.Error = "panic: " + stringify(rv)
		}
	}()

	if r.TimeoutPerCheck > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.TimeoutPerCheck)
		defer cancel()
	}

	start := time.Now()
	res = c.Run(ctx)
	res.Duration = time.Since(start)
	res.Name = c.Name()
	return res
}

// DefaultRunner returns a Runner with sensible defaults: 2 s per-check
// timeout and parallel execution enabled.
func DefaultRunner() Runner {
	return Runner{
		TimeoutPerCheck: 2 * time.Second,
		Parallel:        true,
	}
}

// stringify converts an arbitrary recovered value to a string.
func stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if e, ok := v.(error); ok {
		return e.Error()
	}
	return "unknown panic"
}
