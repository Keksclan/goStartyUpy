package checks

import (
	"context"
	"fmt"
	"time"
)

// FuncCheck adapts a plain function into a Check. It is the simplest way to
// create a custom startup check without defining a new type.
type FuncCheck struct {
	// Label is the human-readable name for this check.
	Label string
	// Fn is the function executed by Run. It should return nil on success.
	Fn func(ctx context.Context) error
}

// Name returns the check label.
func (c FuncCheck) Name() string { return c.Label }

// Run executes Fn, measures its duration, and converts the outcome into a
// Result. If Fn is nil the check always fails. Panics inside Fn are recovered
// and reported as failures.
func (c FuncCheck) Run(ctx context.Context) Result {
	start := time.Now()

	if c.Fn == nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    "nil check function",
		}
	}

	err := safeCall(ctx, c.Fn)

	if err != nil {
		return Result{
			Name:     c.Label,
			OK:       false,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}

	return Result{
		Name:     c.Label,
		OK:       true,
		Duration: time.Since(start),
	}
}

// New creates a Check backed by the given function. If label is empty it
// defaults to "custom".
func New(label string, fn func(ctx context.Context) error) Check {
	if label == "" {
		label = "custom"
	}
	return FuncCheck{Label: label, Fn: fn}
}

// Bool creates a Check from a function that returns a boolean and an error.
// The check passes only when ok is true and err is nil. If label is empty it
// defaults to "custom".
func Bool(label string, fn func(ctx context.Context) (bool, error)) Check {
	if label == "" {
		label = "custom"
	}
	if fn == nil {
		return FuncCheck{Label: label, Fn: nil}
	}
	return FuncCheck{
		Label: label,
		Fn: func(ctx context.Context) error {
			ok, err := fn(ctx)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("check returned false")
			}
			return nil
		},
	}
}

// safeCall executes fn and recovers from panics, returning them as errors.
func safeCall(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	defer func() {
		if rv := recover(); rv != nil {
			err = fmt.Errorf("panic: %s", stringify(rv))
		}
	}()
	return fn(ctx)
}
