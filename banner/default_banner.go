package banner

import "strings"

// DefaultBanner generates a simple box banner from the given service name.
// If serviceName is empty, "SERVICE" is used. When asciiOnly is true the box
// is drawn with plain ASCII characters (+, -, |); otherwise Unicode
// box-drawing characters (┌, ─, ┐, │, └, ┘) are used.
func DefaultBanner(serviceName string, asciiOnly bool) string {
	if serviceName == "" {
		serviceName = "SERVICE"
	}
	name := strings.ToUpper(serviceName)

	const minWidth = 27
	// Inner width = name length + 2*padding (at least 4 each side).
	inner := len(name) + 8
	if inner < minWidth {
		inner = minWidth
	}

	// Centre the name within the inner width.
	pad := inner - len(name)
	left := pad / 2
	right := pad - left
	centred := strings.Repeat(" ", left) + name + strings.Repeat(" ", right)

	var tl, tr, bl, br, h, v string
	if asciiOnly {
		tl, tr, bl, br, h, v = "+", "+", "+", "+", "-", "|"
	} else {
		tl, tr, bl, br, h, v = "┌", "┐", "└", "┘", "─", "│"
	}

	hLine := strings.Repeat(h, inner)
	top := tl + hLine + tr
	mid := v + centred + v
	bot := bl + hLine + br

	return top + "\n" + mid + "\n" + bot
}

// asciiSeparator is the plain-ASCII fallback for the info separator line.
const asciiSeparator = "============================================================"
