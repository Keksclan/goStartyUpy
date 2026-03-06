# Banner System

This document explains how goStartyUpy generates and renders startup banners.

## Banner Styles

goStartyUpy ships with six built-in banner styles. The style is selected via `Options.BannerStyle` (default: `"spring"`). When `Options.Banner` is set to a non-empty string, the custom art is used verbatim and `BannerStyle` is ignored.

### `spring` (default)

A large ASCII-art wordmark rendered with the **big font** (5-row glyphs using underscores, pipes, and slashes). The service name is normalized to uppercase and rendered as multi-line ASCII art.

```
  ___  ____  ____  _____ ____
 / _ \|  _ \|  _ \| ____|  _ \
/ /_\ \ |_) | | | |  _| | |_) |
|  _  |  _ <| |_| | |___|  _ <
|_| |_|_| \_\____/|_____|_| \_\
```

### `classic`

A Spring Boot–inspired wordmark using the **classic font** (7-row glyphs with slashes, backslashes, and underscores). Includes two tagline lines below the wordmark:

- **Tagline 1**: `<ServiceName> <Version>` (overridable via `Options.Tagline1`)
- **Tagline 2**: `Build: <BuildTime>  Commit: <Commit>` (overridable via `Options.Tagline2`)

The `ShowDetails` option controls whether the key-value info section appears below the classic banner.

### `box`

A bordered box drawn around the service name:

```
┌──────────────────────────────┐
│        ORDER-SERVICE         │
└──────────────────────────────┘
```

With `ASCIIOnly: true`:

```
+------------------------------+
|        ORDER-SERVICE         |
+------------------------------+
```

### `mini`

A compact 3-row wordmark using the **mini font**:

```
 _  ___
| |/ _ \
|_|\___/
```

### `block`

A 5-row wordmark using hash-based block glyphs:

```
  ###   ####   ####
 #   # #    # #
 ##### ####   #
 #   # #    # #
 #   #  ####   ####
```

### Custom Banner

When `Options.Banner` is set, the string is used as-is. No font rendering or generation occurs:

```go
opts := banner.Options{
    ServiceName: "my-service",
    Banner: `
   ╔═══════════════════════════╗
   ║    MY CUSTOM BANNER       ║
   ╚═══════════════════════════╝`,
}
```

## Font System

Each banner style that renders text uses a font — a mapping from runes to multi-line glyph arrays.

### Font Files

| File | Font Name | Glyph Height | Style |
|------|-----------|-------------|-------|
| `font_big.go` | `bigFont` | 5 rows | Underscore/pipe/slash |
| `font_classic.go` | `classicFont` | 7 rows | Slash/backslash/underscore |
| `font_mini.go` | `miniFont` | 3 rows | Minimal compact |
| `font_block.go` | `blockFont` | 5 rows | Hash-based blocks |

### Supported Characters

All fonts support the same character set:

- **Letters**: A–Z (input is uppercased automatically)
- **Digits**: 0–9
- **Symbols**: `-`, `_`, ` ` (space), `?` (fallback for unsupported characters)

### Glyph Structure

Each glyph is a fixed-size array of strings, one per row:

```go
// Big font (5 rows):
type glyph [fontHeight]string  // fontHeight = 5

// Classic font (7 rows):
type glyph7 [classicHeight]string  // classicHeight = 7

// Mini font (3 rows):
type glyph3 [miniHeight]string  // miniHeight = 3

// Block font (5 rows):
type glyph5 [blockHeight]string  // blockHeight = 5
```

### Text Rendering

The rendering functions (`renderBigText`, `renderClassicText`, etc.) follow the same algorithm:

1. **Normalize** the input: uppercase, collapse whitespace, replace unsupported characters with `?`.
2. **Look up** each character's glyph in the font map.
3. **Concatenate** glyphs side by side, row by row, with single-space separation between characters.
4. **Trim** trailing whitespace from each row.

## Name Normalization

The `normalizeName()` function prepares the service name for banner rendering:

1. Convert to uppercase.
2. Replace whitespace with dashes.
3. Keep only `A-Z`, `0-9`, `-`, `_`.
4. Drop all other characters.

Example: `"my-awesome-service"` → `"MY-AWESOME-SERVICE"`

## Environment Awareness

The banner can display the deployment environment (e.g., "production", "staging") in two ways:

1. **Explicit**: Set `Options.Environment` directly in code.
2. **Automatic**: Leave `Options.Environment` empty; the renderer reads `GO_STARTYUPY_ENV` at render time.

When the environment is detected from the environment variable, `EnvironmentFromEnv` is set to `true`, which may affect how the environment is displayed in certain styles.

## Width Clamping

When `Options.BannerWidth > 0`, every line of the generated banner is hard-cut to that width. This prevents wide banners from wrapping in narrow terminals:

```go
opts := banner.Options{
    ServiceName: "my-very-long-service-name",
    BannerWidth: 60,
}
```

## ASCII-Only Mode

Setting `Options.ASCIIOnly = true` replaces all Unicode box-drawing characters with plain ASCII equivalents:

| Unicode | ASCII Replacement |
|---------|------------------|
| `┌` `┐` `└` `┘` | `+` |
| `─` | `-` |
| `│` | `\|` |
| `════...` (separator) | `====...` |

This ensures correct rendering in terminals or log systems that do not support Unicode.

## ANSI Color Support

Colors are controlled by `Options.Color`:

- **`false` (default)**: Plain text output with no escape sequences. Safe for log files and non-terminal outputs.
- **`true`**: ANSI escape sequences are applied to banner art (cyan+bold), separators (dim), check tags (green/red), and other elements.

The `colorize()` helper wraps text in the appropriate ANSI codes only when color is enabled. When disabled, text passes through unchanged.
