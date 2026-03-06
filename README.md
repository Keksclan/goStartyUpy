<div align="center">

# goStartyUpy

**Zero-dependency Go library for production-ready startup banners with build metadata, runtime information, and structured health checks.**

[![Go Reference](https://pkg.go.dev/badge/github.com/keksclan/goStartyUpy.svg)](https://pkg.go.dev/github.com/keksclan/goStartyUpy)
[![Go Report Card](https://goreportcard.com/badge/github.com/keksclan/goStartyUpy)](https://goreportcard.com/report/github.com/keksclan/goStartyUpy)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Keksclan/goStartyUpy)](https://github.com/Keksclan/goStartyUpy/blob/master/go.mod)
[![License](https://img.shields.io/github/license/Keksclan/goStartyUpy)](https://github.com/Keksclan/goStartyUpy/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/Keksclan/goStartyUpy)](https://github.com/Keksclan/goStartyUpy/releases)

<!-- TODO: Add CI badge when GitHub Actions workflow is configured -->
<!-- [![CI](https://github.com/Keksclan/goStartyUpy/actions/workflows/ci.yml/badge.svg)](https://github.com/Keksclan/goStartyUpy/actions/workflows/ci.yml) -->

</div>

---

### Why goStartyUpy?

goStartyUpy gives every Go service a **Spring Boot–style startup banner** — instantly showing build metadata, runtime info, and dependency health in one glance.

🎨 **6 Banner Styles** — Spring, Classic, Box, Mini, Block, or your own ASCII art \
🔧 **Build Metadata** — Version, commit, branch, build time — injected via `-ldflags` \
🖥️ **Runtime Info** — Go version, OS/Arch, PID — captured automatically \
✅ **Health Checks** — SQL, TCP, HTTP, Redis — parallel or sequential, with timeout \
📦 **Zero Dependencies** — Pure `stdlib`. Adds nothing to your `go.sum` \
🛡️ **Panic-Safe & Deterministic** — Stable output, no side effects, all errors caught

---

## Feature Overview

| Feature | Description |
|---------|-------------|
| **6 Banner Styles** | `spring` (default), `classic`, `box`, `mini`, `block`, or custom ASCII art via the `Banner` field |
| **Build Metadata** | Version, BuildTime, Commit, Branch, Dirty — injected at compile time via `-ldflags` |
| **Runtime Info** | Go version, OS/Arch, PID — automatically captured at runtime |
| **4 Built-in Checks** | `SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck` — all without external dependencies |
| **Custom Checks** | `checks.New()`, `checks.Bool()`, `checks.NewGroup()` — or implement the `Check` interface |
| **Parallel & Sequential** | `Runner` supports both modes with configurable per-check timeout |
| **Environment Detection** | Automatic via the `GO_STARTYUPY_ENV` environment variable when not explicitly set |
| **ANSI Colors** | Optional via `Color: true` — plain text without escape sequences by default |
| **ASCII-Only Mode** | `ASCIIOnly: true` replaces Unicode box-drawing characters with plain ASCII (`+`, `-`, `\|`) |
| **Banner Width** | `BannerWidth` truncates each line to a maximum width |
| **Deterministic** | Stable output order, no randomness, no side effects |
| **Panic-Safe** | All errors are caught and returned as structured `Result` objects |

---

## Installation

```bash
go get github.com/keksclan/goStartyUpy
```

**Prerequisite:** Go 1.24 or newer.

The module has **no transitive dependencies**. After `go get`, your `go.sum` will contain only the single entry for `goStartyUpy` itself.

### Import Paths

```go
import (
    "github.com/keksclan/goStartyUpy/banner"       // Banner rendering, Options, BuildInfo
    "github.com/keksclan/goStartyUpy/checks"       // Health checks, Runner, Check interface
    "github.com/keksclan/goStartyUpy/configcheck"  // Configuration validation for goConfy structs
    "github.com/keksclan/goStartyUpy/version"      // Module version (e.g., "0.2.0")
)
```

- **`banner`** — Main package. Contains `Options`, `BuildInfo`, `Render()`, `RenderWithChecks()`, all banner functions, and the build metadata variables.
- **`checks`** — Health check system. Contains the `Check` interface, `Runner`, all built-in checks, and helper constructors.
- **`configcheck`** — Configuration validation. Validates that all required fields in a goConfy config struct are populated.
- **`version`** — Exposes the module version as the `ModuleVersion` constant.

---

## Quickstart

### Minimal Banner (Without Checks)

```go
package main

import (
    "fmt"

    "github.com/keksclan/goStartyUpy/banner"
)

func main() {
    opts := banner.Options{
        ServiceName: "my-service",
    }
    info := banner.CurrentBuildInfo()
    fmt.Print(banner.Render(opts, info))
}
```

This produces a Spring Boot–style ASCII art wordmark with all detected build/runtime metadata.

### Full Example (Banner + Checks)

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/keksclan/goStartyUpy/banner"
    "github.com/keksclan/goStartyUpy/checks"
)

func main() {
    opts := banner.Options{
        ServiceName: "order-service",
        Environment: "production",
        Extra: map[string]string{
            "HTTP":  ":8080",
            "gRPC":  ":9090",
        },
    }
    info := banner.CurrentBuildInfo()

    runner := checks.DefaultRunner() // 2s timeout, parallel
    results := runner.Run(context.Background(),
        checks.New("env-DATABASE_URL", func(ctx context.Context) error {
            if os.Getenv("DATABASE_URL") == "" {
                return fmt.Errorf("DATABASE_URL is not set")
            }
            return nil
        }),
        checks.TCPDialCheck{Address: "localhost:5432", Label: "postgres-tcp"},
        checks.HTTPGetCheck{URL: "http://localhost:8080/healthz", Label: "self-http"},
    )

    fmt.Print(banner.RenderWithChecks(opts, info, results))
}
```

### Compiling with Build Metadata

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

Or directly:

```bash
go build -ldflags "$(./scripts/ldflags.sh)" ./cmd/myservice/
```

---

## Build Metadata (`-ldflags`)

The `banner` package provides five link-time variables that are injected at compile time via `-ldflags`. These values automatically appear in the rendered banner.

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `banner.Version` | `string` | Semantic version or `git describe` output (e.g., `v1.2.3`, `v1.2.3-4-gabcdef1`) | `"dev"` |
| `banner.BuildTime` | `string` | UTC build timestamp in RFC 3339 format (e.g., `2026-03-05T18:00:00+01:00`) | `"unknown"` |
| `banner.Commit` | `string` | Short Git commit hash (e.g., `abcdef1`) | `"unknown"` |
| `banner.Branch` | `string` | Git branch used for the build (e.g., `master`, `feature/foo`) | `"unknown"` |
| `banner.Dirty` | `string` | `"true"` if the working tree had uncommitted changes at build time, otherwise `"false"` | `"false"` |

**How does this work?**

Go's linker allows overriding variable values at compile time. The `-X` flags set the package variables directly in the compiled binary without modifying source code. `CurrentBuildInfo()` then reads these values at runtime.

### Makefile (Recommended)

The included `Makefile` collects Git metadata automatically via `scripts/ldflags.sh`:

```bash
make build-example   # Compiles the example binary with all metadata
make run-example     # Compiles and runs the example
make test            # Runs all unit tests (go test ./...)
make lint            # go vet + gofmt check on all packages
make clean           # Removes build artifacts (bin/)
```

For your own service:

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

### Manual Build

If you prefer not to use the Makefile, you can set the ldflags directly:

```bash
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date -Iseconds)
DIRTY=$(git diff --quiet && echo "false" || echo "true")

go build -ldflags "\
  -X 'github.com/keksclan/goStartyUpy/banner.Version=${VERSION}' \
  -X 'github.com/keksclan/goStartyUpy/banner.BuildTime=${BUILD_TIME}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Commit=${COMMIT}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Branch=${BRANCH}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Dirty=${DIRTY}'" \
  ./cmd/myservice/
```

### `scripts/ldflags.sh` Helper

The POSIX-sh-compatible script outputs the complete ldflags string — ideal for CI/CD pipelines or other build systems:

```bash
# Standard usage:
LDFLAGS="$(./scripts/ldflags.sh)" go build -ldflags "$LDFLAGS" ./cmd/myservice

# Override module path (if your import path differs):
MODULE=github.com/my/repo ./scripts/ldflags.sh
```

---

## Banner Styles

The library supports **6 different banner styles**, controlled via `Options.BannerStyle`. If `Options.Banner` is empty, the style determines the automatically generated banner. If `Options.Banner` is set, it is **used directly** (raw mode) and `BannerStyle` is ignored.

**Style Overview:**

| Style | `BannerStyle` Value | Height | Font Technique | Direct Function |
|-------|---------------------|--------|----------------|-----------------|
| Spring (Default) | `"spring"` or `""` | 5 lines | Underscores / Pipes / Slashes | `SpringLikeBanner(name, asciiOnly)` |
| Classic | `"classic"` | 5 lines | Slashes / Backslashes / Underscores | `ClassicLikeBanner(name, asciiOnly)` |
| Box | `"box"` | 3 lines | Unicode box-drawing characters (or ASCII) | `BoxBanner(name, asciiOnly)` |
| Mini | `"mini"` | 3 lines | Compact ASCII glyphs | `MiniBanner(name, asciiOnly)` |
| Block | `"block"` | 5 lines | Thick `#` characters | `BlockBanner(name, asciiOnly)` |
| Custom (Raw) | — | any | Custom ASCII art | — |

**Character Support (All Built-in Fonts):**

Each built-in font supports the same characters: **A–Z**, **0–9**, **`-`**, **`_`**, and **space**. The `ServiceName` is automatically converted to uppercase. Unsupported characters are replaced by a **`?` fallback glyph**. Spaces are normalized to `-`.

### Spring Style (Default)

`BannerStyle: "spring"` (or empty, since `"spring"` is the default) produces a large ASCII art wordmark inspired by the Spring Boot startup banner. The font uses underscores (`_`), pipes (`|`), and slashes (`/`, `\`).

Below the wordmark is the tagline `:: goStartyUpy ::` with an optional environment suffix (only when detected via `GO_STARTYUPY_ENV`, see [Environment Detection](#environment-detection-go_startyupy_env)).

```go
opts := banner.Options{
    ServiceName: "my-svc",
    // BannerStyle defaults to "spring"
}
```

**Example Output:**

```
 __  __  __   __         ____   __     __  ____
|  \/  | \ \ / /        / ___| \ \   / / / ___|
| |\/| |  \ V /  _____  \___ \  \ \ / / | |
| |  | |   | |  |_____|  ___) |  \ V /  | |___
|_|  |_|   |_|          |____/    \_/    \____|

 :: goStartyUpy ::
```

**Direct call** (without `Options`/`Render()`):

```go
art := banner.SpringLikeBanner("my-svc", false)
fmt.Println(art)
```

### Classic Style

`BannerStyle: "classic"` produces a banner with a slash/backslash/underscore font style reminiscent of traditional Java framework startup banners. Below the wordmark, **two configurable taglines** are printed:

| Tagline | Default Value | Example |
|---------|---------------|---------|
| `Tagline1` | `"<ServiceName> <Version>"` | `"my-service v1.2.3"` |
| `Tagline2` | `"Build: <BuildTime>  Commit: <Commit> [Branch] [Dirty]"` | `"Build: 2026-03-05  Commit: abcdef1  Branch: master"` |

Both taglines can be overridden via `Options.Tagline1` and `Options.Tagline2`:

```go
opts := banner.Options{
    ServiceName: "my-svc",
    BannerStyle: "classic",
    Tagline1:    "My Service v2.0.0",
    Tagline2:    "Powered by goStartyUpy",
}
```

**`ShowDetails` Option:** Controls whether the key/value info block (Service, Version, Go version, etc.) is displayed in classic mode. This is a `*bool` pointer. By default, details are shown (`nil` = true). Set explicitly to `false` to hide them:

```go
hide := false
opts := banner.Options{
    ServiceName: "my-svc",
    BannerStyle: "classic",
    ShowDetails: &hide,   // Details block will not be printed
}
```

**Direct call:**

```go
art := banner.ClassicLikeBanner("my-svc", false)
fmt.Println(art)
```

### Box Style

`BannerStyle: "box"` produces the classic box banner using Unicode box-drawing characters (`┌`, `─`, `┐`, `│`, `└`, `┘`). The service name is centered inside the box.

```go
opts := banner.Options{
    ServiceName: "my-service",
    BannerStyle: "box",
}
```

**Output:**

```
┌───────────────────────────┐
│        MY-SERVICE         │
└───────────────────────────┘
```

**ASCII-Only Mode:** Set `Options.ASCIIOnly = true` to replace Unicode box-drawing characters with plain ASCII (`+`, `-`, `|`):

```
+---------------------------+
|        MY-SERVICE         |
+---------------------------+
```

**Direct call:**

```go
art := banner.BoxBanner("my-service", true) // true = ASCII-only
fmt.Println(art)
```

### Mini Style

`BannerStyle: "mini"` produces a **compact 3-line** ASCII art wordmark. Ideal for narrow terminals or logs with limited vertical space.

```go
opts := banner.Options{
    ServiceName: "go",
    BannerStyle: "mini",
}
```

**Example Output (`"GO"`):**

```
 __  _
| _ | |
|__||_|
```

**Direct call:**

```go
art := banner.MiniBanner("my-svc", false)
fmt.Println(art)
```

### Block Style

`BannerStyle: "block"` produces a **thick 5-line** ASCII art wordmark where each letter is built from `#` characters. Highly visible even in noisy log output.

```go
opts := banner.Options{
    ServiceName: "go",
    BannerStyle: "block",
}
```

**Example Output (`"GO"`):**

```
  ####   ###
 #      #   #
 #  ##  #   #
 #   #  #   #
  ####   ###
```

**Direct call:**

```go
art := banner.BlockBanner("my-svc", false)
fmt.Println(art)
```

### Custom Banner (Raw)

To use your own ASCII art, simply set `Options.Banner`. The value is **used directly** without any processing. `BannerStyle` is ignored in this case.

```go
opts := banner.Options{
    ServiceName: "my-service",
    Banner: `
   ╔═══════════════════════════════════╗
   ║     ★  MY AWESOME SERVICE  ★     ║
   ╚═══════════════════════════════════╝`,
}
```

**Tip:** You can use tools like [patorjk.com/software/taag](http://patorjk.com/software/taag/) to generate custom ASCII art fonts and insert them as the `Banner` string.

### Banner Width (`BannerWidth`)

Set `Options.BannerWidth` to a positive integer to **hard-truncate** each banner line to that maximum width. A value of `0` (default) means no restriction.

```go
opts := banner.Options{
    ServiceName: "my-very-long-service-name",
    BannerWidth: 60, // Each line is truncated after 60 characters
}
```

This is useful when the generated banner is too wide for your terminal or logging system.

---

## Environment Detection (`GO_STARTYUPY_ENV`)

goStartyUpy supports **automatic environment detection** via the `GO_STARTYUPY_ENV` environment variable. The behavior is as follows:

| Scenario | `Options.Environment` | `GO_STARTYUPY_ENV` | Result in Banner |
|----------|----------------------|--------------------|------------------|
| Explicitly set | `"production"` | any | No suffix displayed |
| Detected from env var | `""` (empty) | `"staging"` | Suffix `(staging)` is displayed |
| Nothing set | `""` (empty) | not set / empty | No suffix displayed |

**Rule:** The environment suffix (e.g., `(staging)`, `(dev)`) appears in the banner header **only** when the value originates from the `GO_STARTYUPY_ENV` environment variable. If `Options.Environment` is explicitly set in code, **no** suffix is displayed — the value is used internally but not shown in the banner.

**Why this design?**

- Explicitly set values in code are known to the developer — no visual hint needed.
- Values from environment variables may be unexpected (e.g., incorrect configuration in a CI/CD pipeline) — a visual hint in the banner helps with debugging.

**Example with environment variable:**

```bash
export GO_STARTYUPY_ENV=staging
go run ./cmd/myservice/
```

The banner then shows:

```
 :: goStartyUpy :: (staging)
```

**Example without environment variable (explicit):**

```go
opts := banner.Options{
    ServiceName: "my-service",
    Environment: "production",  // Explicit — no suffix in the banner
}
```

---

## Options Reference (Complete)

The `banner.Options` struct controls all aspects of banner rendering. Each field is documented with its type, default value, and description:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `ServiceName` | `string` | `""` | Name of the service. Displayed in the banner and the info section. If empty, `"SERVICE"` is used as a fallback. |
| `Environment` | `string` | `""` | Runtime environment (e.g., `"production"`, `"staging"`). If empty, `GO_STARTYUPY_ENV` is checked. Appears in the info section. |
| `Banner` | `string` | `""` | Custom ASCII art. If set, automatic banner generation is skipped and this text is used directly. |
| `BannerStyle` | `string` | `"spring"` | Controls the automatically generated banner style: `"spring"`, `"classic"`, `"box"`, `"mini"`, `"block"`. Ignored when `Banner` is set. |
| `BannerWidth` | `int` | `0` | Maximum width per banner line. `0` = no restriction. Positive values hard-truncate each line. |
| `Separator` | `string` | `"═"` (Unicode) | Character for the separator line between the banner and info section. In ASCII-only mode, `"="` is used. |
| `ASCIIOnly` | `bool` | `false` | When `true`, all Unicode characters (box-drawing, separator) are replaced with plain ASCII. |
| `Color` | `bool` | `false` | When `true`, the output is colored with ANSI escape sequences. Plain text by default. |
| `Extra` | `map[string]string` | `nil` | Additional key/value pairs displayed in the info section (e.g., `"HTTP": ":8080"`). |
| `Tagline1` | `string` | `""` | Overrides the first tagline in classic style. If empty, the default is generated. |
| `Tagline2` | `string` | `""` | Overrides the second tagline in classic style. If empty, the default is generated. |
| `ShowDetails` | `*bool` | `nil` | Controls display of the details block in classic style. `nil` = show, `&false` = hide. |

**Internal Fields** (unexported; not part of the public API):

| Field | Type | Description |
|-------|------|-------------|
| `environmentFromEnv` | `bool` | Set internally to `true` when the environment originates from `GO_STARTYUPY_ENV`. Controls the suffix display. |

---

## Check System

The `checks` package provides a complete startup check system for verifying dependencies (databases, caches, HTTP services) before accepting traffic.

### `Check` Interface

Every startup check implements the `Check` interface:

```go
type Check interface {
    Name() string                        // Human-readable name of the check
    Run(ctx context.Context) Result      // Executes the check, returns Result
}
```

The `Result` struct contains the outcome:

```go
type Result struct {
    Name     string        // Name of the check
    OK       bool          // true = passed, false = failed
    Duration time.Duration // Execution duration
    Error    string        // Error message (empty on success)
}
```

**Important:** Checks **never** panic. All panics within check functions are automatically caught and returned as a `Result` with `OK: false` and a corresponding error message.

### `Runner` — Check Execution

The `Runner` executes checks with a configurable timeout. It supports both **parallel** and **sequential** execution:

```go
runner := checks.Runner{
    TimeoutPerCheck: 2 * time.Second,  // Timeout per individual check
    Parallel:        true,              // true = parallel, false = sequential
}
results := runner.Run(ctx, check1, check2, check3)
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `TimeoutPerCheck` | `time.Duration` | `0` | Timeout per check. `0` = no additional timeout (only the provided context). |
| `Parallel` | `bool` | `false` | `true` = all checks run concurrently in their own goroutines. `false` = sequential execution in input order. |

**`DefaultRunner()`** returns a preconfigured runner with a 2-second timeout and parallel execution:

```go
runner := checks.DefaultRunner()
// Equivalent to:
// runner := checks.Runner{TimeoutPerCheck: 2 * time.Second, Parallel: true}
```

**Result Order:** Regardless of the execution mode (parallel or sequential), the results are **always returned in the same order** as the input checks. This makes the output deterministic and testable.

### Built-in Checks (4 Types)

#### `SQLPingCheck` — SQL Database

Pings an `*sql.DB` connection via `PingContext`. Useful for PostgreSQL, MySQL, SQLite, and any other `database/sql`-compatible driver.

```go
check := checks.SQLPingCheck{
    DB:        db,           // *sql.DB handle (must not be nil)
    NameLabel: "postgres",   // Human-readable name
}
```

| Field | Type | Description |
|-------|------|-------------|
| `DB` | `*sql.DB` | The database handle. If `nil`, the check fails with `"sql.DB is nil"`. |
| `NameLabel` | `string` | Name of the check in the output. |

#### `TCPDialCheck` — TCP Port

Checks whether a TCP endpoint is reachable by establishing a connection and immediately closing it. Ideal for databases, caches, or other TCP-based services.

```go
check := checks.TCPDialCheck{
    Address: "localhost:5432",   // host:port format
    Label:   "postgres-tcp",     // Human-readable name
}
```

| Field | Type | Description |
|-------|------|-------------|
| `Address` | `string` | TCP address in `host:port` format (e.g., `"localhost:5432"`, `"redis:6379"`). |
| `Label` | `string` | Name of the check in the output. |

#### `HTTPGetCheck` — HTTP Endpoint

Performs an HTTP GET request and checks whether the status code falls within an expected range. Useful for health endpoints of other services.

```go
check := checks.HTTPGetCheck{
    URL:               "http://localhost:8080/healthz",
    Label:             "api-health",
    ExpectedStatusMin: 200,   // Optional, default: 200
    ExpectedStatusMax: 299,   // Optional, default: 399
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `URL` | `string` | — | Full URL to probe (e.g., `"http://localhost:8080/healthz"`). |
| `Label` | `string` | — | Name of the check in the output. |
| `ExpectedStatusMin` | `int` | `200` | Lower bound (inclusive) of the acceptable status code range. |
| `ExpectedStatusMax` | `int` | `399` | Upper bound (inclusive) of the acceptable status code range. |

**Note:** The HTTP client does not use its own timeout — it relies on the runner's context deadline so that behavior is consistent across all check types.

#### `RedisPingCheck` — Redis via Raw TCP

Sends a RESP-encoded `PING` command over a raw TCP connection and expects `+PONG` as the response. **No Redis client or external dependency required** — works with any Redis-compatible server.

```go
check := checks.RedisPingCheck{
    Address: "localhost:6379",  // host:port of the Redis server
    Label:   "redis-ping",     // Human-readable name
}
```

| Field | Type | Description |
|-------|------|-------------|
| `Address` | `string` | TCP address of the Redis server in `host:port` format. |
| `Label` | `string` | Name of the check in the output. |

**Technical Detail:** The check sends the RESP array `*1\r\n$4\r\nPING\r\n` and expects `+PONG\r\n`. On unexpected responses, the check fails with the actual reply in the error text.

### Creating Custom Checks

#### Function-Based Check (`checks.New`)

The simplest way to create a custom check. Pass a label string and a function that returns `error` (`nil` = passed):

```go
envCheck := checks.New("env-DATABASE_URL", func(ctx context.Context) error {
    if os.Getenv("DATABASE_URL") == "" {
        return fmt.Errorf("DATABASE_URL is not set")
    }
    return nil
})
```

If `label` is empty, `"custom"` is used as a fallback. If `fn` is `nil`, the check always fails with `"nil check function"`.

#### Boolean Check (`checks.Bool`)

For checks that return a boolean result plus an optional error:

```go
featureFlag := checks.Bool("feature-flag", func(ctx context.Context) (bool, error) {
    return os.Getenv("ENABLE_NEW_UI") == "true", nil
})
```

The check passes **only** when `ok == true` **and** `err == nil`. If `ok == false` and `err == nil`, the error `"check returned false"` is produced.

#### Grouped Check (`checks.NewGroup`)

Combines multiple checks into a single one. The group passes **only** when **all** children pass:

```go
deps := checks.NewGroup("dependencies", checks.GroupOptions{
    Parallel:        true,                  // Run children in parallel
    TimeoutPerCheck: 3 * time.Second,       // Timeout per child check
},
    checks.SQLPingCheck{DB: db, NameLabel: "postgres"},
    checks.TCPDialCheck{Address: "localhost:6379", Label: "redis-tcp"},
    checks.RedisPingCheck{Address: "localhost:6379", Label: "redis-ping"},
)
```

| `GroupOptions` Field | Type | Default | Description |
|----------------------|------|---------|-------------|
| `Parallel` | `bool` | `false` | `true` = run child checks in parallel. |
| `TimeoutPerCheck` | `time.Duration` | `0` | Timeout per child check. `0` = no additional timeout. |

On failures, the error string contains a compact summary: `"2 failing: postgres: connection refused; redis-tcp: dial timeout"`.

#### Implementing the `Check` Interface

For more complex scenarios, you can implement the `Check` interface directly:

```go
type MyCustomCheck struct {
    // custom fields
}

func (c MyCustomCheck) Name() string { return "my-custom" }

func (c MyCustomCheck) Run(ctx context.Context) checks.Result {
    start := time.Now()
    // ... your logic ...
    return checks.Result{
        Name:     c.Name(),
        OK:       true,
        Duration: time.Since(start),
    }
}
```

### Parallel vs. Sequential Execution

| Mode | `Runner.Parallel` | Behavior |
|------|-------------------|----------|
| Parallel | `true` | Each check runs in its own goroutine. The runner waits until all are finished. Fastest overall throughput. |
| Sequential | `false` | Checks run one after another in input order. A slow check blocks subsequent ones. |

In **both modes**, the results are returned in the **same order** as the input.

---

## Configuration Validation

The `configcheck` package provides **startup-time configuration validation** for structs loaded by [goConfy](https://github.com/Keksclan/goConfy). It inspects the config struct via reflection and reports any required fields that are missing or empty — catching configuration mistakes before the service starts serving traffic.

### Why Validation?

Missing or empty configuration values often cause cryptic runtime errors (nil pointer dereferences, empty connection strings, silent fallbacks). Running a validation step at startup ensures that **all required values are present** before any component is initialized.

### Enabling Validation

Validation is **optional and off by default**. Enable it by passing `configcheck.Options{Enabled: true}` before printing the startup banner:

```go
package main

import (
    "fmt"
    "log"

    goconfy "github.com/keksclan/goConfy"
    "github.com/keksclan/goStartyUpy/banner"
    "github.com/keksclan/goStartyUpy/configcheck"
)

type AppConfig struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    Redis struct {
        Address string `yaml:"address"`
    } `yaml:"redis"`
    LogLevel string `yaml:"log_level" required:"false"`
}

func main() {
    // 1. Load config via goConfy
    cfg, err := goconfy.Load[AppConfig](goconfy.WithFile("config.yml"))
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // 2. Run config validation
    msg, err := configcheck.RunStartupCheck(configcheck.Options{
        Enabled: true,
        Config:  cfg,
    })
    if err != nil {
        fmt.Println(msg)
        log.Fatalf("Aborting startup due to configuration errors: %v", err)
    }

    // 3. Print startup banner
    fmt.Print(banner.Render(banner.Options{ServiceName: "my-service"}, banner.CurrentBuildInfo()))

    // 4. Continue service initialization...
}
```

### Startup Order

The recommended startup sequence when using config validation:

1. Load config via goConfy (`goconfy.Load[T]`)
2. Run config validation (`configcheck.RunStartupCheck`)
3. Print startup banner (`banner.Render` / `banner.RenderWithChecks`)
4. Continue service initialization

### Required vs. Optional Fields

By default, **all exported struct fields are required**. Mark a field as optional with the `required:"false"` struct tag:

```go
type Config struct {
    Host     string `yaml:"host"`                       // required (default)
    Port     int    `yaml:"port"`                       // required
    LogLevel string `yaml:"log_level" required:"false"` // optional
}
```

The validator uses the `yaml` struct tag to determine the YAML key name. If no `yaml` tag is present, the Go field name is used. Fields tagged `yaml:"-"` are skipped entirely.

### What Gets Checked

| Check | Description |
|-------|-------------|
| **Zero-value scalars** | Empty strings, zero ints/floats, false bools (when required) |
| **Nil pointers** | Pointer fields that are nil |
| **Empty slices/maps** | Nil or zero-length slices and maps |
| **Nested structs** | Recursively validated with dot-separated paths |
| **Leaf structs** | Types implementing `encoding.TextUnmarshaler` (e.g., `time.Time`) are checked as single values |

### Example Error Output

When validation fails, the output is a clear diagnostic block:

```
Config validation failed:

Missing configuration keys:
  - database.port
  - database.password
  - redis.address
```

This helps developers quickly identify and fix configuration problems.

### Convenience: `MustPassStartupCheck`

For a one-liner that logs and exits on failure:

```go
configcheck.MustPassStartupCheck(configcheck.Options{
    Enabled: true,
    Config:  cfg,
}, log.Fatalf)
```

### Full Configuration Reference

A complete example YAML file containing **every supported configuration option** is available at:

```
examples/config-validation/config.full.example.yml
```

Use it as a starting point for your own configuration. The file includes all
banner, config-check, and application-level settings with explanatory comments.

Additional example files in the same directory:

| File | Purpose |
|------|---------|
| `config.full.example.yml` | All supported options with comments |
| `config.minimal.yml` | Smallest valid configuration |
| `config_invalid.yml` | Intentionally broken config for testing |

Run the example tests to verify all configs:

```bash
cd examples/config-validation && go test ./... -v
```

### Best Practices

- **Enable validation in all environments** — it adds negligible overhead and catches drift between config files.
- **Use `required:"false"`** sparingly — only for truly optional fields with sensible zero-value behavior.
- **Run validation before the banner** — so failures are visible immediately, not buried after startup output.
- **Use goConfy's strict mode** (default) together with configcheck — goConfy rejects unknown YAML keys, configcheck catches missing ones.
- **Keep example YAML files in sync** — the full example config acts as living documentation; update it when adding new config options.

---

## Module Version vs. Build Version

goStartyUpy strictly distinguishes between **two different version values** that are independent of each other:

| Value | Package | Purpose | Set by |
|-------|---------|---------|--------|
| `version.ModuleVersion` | `version` | Release version of the **library** itself (e.g., `"0.2.0"`) | In the source code (`version/version.go`) |
| `banner.Version` | `banner` | Build version of the **service binary** (e.g., `"v1.2.3"`) | `-ldflags` at compile time |

**Why two versions?**

- `ModuleVersion` tells you which version of goStartyUpy you are using as a dependency.
- `banner.Version` tells you which version of your own service is currently running.

Both are independent values. Your service can use `goStartyUpy@v0.2.0` and still be tagged as `v3.7.2`.

```go
import "github.com/keksclan/goStartyUpy/version"

fmt.Println("goStartyUpy Library:", version.ModuleVersion) // "0.2.0"
```

---

## `BuildInfo` Struct

`CurrentBuildInfo()` creates a `BuildInfo` struct from the package-level link-time variables. This struct is passed to `Render()` and `RenderWithChecks()`:

```go
info := banner.CurrentBuildInfo()
```

### Struct Fields

These are the actual fields of the `BuildInfo` struct, set via `-ldflags` at build time:

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `Version` | `string` | `-ldflags` | Build version of the service (default: `"dev"`) |
| `BuildTime` | `string` | `-ldflags` | UTC build timestamp (default: `"unknown"`) |
| `Commit` | `string` | `-ldflags` | Short Git commit hash (default: `"unknown"`) |
| `Branch` | `string` | `-ldflags` | Git branch (default: `"unknown"`) |
| `Dirty` | `string` | `-ldflags` | `"true"` / `"false"` for uncommitted changes |

### Rendered Values

The following values are **not** fields on `BuildInfo`. They are collected internally by the renderer at render time from the `runtime` package and `os.Getpid()`, and appear in the banner output alongside the struct fields:

| Value | Type | Source | Description |
|-------|------|--------|-------------|
| `GoVersion` | `string` | `runtime.Version()` | Go version (e.g., `"go1.26"`) |
| `OS` | `string` | `runtime.GOOS` | Operating system (e.g., `"linux"`, `"darwin"`) |
| `Arch` | `string` | `runtime.GOARCH` | CPU architecture (e.g., `"amd64"`, `"arm64"`) |
| `PID` | `int` | `os.Getpid()` | Process ID of the running binary |

---

## Render Functions

The `banner` package provides two main render functions:

### `Render(opts Options, info BuildInfo) string`

Produces the complete startup banner **without checks**. Returns the entire banner as a string (ready for `fmt.Print`):

```go
output := banner.Render(opts, info)
fmt.Print(output)
```

### `RenderWithChecks(opts Options, info BuildInfo, results []checks.Result) string`

Produces the complete startup banner **with check results**. The check results are appended as a list at the end of the banner:

```go
output := banner.RenderWithChecks(opts, info, results)
fmt.Print(output)
```

**Check Output Format:**

```
Checks:
  [OK]   postgres (12ms)
  [OK]   redis-tcp (3ms)
  [FAIL] kafka: connection refused (1.2s)

Startup Complete
```

- `[OK]` = Check passed (green when `Color: true`)
- `[FAIL]` = Check failed (red when `Color: true`), with error message

---

## Public API Stability

The following elements are considered **public API** and are subject to versioning guarantees:

- All exported types, functions, variables, and constants in the `banner`, `checks`, and `version` packages.
- The `Check` interface and its contract.
- The fields of the `Options`, `BuildInfo`, `Result`, `Runner`, and `GroupOptions` structs.
- The fields of the built-in check structs (`SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck`).

**Not part of the public API** (may change without notice):

- All unexported (lowercase) identifiers.
- The `example/` directory.
- The `scripts/` directory.
- Internal font data and render helper functions.

---

## Versioning Strategy

This project follows [Semantic Versioning 2.0.0](https://semver.org/):

| Version Part | When? | Example |
|--------------|-------|---------|
| **MAJOR** (`X.0.0`) | Incompatible API changes (removing/renaming exported symbols, changing function signatures, breaking changes to the `Check` interface) | `1.0.0` → `2.0.0` |
| **MINOR** (`0.X.0`) | New features, backward-compatible (new check types, new `Options` fields, new helper functions) | `0.1.0` → `0.2.0` |
| **PATCH** (`0.0.X`) | Backward-compatible bug fixes and documentation corrections | `0.1.0` → `0.1.1` |

**Note:** As long as the module is at `0.x.y`, the API may change between minor versions. A `1.0.0` release signals a stable API commitment.

---

## Release Process

1. **Update version:** Set `ModuleVersion` in `version/version.go` to the new version.
2. **Update CHANGELOG:** Move entries from `[Unreleased]` into a new version section with date.
3. **Commit:**
   ```bash
   git add -A
   git commit -m "release: v0.2.0"
   ```
4. **Tag and push:**
   ```bash
   git tag v0.2.0
   git push origin master v0.2.0
   ```
5. **Consumers can pin the version:**
   ```bash
   go get github.com/keksclan/goStartyUpy@v0.2.0
   ```

---

## Example Programs

The `example/` directory contains runnable programs for various use cases:

| Example | Description | Run with |
|---------|-------------|----------|
| `example/` | Full demo: Custom checks, groups, built-in checks | `make run-example` |
| `example/simple/` | Minimal banner without checks | `go run ./example/simple/` |
| `example/basic_start/` | Simplest possible usage — banner with defaults | `go run ./example/basic_start/` |
| `example/custom_banner/` | Custom ASCII art as banner | `go run ./example/custom_banner/` |
| `example/env_aware_start/` | Automatic environment detection via `GO_STARTYUPY_ENV` | `go run ./example/env_aware_start/` |
| `example/ascii_only/` | ASCII-only mode for terminals without Unicode | `go run ./example/ascii_only/` |
| `example/checks_demo/` | All built-in check types (SQL, TCP, HTTP, Redis) | `go run ./example/checks_demo/` |
| `example/custom_checks/` | Function-based, boolean, and grouped checks | `go run ./example/custom_checks/` |
| `example/font_preview/` | Prints the big-font ASCII wordmark for a service name | `go run ./example/font_preview/` |
| `example/config_validation/` | Configuration validation with configcheck | `go run ./example/config_validation/` |

```bash
# Simplest start:
go run ./example/simple/

# Full demo with build metadata:
make run-example
```

---

## Output Example (Box Style with Checks)

The following example shows the complete output in box style with build metadata, extra fields, and check results:

```
┌──────────────────────────────┐
│        ORDER-SERVICE         │
└──────────────────────────────┘
════════════════════════════════════════════════════════════
  Service     : order-service
  Environment : staging
  Version     : v1.2.3
  BuildTime   : 2026-02-24T09:00:00Z
  Commit      : abcdef1
  Branch      : master
  Dirty       : false
  Go          : go1.26
  OS/Arch     : linux/amd64
  PID         : 12345
  HTTP        : :8080

Checks:
  [OK]   postgres (12ms)
  [OK]   redis-tcp (3ms)
  [OK]   self-http (8ms)
  [OK]   redis-ping (2ms)

Startup Complete
```

**Output Structure:**

1. **Banner** — ASCII art or box (depending on style)
2. **Separator** — Separator line (`═══...` or `===...` in ASCII mode)
3. **Info Section** — Key/value pairs (Service, Environment, Version, BuildTime, Commit, Branch, Dirty, Go, OS/Arch, PID, plus all `Extra` entries)
4. **Checks** (only with `RenderWithChecks`) — Results with `[OK]`/`[FAIL]` status and duration
5. **Footer** — `"Startup Complete"` or check summary

---

## Security Notice

The banner prints exclusively **safe, non-secret** information (version, addresses, PID, etc.).

⚠️ **Never pass secrets** (passwords, tokens, API keys) via `Options.Extra` or other fields. The caller is responsible for what is printed.

---

## Tests

```bash
# Run all unit tests:
go test ./...

# Via Makefile (identical):
make test

# Linting (go vet + gofmt):
make lint
```

The project includes tests for:
- Banner rendering of all 6 styles
- Font rendering and fallback glyphs
- Build metadata snapshot
- Formatting and separator
- Environment detection (explicit, from env var, not set)
- Check runner (parallel and sequential)
- All built-in check types
- FuncCheck, Bool check, Group check
- Panic recovery
- Unicode and edge-case safety

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE) for the full text.

**Summary:**
- Open-source usage: Free under the MIT License, no attribution required.
- Commercial/corporate usage: MIT-permissive, no attribution required.

---

## Documentation

Detailed documentation is available in the [`docs/`](docs/) directory:

| Document | Description |
|----------|-------------|
| [`docs/architecture.md`](docs/architecture.md) | Module layout, design principles, package responsibilities, data flow |
| [`docs/startup-flow.md`](docs/startup-flow.md) | Step-by-step rendering sequence, check execution phases |
| [`docs/banner-system.md`](docs/banner-system.md) | Banner styles, font system, name normalization, ASCII/color modes |

---

## Used by

This project is used by:
- [Keksclan](https://github.com/Keksclan) — Creator of goStartyUpy
- Internal microservices — Production startup banners and health checks
- Community projects — Open-source Go services using goStartyUpy for boot diagnostics

➕ **Add your project/organization here?** Open a pull request and edit [`USED_BY.md`](USED_BY.md). Rules:
- Sort alphabetically
- 1 line per entry, no marketing
- Format: `- [Name](URL) - Short description [tags]`

---

## Topics

`golang` · `go-library` · `banner` · `startup-banner` · `cli` · `microservices` · `health-check` · `startup-checks` · `zero-dependencies` · `devops` · `ascii-art` · `build-metadata` · `spring-boot-style` · `production-ready` · `deterministic`
