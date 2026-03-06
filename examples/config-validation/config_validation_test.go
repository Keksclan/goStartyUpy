package main

import (
	"slices"
	"testing"

	goconfy "github.com/keksclan/goConfy"
	"github.com/keksclan/goStartyUpy/configcheck"
)

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

	if err := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	}); err != nil {
		t.Fatalf("startup check should pass for valid config:\n%s", err)
	}
}

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

	checkErr := configcheck.RunStartupCheck(configcheck.Options{
		Enabled: true,
		Config:  cfg,
	})
	if checkErr == nil {
		t.Fatal("expected startup check to fail for invalid config")
	}

	msg := checkErr.Error()
	for _, key := range []string{"database.port", "database.password", "redis.address"} {
		if !containsSubstring(msg, key) {
			t.Errorf("error message should mention %q, got:\n%s", key, msg)
		}
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := range len(s) - len(sub) + 1 {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
