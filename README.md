# goStartyUpy

A zero-dependency Go library that generates a production-ready startup banner with build metadata, runtime info, and optional health checks.

## Features

- **Startup banner** with ASCII art, build metadata, and runtime details.
- **Startup checks** – verify databases, TCP endpoints, HTTP services, and Redis before accepting traffic.
- **No external dependencies** – standard library only.
- **Never panics** – all errors are captured and returned as structured results.
- **Deterministic & testable** – stable output ordering, no randomness.

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
        Extra: map[string]string{
            "HTTP": ":8080",
            "gRPC": ":9090",
        },
    }
    info := banner.CurrentBuildInfo()
    fmt.Print(banner.Render(opts, info))
}
```

## Build with ldflags

Inject git metadata at build time:

```bash
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
DIRTY=$(test -z "$(git status --porcelain)" && echo "false" || echo "true")

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
   _____ _                 _         _    _
  / ____| |               | |       | |  | |
 | (___ | |_ __ _ _ __ | |_ _   _| |  | |_ __  _   _
  \___ \| __/ _' | '__| __| | | | |  | | '_ \| | | |
  ____) | || (_| | |  | |_| |_| | |__| | |_) | |_| |
 |_____/ \__\__,_|_|   \__|\__, |\____/| .__/ \__, |
                             __/ |      | |     __/ |
                            |___/       |_|    |___/
════════════════════════════════════════════════════════════
  Service     : order-service
  Environment : staging
  Version     : v1.2.3
  BuildTime   : 2026-02-24T09:00:00Z
  Commit      : abcdef1234567890
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
