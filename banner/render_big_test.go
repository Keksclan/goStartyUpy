package banner

import (
	"strings"
	"testing"
)

func TestRenderBigText_ReturnsFiveLines(t *testing.T) {
	lines := renderBigText("HELLO")
	if len(lines) != fontHeight {
		t.Errorf("expected %d lines, got %d", fontHeight, len(lines))
	}
}

func TestRenderBigText_SingleChar_A(t *testing.T) {
	lines := renderBigText("A")
	expected := bigFont['A']
	for i, line := range lines {
		want := strings.TrimRight(expected[i], " ")
		if line != want {
			t.Errorf("row %d: got %q, want %q", i, line, want)
		}
	}
}

func TestRenderBigText_UnknownCharBecomesFallback(t *testing.T) {
	lines := renderBigText("@")
	fallback := bigFont['?']
	for i, line := range lines {
		want := strings.TrimRight(fallback[i], " ")
		if line != want {
			t.Errorf("row %d: got %q, want %q", i, line, want)
		}
	}
}

func TestRenderBigText_HyphenNonEmpty(t *testing.T) {
	lines := renderBigText("-")
	nonEmpty := false
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmpty = true
			break
		}
	}
	if !nonEmpty {
		t.Error("hyphen glyph should have non-empty content")
	}
}

func TestRenderBigText_UnderscoreNonEmpty(t *testing.T) {
	lines := renderBigText("_")
	nonEmpty := false
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmpty = true
			break
		}
	}
	if !nonEmpty {
		t.Error("underscore glyph should have non-empty content")
	}
}

func TestRenderBigText_NoTrailingSpaces(t *testing.T) {
	lines := renderBigText("AB")
	for i, line := range lines {
		if line != strings.TrimRight(line, " ") {
			t.Errorf("row %d has trailing spaces: %q", i, line)
		}
	}
}

func TestRenderBigText_Deterministic(t *testing.T) {
	a := renderBigText("TEST")
	b := renderBigText("TEST")
	for i := range fontHeight {
		if a[i] != b[i] {
			t.Errorf("row %d differs between calls", i)
		}
	}
}

func TestRenderBigText_EmptyInput(t *testing.T) {
	lines := renderBigText("")
	if len(lines) != fontHeight {
		t.Errorf("expected %d lines for empty input, got %d", fontHeight, len(lines))
	}
	for i, line := range lines {
		if line != "" {
			t.Errorf("row %d should be empty for empty input, got %q", i, line)
		}
	}
}

func TestRenderBigText_LowercaseNormalized(t *testing.T) {
	lower := renderBigText("abc")
	upper := renderBigText("ABC")
	for i := range fontHeight {
		if lower[i] != upper[i] {
			t.Errorf("row %d: lowercase and uppercase should render the same", i)
		}
	}
}

func TestRenderBigText_WhitespaceCollapsed(t *testing.T) {
	single := renderBigText("A B")
	multi := renderBigText("A   B")
	for i := range fontHeight {
		if single[i] != multi[i] {
			t.Errorf("row %d: multiple spaces should collapse to single space", i)
		}
	}
}

func TestNormalizeBigInput(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"hello", "HELLO"},
		{"a b  c", "A B C"},
		{"café", "CAF?"},
		{"a-b_c", "A-B_C"},
		{"  spaces  ", "SPACES"},
		{"123", "123"},
		{"@#$", "???"},
		{"", ""},
	}
	for _, tt := range tests {
		got := normalizeBigInput(tt.in)
		if got != tt.want {
			t.Errorf("normalizeBigInput(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBigFont_AllRequiredChars(t *testing.T) {
	required := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_ ?"
	for _, r := range required {
		if _, ok := bigFont[r]; !ok {
			t.Errorf("bigFont missing required character %q", r)
		}
	}
}

func TestBigFont_AllGlyphsHaveCorrectHeight(t *testing.T) {
	for r, g := range bigFont {
		if len(g) != fontHeight {
			t.Errorf("glyph %q has %d rows, want %d", r, len(g), fontHeight)
		}
	}
}
