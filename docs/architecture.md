# Architecture

This document describes the high-level architecture of the goStartyUpy library.

## Module Layout

```
goStartyUpy/
├── banner/          # Core package – banner rendering, options, build metadata
│   ├── banner.go         # Render() and RenderWithChecks() entry points
│   ├── buildinfo.go      # BuildInfo struct and link-time variables
│   ├── options.go        # Options struct (all banner configuration)
│   ├── default_banner.go # Banner style resolution and line clamping
│   ├── springlike.go     # "spring" style generator (default)
│   ├── classiclike.go    # "classic" style generator (Spring Boot–like)
│   ├── render_big.go     # Big-font text renderer (used by spring/box)
│   ├── render_classic.go # Classic-font text renderer (used by classic)
│   ├── font_big.go       # 5-row ASCII glyphs (A–Z, 0–9, symbols)
│   ├── font_block.go     # 5-row block glyphs (hash-style)
│   ├── font_classic.go   # 7-row classic glyphs (slash/backslash style)
│   ├── font_mini.go      # 3-row mini glyphs
│   └── color.go          # ANSI escape helpers
├── checks/          # Health-check framework
│   ├── checks.go         # Check interface, Result, Runner
│   ├── func.go           # FuncCheck – wrap any function as a check
│   ├── group.go          # Group – bundle multiple checks into one
│   ├── db_sql.go         # SQLPingCheck – database/sql PingContext
│   ├── tcp.go            # TCPDialCheck – raw TCP dial
│   ├── http.go           # HTTPGetCheck – HTTP GET with status range
│   └── redis_like.go     # RedisPingCheck – inline RESP PING
├── version/         # Module version constant
│   └── version.go
├── example/         # Runnable example programs
├── scripts/         # Build helper scripts (ldflags.sh)
└── docs/            # Documentation (this directory)
```

## Design Principles

1. **Zero Dependencies** — The entire library uses only the Go standard library. The `go.sum` file will never grow because of goStartyUpy.

2. **Deterministic Output** — Banner rendering is a pure function of its inputs (`Options` + `BuildInfo` + check results). There is no randomness, no global state mutation, and no I/O beyond reading the `GO_STARTYUPY_ENV` environment variable.

3. **Panic Safety** — The check `Runner` recovers from panics inside individual checks and reports them as failed `Result` values. The banner renderer never panics.

4. **Composability** — Checks implement a simple `Check` interface (`Name() string`, `Run(ctx) Result`). Built-in checks, function-based checks, boolean checks, and grouped checks can all be freely combined.

## Package Responsibilities

### `banner`

The `banner` package is the primary entry point. It owns:

- **Rendering** — `Render()` and `RenderWithChecks()` produce a complete startup message as a string. The caller decides where to print it (`fmt.Print`, a logger, etc.).
- **Build Metadata** — Package-level variables (`Version`, `BuildTime`, `Commit`, `Branch`, `Dirty`) are set at link time via `-ldflags`. `CurrentBuildInfo()` captures them as a snapshot.
- **Banner Generation** — Six styles are supported: `spring` (default), `classic`, `box`, `mini`, `block`, and raw custom art. Style selection happens in `resolveBanner()`.
- **Formatting** — Key-value pairs are aligned with consistent padding. ANSI colors are applied only when `Color: true`.

### `checks`

The `checks` package provides a standalone health-check framework:

- **Runner** — Executes checks sequentially or in parallel with per-check timeouts. Results are returned in input order regardless of execution mode.
- **Built-in Checks** — `SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck` cover common infrastructure probes without external drivers.
- **Custom Checks** — `New()` wraps any `func(context.Context) error`. `Bool()` wraps a boolean probe. `NewGroup()` bundles checks into a single composite result.

### `version`

Exposes `ModuleVersion` — the semantic version of the goStartyUpy library itself (not the consuming service). Updated with each tagged release.

## Data Flow

```
Options + BuildInfo
       │
       ▼
  resolveBanner()     ← selects style, generates ASCII art
       │
       ▼
  buildKVs()          ← collects key-value metadata pairs
       │
       ▼
  writeAligned()      ← formats pairs with consistent padding
       │
       ▼
  (optional) checks   ← Runner.Run() produces []Result
       │
       ▼
  RenderWithChecks()  ← assembles banner + separator + info + checks
       │
       ▼
  string              ← complete startup message
```

## Threading Model

- Banner rendering is single-threaded and safe to call from any goroutine.
- The check `Runner` with `Parallel: true` spawns one goroutine per check. Results are written to pre-allocated slots with no shared mutation, so no additional synchronization is needed beyond the `WaitGroup`.
