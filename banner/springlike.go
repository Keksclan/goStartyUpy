package banner

import (
	"strings"
	"unicode"
)

// normalizeName sanitizes a service name for banner rendering.
// It uppercases the name, keeps only A-Z 0-9 - _, and replaces spaces with
// dashes.
func normalizeName(name string) string {
	name = strings.ToUpper(name)
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune('-')
		default:
			// skip unknown characters
		}
	}
	return b.String()
}

// renderWordmark renders the given text as multi-line ASCII art using bigFont.
func renderWordmark(text string) string {
	if text == "" {
		return ""
	}
	lines := renderBigText(text)
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// SpringLikeBanner generates a Spring Boot–style startup banner from the
// service name. It produces a large ASCII-art wordmark followed by a tagline.
// The asciiOnly parameter is accepted for API consistency but does not
// currently alter the output (the wordmark uses plain ASCII by default).
func SpringLikeBanner(serviceName string, asciiOnly bool) string {
	return springLikeBannerInternal(serviceName, "", false, asciiOnly)
}

func springLikeBannerInternal(serviceName string, env string, fromEnv bool, _ bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	name := normalizeName(serviceName)
	if name == "" {
		name = "SERVICE"
	}

	var b strings.Builder
	b.WriteString(renderWordmark(name))
	b.WriteByte('\n')
	b.WriteString(" :: goStartyUpy ::")
	if fromEnv && env != "" {
		b.WriteString(" (" + env + ")")
	}
	return b.String()
}
