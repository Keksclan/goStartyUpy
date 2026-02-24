package checks

import (
	"context"
	"database/sql"
	"time"
)

// SQLPingCheck verifies connectivity to a SQL database by calling PingContext.
type SQLPingCheck struct {
	// DB is the *sql.DB handle to ping. Must not be nil.
	DB *sql.DB
	// NameLabel is the human-readable name for this check (e.g. "postgres").
	NameLabel string
}

// Name returns the check label.
func (c SQLPingCheck) Name() string { return c.NameLabel }

// Run pings the database. If DB is nil the check fails gracefully.
func (c SQLPingCheck) Run(ctx context.Context) Result {
	start := time.Now()

	if c.DB == nil {
		return Result{
			Name:     c.NameLabel,
			OK:       false,
			Duration: time.Since(start),
			Error:    "sql.DB is nil",
		}
	}

	if err := c.DB.PingContext(ctx); err != nil {
		return Result{
			Name:     c.NameLabel,
			OK:       false,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}

	return Result{
		Name:     c.NameLabel,
		OK:       true,
		Duration: time.Since(start),
	}
}
