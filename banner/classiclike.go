package banner

import "strings"

// ClassicLikeBanner generates a classic Spring Boot–style startup banner from
// the service name. It produces a large ASCII-art wordmark using the classic
// font (slashes, backslashes, underscores) followed by two tagline lines.
//
// The asciiOnly parameter controls whether backticks are replaced with
// apostrophes. In classic mode no Unicode box-drawing characters are used
// regardless.
func ClassicLikeBanner(serviceName string, asciiOnly bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	name := normalizeName(serviceName)
	if name == "" {
		name = "SERVICE"
	}

	lines := renderClassicText(name)
	var b strings.Builder
	for _, line := range lines {
		if asciiOnly {
			line = strings.ReplaceAll(line, "`", "'")
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}
