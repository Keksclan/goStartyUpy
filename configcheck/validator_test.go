package configcheck

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

// --- test helpers ----------------------------------------------------------

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type redisConfig struct {
	Address string `yaml:"address"`
}

type fullConfig struct {
	Database dbConfig    `yaml:"database"`
	Redis    redisConfig `yaml:"redis"`
	AppName  string      `yaml:"app_name"`
	Debug    bool        `yaml:"debug" required:"false"`
}

type optionalFieldsConfig struct {
	Name     string `yaml:"name"`
	Nickname string `yaml:"nickname" required:"false"`
	Age      int    `yaml:"age" required:"false"`
}

type configWithSlice struct {
	Hosts []string `yaml:"hosts"`
	Tags  []string `yaml:"tags" required:"false"`
}

type configWithMap struct {
	Labels map[string]string `yaml:"labels"`
	Extra  map[string]string `yaml:"extra" required:"false"`
}

type configWithPointer struct {
	DB *dbConfig `yaml:"db"`
}

type configWithTime struct {
	CreatedAt time.Time `yaml:"created_at"`
	Name      string    `yaml:"name"`
}

type configYAMLDash struct {
	Internal string `yaml:"-"`
	Name     string `yaml:"name"`
}

type noTagConfig struct {
	FieldOne string
	FieldTwo int
}

// --- Validate tests --------------------------------------------------------

func TestValidate_AllFieldsPresent(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "localhost", Port: 5432, Password: "secret"},
		Redis:    redisConfig{Address: "localhost:6379"},
		AppName:  "my-app",
		Debug:    false, // optional — zero value is fine
	}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_MissingRequiredFields(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "localhost"},
		// Redis and AppName are missing
	}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error, got nil")
	}

	// database.port is 0 (zero value for int) → missing
	// database.password is "" → missing
	// redis.address is "" → missing
	// app_name is "" → missing
	// debug is optional → NOT missing
	expected := []string{"database.port", "database.password", "redis.address", "app_name"}
	for _, key := range expected {
		if !slices.Contains(ve.Missing, key) {
			t.Errorf("expected missing key %q, got: %v", key, ve.Missing)
		}
	}
	if slices.Contains(ve.Missing, "debug") {
		t.Error("debug is optional and should not be reported as missing")
	}
}

func TestValidate_OptionalFieldsNotReported(t *testing.T) {
	cfg := optionalFieldsConfig{
		Name: "Alice",
		// Nickname and Age are optional and zero
	}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_NilPointer(t *testing.T) {
	ve := Validate((*fullConfig)(nil))
	if ve == nil {
		t.Fatal("expected error for nil pointer")
	}
	if len(ve.Errors) == 0 {
		t.Fatal("expected error message for nil config")
	}
}

func TestValidate_NonStruct(t *testing.T) {
	ve := Validate("not a struct")
	if ve == nil {
		t.Fatal("expected error for non-struct")
	}
	if len(ve.Errors) == 0 || ve.Errors[0] != "expected struct, got string" {
		t.Fatalf("unexpected error: %v", ve.Errors)
	}
}

func TestValidate_PointerToStruct(t *testing.T) {
	cfg := &fullConfig{
		Database: dbConfig{Host: "h", Port: 1, Password: "p"},
		Redis:    redisConfig{Address: "a"},
		AppName:  "x",
	}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_EmptySlice(t *testing.T) {
	cfg := configWithSlice{
		Hosts: nil,
	}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error for nil slice")
	}
	if !slices.Contains(ve.Missing, "hosts") {
		t.Errorf("expected missing key 'hosts', got: %v", ve.Missing)
	}
	if slices.Contains(ve.Missing, "tags") {
		t.Error("tags is optional and should not be reported")
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	cfg := configWithMap{}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error for nil map")
	}
	if !slices.Contains(ve.Missing, "labels") {
		t.Errorf("expected missing key 'labels', got: %v", ve.Missing)
	}
	if slices.Contains(ve.Missing, "extra") {
		t.Error("extra is optional and should not be reported")
	}
}

func TestValidate_NilPointerField(t *testing.T) {
	cfg := configWithPointer{DB: nil}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error for nil pointer field")
	}
	if !slices.Contains(ve.Missing, "db") {
		t.Errorf("expected missing key 'db', got: %v", ve.Missing)
	}
}

