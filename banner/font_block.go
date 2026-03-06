package banner

import "strings"

// blockHeight is the number of rows each block glyph occupies.
const blockHeight = 5

// glyph5 represents a single character rendered as blockHeight rows of text.
type glyph5 [blockHeight]string

// blockFont maps runes to their block ASCII-art glyphs.
// Supported: A–Z, 0–9, '-', '_', ' ', '?'.
var blockFont = map[rune]glyph5{
	'A': {`  ###  `, ` #   # `, ` ##### `, ` #   # `, ` #   # `},
	'B': {` ####  `, ` #   # `, ` ####  `, ` #   # `, ` ####  `},
	'C': {`  #### `, ` #     `, ` #     `, ` #     `, `  #### `},
	'D': {` ####  `, ` #   # `, ` #   # `, ` #   # `, ` ####  `},
	'E': {` ##### `, ` #     `, ` ####  `, ` #     `, ` ##### `},
	'F': {` ##### `, ` #     `, ` ####  `, ` #     `, ` #     `},
	'G': {`  #### `, ` #     `, ` #  ## `, ` #   # `, `  #### `},
	'H': {` #   # `, ` #   # `, ` ##### `, ` #   # `, ` #   # `},
	'I': {` ##### `, `   #   `, `   #   `, `   #   `, ` ##### `},
	'J': {`  #### `, `     # `, `     # `, ` #   # `, `  ###  `},
	'K': {` #   # `, ` #  #  `, ` ###   `, ` #  #  `, ` #   # `},
	'L': {` #     `, ` #     `, ` #     `, ` #     `, ` ##### `},
	'M': {` #   # `, ` ## ## `, ` # # # `, ` #   # `, ` #   # `},
	'N': {` #   # `, ` ##  # `, ` # # # `, ` #  ## `, ` #   # `},
	'O': {`  ###  `, ` #   # `, ` #   # `, ` #   # `, `  ###  `},
	'P': {` ####  `, ` #   # `, ` ####  `, ` #     `, ` #     `},
	'Q': {`  ###  `, ` #   # `, ` #   # `, ` #  ## `, `  #### `},
	'R': {` ####  `, ` #   # `, ` ####  `, ` #  #  `, ` #   # `},
	'S': {`  #### `, ` #     `, `  ###  `, `     # `, ` ####  `},
	'T': {` ##### `, `   #   `, `   #   `, `   #   `, `   #   `},
	'U': {` #   # `, ` #   # `, ` #   # `, ` #   # `, `  ###  `},
	'V': {` #   # `, ` #   # `, ` #   # `, `  # #  `, `   #   `},
	'W': {` #   # `, ` #   # `, ` # # # `, ` ## ## `, ` #   # `},
	'X': {` #   # `, `  # #  `, `   #   `, `  # #  `, ` #   # `},
	'Y': {` #   # `, `  # #  `, `   #   `, `   #   `, `   #   `},
	'Z': {` ##### `, `    #  `, `   #   `, `  #    `, ` ##### `},
	'0': {`  ###  `, ` #   # `, ` #   # `, ` #   # `, `  ###  `},
	'1': {`   #   `, `  ##   `, `   #   `, `   #   `, `  ###  `},
	'2': {`  ###  `, `     # `, `   ##  `, `  #    `, `  #### `},
	'3': {`  ###  `, `     # `, `   ##  `, `     # `, `  ###  `},
	'4': {` #   # `, ` #   # `, `  #### `, `     # `, `     # `},
	'5': {` ##### `, ` #     `, ` ####  `, `     # `, ` ####  `},
	'6': {`  ###  `, ` #     `, ` ####  `, ` #   # `, `  ###  `},
	'7': {` ##### `, `     # `, `    #  `, `   #   `, `  #    `},
	'8': {`  ###  `, ` #   # `, `  ###  `, ` #   # `, `  ###  `},
	'9': {`  ###  `, ` #   # `, `  #### `, `     # `, `  ###  `},
	'-': {`       `, `       `, `  ###  `, `       `, `       `},
	'_': {`       `, `       `, `       `, `       `, ` ##### `},
	' ': {`       `, `       `, `       `, `       `, `       `},
	'?': {`  ###  `, ` #   # `, `    #  `, `       `, `    #  `},
}

// renderBlockText renders the input string as multi-line ASCII art using blockFont.
func renderBlockText(input string) []string {
	normalized := normalizeBigInput(input)
	if normalized == "" {
		normalized = "SERVICE"
	}

	rows := make([]strings.Builder, blockHeight)
	for i, r := range normalized {
		g, ok := blockFont[r]
		if !ok {
			g = blockFont['?']
		}
		for row := range blockHeight {
			if i > 0 {
				rows[row].WriteByte(' ')
			}
			rows[row].WriteString(g[row])
		}
	}

	result := make([]string, blockHeight)
	for i := range blockHeight {
		result[i] = strings.TrimRight(rows[i].String(), " ")
	}
	return result
}

// BlockBanner generates a blocky startup banner from the given service name.
func BlockBanner(serviceName string, _ bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	lines := renderBlockText(serviceName)
	return strings.Join(lines, "\n")
}
