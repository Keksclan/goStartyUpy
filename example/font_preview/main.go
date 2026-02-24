// Command font_preview prints the big ASCII-art wordmark for a service name.
//
// Usage:
//
//	go run ./example/font_preview
//	SERVICE_NAME="my-app" go run ./example/font_preview
package main

import (
	"fmt"
	"os"

	"github.com/keksclan/goStartyUpy/banner"
)

func main() {
	name := os.Getenv("SERVICE_NAME")
	if name == "" {
		name = "goStartyUpy"
	}
	fmt.Println(banner.SpringLikeBanner(name, false))
}
