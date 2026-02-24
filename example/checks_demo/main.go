// Command checks_demo showcases every built-in check type that ships with
// goStartyUpy: SQLPingCheck, TCPDialCheck, HTTPGetCheck, and RedisPingCheck.
//
// The checks will fail in most environments because the target services are
// unlikely to be running – this is intentional and demonstrates how failures
// are rendered in the startup banner.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/keksclan/goStartyUpy/banner"
	"github.com/keksclan/goStartyUpy/checks"
)

func main() {
	opts := banner.Options{
		ServiceName: "checks-demo",
		Environment: "local",
		Color:       true,
		Extra: map[string]string{
			"HTTP": ":8080",
		},
	}
	info := banner.CurrentBuildInfo()

	// SQL database ping – nil DB demonstrates graceful failure.
	var db *sql.DB

	builtinChecks := []checks.Check{
		// 1. SQL ping – passes when the database is reachable.
		checks.SQLPingCheck{DB: db, NameLabel: "postgres"},

		// 2. TCP dial – verifies a raw TCP connection can be established.
		checks.TCPDialCheck{Address: "localhost:5432", Label: "postgres-tcp"},

		// 3. HTTP GET – sends a GET and checks the status code range.
		checks.HTTPGetCheck{
			URL:               "http://localhost:8080/healthz",
			Label:             "api-health",
			ExpectedStatusMin: 200,
			ExpectedStatusMax: 299,
		},

		// 4. Redis PING – speaks the RESP protocol, no client library needed.
		checks.RedisPingCheck{Address: "localhost:6379", Label: "redis"},
	}

	// Run with an explicit runner to show manual configuration.
	runner := checks.Runner{
		TimeoutPerCheck: 3 * time.Second,
		Parallel:        true,
	}
	results := runner.Run(context.Background(), builtinChecks...)

	fmt.Print(banner.RenderWithChecks(opts, info, results))
}
