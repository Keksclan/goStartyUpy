// Command env_aware_start demonstrates automatic environment detection.
//
// goStartyUpy reads the GO_STARTYUPY_ENV environment variable when
// Options.Environment is not set explicitly. This lets operators control the
// displayed environment without code changes.
//
// Run with an explicit environment:
//
//	GO_STARTYUPY_ENV=staging go run ./example/env_aware_start/
//
// Run without the variable (falls back to empty):
//
//	go run ./example/env_aware_start/
package main

import (
	"fmt"

	"github.com/keksclan/goStartyUpy/banner"
)

func main() {
	// Environment is intentionally left empty so that the library reads
	// the GO_STARTYUPY_ENV environment variable at render time.
	opts := banner.Options{
		ServiceName: "env-aware-service",
		Color:       true,
		Extra: map[string]string{
			"HTTP": ":8080",
		},
	}
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(opts, info))
}
