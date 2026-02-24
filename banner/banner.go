package banner

import (
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/keksclan/goStartyUpy/checks"
	"maps"
)

// defaultBanner is the built-in ASCII art shown when Options.Banner is empty.
const defaultBanner = `
   _____ _                 _         _    _
  / ____| |               | |       | |  | |
 | (___ | |_ __ _ _ __ | |_ _   _| |  | |_ __  _   _
  \___ \| __/ _' | '__| __| | | | |  | | '_ \| | | |
  ____) | || (_| | |  | |_| |_| | |__| | |_) | |_| |
 |_____/ \__\__,_|_|   \__|\__, |\____/| .__/ \__, |
                             __/ |      | |     __/ |
                            |___/       |_|    |___/
`

// defaultSeparator is a Unicode box-drawing line used between banner and info.
const defaultSeparator = "════════════════════════════════════════════════════════════"

// Render produces a startup message string containing the ASCII banner,
// a separator, and aligned key-value metadata. The output always ends with
// a newline.
func Render(opts Options, info BuildInfo) string {
	return RenderWithChecks(opts, info, nil)
}

// RenderWithChecks is like Render but appends a "Checks:" section when
// results is non-empty.
func RenderWithChecks(opts Options, info BuildInfo, results []checks.Result) string {
	var b strings.Builder
	c := opts.Color

	// --- banner art ---
	art := opts.Banner
	if art == "" {
		art = defaultBanner
	}
	// Ensure exactly one trailing newline after the art block.
	b.WriteString(colorize(strings.TrimRight(art, "\n"), ansiCyan+ansiBold, c))
	b.WriteByte('\n')

	// --- separator ---
	sep := opts.Separator
	if sep == "" {
		sep = defaultSeparator
	}
	b.WriteString(colorize(sep, ansiDim, c))
	b.WriteByte('\n')

	// --- key/value info lines ---
	kvs := buildKVs(opts, info)
	writeAligned(&b, kvs, c)

	// --- checks section ---
	if len(results) > 0 {
		b.WriteByte('\n')
		b.WriteString(colorize("Checks:", ansiBold, c))
		b.WriteByte('\n')
		allOK := true
		for _, r := range results {
			if r.OK {
				tag := colorize("[OK]", ansiBrightGreen+ansiBold, c)
				fmt.Fprintf(&b, "  %s   %s (%s)\n", tag, r.Name, r.Duration)
			} else {
				allOK = false
				tag := colorize("[FAIL]", ansiBrightRed+ansiBold, c)
				errMsg := colorize(r.Error, ansiRed, c)
				fmt.Fprintf(&b, "  %s %s (%s): %s\n", tag, r.Name, r.Duration, errMsg)
			}
		}
		b.WriteByte('\n')
		if allOK {
			b.WriteString(colorize("Startup Complete", ansiBrightGreen+ansiBold, c))
			b.WriteByte('\n')
		} else {
			b.WriteString(colorize("Startup Failed", ansiBrightRed+ansiBold, c))
			b.WriteByte('\n')
		}
	}

	// Guarantee trailing newline.
	s := b.String()
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}

// kv is an internal key-value pair used for aligned formatting.
type kv struct {
	Key   string
	Value string
}

// buildKVs assembles the ordered list of key-value pairs to display.
func buildKVs(opts Options, info BuildInfo) []kv {
	pairs := []kv{
		{"Service", opts.ServiceName},
		{"Environment", opts.Environment},
		{"Version", info.Version},
		{"BuildTime", info.BuildTime},
		{"Commit", info.Commit},
		{"Branch", info.Branch},
		{"Dirty", info.Dirty},
		{"Go", runtime.Version()},
		{"OS/Arch", runtime.GOOS + "/" + runtime.GOARCH},
		{"PID", fmt.Sprintf("%d", os.Getpid())},
	}

	// Append extra entries sorted by key for deterministic output.
	if len(opts.Extra) > 0 {
		sortedKeys := slices.Sorted(maps.Keys(opts.Extra))
		for _, k := range sortedKeys {
			pairs = append(pairs, kv{k, opts.Extra[k]})
		}
	}

	return pairs
}

// writeAligned writes key-value pairs to b with all colons aligned to the
// longest key. When color is true, keys and separators are colorized.
func writeAligned(b *strings.Builder, pairs []kv, color bool) {
	maxKey := 0
	for _, p := range pairs {
		if len(p.Key) > maxKey {
			maxKey = len(p.Key)
		}
	}
	for _, p := range pairs {
		key := fmt.Sprintf("%-*s", maxKey, p.Key)
		fmt.Fprintf(b, "  %s %s %s\n",
			colorize(key, ansiYellow, color),
			colorize(":", ansiDim, color),
			p.Value,
		)
	}
}
