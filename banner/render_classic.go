package banner

import "strings"

// renderClassicText renders the input string as multi-line ASCII art using
// classicFont. The input is normalized: uppercased, whitespace collapsed to
// single spaces, characters outside [A-Z0-9-_ ] replaced with '?'. The result
// is exactly classicHeight lines with trailing spaces trimmed.
func renderClassicText(input string) []string {
	normalized := normalizeBigInput(input) // reuse existing normalizer

	rows := make([]strings.Builder, classicHeight)
	for i, r := range normalized {
		g, ok := classicFont[r]
		if !ok {
			g = classicFont['?']
		}
		for row := range classicHeight {
			if i > 0 {
				rows[row].WriteByte(' ')
			}
			rows[row].WriteString(g[row])
		}
	}

	result := make([]string, classicHeight)
	for i := range classicHeight {
		result[i] = strings.TrimRight(rows[i].String(), " ")
	}
	return result
}
