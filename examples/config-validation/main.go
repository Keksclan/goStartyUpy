// Command config-validation demonstrates the configuration validation
// feature provided by goStartyUpy's configcheck package together with
// goConfy for YAML loading.
//
// It loads a YAML configuration file via goConfy, then runs the
// goStartyUpy config validation to check for missing required fields.
// When the configuration is valid, the startup banner is printed.
//
// Run from the repository root:
//
//	go run ./examples/config-validation/
package main

import (
	"fmt"
	"log"

	goconfy "github.com/keksclan/goConfy"
	"github.com/keksclan/goStartyUpy/banner"
	"github.com/keksclan/goStartyUpy/configcheck"
)

// AppConfig defines the application configuration structure.
// Fields are required by default. Use `required:"false"` to mark
// a field as optional.
type AppConfig struct {
	AppName  string   `yaml:"app_name"`
	LogLevel string   `yaml:"log_level" required:"false"`
	Database DBConfig `yaml:"database"`
	Redis    struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func main() {
	// Load configuration from YAML using goConfy.
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("examples/config-validation/config.yml"),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Run goStartyUpy config validation.
	configcheck.MustPassStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}, log.Fatalf)

	// Print the startup banner (only reached when config is valid).
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(banner.Options{
		ServiceName: cfg.AppName,
		Color:       true,
	}, info))

	fmt.Printf("Database: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("Redis:    %s\n", cfg.Redis.Address)
}
