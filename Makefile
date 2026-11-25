# Moses development and release targets.
BINARY      ?= moses
GO          ?= go
GOLANGCI    ?= golangci-lint
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_TIME  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_TIME)

.DEFAULT_GOAL := help

.PHONY: all build run test test-race bench lint lint-fix lint-config format format-check mod-check verify clean help

all: build

build: ## Build the moses binary with version metadata.
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/moses

run: ## Run moses from source.
	$(GO) run ./cmd/moses

test: ## Run the test suite.
	$(GO) test ./...

test-race: ## Run the test suite with the race detector.
	$(GO) test -race -shuffle=on -count=1 ./...

bench: ## Run the pipeline benchmarks.
	$(GO) test -run '^$$' -bench 'Benchmark(KnapsackPack|Pack|PackBestFit|Format|ParseFile|Convert10510Gen)$$' -benchmem ./pkg/binpack ./internal/nesting ./internal/mgf ./internal/source/txt ./internal/app

lint: ## Run go vet and golangci-lint with the project configuration.
	$(GO) vet ./...
	$(GOLANGCI) run ./...

lint-fix: ## Auto-fix lint findings where possible.
	$(GOLANGCI) run --fix ./...

lint-config: ## Validate the golangci-lint configuration.
	$(GOLANGCI) config verify

format: ## Format Go sources.
	gofmt -w cmd internal pkg

format-check: ## Fail when Go sources are not formatted.
	@test -z "$$(gofmt -l cmd internal pkg)" || (gofmt -l cmd internal pkg; exit 1)

mod-check: ## Fail when go.mod/go.sum are not tidy.
	@before="$$(cksum go.mod go.sum)"; \
	$(GO) mod tidy; \
	after="$$(cksum go.mod go.sum)"; \
	if [ "$$before" != "$$after" ]; then \
		echo "go.mod/go.sum changed by go mod tidy; commit the tidy result"; \
		exit 1; \
	fi

verify: format-check mod-check lint test-race ## Run the local CI gate.

clean: ## Remove local build artifacts.
	rm -f $(BINARY)
	rm -rf build

help: ## Show this help.
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build           Build the moses binary with version metadata"
	@echo "  run             Run moses from source"
	@echo "  test            Run the test suite"
	@echo "  test-race       Run tests with the race detector"
	@echo "  bench           Run the pipeline benchmarks"
	@echo "  lint            Run go vet and golangci-lint"
	@echo "  lint-fix        Auto-fix lint findings"
	@echo "  lint-config     Validate the golangci-lint configuration"
	@echo "  format          Format Go sources"
	@echo "  format-check    Fail when Go sources are not formatted"
	@echo "  mod-check       Fail when go.mod/go.sum are not tidy"
	@echo "  verify          Run format, mod, lint and race tests"
	@echo "  clean           Remove local build artifacts"
	@echo "  help            Show this help"
