package banner

import (
	"strings"
	"testing"
)

// --- BoxBanner tests (legacy box style) ---

func TestBoxBanner_UsesServiceName(t *testing.T) {
	out := BoxBanner("my-service", false)
	if !strings.Contains(out, "MY-SERVICE") {
		t.Errorf("expected banner to contain upper-cased service name, got:\n%s", out)
	}
}

func TestBoxBanner_EmptyNameFallback(t *testing.T) {
	out := BoxBanner("", false)
	if !strings.Contains(out, "SERVICE") {
		t.Errorf("expected banner to contain SERVICE for empty name, got:\n%s", out)
	}
}

func TestBoxBanner_UnicodeBox(t *testing.T) {
	out := BoxBanner("test", false)
	if !strings.Contains(out, "┌") || !strings.Contains(out, "┘") {
		t.Errorf("expected Unicode box-drawing chars, got:\n%s", out)
	}
	if strings.Contains(out, "+") {
		t.Error("Unicode mode should not contain '+' characters")
	}
}

func TestBoxBanner_ASCIIOnly(t *testing.T) {
	out := BoxBanner("test", true)
	if !strings.Contains(out, "+") || !strings.Contains(out, "-") || !strings.Contains(out, "|") {
		t.Errorf("expected ASCII box chars, got:\n%s", out)
	}
	if strings.Contains(out, "┌") || strings.Contains(out, "─") {
		t.Error("ASCII mode should not contain Unicode box-drawing chars")
	}
}

func TestBoxBanner_MinWidth(t *testing.T) {
	out := BoxBanner("x", false)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	topRunes := []rune(lines[0])
	if len(topRunes) < 29 {
		t.Errorf("top line too short (%d runes): %s", len(topRunes), lines[0])
	}
}

func TestBoxBanner_Deterministic(t *testing.T) {
	a := BoxBanner("svc", false)
	b := BoxBanner("svc", false)
	if a != b {
		t.Error("BoxBanner should be deterministic")
	}
}

// --- DefaultBanner tests (now defaults to spring style) ---

func TestDefaultBanner_UsesSpringStyle(t *testing.T) {
	out := DefaultBanner("my-svc", false)
	// Spring-style banner must be multi-line (wordmark) and contain tagline.
	if !strings.Contains(out, ":: goStartyUpy ::") {
		t.Errorf("expected spring-style tagline, got:\n%s", out)
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 5 {
		t.Errorf("expected multi-line wordmark (>= 5 lines), got %d lines", len(lines))
	}
}

func TestDefaultBanner_EmptyNameFallback(t *testing.T) {
	out := DefaultBanner("", false)
	// Should render "SERVICE" wordmark — check the tagline is present.
	if !strings.Contains(out, ":: goStartyUpy ::") {
		t.Errorf("expected spring-style tagline for empty name, got:\n%s", out)
	}
}

// --- Render integration tests ---

func TestRender_UsesSpringBannerByDefault(t *testing.T) {
	opts := Options{
		ServiceName: "auto-svc",
	}
	out := Render(opts, BuildInfo{})
	if !strings.Contains(out, ":: goStartyUpy ::") {
		t.Errorf("expected spring-style tagline in default render, got:\n%s", out)
	}
}

func TestRender_BoxStyleWorks(t *testing.T) {
	opts := Options{
		ServiceName: "auto-svc",
		BannerStyle: "box",
	}
	out := Render(opts, BuildInfo{})
	if !strings.Contains(out, "AUTO-SVC") {
		t.Errorf("expected box banner with service name, got:\n%s", out)
	}
}

func TestRender_ASCIIOnlySwitchesSeparator(t *testing.T) {
	opts := Options{
		ServiceName: "test",
		ASCIIOnly:   true,
	}
	out := Render(opts, BuildInfo{})
	if strings.Contains(out, "═") {
		t.Error("ASCIIOnly should not use Unicode separator")
	}
	if !strings.Contains(out, "====") {
		t.Error("ASCIIOnly should use ASCII separator")
	}
}

func TestRender_ASCIIOnlyBoxBannerChars(t *testing.T) {
	opts := Options{
		ServiceName: "test",
		ASCIIOnly:   true,
		BannerStyle: "box",
	}
	out := Render(opts, BuildInfo{})
	if strings.Contains(out, "┌") || strings.Contains(out, "─") {
		t.Error("ASCIIOnly should not use Unicode box-drawing in banner")
	}
	if !strings.Contains(out, "+") {
		t.Error("ASCIIOnly banner should use '+' corners")
	}
}
