package banner

import (
	"strings"
	"testing"
	"time"

	"github.com/keksclan/goStartyUpy/checks"
)

func TestRender_ContainsVersionBuildCommit(t *testing.T) {
	info := BuildInfo{
		Version:   "v1.0.0",
		BuildTime: "2026-01-01T00:00:00Z",
		Commit:    "abc1234",
		Branch:    "main",
		Dirty:     "false",
	}
	opts := Options{
		ServiceName: "test-svc",
		Environment: "test",
		Extra: map[string]string{
			"HTTP": ":8080",
			"gRPC": ":9090",
		},
	}

	out := Render(opts, info)

	for _, want := range []string{"v1.0.0", "2026-01-01T00:00:00Z", "abc1234", "main"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestRender_EndsWithNewline(t *testing.T) {
	out := Render(Options{}, BuildInfo{})
	if !strings.HasSuffix(out, "\n") {
		t.Error("output should end with newline")
	}
}

func TestRender_NilExtra_NoPanic(t *testing.T) {
	// Must not panic when Extra is nil.
	out := Render(Options{Extra: nil}, BuildInfo{})
	if out == "" {
		t.Error("output should not be empty")
	}
}

func TestRender_CustomBannerAndSeparator(t *testing.T) {
	opts := Options{
		Banner:    "=== MY BANNER ===",
		Separator: "---",
	}
	out := Render(opts, BuildInfo{})
	if !strings.Contains(out, "=== MY BANNER ===") {
		t.Error("custom banner not found")
	}
	if !strings.Contains(out, "---") {
		t.Error("custom separator not found")
	}
}

func TestRenderWithChecks_OKAndFail(t *testing.T) {
	results := []checks.Result{
		{Name: "postgres", OK: true, Duration: 12 * time.Millisecond},
		{Name: "redis", OK: false, Duration: 2 * time.Second, Error: "dial tcp: timeout"},
	}

	out := RenderWithChecks(Options{}, BuildInfo{}, results)

	if !strings.Contains(out, "Checks:") {
		t.Error("missing Checks: header")
	}
	if !strings.Contains(out, "[OK]   postgres") {
		t.Error("missing OK check line")
	}
	if !strings.Contains(out, "[FAIL] redis") {
		t.Error("missing FAIL check line")
	}
	if !strings.Contains(out, "dial tcp: timeout") {
		t.Error("missing error detail in FAIL line")
	}
	if !strings.Contains(out, "Startup Failed") {
		t.Error("missing Startup Failed message when a check fails")
	}
}

func TestRenderWithChecks_AllOK_StartupComplete(t *testing.T) {
	results := []checks.Result{
		{Name: "postgres", OK: true, Duration: 5 * time.Millisecond},
		{Name: "redis", OK: true, Duration: 3 * time.Millisecond},
	}

	out := RenderWithChecks(Options{}, BuildInfo{}, results)

	if !strings.Contains(out, "Startup Complete") {
		t.Error("missing Startup Complete when all checks pass")
	}
}

func TestRenderWithChecks_NoChecks(t *testing.T) {
	out := RenderWithChecks(Options{}, BuildInfo{}, nil)
	if strings.Contains(out, "Checks:") {
		t.Error("should not contain Checks: section when no results")
	}
}

func TestRender_ColorEnabled(t *testing.T) {
	opts := Options{
		ServiceName: "color-svc",
		Environment: "test",
		Color:       true,
	}
	info := BuildInfo{Version: "v1.0.0", Commit: "abc123"}
	out := Render(opts, info)

	// Output must contain ANSI escape sequences when Color is true.
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in colored output")
	}
	// Content must still be present.
	if !strings.Contains(out, "color-svc") {
		t.Error("service name missing in colored output")
	}
	if !strings.Contains(out, "v1.0.0") {
		t.Error("version missing in colored output")
	}
}

func TestRender_ColorDisabled_NoEscapes(t *testing.T) {
	opts := Options{
		ServiceName: "plain-svc",
		Color:       false,
	}
	out := Render(opts, BuildInfo{})

	if strings.Contains(out, "\033[") {
		t.Error("plain output must not contain ANSI escape codes")
	}
}

func TestRenderWithChecks_ColorChecks(t *testing.T) {
	results := []checks.Result{
		{Name: "db", OK: true, Duration: 5 * time.Millisecond},
		{Name: "cache", OK: false, Duration: 1 * time.Second, Error: "refused"},
	}
	out := RenderWithChecks(Options{Color: true}, BuildInfo{}, results)

	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI codes in colored checks output")
	}
	if !strings.Contains(out, "db") || !strings.Contains(out, "cache") {
		t.Error("check names missing in colored output")
	}
	if !strings.Contains(out, "Startup Failed") {
		t.Error("Startup Failed missing in colored output")
	}
}
