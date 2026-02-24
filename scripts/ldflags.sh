#!/bin/sh
# ldflags.sh — prints Go ldflags that inject build metadata into the banner
# package. POSIX sh compatible; no bash-only features.
#
# Usage:
#   LDFLAGS="$(./scripts/ldflags.sh)" go build -ldflags "$LDFLAGS" ./cmd/myservice
#
# You can override MODULE if your import path differs:
#   MODULE=github.com/my/repo ./scripts/ldflags.sh

set -e

MODULE="${MODULE:-github.com/keksclan/goStartyUpy}"
PKG="${MODULE}/banner"

# --- Collect values ---

COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"

if git diff --quiet 2>/dev/null; then
    DIRTY="false"
else
    DIRTY="true"
fi

# RFC 3339 timestamp in local time.
# Try GNU/BSD `date -Iseconds` first; fall back to explicit format string.
BUILDTIME="$(date -Iseconds 2>/dev/null || date +"%Y-%m-%dT%H:%M:%S%z" 2>/dev/null || echo unknown)"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"

# --- Print ldflags ---

printf "%s" \
  "-X '${PKG}.Version=${VERSION}'" \
  " -X '${PKG}.BuildTime=${BUILDTIME}'" \
  " -X '${PKG}.Commit=${COMMIT}'" \
  " -X '${PKG}.Branch=${BRANCH}'" \
  " -X '${PKG}.Dirty=${DIRTY}'"
