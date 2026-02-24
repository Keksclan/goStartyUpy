# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Built-in big ASCII font for spring-like banners using an "underscore / pipe /
  slash" style (`banner/font_big.go`). Supports A–Z, 0–9, `-`, `_`, space, and
  a `?` fallback glyph. No external dependencies.
- `renderBigText` internal renderer (`banner/render_big.go`) that normalizes
  input and produces aligned 5-row ASCII-art output.
- Font preview example (`example/font_preview/main.go`) for visual verification.
- `SpringLikeBanner` function that generates a large ASCII-art wordmark from
  `ServiceName`, inspired by Spring Boot startup banners.
- `BoxBanner` function — the previous box-style banner extracted as a named
  public function.
- `Options.BannerStyle` field: `"spring"` (default) for the new wordmark
  style, `"box"` for the legacy box style.
- `Options.BannerWidth` field: when > 0, hard-cuts every banner line to that
  maximum width.
- Tests for `SpringLikeBanner`, `BoxBanner`, `BannerStyle`, `BannerWidth`,
  name normalization, and Unicode/edge-case input safety.

### Changed

- `DefaultBanner` now returns a Spring Boot–style wordmark banner instead of
  the box banner. Use `BoxBanner` or `BannerStyle: "box"` for the old style.
- Default banner style changed from box to `"spring"`.

### Fixed

## [0.1.0] - 2026-02-24

### Added

- `banner` package with `Render` and `RenderWithChecks` functions.
- `BuildInfo` struct and `CurrentBuildInfo` for build metadata snapshots.
- `Options` struct for configuring service name, environment, banner art,
  separator, extra key-value pairs, color, and ASCII-only mode.
- Build metadata injection via `-ldflags` (`Version`, `BuildTime`, `Commit`,
  `Branch`, `Dirty`).
- Auto-generated box banner from `ServiceName` via `DefaultBanner`.
- ASCII-only fallback mode (`ASCIIOnly` option).
- Optional ANSI color output (`Color` option).
- `checks` package with `Check` interface, `Result` struct, and `Runner`.
- Built-in checks: `SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`,
  `RedisPingCheck`.
- `FuncCheck` adapter and `New` constructor for function-based checks.
- `Bool` constructor for boolean-returning check functions.
- `NewGroup` and `GroupOptions` for composite grouped checks.
- `DefaultRunner` with 2 s timeout and parallel execution.
- Parallel and sequential execution modes in `Runner`.
- Panic recovery in `Runner` and `FuncCheck`.
- `version` package with `ModuleVersion` constant.
- `Makefile` with `build`, `build-example`, `run-example`, `test`, `lint`,
  and `clean` targets.
- `scripts/ldflags.sh` helper for build metadata collection.
- Example programs: simple, custom banner, ASCII-only, checks demo,
  custom checks, and full demo.
- Unit tests for banner rendering, formatting, default banner generation,
  checks runner, function checks, and group checks.
- README with installation, quickstart, API reference, and release process.
- `.gitignore`, `.editorconfig`, `.gitattributes` for repository hygiene.

[Unreleased]: https://github.com/keksclan/goStartyUpy/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/keksclan/goStartyUpy/releases/tag/v0.1.0
