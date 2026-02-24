package banner

import (
	"strings"
	"unicode"
)

// renderBigText renders the input string as multi-line ASCII art using bigFont.
// The input is normalized: uppercased, whitespace collapsed to single spaces,
// characters outside [A-Z0-9-_ ] replaced with '?'. The result is exactly
// fontHeight lines with trailing spaces trimmed.
func renderBigText(input string) []string {
	normalized := normalizeBigInput(input)

	rows := make([]strings.Builder, fontHeight)
	for i, r := range normalized {
		g, ok := bigFont[r]
		if !ok {
			g = bigFont['?']
		}
		for row := range fontHeight {
			if i > 0 {
				rows[row].WriteByte(' ')
			}
			rows[row].WriteString(g[row])
		}
	}

	result := make([]string, fontHeight)
	for i := range fontHeight {
		result[i] = strings.TrimRight(rows[i].String(), " ")
	}
	return result
}

// normalizeBigInput prepares input for big-font rendering:
// - converts to uppercase
// - collapses whitespace runs to a single space
// - keeps only A-Z, 0-9, '-', '_', ' '
// - replaces anything else with '?'
func normalizeBigInput(input string) string {
	input = strings.ToUpper(input)

	var b strings.Builder
	b.Grow(len(input))
	prevSpace := false
	for _, r := range input {
		switch {
		case r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
			prevSpace = false
		case unicode.IsSpace(r):
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
		default:
			b.WriteByte('?')
			prevSpace = false
		}
	}
	return strings.TrimSpace(b.String())
}
