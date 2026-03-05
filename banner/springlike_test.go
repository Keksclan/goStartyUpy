package banner

import (
	"strings"
	"testing"
)

func TestSpringLikeBanner_DefaultWhenEmpty(t *testing.T) {
	opts := Options{ServiceName: "test-svc"}
	art := resolveBanner(opts)
	spring := SpringLikeBanner("test-svc", false)
	if art != spring {
		t.Error("default style should use SpringLikeBanner")
	}
}

func TestSpringLikeBanner_ContainsTagline(t *testing.T) {
	out := SpringLikeBanner("hello", false)
	if !strings.Contains(out, ":: goStartyUpy ::") {
		t.Errorf("expected tagline, got:\n%s", out)
	}
}

func TestSpringLikeBanner_NoDefaultSuffix(t *testing.T) {
	out := SpringLikeBanner("x", false)
	// Version/env suffix is no longer shown by default.
	if strings.Contains(out, "(") || strings.Contains(out, ")") {
		t.Errorf("did not expect suffix in default tagline, got:\n%s", out)
	}
}

func TestSpringLikeBanner_MultiLine(t *testing.T) {
	out := SpringLikeBanner("AB", false)
	lines := strings.Split(out, "\n")
	// At least 5 wordmark lines + 1 empty + 1 tagline = 7
	if len(lines) < 7 {
		t.Errorf("expected at least 7 lines, got %d:\n%s", len(lines), out)
	}
}

func TestSpringLikeBanner_EmptyNameUsesService(t *testing.T) {
	out := SpringLikeBanner("", false)
	// "SERVICE" wordmark includes the S glyph which contains "/ ___|"
	if !strings.Contains(out, "/ ___|") {
		t.Errorf("expected SERVICE wordmark for empty name, got:\n%s", out)
	}
}

func TestSpringLikeBanner_Deterministic(t *testing.T) {
	a := SpringLikeBanner("det", false)
	b := SpringLikeBanner("det", false)
	if a != b {
		t.Error("SpringLikeBanner should be deterministic")
	}
}

func TestNormalizeName_Basic(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"hello", "HELLO"},
		{"My Service", "MY-SERVICE"},
		{"a_b-c", "A_B-C"},
		{"café", "CAF"},
		{"", ""},
		{"123", "123"},
		{"a b  c", "A-B--C"},
		{"★unicode★", "UNICODE"},
	}
	for _, tt := range tests {
		got := normalizeName(tt.in)
		if got != tt.want {
			t.Errorf("normalizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSpringLikeBanner_NoPanicOnUnicode(t *testing.T) {
	// Must not panic on weird input.
	_ = SpringLikeBanner("日本語サービス", false)
	_ = SpringLikeBanner("🚀🎉", false)
	_ = SpringLikeBanner(strings.Repeat("x", 1000), false)
}

func TestBannerStyle_Box(t *testing.T) {
	opts := Options{
		ServiceName: "my-svc",
		BannerStyle: "box",
	}
	art := resolveBanner(opts)
	if !strings.Contains(art, "MY-SVC") {
		t.Errorf("box style should contain service name, got:\n%s", art)
	}
}

func TestBannerWidth_ClampsLines(t *testing.T) {
	opts := Options{
		ServiceName: "LONGNAME",
		BannerWidth: 20,
	}
	art := resolveBanner(opts)
	for i, line := range strings.Split(art, "\n") {
		if len(line) > 20 {
			t.Errorf("line %d exceeds BannerWidth 20 (%d chars): %q", i, len(line), line)
		}
	}
}

func TestBannerWidth_ZeroNoClamping(t *testing.T) {
	opts := Options{
		ServiceName: "LONGNAME",
		BannerWidth: 0,
	}
	art := resolveBanner(opts)
	// Should have long wordmark lines (LONGNAME = 8 chars * 6 wide = ~53).
	maxLen := 0
	for _, line := range strings.Split(art, "\n") {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	if maxLen < 40 {
		t.Errorf("expected long lines without clamping, max was %d", maxLen)
	}
}

func TestRawBanner_IgnoresStyle(t *testing.T) {
	opts := Options{
		Banner:      "MY RAW BANNER",
		BannerStyle: "box",
	}
	art := resolveBanner(opts)
	if art != "MY RAW BANNER" {
		t.Errorf("raw banner should be used as-is, got: %q", art)
	}
}
