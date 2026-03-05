package banner

import "strings"

// resolveBanner selects the banner art based on Options. If opts.Banner is
// set it is used as-is (raw). Otherwise the BannerStyle field controls the
// generator: "spring" (default) uses SpringLikeBanner, "classic" uses
// ClassicLikeBanner, "box" uses BoxBanner.
// If BannerWidth > 0 every line is hard-cut to that width.
func resolveBanner(opts Options) string {
	var art string
	if opts.Banner != "" {
		art = opts.Banner
	} else {
		style := opts.BannerStyle
		if style == "" {
			style = "spring"
		}
		switch style {
		case "classic":
			art = ClassicLikeBanner(opts.ServiceName, opts.ASCIIOnly)
		case "box":
			art = BoxBanner(opts.ServiceName, opts.ASCIIOnly)
		case "mini":
			art = MiniBanner(opts.ServiceName, opts.ASCIIOnly)
		case "block":
			art = BlockBanner(opts.ServiceName, opts.ASCIIOnly)
		default: // "spring"
			art = springLikeBannerInternal(opts.ServiceName, opts.Environment, opts.EnvironmentFromEnv, opts.ASCIIOnly)
		}
	}
	if opts.BannerWidth > 0 {
		art = clampLines(art, opts.BannerWidth)
	}
	return art
}

// clampLines hard-cuts every line in s to at most width characters.
func clampLines(s string, width int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if len(line) > width {
			lines[i] = line[:width]
		}
	}
	return strings.Join(lines, "\n")
}

// DefaultBanner generates a startup banner from the service name using the
// default style ("spring"). For the legacy box style use BoxBanner.
func DefaultBanner(serviceName string, asciiOnly bool) string {
	return SpringLikeBanner(serviceName, asciiOnly)
}

// BoxBanner generates a simple box banner from the given service name.
// If serviceName is empty, "SERVICE" is used. When asciiOnly is true the box
// is drawn with plain ASCII characters (+, -, |); otherwise Unicode
// box-drawing characters (┌, ─, ┐, │, └, ┘) are used.
func BoxBanner(serviceName string, asciiOnly bool) string {
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
