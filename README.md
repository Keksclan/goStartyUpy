# goStartyUpy

A zero-dependency Go library that generates a production-ready startup banner with build metadata, runtime info, and optional health checks.

## Features

- **Startup banner** with ASCII art, build metadata, and runtime details.
- **Auto-generated banner** from `ServiceName` when no custom banner is set.
- **Startup checks** – verify databases, TCP endpoints, HTTP services, and Redis before accepting traffic.
- **No external dependencies** – standard library only.
- **Never panics** – all errors are captured and returned as structured results.
- **Deterministic & testable** – stable output ordering, no randomness.
- **Optional ANSI colors** – set `Color: true` for colorized terminal output; plain text by default.
- **ASCII-only mode** – set `ASCIIOnly: true` to avoid Unicode box-drawing characters.

## Quickstart

```go
package main

import (
    "fmt"

    "github.com/keksclan/goStartyUpy/banner"
)

func main() {
    opts := banner.Options{
        ServiceName: "my-service",
        Environment: "production",
        Color:       true, // enable ANSI colors (set false for plain text)
        Extra: map[string]string{
            "HTTP": ":8080",
            "gRPC": ":9090",
        },
    }
    info := banner.CurrentBuildInfo()
    fmt.Print(banner.Render(opts, info))
}
```

## Banner Customization

### Auto-generated banner (default)

When `Options.Banner` is empty the library generates a box banner from
`Options.ServiceName` automatically:

```
┌───────────────────────────┐
│        MY-SERVICE         │
└───────────────────────────┘
```

Set `Options.ASCIIOnly = true` to use plain ASCII characters instead of
Unicode box-drawing:

```
+---------------------------+
|        MY-SERVICE         |
+---------------------------+
```

You can also call `banner.DefaultBanner(name, asciiOnly)` directly if you need
the generated string elsewhere.

### Custom banner

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

## Build with Makefile

The included `Makefile` automatically collects git metadata and injects it via
`-ldflags`:

```bash
make build-example   # compile example binary with metadata
make run-example     # build and run example
make test            # run all unit tests
make lint            # go vet + gofmt check
make clean           # remove build artifacts
```

For your own service you can use the generic `build` target:

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

### Using `scripts/ldflags.sh` from another repo

The helper script `scripts/ldflags.sh` prints the ldflags string so you can
integrate it into any build system. It is POSIX sh compatible.

```bash
# From the goStartyUpy directory:
LDFLAGS="$(./scripts/ldflags.sh)" go build -ldflags "$LDFLAGS" ./cmd/myservice

# Override the module path if your import path differs:
MODULE=github.com/my/repo ./scripts/ldflags.sh
```

## Build with ldflags (manual)

Inject git metadata at build time:

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
  ./example/
```

## Running Checks

Use the `checks` package to verify dependencies at startup:

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib" // your driver

    "github.com/keksclan/goStartyUpy/banner"
    "github.com/keksclan/goStartyUpy/checks"
)

func main() {
    db, _ := sql.Open("pgx", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")
    defer db.Close()

    opts := banner.Options{
        ServiceName: "order-service",
        Environment: "staging",
        Color:       true,
        Extra: map[string]string{
            "HTTP": ":8080",
        },
    }
    info := banner.CurrentBuildInfo()

    runner := checks.Runner{
        TimeoutPerCheck: 2 * time.Second,
        Parallel:        true,
    }

    results := runner.Run(context.Background(),
        checks.SQLPingCheck{DB: db, NameLabel: "postgres"},
        checks.TCPDialCheck{Address: "localhost:6379", Label: "redis-tcp"},
        checks.HTTPGetCheck{URL: "http://localhost:8080/healthz", Label: "self-http"},
        checks.RedisPingCheck{Address: "localhost:6379", Label: "redis-ping"},
    )

    fmt.Print(banner.RenderWithChecks(opts, info, results))
}
```

### Available Checks

| Check | Description |
|---|---|
| `SQLPingCheck` | Pings a `*sql.DB` via `PingContext`. |
| `TCPDialCheck` | Dials a TCP `host:port` and closes the connection. |
| `HTTPGetCheck` | Sends an HTTP GET and checks the status code range. |
| `RedisPingCheck` | Sends a RESP `PING` command over TCP (no Redis client needed). |

### Custom Checks

Implement the `checks.Check` interface:

```go
type Check interface {
    Name() string
    Run(ctx context.Context) checks.Result
}
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

The banner only prints **safe, non-secret** information (version, addresses, PID, etc.). **Do not** pass secrets (passwords, tokens) via `Options.Extra` or any other field. The caller is responsible for what gets printed.

## Testing

```bash
go test ./...
```

## License

[MIT](LICENSE)
