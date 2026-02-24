MODULE   := github.com/keksclan/goStartyUpy

# --- Build metadata (collected automatically) ---
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH   := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
# RFC 3339 timestamp in local time. Try `date -Iseconds` first (GNU/BSD),
# fall back to explicit format string for broader portability.
BUILDTIME := $(shell date -Iseconds 2>/dev/null || date +"%Y-%m-%dT%H:%M:%S%z" 2>/dev/null || echo "unknown")
DIRTY    := $(shell git diff --quiet 2>/dev/null && echo "false" || echo "true")

LDFLAGS := -X '$(MODULE)/banner.Version=$(VERSION)' \
           -X '$(MODULE)/banner.BuildTime=$(BUILDTIME)' \
           -X '$(MODULE)/banner.Commit=$(COMMIT)' \
           -X '$(MODULE)/banner.Branch=$(BRANCH)' \
           -X '$(MODULE)/banner.Dirty=$(DIRTY)'

.PHONY: build build-example run-example test lint clean

## build: compile a binary from PKG (override PKG and BIN as needed)
##   make build PKG=./cmd/myservice BIN=bin/myservice
PKG ?= ./example/
BIN ?= bin/example
build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)

## build-example: compile the example binary with embedded metadata
build-example:
	go build -ldflags "$(LDFLAGS)" -o bin/example ./example/

## run-example: build and run the example binary
run-example: build-example
	./bin/example

## test: run all unit tests
test:
	go test -v -count=1 ./...

## lint: run go vet and check gofmt
lint:
	go vet ./...
	@test -z "$$(gofmt -l .)" || { echo "gofmt: the following files need formatting:"; gofmt -l .; exit 1; }

## clean: remove build artifacts
clean:
	rm -rf bin/
