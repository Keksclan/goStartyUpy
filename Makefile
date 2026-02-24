MODULE   := github.com/keksclan/goStartyUpy
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BRANCH   := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
DIRTY    := $(shell test -z "$$(git status --porcelain 2>/dev/null)" && echo "false" || echo "true")

LDFLAGS := -X '$(MODULE)/banner.Version=$(VERSION)' \
           -X '$(MODULE)/banner.BuildTime=$(BUILD_TIME)' \
           -X '$(MODULE)/banner.Commit=$(COMMIT)' \
           -X '$(MODULE)/banner.Branch=$(BRANCH)' \
           -X '$(MODULE)/banner.Dirty=$(DIRTY)'

.PHONY: build test vet fmt clean

## build: compile the example binary with embedded metadata
build:
	go build -ldflags "$(LDFLAGS)" -o bin/example ./example/

## test: run all unit tests
test:
	go test -v -count=1 ./...

## vet: run go vet
vet:
	go vet ./...

## fmt: run gofmt check
fmt:
	gofmt -l .

## clean: remove build artifacts
clean:
	rm -rf bin/
