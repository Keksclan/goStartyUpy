package checks

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// GroupOptions configures how a Group executes its child checks.
type GroupOptions struct {
	// Parallel, when true, runs child checks concurrently.
	Parallel bool
	// TimeoutPerCheck is applied to each child check individually.
	// Zero means no extra timeout beyond the parent context.
	TimeoutPerCheck time.Duration
}

// Group bundles multiple checks into a single composite Check. When run it
// executes all children and aggregates their results: the group passes only
// when every child passes.
type Group struct {
	// Label is the human-readable name for the group.
	Label string
	// Checks is the ordered list of child checks.
	Checks []Check
}

// Name returns the group label, defaulting to "group" when empty.
func (g Group) Name() string {
	if g.Label == "" {
		return "group"
	}
	return g.Label
}

// Run executes all child checks sequentially using a local Runner, then
// aggregates the results into a single Result. OK is true only if every child
// passed; Error contains a compact summary of failing children.
func (g Group) Run(ctx context.Context) Result {
	start := time.Now()

	if len(g.Checks) == 0 {
		return Result{
			Name:     g.Name(),
			OK:       true,
			Duration: time.Since(start),
		}
	}

	runner := Runner{Parallel: false}
	results := runner.Run(ctx, g.Checks...)

	allOK := true
	var failures []string
	for _, r := range results {
		if !r.OK {
			allOK = false
			msg := r.Name
			if r.Error != "" {
				msg += ": " + r.Error
			}
			failures = append(failures, msg)
		}
	}

	res := Result{
		Name:     g.Name(),
		OK:       allOK,
		Duration: time.Since(start),
	}
	if !allOK {
		res.Error = fmt.Sprintf("%d failing: %s", len(failures), strings.Join(failures, "; "))
	}
	return res
}

// NewGroup creates a composite Check that bundles the given children. If label
// is empty it defaults to "group". The opts parameter controls parallelism and
// per-child timeouts.
func NewGroup(label string, opts GroupOptions, chks ...Check) Check {
	if label == "" {
		label = "group"
	}
	if opts.Parallel || opts.TimeoutPerCheck > 0 {
		// Wrap children so the group's Run uses the caller-specified options.
		return &configuredGroup{
			label:  label,
			checks: chks,
			opts:   opts,
		}
	}
	return Group{Label: label, Checks: chks}
}

// configuredGroup is an internal variant of Group that honours GroupOptions.
type configuredGroup struct {
	label  string
	checks []Check
	opts   GroupOptions
}

func (g *configuredGroup) Name() string { return g.label }

func (g *configuredGroup) Run(ctx context.Context) Result {
	start := time.Now()

	if len(g.checks) == 0 {
		return Result{
			Name:     g.label,
			OK:       true,
			Duration: time.Since(start),
		}
	}

	runner := Runner{
		Parallel:        g.opts.Parallel,
		TimeoutPerCheck: g.opts.TimeoutPerCheck,
	}
	results := runner.Run(ctx, g.checks...)

	allOK := true
	var failures []string
	for _, r := range results {
		if !r.OK {
			allOK = false
			msg := r.Name
			if r.Error != "" {
				msg += ": " + r.Error
			}
			failures = append(failures, msg)
		}
	}

	res := Result{
		Name:     g.label,
		OK:       allOK,
		Duration: time.Since(start),
	}
	if !allOK {
		res.Error = fmt.Sprintf("%d failing: %s", len(failures), strings.Join(failures, "; "))
	}
	return res
}
