# Startup Flow

This document explains the step-by-step sequence that occurs when a service calls `banner.Render()` or `banner.RenderWithChecks()`.

## Overview

The startup flow has two independent phases:

1. **Banner Rendering** — Assembles the visual output (banner art, metadata, optional check results) into a single string.
2. **Check Execution** (optional) — Runs health probes against infrastructure dependencies and produces `[]checks.Result`.

The caller controls both phases explicitly. goStartyUpy never performs I/O on its own (except reading one environment variable).

## Phase 1: Banner Rendering

### Step 1 — Environment Detection

If `Options.Environment` is empty, the renderer checks `os.Getenv("GO_STARTYUPY_ENV")`. When the variable is set, its value becomes the environment and `EnvironmentFromEnv` is set to `true` internally.

```go
// Automatic detection:
opts := banner.Options{ServiceName: "my-svc"}
// → reads GO_STARTYUPY_ENV at render time

// Explicit override (skips env var):
opts := banner.Options{ServiceName: "my-svc", Environment: "production"}
```

### Step 2 — Banner Art Generation

`resolveBanner(opts)` selects the banner art:

| Condition | Action |
|-----------|--------|
| `opts.Banner` is set | Use as-is (raw custom art) |
| `opts.BannerStyle == "spring"` (or empty) | Call `SpringLikeBanner()` — big-font wordmark |
| `opts.BannerStyle == "classic"` | Call `ClassicLikeBanner()` — slash/backslash wordmark |
| `opts.BannerStyle == "box"` | Call `BoxBanner()` — bordered box |
| `opts.BannerStyle == "mini"` | Generate 3-row mini wordmark |
| `opts.BannerStyle == "block"` | Generate 5-row block wordmark |

If `BannerWidth > 0`, each line is hard-cut to that width.

### Step 3 — Separator

A horizontal line separates the banner art from the info section:

- Default: Unicode double-line (`════...`)
- With `ASCIIOnly: true`: Plain equals signs (`====...`)
- Custom: `opts.Separator` overrides both defaults

### Step 4 — Key-Value Info Block

`buildKVs()` collects metadata pairs in a fixed order:

1. Service, Environment, Version, BuildTime, Commit, Branch, Dirty
2. Go version, OS/Arch, PID
3. All `Extra` entries (sorted alphabetically by key)

`writeAligned()` formats them with consistent label padding:

```
  Service     : my-awesome-service
  Environment : production
  Version     : v1.2.3
  Go          : go1.24
  PID         : 12345
```

### Step 5 — Check Results (only with `RenderWithChecks`)

When `results` is non-empty, a "Checks:" section is appended:

```
Checks:
  [OK]   postgres (12ms)
  [FAIL] redis-tcp (2.001s): dial tcp: connection refused
```

Followed by either `"Startup Complete"` (all OK) or `"Startup Failed"`.

### Step 6 — Final Assembly

The complete string is returned with a guaranteed trailing newline. The caller prints it:

```go
fmt.Print(banner.Render(opts, info))
// or
fmt.Print(banner.RenderWithChecks(opts, info, results))
```

## Phase 2: Check Execution

Check execution is independent of banner rendering. The caller creates checks, runs them, and passes the results to `RenderWithChecks`.

### Creating Checks

```go
// Built-in checks:
sqlCheck := checks.SQLPingCheck{DB: db, NameLabel: "postgres"}
tcpCheck := checks.TCPDialCheck{Address: "localhost:6379", Label: "redis"}
httpCheck := checks.HTTPGetCheck{URL: "http://localhost/healthz", Label: "self"}
redisCheck := checks.RedisPingCheck{Address: "localhost:6379", Label: "redis-ping"}

// Custom function check:
envCheck := checks.New("env-check", func(ctx context.Context) error {
    if os.Getenv("DATABASE_URL") == "" {
        return fmt.Errorf("DATABASE_URL not set")
    }
    return nil
})

// Boolean check:
flagCheck := checks.Bool("feature-flag", func(ctx context.Context) (bool, error) {
    return os.Getenv("ENABLE_FEATURE") == "true", nil
})

// Grouped checks:
group := checks.NewGroup("infra", checks.GroupOptions{Parallel: true},
    sqlCheck, tcpCheck, redisCheck,
)
```

### Running Checks

```go
runner := checks.DefaultRunner() // parallel, 2s per-check timeout
// or custom:
runner := checks.Runner{
    TimeoutPerCheck: 5 * time.Second,
    Parallel:        false,
}

ctx := context.Background()
results := runner.Run(ctx, sqlCheck, tcpCheck, httpCheck)
```

### Execution Modes

| Mode | Behavior |
|------|----------|
| `Parallel: false` | Checks run sequentially in input order |
| `Parallel: true` | Checks run concurrently; results are still returned in input order |

Each check receives its own timeout context derived from `TimeoutPerCheck`. If a check panics, the panic is recovered and reported as a failed `Result`.

## Complete Example

```go
func main() {
    opts := banner.Options{
        ServiceName: "order-service",
        Environment: "staging",
        Color:       true,
        Extra: map[string]string{"HTTP": ":8080"},
    }
    info := banner.CurrentBuildInfo()

    runner := checks.DefaultRunner()
    results := runner.Run(context.Background(),
        checks.TCPDialCheck{Address: "localhost:5432", Label: "postgres"},
        checks.HTTPGetCheck{URL: "http://localhost:8080/healthz", Label: "self"},
    )

    fmt.Print(banner.RenderWithChecks(opts, info, results))
}
```
