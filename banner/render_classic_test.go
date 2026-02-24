package banner

import (
	"strings"
	"testing"
)

func TestRenderClassicText_ReturnsSevenLines(t *testing.T) {
	lines := renderClassicText("HELLO")
	if len(lines) != classicHeight {
		t.Errorf("expected %d lines, got %d", classicHeight, len(lines))
	}
}

func TestRenderClassicText_SingleChar_A(t *testing.T) {
	lines := renderClassicText("A")
	expected := classicFont['A']
	for i, line := range lines {
		want := strings.TrimRight(expected[i], " ")
		if line != want {
			t.Errorf("row %d: got %q, want %q", i, line, want)
		}
	}
}

func TestRenderClassicText_UnknownCharBecomesFallback(t *testing.T) {
	lines := renderClassicText("@")
	fallback := classicFont['?']
	for i, line := range lines {
		want := strings.TrimRight(fallback[i], " ")
		if line != want {
			t.Errorf("row %d: got %q, want %q", i, line, want)
		}
	}
}

func TestRenderClassicText_NoTrailingSpaces(t *testing.T) {
	lines := renderClassicText("AB")
	for i, line := range lines {
		if line != strings.TrimRight(line, " ") {
			t.Errorf("row %d has trailing spaces: %q", i, line)
		}
	}
}

func TestRenderClassicText_Deterministic(t *testing.T) {
	a := renderClassicText("TEST")
	b := renderClassicText("TEST")
	for i := range classicHeight {
		if a[i] != b[i] {
			t.Errorf("row %d differs between calls", i)
		}
	}
}

func TestRenderClassicText_EmptyInput(t *testing.T) {
	lines := renderClassicText("")
	if len(lines) != classicHeight {
		t.Errorf("expected %d lines for empty input, got %d", classicHeight, len(lines))
	}
	for i, line := range lines {
		if line != "" {
			t.Errorf("row %d should be empty for empty input, got %q", i, line)
		}
	}
}

func TestRenderClassicText_LowercaseNormalized(t *testing.T) {
	lower := renderClassicText("abc")
	upper := renderClassicText("ABC")
	for i := range classicHeight {
		if lower[i] != upper[i] {
			t.Errorf("row %d: lowercase and uppercase should render the same", i)
		}
	}
}

func TestClassicFont_AllRequiredChars(t *testing.T) {
	required := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_ ?"
	for _, r := range required {
		if _, ok := classicFont[r]; !ok {
			t.Errorf("classicFont missing required character %q", r)
		}
	}
}

func TestClassicFont_AllGlyphsHaveCorrectHeight(t *testing.T) {
	for r, g := range classicFont {
		if len(g) != classicHeight {
			t.Errorf("glyph %q has %d rows, want %d", r, len(g), classicHeight)
		}
	}
}

