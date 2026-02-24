package banner

import (
	"strings"
	"unicode"
)

// glyphHeight is the number of rows each glyph occupies.
const glyphHeight = 5

// blockFont is a minimal 5-row block font for A-Z, 0-9, '-', '_'.
// Each glyph is represented as [glyphHeight]string.
var blockFont = map[rune][glyphHeight]string{
	'A': {
		" AAA ",
		"A   A",
		"AAAAA",
		"A   A",
		"A   A",
	},
	'B': {
		"BBBB ",
		"B   B",
		"BBBB ",
		"B   B",
		"BBBB ",
	},
	'C': {
		" CCCC",
		"C    ",
		"C    ",
		"C    ",
		" CCCC",
	},
	'D': {
		"DDDD ",
		"D   D",
		"D   D",
		"D   D",
		"DDDD ",
	},
	'E': {
		"EEEEE",
		"E    ",
		"EEE  ",
		"E    ",
		"EEEEE",
	},
	'F': {
		"FFFFF",
		"F    ",
		"FFF  ",
		"F    ",
		"F    ",
	},
	'G': {
		" GGG ",
		"G    ",
		"G  GG",
		"G   G",
		" GGG ",
	},
	'H': {
		"H   H",
		"H   H",
		"HHHHH",
		"H   H",
		"H   H",
	},
	'I': {
		"IIIII",
		"  I  ",
		"  I  ",
		"  I  ",
		"IIIII",
	},
	'J': {
		"JJJJJ",
		"    J",
		"    J",
		"J   J",
		" JJJ ",
	},
	'K': {
		"K   K",
		"K  K ",
		"KKK  ",
		"K  K ",
		"K   K",
	},
	'L': {
		"L    ",
		"L    ",
		"L    ",
		"L    ",
		"LLLLL",
	},
	'M': {
		"M   M",
		"MM MM",
		"M M M",
		"M   M",
		"M   M",
	},
	'N': {
		"N   N",
		"NN  N",
		"N N N",
		"N  NN",
		"N   N",
	},
	'O': {
		" OOO ",
		"O   O",
		"O   O",
		"O   O",
		" OOO ",
	},
	'P': {
		"PPPP ",
		"P   P",
		"PPPP ",
		"P    ",
		"P    ",
	},
	'Q': {
		" QQQ ",
		"Q   Q",
		"Q   Q",
		"Q  Q ",
		" QQ Q",
	},
	'R': {
		"RRRR ",
		"R   R",
		"RRRR ",
		"R  R ",
		"R   R",
	},
	'S': {
		" SSS ",
		"S    ",
		" SSS ",
		"    S",
		" SSS ",
	},
	'T': {
		"TTTTT",
		"  T  ",
		"  T  ",
		"  T  ",
		"  T  ",
	},
	'U': {
		"U   U",
		"U   U",
		"U   U",
		"U   U",
		" UUU ",
	},
	'V': {
		"V   V",
		"V   V",
		"V   V",
		" V V ",
		"  V  ",
	},
	'W': {
		"W   W",
		"W   W",
		"W W W",
		"WW WW",
		"W   W",
	},
	'X': {
		"X   X",
		" X X ",
		"  X  ",
		" X X ",
		"X   X",
	},
	'Y': {
		"Y   Y",
		" Y Y ",
		"  Y  ",
		"  Y  ",
		"  Y  ",
	},
	'Z': {
		"ZZZZZ",
		"   Z ",
		"  Z  ",
		" Z   ",
		"ZZZZZ",
	},
	'0': {
		" 000 ",
		"0   0",
		"0   0",
		"0   0",
		" 000 ",
	},
	'1': {
		" 1   ",
		"11   ",
		" 1   ",
		" 1   ",
		"11111",
	},
	'2': {
		" 222 ",
		"2   2",
		"  22 ",
		" 2   ",
		"22222",
	},
	'3': {
		" 333 ",
		"    3",
		"  33 ",
		"    3",
		" 333 ",
	},
	'4': {
		"4   4",
		"4   4",
		"44444",
		"    4",
		"    4",
	},
	'5': {
		"55555",
		"5    ",
		"5555 ",
		"    5",
		"5555 ",
	},
	'6': {
		" 666 ",
		"6    ",
		"6666 ",
		"6   6",
		" 666 ",
	},
	'7': {
		"77777",
		"   7 ",
		"  7  ",
		" 7   ",
		"7    ",
	},
	'8': {
		" 888 ",
		"8   8",
		" 888 ",
		"8   8",
		" 888 ",
	},
	'9': {
		" 999 ",
		"9   9",
		" 9999",
		"    9",
		" 999 ",
	},
	'-': {
		"     ",
		"     ",
		"-----",
		"     ",
		"     ",
	},
	'_': {
		"     ",
		"     ",
		"     ",
		"     ",
		"_____",
	},
	'?': {
		" ??? ",
		"    ?",
		"  ?? ",
		"     ",
		"  ?  ",
	},
}

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

// renderWordmark renders the given text as multi-line ASCII art using blockFont.
func renderWordmark(text string) string {
	if text == "" {
		return ""
	}
	rows := [glyphHeight]strings.Builder{}
	for i, r := range text {
		glyph, ok := blockFont[r]
		if !ok {
			glyph = blockFont['?']
		}
		for row := range glyphHeight {
			if i > 0 {
				rows[row].WriteByte(' ')
			}
			rows[row].WriteString(glyph[row])
		}
	}
	var b strings.Builder
	for row := range glyphHeight {
		b.WriteString(rows[row].String())
		b.WriteByte('\n')
	}
	return b.String()
}

// SpringLikeBanner generates a Spring Boot–style startup banner from the
// service name. It produces a large ASCII-art wordmark followed by a tagline.
// The asciiOnly parameter is accepted for API consistency but does not
// currently alter the output (the wordmark uses plain ASCII by default).
func SpringLikeBanner(serviceName string, asciiOnly bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	name := normalizeName(serviceName)
	if name == "" {
		name = "SERVICE"
	}

	ver := Version
	if ver == "" {
		ver = "dev"
	}

	var b strings.Builder
	b.WriteString(renderWordmark(name))
	b.WriteByte('\n')
	b.WriteString(" :: goStartyUpy :: (" + ver + ")")
	return b.String()
}
