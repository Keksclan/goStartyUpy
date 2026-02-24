# goStartyUpy

A zero-dependency Go library that renders a production-ready startup banner
with build metadata, runtime information, and structured health checks.
It is designed as a reusable module for any Go service that needs a clear,
deterministic startup message with optional dependency verification.

## Feature Overview

- **Build metadata injection** — inject Version, BuildTime, Commit, Branch,
  and Dirty flag via `-ldflags` at compile time.
- **Customizable banner** — auto-generated Spring Boot–style ASCII-art
  wordmark from `ServiceName` (default), classic box style, or your own
  multiline ASCII art.
- **Built-in big font (no deps)** — underscore / pipe / slash style glyphs
  for A–Z, 0–9, `-`, `_`, and space. Unknown characters render as a `?`
  fallback glyph. Deterministic, zero external dependencies.
- **ASCII fallback mode** — `ASCIIOnly: true` replaces Unicode box-drawing
  characters with plain ASCII for restricted terminals.
- **Startup checks** — verify SQL databases, TCP endpoints, HTTP services,
  and Redis connectivity before accepting traffic.
- **Custom checks** — implement the `Check` interface, or use `checks.New`,
  `checks.Bool`, and `checks.NewGroup` helpers.
- **Parallel & sequential execution** — `Runner` supports both modes with
  per-check timeouts.
- **Deterministic & testable** — stable output ordering, no randomness,
  no external dependencies.
- **Optional ANSI colors** — `Color: true` for colorized terminal output;
  plain text by default.
- **Never panics** — all errors are captured and returned as structured results.

## Installation

```bash
go get github.com/keksclan/goStartyUpy
```

Requires Go 1.26 or later.

## Quickstart

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
        ServiceName: "my-service",
        Environment: "production",
        Extra: map[string]string{
            "HTTP": ":8080",
        },
    }
    info := banner.CurrentBuildInfo()

    runner := checks.DefaultRunner()
    results := runner.Run(context.Background(),
        checks.New("env-check", func(ctx context.Context) error {
            if os.Getenv("APP_SECRET") == "" {
                return fmt.Errorf("APP_SECRET not set")
            }
            return nil
        }),
        checks.TCPDialCheck{Address: "localhost:5432", Label: "postgres-tcp"},
    )

    fmt.Print(banner.RenderWithChecks(opts, info, results))
}
```

Build with metadata:

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

## Build Metadata (ldflags)

The `banner` package exposes five link-time variables that are injected via
`-ldflags` during `go build`:

| Variable    | Description                                          | Default     |
|-------------|------------------------------------------------------|-------------|
| `Version`   | Semantic version or git describe output              | `"dev"`     |
| `BuildTime` | UTC timestamp of the build (RFC 3339)                | `"unknown"` |
| `Commit`    | Short git commit hash                                | `"unknown"` |
| `Branch`    | Git branch the binary was built from                 | `"unknown"` |
| `Dirty`     | `"true"` if the working tree had uncommitted changes | `"false"`   |

### Using the Makefile

The included `Makefile` collects git metadata automatically:

```bash
make build-example   # compile example binary with metadata
make run-example     # build and run example
make test            # run all unit tests
make lint            # go vet + gofmt check
make clean           # remove build artifacts
```

For your own service:

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

### Manual build

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

### Using `scripts/ldflags.sh`

The helper script prints the full ldflags string for integration into any
build system (POSIX sh compatible):

```bash
LDFLAGS="$(./scripts/ldflags.sh)" go build -ldflags "$LDFLAGS" ./cmd/myservice

# Override the module path if your import path differs:
MODULE=github.com/my/repo ./scripts/ldflags.sh
```

## Banner Styles

The library supports multiple banner styles controlled by `Options.BannerStyle`.
When `Options.Banner` is empty, the style determines the auto-generated banner.
When `Options.Banner` is set, the value is used as-is ("raw" mode) and
`BannerStyle` is ignored.

### Spring style (default)

`BannerStyle: "spring"` (or empty, which defaults to `"spring"`) generates a
large ASCII-art wordmark derived from `ServiceName`, followed by a tagline.
This is inspired by the Spring Boot startup banner structure — no external
dependencies, fully deterministic, built-in big font using an "underscore /
pipe / slash" style.

Supported characters: **A–Z**, **0–9**, **`-`**, **`_`**, and **space**.
Any unsupported character is replaced with a **`?`** fallback glyph.

```go
opts := banner.Options{
    ServiceName: "my-svc",
    // BannerStyle defaults to "spring"
}
```

Example output (shape):

```
 __  __  __   __         ____   __     __  ____
|  \/  | \ \ / /        / ___| \ \   / / / ___|
| |\/| |  \ V /  _____  \___ \  \ \ / / | |
| |  | |   | |  |_____|  ___) |  \ V /  | |___
|_|  |_|   |_|          |____/    \_/    \____|

 :: goStartyUpy :: (dev)
