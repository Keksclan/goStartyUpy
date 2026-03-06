package banner

import (
	"strings"
	"testing"
)

func TestMiniBanner(t *testing.T) {
	out := MiniBanner("ABC", false)
	lines := strings.Split(out, "\n")
	if len(lines) != miniHeight {
		t.Errorf("expected %d lines, got %d", miniHeight, len(lines))
	}
	// Check if "ABC" is rendered. Just a smoke test.
	if !strings.Contains(out, "_") {
		t.Error("MiniBanner output looks wrong")
	}
}

func TestBlockBanner(t *testing.T) {
	out := BlockBanner("ABC", false)
	lines := strings.Split(out, "\n")
	if len(lines) != blockHeight {
		t.Errorf("expected %d lines, got %d", blockHeight, len(lines))
	}
	// Check if "ABC" is rendered.
	if !strings.Contains(out, "#") {
		t.Error("BlockBanner output looks wrong")
	}
}

func TestResolveBanner_NewStyles(t *testing.T) {
	optsMini := Options{BannerStyle: "mini", ServiceName: "TEST"}
	outMini := resolveBanner(optsMini)
	if len(strings.Split(outMini, "\n")) != miniHeight {
		t.Error("resolveBanner failed for style 'mini'")
	}

	optsBlock := Options{BannerStyle: "block", ServiceName: "TEST"}
	outBlock := resolveBanner(optsBlock)
	if len(strings.Split(outBlock, "\n")) != blockHeight {
		t.Error("resolveBanner failed for style 'block'")
	}
}
