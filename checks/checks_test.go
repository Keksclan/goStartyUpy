package checks

import (
	"testing"
	"time"
)

func TestRunner_SequentialPreservesOrder(t *testing.T) {
	runner := Runner{TimeoutPerCheck: 5 * time.Second, Parallel: false}
	chks := []Check{
		fakeCheck{name: "alpha", ok: true, latency: 1 * time.Millisecond},
		fakeCheck{name: "beta", ok: false, err: "down", latency: 1 * time.Millisecond},
		fakeCheck{name: "gamma", ok: true, latency: 1 * time.Millisecond},
	}

	results := runner.Run(t.Context(), chks...)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, want := range []string{"alpha", "beta", "gamma"} {
		if results[i].Name != want {
			t.Errorf("result[%d].Name = %q, want %q", i, results[i].Name, want)
		}
	}
	if results[0].OK != true {
		t.Error("alpha should be OK")
	}
	if results[1].OK != false {
		t.Error("beta should be FAIL")
	}
	if results[1].Error != "down" {
		t.Errorf("beta error = %q, want %q", results[1].Error, "down")
	}
}

func TestRunner_ParallelPreservesOrder(t *testing.T) {
	runner := Runner{TimeoutPerCheck: 5 * time.Second, Parallel: true}
	chks := []Check{
		fakeCheck{name: "first", ok: true, latency: 20 * time.Millisecond},
		fakeCheck{name: "second", ok: true, latency: 1 * time.Millisecond},
		fakeCheck{name: "third", ok: true, latency: 10 * time.Millisecond},
	}

	results := runner.Run(t.Context(), chks...)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// Order must match input, not completion order.
	for i, want := range []string{"first", "second", "third"} {
		if results[i].Name != want {
			t.Errorf("result[%d].Name = %q, want %q", i, results[i].Name, want)
		}
	}
}

func TestRunner_TimeoutProducesFail(t *testing.T) {
	runner := Runner{TimeoutPerCheck: 10 * time.Millisecond, Parallel: false}
	chks := []Check{
		// This check sleeps longer than the per-check timeout.
		fakeCheck{name: "slow", ok: true, latency: 5 * time.Second},
	}

	results := runner.Run(t.Context(), chks...)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OK {
		t.Error("expected FAIL due to timeout")
	}
	if results[0].Error == "" {
		t.Error("expected non-empty error on timeout")
	}
}

func TestRunner_PanicRecovery(t *testing.T) {
	runner := Runner{TimeoutPerCheck: time.Second, Parallel: false}
	chks := []Check{
		panicCheck{name: "panicker"},
	}

	results := runner.Run(t.Context(), chks...)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OK {
		t.Error("panicking check should be FAIL")
	}
	if results[0].Error == "" {
		t.Error("expected non-empty error on panic")
	}
}

func TestRunner_EmptyChecks(t *testing.T) {
	runner := Runner{TimeoutPerCheck: time.Second}
	results := runner.Run(t.Context())
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSQLPingCheck_NilDB(t *testing.T) {
	c := SQLPingCheck{DB: nil, NameLabel: "pg"}
	res := c.Run(t.Context())
	if res.OK {
		t.Error("expected FAIL for nil DB")
	}
	if res.Error != "sql.DB is nil" {
		t.Errorf("unexpected error: %q", res.Error)
	}
}