```

You can also call `banner.SpringLikeBanner(name, asciiOnly)` directly.

### Box style

`BannerStyle: "box"` generates the classic box banner:

```go
opts := banner.Options{
    ServiceName: "my-service",
    BannerStyle: "box",
}
```

```
┌───────────────────────────┐
│        MY-SERVICE         │
└───────────────────────────┘
```

Set `Options.ASCIIOnly = true` to use plain ASCII box characters (`+`, `-`, `|`).
You can also call `banner.BoxBanner(name, asciiOnly)` directly.

### Custom banner (raw)

Provide your own multiline ASCII art via `Options.Banner`:

```go
opts := banner.Options{
    ServiceName: "my-service",
    Banner: `
   ╔═══════════════════════════════════╗
   ║     ★  MY AWESOME SERVICE  ★     ║
   ╚═══════════════════════════════════╝`,
}
```

When `Banner` is set the auto-generation is skipped entirely.

### Banner width clamping

Set `Options.BannerWidth` to a positive integer to hard-cut every banner line
to that maximum width. A value of `0` (default) means no clamping.

```go
opts := banner.Options{
    ServiceName: "my-service",
    BannerWidth: 60,
}
```

## Checks System

### Check interface

Every startup probe implements the `Check` interface:

```go
type Check interface {
    Name() string
    Run(ctx context.Context) checks.Result
}
```

### Runner

`Runner` executes checks with a configurable per-check timeout. When
`Parallel` is true all checks run concurrently; results are returned in
input order regardless.

```go
runner := checks.Runner{
    TimeoutPerCheck: 2 * time.Second,
    Parallel:        true,
}
results := runner.Run(ctx, check1, check2, check3)
```

`checks.DefaultRunner()` returns a runner with 2 s timeout and parallel
execution enabled.

### Built-in checks

| Check            | Description                                                    |
|------------------|----------------------------------------------------------------|
| `SQLPingCheck`   | Pings a `*sql.DB` via `PingContext`.                           |
| `TCPDialCheck`   | Dials a TCP `host:port` and closes the connection.             |
| `HTTPGetCheck`   | Sends an HTTP GET and checks the status code range.            |
| `RedisPingCheck` | Sends a RESP `PING` command over TCP (no Redis client needed). |

### Custom checks

#### Function-based check

```go
envCheck := checks.New("env-DATABASE_URL", func(ctx context.Context) error {
    if os.Getenv("DATABASE_URL") == "" {
        return fmt.Errorf("DATABASE_URL is not set")
    }
    return nil
})
```

#### Boolean check

```go
featureFlag := checks.Bool("feature-flag", func(ctx context.Context) (bool, error) {
    return os.Getenv("ENABLE_NEW_UI") == "true", nil
})
```

#### Grouped checks

```go
deps := checks.NewGroup("dependencies", checks.GroupOptions{},
    checks.SQLPingCheck{DB: db, NameLabel: "postgres"},
    checks.TCPDialCheck{Address: "localhost:6379", Label: "redis-tcp"},
)
```

The group passes only when every child passes; the error summary lists
which children failed.

### Parallel vs sequential execution

Set `Runner.Parallel = true` to run all checks concurrently. Each check
gets its own goroutine, and the runner waits for all to complete. When
`Parallel` is false checks execute sequentially in input order. In both
modes results are returned in the same order as the input slice.

## Module Version vs Build Version

This module exposes two distinct version values:

| Value                    | Purpose                             | Set by            |
|--------------------------|-------------------------------------|-------------------|
| `version.ModuleVersion`  | Library release version (`0.1.0`)   | Source code        |
| `banner.Version`         | Service build version (`v1.2.3`)    | `-ldflags` at build time |

`ModuleVersion` tracks the goStartyUpy library release. `banner.Version`
is injected by the service that imports the library and represents the
service binary version. They are independent and serve different purposes.

```go
import "github.com/keksclan/goStartyUpy/version"

fmt.Println("Library version:", version.ModuleVersion)
```

## Public API Stability

The following are considered **public API** and are subject to the
versioning guarantees described below:

- All exported types, functions, variables, and constants in the `banner`,
  `checks`, and `version` packages.
- The `Check` interface contract.
- The `Options` and `BuildInfo` struct fields.
- The `Result` struct fields.

Internal (unexported) identifiers, the `example/` directory, and
`scripts/` are not part of the public API and may change without notice.

## Versioning Policy

This project follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** — incompatible API changes (removing/renaming exported symbols,
  changing function signatures, breaking the `Check` interface).
- **MINOR** — new features added in a backwards-compatible manner (new check
  types, new `Options` fields, new helper functions).
- **PATCH** — backwards-compatible bug fixes and documentation corrections.

While the module is at `0.x.y`, the API may change between minor versions.
A `1.0.0` release will signal a stable public API commitment.

## Release Process

1. Update `version/version.go` — set `ModuleVersion` to the new version.
2. Update `CHANGELOG.md` — move items from `[Unreleased]` to a new
   version section with the release date.
3. Commit the changes:
   ```bash
   git add -A
   git commit -m "release: v0.1.0"
   ```
4. Tag and push:
   ```bash
   git tag v0.1.0
   git push origin main v0.1.0
   ```
5. Consumers can then pin the version:
   ```bash
   go get github.com/keksclan/goStartyUpy@v0.1.0
   ```

## Examples

The `example/` directory contains runnable programs:

| Example                  | What it shows                                       |
|--------------------------|-----------------------------------------------------|
| `example/`               | Full demo: custom checks, groups, built-in checks   |
| `example/simple/`        | Minimal banner-only usage (no checks)               |
| `example/custom_banner/` | Supplying your own ASCII art banner                 |
| `example/ascii_only/`    | ASCII-only mode for terminals without Unicode       |
| `example/checks_demo/`   | All built-in check types (SQL, TCP, HTTP, Redis)    |
| `example/custom_checks/` | Function-based, boolean, and grouped custom checks  |
| `example/font_preview/`  | Prints the big-font ASCII wordmark for a service name |

```bash
go run ./example/simple/
make run-example
```

## Output Example

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
  Branch      : main
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

## Security Note

The banner only prints **safe, non-secret** information (version, addresses,
PID, etc.). **Do not** pass secrets (passwords, tokens) via `Options.Extra`
or any other field. The caller is responsible for what gets printed.

## Testing

```bash
go test ./...
```

## License

[MIT](LICENSE)
