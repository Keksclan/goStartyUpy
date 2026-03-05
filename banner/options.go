package banner

// Options configures the banner output. All fields are optional; sensible
// defaults are applied during rendering.
//
// Service-specific addresses (HTTP, gRPC, etc.) are not dedicated fields;
// pass them via the Extra map instead.
type Options struct {
	// ServiceName is the human-readable name of the service.
	ServiceName string
	// Environment is the deployment environment (e.g. "production", "staging").
	Environment string
	// Banner is the multiline ASCII art banner. If empty, a built-in default
	// is used.
	Banner string
	// Separator is the line drawn between the banner and the info block.
	// If empty, a default Unicode line is used.
	Separator string
	// Extra holds additional key-value pairs that are printed after the
	// standard fields, sorted alphabetically by key. Use this for addresses
	// like "HTTP" or "gRPC".
	Extra map[string]string
	// Color enables ANSI color escape sequences in the output. Set to true
	// when writing to a terminal that supports colors. When false (default),
	// the output is plain text with no escape sequences.
	Color bool
	// ASCIIOnly forces plain ASCII characters for the auto-generated banner
	// box and the default separator line. Use this when the output terminal
	// does not support Unicode box-drawing characters.
	ASCIIOnly bool
	// BannerStyle selects the auto-generated banner style when Banner is
	// empty. Supported values: "spring" (default), "classic", "box".
	// When Banner is set, this field is ignored and the banner is treated
	// as raw text. The "classic" style renders a Spring Boot–like banner
	// with slashes, backslashes and underscores plus two tagline lines.
	BannerStyle string
	// BannerWidth is the optional maximum line width. Lines longer than
	// this value are hard-cut. A value of 0 (default) means no clamping.
	BannerWidth int
	// Tagline1 overrides the first line printed below the classic-style
	// wordmark. Defaults to "<ServiceName> <Version>".
	Tagline1 string
	// Tagline2 overrides the second line printed below the classic-style
	// wordmark. Defaults to "Build: <BuildTime>  Commit: <Commit>".
	Tagline2 string
	// ShowDetails controls whether the key/value info section is printed
	// after the banner. Defaults to true when zero-valued (false means
	// the details are hidden). Only affects "classic" style; other styles
	// always show details.
	ShowDetails *bool
	// EnvironmentFromEnv is an internal flag (additive) that indicates
	// whether the Environment value was sourced from an environment
	// variable at runtime. When true, the banner header may display
	// the environment as a suffix.
	EnvironmentFromEnv bool
}
