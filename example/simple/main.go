// Command simple demonstrates the most minimal goStartyUpy usage: just a
// startup banner with build metadata, no health checks.
package main

import (
	"fmt"

	"github.com/keksclan/goStartyUpy/banner"
)

func main() {
	opts := banner.Options{
		ServiceName: "hello-service",
		Environment: "development",
	}
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(opts, info))
}
