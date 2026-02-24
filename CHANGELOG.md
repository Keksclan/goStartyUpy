# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.1.0] – 2026-02-24

### Added
- `banner` package: `Render`, `RenderWithChecks`, `BuildInfo`, `Options`.
- `checks` package: `Runner`, `Check` interface, `Result` struct.
- Built-in checks: `SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck`.
- Example program in `example/main.go`.
- Unit tests for banner rendering/formatting and checks runner.
