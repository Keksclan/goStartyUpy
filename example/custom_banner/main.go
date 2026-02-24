// Command custom_banner demonstrates how to supply your own ASCII art banner
// instead of using the auto-generated one. When Options.Banner is set, the
// library skips auto-generation entirely.
package main

import (
	"fmt"

	"github.com/keksclan/goStartyUpy/banner"
)

func main() {
	opts := banner.Options{
		ServiceName: "payment-service",
		Environment: "production",
		Color:       true,
		Banner: `
   ╔═══════════════════════════════════╗
   ║    💳  PAYMENT SERVICE  💳        ║
   ╚═══════════════════════════════════╝`,
		Extra: map[string]string{
			"HTTP":    ":443",
			"Webhook": "/api/v1/webhook",
		},
	}
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(opts, info))
}
