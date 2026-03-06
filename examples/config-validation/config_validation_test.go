package main

import (
	"slices"
	"strings"
	"testing"

	goconfy "github.com/keksclan/goConfy"
	"github.com/keksclan/goStartyUpy/configcheck"
)

// ---------------------------------------------------------------------------
// Minimal config
// ---------------------------------------------------------------------------

func TestMinimalConfig(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.minimal.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load minimal config: %v", err)
	}

	ve := configcheck.Validate(cfg)
	if ve != nil {
		t.Fatalf("expected minimal config to be valid, got errors:\n%s", ve.Error())
	}
}

func TestMinimalConfig_RunStartupCheck(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.minimal.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load minimal config: %v", err)
	}

	if _, err := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}); err != nil {
		t.Fatalf("startup check should pass for minimal config:\n%s", err)
	}
}

// ---------------------------------------------------------------------------
// Full example config
// ---------------------------------------------------------------------------

func TestFullExampleConfig(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.full.example.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load full example config: %v", err)
	}

	ve := configcheck.Validate(cfg)
	if ve != nil {
		t.Fatalf("expected full example config to be valid, got errors:\n%s", ve.Error())
	}
}

func TestFullExampleConfig_RunStartupCheck(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.full.example.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load full example config: %v", err)
	}

	if _, err := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}); err != nil {
		t.Fatalf("startup check should pass for full example config:\n%s", err)
	}
}

func TestFullExampleConfig_AllFieldsPopulated(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.full.example.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load full example config: %v", err)
	}

	// Verify key fields are actually populated (not just zero values).
	if cfg.AppName == "" {
		t.Error("app_name should be set in the full example")
	}
	if cfg.Environment == "" {
		t.Error("environment should be set in the full example")
	}
	if cfg.Database.Host == "" {
		t.Error("database.host should be set in the full example")
	}
	if cfg.Database.Port == 0 {
		t.Error("database.port should be set in the full example")
	}
	if cfg.Database.Password == "" {
		t.Error("database.password should be set in the full example")
	}
	if cfg.Redis.Address == "" {
		t.Error("redis.address should be set in the full example")
	}
	if cfg.Banner.BannerStyle == "" {
		t.Error("banner.banner_style should be set in the full example")
	}
	if cfg.Banner.BannerWidth == 0 {
		t.Error("banner.banner_width should be set in the full example")
	}
	if cfg.Banner.Extra == nil || len(cfg.Banner.Extra) == 0 {
		t.Error("banner.extra should contain at least one entry in the full example")
	}
	if !cfg.ConfigCheck.Enabled {
		t.Error("config_check.enabled should be true in the full example")
	}
}

// ---------------------------------------------------------------------------
// Valid config (original config.yml)
// ---------------------------------------------------------------------------

func TestValidConfig(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load valid config: %v", err)
	}

	ve := configcheck.Validate(cfg)
	if ve != nil {
		t.Fatalf("expected valid config, got validation errors:\n%s", ve.Error())
	}
}

func TestValidConfig_RunStartupCheck(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load valid config: %v", err)
	}

	if _, err := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}); err != nil {
		t.Fatalf("startup check should pass for valid config:\n%s", err)
	}
}

// ---------------------------------------------------------------------------
// Invalid config
// ---------------------------------------------------------------------------

func TestInvalidConfig_MissingFields(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config_invalid.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load invalid config: %v", err)
	}

	ve := configcheck.Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation errors for invalid config, got nil")
	}

	// config_invalid.yml is missing: database.port, database.password,
	// and the entire redis.address.
	expectedMissing := []string{"database.port", "database.password", "redis.address"}
	for _, key := range expectedMissing {
		if !slices.Contains(ve.Missing, key) {
			t.Errorf("expected missing key %q, got missing: %v", key, ve.Missing)
		}
	}

	// app_name is set in the invalid config, so it must not appear.
	if slices.Contains(ve.Missing, "app_name") {
		t.Error("app_name is set and should not be reported as missing")
	}
}

func TestInvalidConfig_ErrorMessageContainsKeys(t *testing.T) {
	cfg, err := goconfy.Load[AppConfig](
		goconfy.WithFile("config_invalid.yml"),
	)
	if err != nil {
		t.Fatalf("failed to load invalid config: %v", err)
	}

	_, checkErr := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	})
	if checkErr == nil {
		t.Fatal("expected startup check to fail for invalid config")
	}

	msg := checkErr.Error()
	for _, key := range []string{"database.port", "database.password", "redis.address"} {
		if !strings.Contains(msg, key) {
			t.Errorf("error message should mention %q, got:\n%s", key, msg)
		}
	}
}
