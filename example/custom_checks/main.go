// Command custom_checks shows how to create custom startup checks without
// implementing the Check interface by hand. It demonstrates checks.New,
// checks.Bool, and checks.NewGroup.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/keksclan/goStartyUpy/banner"
	"github.com/keksclan/goStartyUpy/checks"
)

func main() {
	opts := banner.Options{
		ServiceName: "custom-checks-demo",
		Environment: "staging",
		Color:       true,
		Extra: map[string]string{
			"HTTP": ":8080",
		},
	}
	info := banner.CurrentBuildInfo()

	// --- Function-based checks ---

	// Validate that a required environment variable is present.
	envCheck := checks.New("env-APP_SECRET", func(_ context.Context) error {
		if os.Getenv("APP_SECRET") == "" {
			return fmt.Errorf("APP_SECRET is not set")
		}
		return nil
	})

	// Check that PORT is a valid number.
	portCheck := checks.New("env-PORT", func(_ context.Context) error {
		raw := os.Getenv("PORT")
		if raw == "" {
			return fmt.Errorf("PORT is not set")
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("PORT is not a number: %s", raw)
		}
		if n < 1 || n > 65535 {
			return fmt.Errorf("PORT out of range: %d", n)
		}
		return nil
	})

	// --- Boolean checks ---

	// Feature flag: the check passes when the flag is "true".
	featureFlag := checks.Bool("feature-new-ui", func(_ context.Context) (bool, error) {
		return os.Getenv("ENABLE_NEW_UI") == "true", nil
	})

	// Debug mode: reports whether debug logging is active.
	debugMode := checks.Bool("debug-mode", func(_ context.Context) (bool, error) {
		return os.Getenv("DEBUG") == "1", nil
	})

	// --- Grouped check ---

	// Bundle the environment checks into a single composite check.
	envGroup := checks.NewGroup("env-vars", checks.GroupOptions{},
		envCheck,
		portCheck,
	)

	// Run everything with the default runner (parallel, 2 s timeout).
	runner := checks.DefaultRunner()
	results := runner.Run(context.Background(),
		envGroup,
		featureFlag,
		debugMode,
	)

	fmt.Print(banner.RenderWithChecks(opts, info, results))
}
