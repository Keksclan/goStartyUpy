// Command basic_start demonstrates the simplest possible goStartyUpy usage.
//
// It renders a startup banner with the default Spring style, showing build
// metadata and runtime information. No health checks are configured.
//
// Run:
//
//	go run ./example/basic_start/
package main

import (
	"fmt"

	"github.com/keksclan/goStartyUpy/banner"
)

func main() {
	opts := banner.Options{
		ServiceName: "basic-service",
	}
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(opts, info))
}
