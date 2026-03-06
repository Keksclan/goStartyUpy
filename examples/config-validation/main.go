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

// BannerConfig holds all user-facing banner rendering options.
type BannerConfig struct {
	Color        bool              `yaml:"color"         required:"false"`
	ASCIIOnly    bool              `yaml:"ascii_only"    required:"false"`
	BannerStyle  string            `yaml:"banner_style"  required:"false"`
	BannerWidth  int               `yaml:"banner_width"  required:"false"`
	CustomBanner string            `yaml:"custom_banner" required:"false"`
	Separator    string            `yaml:"separator"     required:"false"`
	Tagline1     string            `yaml:"tagline1"      required:"false"`
	Tagline2     string            `yaml:"tagline2"      required:"false"`
	ShowDetails  *bool             `yaml:"show_details"  required:"false"`
	Extra        map[string]string `yaml:"extra"         required:"false"`
}

// ConfigCheckConfig controls startup configuration validation.
type ConfigCheckConfig struct {
	Enabled bool `yaml:"enabled" required:"false"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

// AppConfig defines the full application configuration structure.
// Fields are required by default. Use `required:"false"` to mark
// a field as optional.
type AppConfig struct {
	AppName     string            `yaml:"app_name"`
	Environment string            `yaml:"environment"   required:"false"`
	LogLevel    string            `yaml:"log_level"     required:"false"`
	Banner      BannerConfig      `yaml:"banner"        required:"false"`
	ConfigCheck ConfigCheckConfig `yaml:"config_check"  required:"false"`
	Database    DBConfig          `yaml:"database"`
	Redis       struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
}

func main() {
	// Load configuration from YAML using goConfy.
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("examples/config-validation/config.full.example.yml"),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Run goStartyUpy config validation.
	configcheck.MustPassStartupCheck(configcheck.Options{
		Enabled: cfg.ConfigCheck.Enabled,
		Config:  cfg,
	}, log.Fatalf)

	// Print the startup banner (only reached when config is valid).
	info := banner.CurrentBuildInfo()
	fmt.Print(banner.Render(banner.Options{
		ServiceName: cfg.AppName,
		Environment: cfg.Environment,
		Color:       cfg.Banner.Color,
		ASCIIOnly:   cfg.Banner.ASCIIOnly,
		BannerStyle: cfg.Banner.BannerStyle,
		BannerWidth: cfg.Banner.BannerWidth,
		Banner:      cfg.Banner.CustomBanner,
		Separator:   cfg.Banner.Separator,
		Tagline1:    cfg.Banner.Tagline1,
		Tagline2:    cfg.Banner.Tagline2,
		ShowDetails: cfg.Banner.ShowDetails,
		Extra:       cfg.Banner.Extra,
	}, info))

	fmt.Printf("Database: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("Redis:    %s\n", cfg.Redis.Address)
}
