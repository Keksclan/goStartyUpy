package banner

// ANSI escape sequences for terminal colors. These are no-ops when color is
// disabled (the render functions simply skip them).
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiItalic = "\033[3m"

	// Foreground colors.
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiWhite   = "\033[37m"

	// Bright foreground colors.
	ansiBrightGreen = "\033[92m"
	ansiBrightRed   = "\033[91m"
	ansiBrightCyan  = "\033[96m"
)

// colorize wraps text in ANSI codes if color is enabled.
func colorize(text, code string, color bool) string {
	if !color {
		return text
	}
	return code + text + ansiReset
}
