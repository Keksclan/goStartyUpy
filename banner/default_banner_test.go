package banner

import (
	"strings"
	"testing"
)

func TestDefaultBanner_UsesServiceName(t *testing.T) {
	out := DefaultBanner("my-service", false)
	if !strings.Contains(out, "MY-SERVICE") {
		t.Errorf("expected banner to contain upper-cased service name, got:\n%s", out)
	}
}

func TestDefaultBanner_EmptyNameFallback(t *testing.T) {
	out := DefaultBanner("", false)
	if !strings.Contains(out, "SERVICE") {
		t.Errorf("expected banner to contain SERVICE for empty name, got:\n%s", out)
	}
}

func TestDefaultBanner_UnicodeBox(t *testing.T) {
	out := DefaultBanner("test", false)
	if !strings.Contains(out, "┌") || !strings.Contains(out, "┘") {
		t.Errorf("expected Unicode box-drawing chars, got:\n%s", out)
	}
	if strings.Contains(out, "+") {
		t.Error("Unicode mode should not contain '+' characters")
	}
}

func TestDefaultBanner_ASCIIOnly(t *testing.T) {
	out := DefaultBanner("test", true)
	if !strings.Contains(out, "+") || !strings.Contains(out, "-") || !strings.Contains(out, "|") {
		t.Errorf("expected ASCII box chars, got:\n%s", out)
	}
	if strings.Contains(out, "┌") || strings.Contains(out, "─") {
		t.Error("ASCII mode should not contain Unicode box-drawing chars")
	}
}

func TestDefaultBanner_MinWidth(t *testing.T) {
	out := DefaultBanner("x", false)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// Top line includes border chars; inner width must be >= 27.
	// "┌" + 27×"─" + "┐" = 29 runes minimum.
	topRunes := []rune(lines[0])
	if len(topRunes) < 29 {
		t.Errorf("top line too short (%d runes): %s", len(topRunes), lines[0])
	}
}

func TestDefaultBanner_Deterministic(t *testing.T) {
	a := DefaultBanner("svc", false)
	b := DefaultBanner("svc", false)
	if a != b {
		t.Error("DefaultBanner should be deterministic")
	}
}

func TestRender_UsesDefaultBannerWhenEmpty(t *testing.T) {
	opts := Options{
		ServiceName: "auto-svc",
	}
	out := Render(opts, BuildInfo{})
	if !strings.Contains(out, "AUTO-SVC") {
		t.Errorf("expected auto-generated banner with service name, got:\n%s", out)
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

func TestRender_ASCIIOnlyBannerChars(t *testing.T) {
	opts := Options{
		ServiceName: "test",
		ASCIIOnly:   true,
	}
	out := Render(opts, BuildInfo{})
	if strings.Contains(out, "┌") || strings.Contains(out, "─") {
		t.Error("ASCIIOnly should not use Unicode box-drawing in banner")
	}
	if !strings.Contains(out, "+") {
		t.Error("ASCIIOnly banner should use '+' corners")
	}
}