func TestValidate_PopulatedPointerField(t *testing.T) {
	cfg := configWithPointer{DB: &dbConfig{Host: "h", Port: 1, Password: "p"}}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_LeafStructZeroTime(t *testing.T) {
	cfg := configWithTime{Name: "test"}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error for zero time.Time")
	}
	if !slices.Contains(ve.Missing, "created_at") {
		t.Errorf("expected missing key 'created_at', got: %v", ve.Missing)
	}
}

func TestValidate_LeafStructNonZeroTime(t *testing.T) {
	cfg := configWithTime{Name: "test", CreatedAt: time.Now()}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_YAMLDashSkipped(t *testing.T) {
	cfg := configYAMLDash{Name: "hello"}
	if ve := Validate(cfg); ve != nil {
		t.Fatalf("expected no error, got: %v", ve)
	}
}

func TestValidate_NoTagsFallbackToFieldName(t *testing.T) {
	cfg := noTagConfig{}
	ve := Validate(cfg)
	if ve == nil {
		t.Fatal("expected validation error")
	}
	if !slices.Contains(ve.Missing, "FieldOne") {
		t.Errorf("expected missing key 'FieldOne', got: %v", ve.Missing)
	}
	if !slices.Contains(ve.Missing, "FieldTwo") {
		t.Errorf("expected missing key 'FieldTwo', got: %v", ve.Missing)
	}
}

// --- ValidationError tests -------------------------------------------------

func TestValidationError_Error(t *testing.T) {
	ve := &ValidationError{
		Missing: []string{"database.host", "redis.address"},
	}
	s := ve.Error()
	if s == "" {
		t.Fatal("expected non-empty error string")
	}
	for _, key := range ve.Missing {
		if !strings.Contains(s, key) {
			t.Errorf("error string should contain %q", key)
		}
	}
}

func TestValidationError_ErrorEmpty(t *testing.T) {
	ve := &ValidationError{}
	if ve.Error() != "" {
		t.Fatal("expected empty error string for no errors")
	}
}

func TestValidationError_HasErrors(t *testing.T) {
	if (&ValidationError{}).HasErrors() {
		t.Error("empty ValidationError should not have errors")
	}
	if !(&ValidationError{Missing: []string{"x"}}).HasErrors() {
		t.Error("ValidationError with missing keys should have errors")
	}
	if !(&ValidationError{Errors: []string{"y"}}).HasErrors() {
		t.Error("ValidationError with errors should have errors")
	}
}

// --- RunStartupCheck tests -------------------------------------------------

func TestRunStartupCheck_Disabled(t *testing.T) {
	msg, err := RunStartupCheck(Options{Enabled: false, Config: "not even a struct"})
	if err != nil {
		t.Fatalf("expected nil when disabled, got: %v", err)
	}
	if msg != "" {
		t.Fatalf("expected empty message when disabled, got: %q", msg)
	}
}

func TestRunStartupCheck_NilConfig(t *testing.T) {
	msg, err := RunStartupCheck(Options{Enabled: true, Config: nil})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
	if msg != "" {
		t.Fatalf("expected empty message on error, got: %q", msg)
	}
}

func TestRunStartupCheck_Valid(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "h", Port: 1, Password: "p"},
		Redis:    redisConfig{Address: "a"},
		AppName:  "x",
	}
	msg, err := RunStartupCheck(Options{Enabled: true, Config: cfg})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if !strings.Contains(msg, "Config Check") {
		t.Fatalf("expected success message containing 'Config Check', got: %q", msg)
	}
}

