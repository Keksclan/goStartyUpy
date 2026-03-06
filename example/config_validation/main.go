// Command config_validation demonstrates the configuration validation
// feature provided by the configcheck package.
//
// It shows how to validate a configuration struct (typically loaded via
// goConfy) before printing the startup banner. When required fields are
// missing, the program outputs a clear diagnostic and exits.
//
// Run with:
//
//	go run ./example/config_validation/
package main

import (
	"fmt"
	"log"

	"github.com/keksclan/goStartyUpy/banner"
	"github.com/keksclan/goStartyUpy/configcheck"
)

// AppConfig mirrors a typical goConfy configuration struct. Fields tagged
// with `required:"false"` are optional; all others must be populated.
type AppConfig struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Redis struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
	AppName  string `yaml:"app_name"`
	LogLevel string `yaml:"log_level" required:"false"`
}

func main() {
	// In a real application this would come from goConfy:
	//   cfg, err := goconfy.Load[AppConfig](goconfy.WithFile("config.yml"))
	//
	// Here we simulate a partially loaded config with missing fields.
	cfg := AppConfig{
		AppName: "demo-service",
	}
	cfg.Database.Host = "localhost"
	// Database.Port, Database.Password, and Redis.Address are intentionally
	// left empty to trigger validation errors.

	// 1. Run config validation (before printing the banner).
	// MustPassStartupCheck prints a green ✔ on success and calls
	// log.Fatalf with a diagnostic when required fields are missing.
	configcheck.MustPassStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}, log.Fatalf)

	// 2. Print the startup banner (only reached when config is valid).
	opts := banner.Options{
		ServiceName: cfg.AppName,
		Color:       true,
	}
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(opts, info))
}
