// Command ascii_only demonstrates the ASCIIOnly option that replaces Unicode
// box-drawing characters with plain ASCII. This is useful for terminals or log
// systems that do not support Unicode.
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
		ServiceName: "ascii-service",
		Environment: "ci",
		ASCIIOnly:   true, // no Unicode box-drawing characters
		Extra: map[string]string{
			"HTTP": ":8080",
		},
	}
	info := banner.CurrentBuildInfo()

	// A simple environment check to show how checks render in ASCII mode.
	homeCheck := checks.New("env-HOME", func(_ context.Context) error {
		if os.Getenv("HOME") == "" && os.Getenv("USERPROFILE") == "" {
			return fmt.Errorf("neither HOME nor USERPROFILE is set")
		}
		return nil
	})

	runner := checks.DefaultRunner()
	results := runner.Run(context.Background(), homeCheck)

	fmt.Print(banner.RenderWithChecks(opts, info, results))
}