func TestClassicLikeBanner_EmptyNameUsesService(t *testing.T) {
	out := ClassicLikeBanner("", false)
	// Should render "SERVICE" and not panic.
	if out == "" {
		t.Error("ClassicLikeBanner should not return empty string")
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != classicHeight {
		t.Errorf("expected %d wordmark lines, got %d", classicHeight, len(lines))
	}
}

func TestClassicLikeBanner_NoPanicOnUnicode(t *testing.T) {
	// Must not panic on weird input.
	_ = ClassicLikeBanner("日本語サービス", false)
	_ = ClassicLikeBanner("🚀🎉", false)
	_ = ClassicLikeBanner(strings.Repeat("x", 200), false)
}

func TestClassicLikeBanner_Deterministic(t *testing.T) {
	a := ClassicLikeBanner("det", false)
	b := ClassicLikeBanner("det", false)
	if a != b {
		t.Error("ClassicLikeBanner should be deterministic")
	}
}

func TestBannerStyle_Classic(t *testing.T) {
	opts := Options{
		ServiceName: "my-svc",
		BannerStyle: "classic",
	}
	art := resolveBanner(opts)
	if art == "" {
		t.Error("classic style should produce non-empty banner")
	}
	// Classic banners use slashes/backslashes.
	if !strings.Contains(art, `\`) && !strings.Contains(art, `/`) {
		t.Errorf("classic style should contain slashes, got:\n%s", art)
	}
}

func TestBannerStyle_Switching(t *testing.T) {
	spring := resolveBanner(Options{ServiceName: "test", BannerStyle: "spring"})
	classic := resolveBanner(Options{ServiceName: "test", BannerStyle: "classic"})
	box := resolveBanner(Options{ServiceName: "test", BannerStyle: "box"})

	if spring == classic {
		t.Error("spring and classic styles should produce different output")
	}
	if spring == box {
		t.Error("spring and box styles should produce different output")
	}
	if classic == box {
		t.Error("classic and box styles should produce different output")
	}
}

func TestClassicBanner_DefaultTaglines(t *testing.T) {
	opts := Options{
		ServiceName: "my-svc",
		BannerStyle: "classic",
	}
	info := BuildInfo{
		Version:   "v1.0.0",
		BuildTime: "2024-01-01",
		Commit:    "abc123",
		Branch:    "main",
		Dirty:     "false",
	}
	out := Render(opts, info)
	if !strings.Contains(out, "my-svc v1.0.0") {
		t.Errorf("expected tagline1 with service name and version, got:\n%s", out)
	}
	if !strings.Contains(out, "Build: 2024-01-01") {
		t.Errorf("expected tagline2 with build time, got:\n%s", out)
	}
	if !strings.Contains(out, "Commit: abc123") {
		t.Errorf("expected tagline2 with commit, got:\n%s", out)
	}
	if !strings.Contains(out, "Branch: main") {
		t.Errorf("expected tagline2 with branch, got:\n%s", out)
	}
}

func TestClassicBanner_CustomTaglines(t *testing.T) {
	opts := Options{
		ServiceName: "svc",
		BannerStyle: "classic",
		Tagline1:    "Custom Line 1",
		Tagline2:    "Custom Line 2",
	}
	info := BuildInfo{}
	out := Render(opts, info)
	if !strings.Contains(out, "Custom Line 1") {
		t.Errorf("expected custom tagline1, got:\n%s", out)
	}
	if !strings.Contains(out, "Custom Line 2") {
		t.Errorf("expected custom tagline2, got:\n%s", out)
	}
}

func TestClassicBanner_DefaultTaglineEmptyService(t *testing.T) {
	opts := Options{BannerStyle: "classic"}
	info := BuildInfo{}
	out := Render(opts, info)
	// Default: "SERVICE dev"
	if !strings.Contains(out, "SERVICE dev") {
		t.Errorf("expected default tagline 'SERVICE dev', got:\n%s", out)
	}
}

func TestClassicBanner_ShowDetailsFalse(t *testing.T) {
	hide := false
	opts := Options{
		ServiceName: "svc",
		BannerStyle: "classic",
		ShowDetails: &hide,
	}
	info := BuildInfo{Version: "v1.0.0"}
	out := Render(opts, info)
	// When ShowDetails is false, key/value section should be absent.
	if strings.Contains(out, "Service") && strings.Contains(out, ":") && strings.Contains(out, "svc") {
		// Check it doesn't have the aligned "Service : svc" line.
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "Service") && strings.Contains(line, ":") && strings.Contains(line, "svc") {
				t.Errorf("ShowDetails=false should hide key/value section, but found: %q", line)
			}
		}
	}
}

func TestClassicBanner_DirtyTagline(t *testing.T) {
	opts := Options{
		ServiceName: "svc",
		BannerStyle: "classic",
	}
	info := BuildInfo{
		Version:   "v1.0.0",
		BuildTime: "2024-01-01",
		Commit:    "abc",
		Branch:    "dev",
		Dirty:     "true",
	}
	out := Render(opts, info)
	if !strings.Contains(out, "Dirty: true") {
		t.Errorf("expected dirty flag in tagline2, got:\n%s", out)
	}
}

func TestResolveStyle(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want string
	}{
		{"default", Options{}, "spring"},
		{"spring", Options{BannerStyle: "spring"}, "spring"},
		{"classic", Options{BannerStyle: "classic"}, "classic"},
		{"box", Options{BannerStyle: "box"}, "box"},
		{"raw override", Options{Banner: "RAW", BannerStyle: "classic"}, "raw"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveStyle(tt.opts)
			if got != tt.want {
				t.Errorf("resolveStyle() = %q, want %q", got, tt.want)
			}
		})
	}
}
