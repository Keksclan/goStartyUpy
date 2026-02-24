// Command example demonstrates how to use goStartyUpy to print a startup
// banner and run optional health checks.
//
// By default the banner is auto-generated from ServiceName. Set the
// environment variable CUSTOM_BANNER=1 to see an explicitly provided banner
// instead.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/keksclan/goStartyUpy/banner"
	"github.com/keksclan/goStartyUpy/checks"
)

func main() {
	// Build options – customise per service.
	// HTTP and gRPC addresses are passed via Extra (not dedicated fields).
	opts := banner.Options{
		ServiceName: "my-awesome-service",
		Environment: "production",
		Color:       true, // enable ANSI colors for terminal output
		Extra: map[string]string{
			"HTTP":   ":8080",
			"gRPC":   ":9090",
			"Region": "eu-central-1",
		},
	}

	// When CUSTOM_BANNER=1 is set, use an explicitly provided banner string.
	// Otherwise the banner is auto-generated from ServiceName.
	if os.Getenv("CUSTOM_BANNER") == "1" {
		opts.Banner = `
   ╔═══════════════════════════════════╗
   ║     ★  MY AWESOME SERVICE  ★     ║
   ╚═══════════════════════════════════╝`
	}

	info := banner.CurrentBuildInfo()

	// --- Setup checks (examples; adjust to your infrastructure) ---

	// SQL database ping – open your own *sql.DB first.
	// Replace the placeholder DSN with your actual connection string.
	// db, err := sql.Open("pgx", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")
	var db *sql.DB // placeholder – will gracefully fail as nil

	checkList := []checks.Check{
		checks.SQLPingCheck{DB: db, NameLabel: "postgres"},
		checks.TCPDialCheck{Address: "localhost:6379", Label: "redis-tcp"},
		checks.HTTPGetCheck{URL: "http://localhost:8080/healthz", Label: "self-http"},
		checks.RedisPingCheck{Address: "localhost:6379", Label: "redis-ping"},
	}

	// Run all checks in parallel with a 2-second per-check timeout.
	runner := checks.Runner{
		TimeoutPerCheck: 2 * time.Second,
		Parallel:        true,
	}

	ctx := context.Background()
	results := runner.Run(ctx, checkList...)

	// Render and print the full startup message.
	// The output includes "[OK]" per successful check and ends with
	// "Startup Complete" when all checks pass, or "Startup Failed" otherwise.
	fmt.Print(banner.RenderWithChecks(opts, info, results))
}
