package banner

import "strings"

// miniHeight is the number of rows each mini glyph occupies.
const miniHeight = 3

// glyph3 represents a single character rendered as miniHeight rows of text.
type glyph3 [miniHeight]string

// miniFont maps runes to their mini ASCII-art glyphs.
// Supported: A–Z, 0–9, '-', '_', ' ', '?'.
var miniFont = map[rune]glyph3{
	'A': {` _ `, `/_\`, `| |`},
	'B': {`|_ `, `|_)`, `|_|`},
	'C': {` _ `, `|  `, `|_ `},
	'D': {`| \`, `| |`, `|_/`},
	'E': {`___`, `|_ `, `|__`},
	'F': {`___`, `|_ `, `|  `},
	'G': {` __`, `| _`, `|__|`},
	'H': {`| |`, `|_|`, `| |`},
	'I': {`_`, `|`, `_`},
	'J': {` _`, ` |`, `_|`},
	'K': {`|/`, `|\`, `| \`},
	'L': {`|  `, `|  `, `|__`},
	'M': {`|\/|`, `|  |`, `|  |`},
	'N': {`|\ |`, `| \|`, `|  |`},
	'O': {` _ `, `| |`, `|_|`},
	'P': {`|_ `, `|_)`, `|  `},
	'Q': {` _ `, `| |`, `|_\|`},
	'R': {`|_ `, `| \`, `|  \`},
	'S': {` _ `, `|_ `, ` _|`},
	'T': {`___`, ` | `, ` | `},
	'U': {`| |`, `| |`, `|_|`},
	'V': {`| |`, `| |`, ` \/ `},
	'W': {`|  |`, `|/\|`, `|  |`},
	'X': {`\ /`, ` X `, `/ \`},
	'Y': {`\ /`, ` | `, ` | `},
	'Z': {`__ `, ` / `, `/__`},
	'0': {` _ `, `| |`, `|_|`},
	'1': {` | `, ` | `, ` | `},
	'2': {` _ `, ` _|`, `|_ `},
	'3': {` _ `, ` _|`, ` _|`},
	'4': {`|_|`, `  |`, `  |`},
	'5': {` _ `, `|_ `, ` _|`},
	'6': {` _ `, `|_ `, `|_|`},
	'7': {`__ `, `  /`, ` / `},
	'8': {` _ `, `|_|`, `|_|`},
	'9': {` _ `, `|_|`, ` _|`},
	'-': {`   `, `---`, `   `},
	'_': {`   `, `   `, `___`},
	' ': {`   `, `   `, `   `},
	'?': {` _ `, ` _|`, ` . `},
}

// renderMiniText renders the input string as multi-line ASCII art using miniFont.
func renderMiniText(input string) []string {
	normalized := normalizeBigInput(input)

	rows := make([]strings.Builder, miniHeight)
	for i, r := range normalized {
		g, ok := miniFont[r]
		if !ok {
			g = miniFont['?']
		}
		for row := range miniHeight {
			if i > 0 {
				rows[row].WriteByte(' ')
			}
			rows[row].WriteString(g[row])
		}
	}

	result := make([]string, miniHeight)
	for i := range miniHeight {
		result[i] = strings.TrimRight(rows[i].String(), " ")
	}
	return result
}

// MiniBanner generates a mini startup banner from the given service name.
func MiniBanner(serviceName string, _ bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	lines := renderMiniText(serviceName)
	return strings.Join(lines, "\n")
}