func TestRunStartupCheck_ValidPlain(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "h", Port: 1, Password: "p"},
		Redis:    redisConfig{Address: "a"},
		AppName:  "x",
	}
	msg, err := RunStartupCheck(Options{Enabled: true, Config: cfg})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if !strings.Contains(msg, "[OK]") {
		t.Fatalf("expected plain [OK] tag, got: %q", msg)
	}
}

func TestRunStartupCheck_ValidColor(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "h", Port: 1, Password: "p"},
		Redis:    redisConfig{Address: "a"},
		AppName:  "x",
	}
	msg, err := RunStartupCheck(Options{Enabled: true, Config: cfg, Color: true})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if !strings.Contains(msg, "✔") {
		t.Fatalf("expected colored output with checkmark, got: %q", msg)
	}
	if strings.Contains(msg, "[OK]") {
		t.Fatalf("colored output should not contain [OK] tag, got: %q", msg)
	}
}

func TestRunStartupCheck_Invalid(t *testing.T) {
	cfg := fullConfig{}
	msg, err := RunStartupCheck(Options{Enabled: true, Config: cfg})
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
	if msg != "" {
		t.Fatalf("expected empty message on error, got: %q", msg)
	}
}

// --- MustPassStartupCheck tests --------------------------------------------

func TestMustPassStartupCheck_NoFatal(t *testing.T) {
	cfg := fullConfig{
		Database: dbConfig{Host: "h", Port: 1, Password: "p"},
		Redis:    redisConfig{Address: "a"},
		AppName:  "x",
	}
	called := false
	MustPassStartupCheck(Options{Enabled: true, Config: cfg}, func(string, ...any) {
		called = true
	})
	if called {
		t.Fatal("fatalf should not have been called for valid config")
	}
}

func TestMustPassStartupCheck_Fatal(t *testing.T) {
	cfg := fullConfig{}
	called := false
	MustPassStartupCheck(Options{Enabled: true, Config: cfg}, func(string, ...any) {
		called = true
	})
	if !called {
		t.Fatal("fatalf should have been called for invalid config")
	}
}

// --- FormatValidationError tests -------------------------------------------

func TestFormatValidationError_Nil(t *testing.T) {
	if s := FormatValidationError(nil); s != "" {
		t.Fatalf("expected empty string, got: %q", s)
	}
}

func TestFormatValidationError_ValidationError(t *testing.T) {
	ve := &ValidationError{Missing: []string{"a.b"}}
	s := FormatValidationError(ve)
	if !strings.Contains(s, "a.b") {
		t.Fatalf("expected formatted output to contain 'a.b', got: %q", s)
	}
}

func TestFormatValidationError_OtherError(t *testing.T) {
	s := FormatValidationError(fmt.Errorf("boom"))
	if !strings.Contains(s, "boom") {
		t.Fatalf("expected formatted output to contain 'boom', got: %q", s)
	}
}

// --- Easter Egg tests ------------------------------------------------------

func TestValidationError_ErrorIsDeterministic(t *testing.T) {
	ve := &ValidationError{Missing: []string{"some.key"}}

	const kimPhrase = "Kim mag dich nicht"

	// Error() must never contain the Easter egg message.
	for range 1000 {
		s := ve.Error()
		if strings.Contains(s, kimPhrase) {
			t.Fatal("Error() must not contain the Easter egg message")
		}
		if !strings.Contains(s, "Config validation failed") {
			t.Fatal("error string must always contain 'Config validation failed'")
		}
		if !strings.Contains(s, "some.key") {
			t.Fatal("error string must always contain the missing key")
		}
	}
}

func TestShowEasterEgg(t *testing.T) {
	msg := ShowEasterEgg()
	if !strings.Contains(msg, "Kim mag dich nicht") {
		t.Fatalf("expected Easter egg message, got: %q", msg)
	}
}
