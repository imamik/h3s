# Makefile for h3s project

# Go command and parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=h3s

# Main package path
MAIN_PACKAGE=.

# Build variables
VERSION ?= dev
LDFLAGS := -ldflags="\
	-X h3s/internal/version.Version=$(VERSION) \
	-X h3s/internal/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown') \
	-X 'h3s/internal/version.GoVersion=$(shell go version | cut -d" " -f3)' \
	-s -w"

# Primary targets
.PHONY: all
all: test build ## Run tests and build binary

.PHONY: build
build: ## Build binary for current platform
	@mkdir -p bin
	CGO_ENABLED=0 $(GOBUILD) -trimpath $(LDFLAGS) -o bin/$(BINARY_NAME) -v $(MAIN_PACKAGE)
	@echo "Build complete: bin/$(BINARY_NAME)"

.PHONY: test
test: ## Run tests
	$(GOTEST) -v -coverprofile=coverage.txt ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run --timeout 5m

.PHONY: lint\:fix
lint\:fix: ## Run linter with auto-fix enabled
	golangci-lint run --timeout 5m --fix

.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf bin/ dist/ coverage.txt

.PHONY: run
run: build ## Build and run the project
	./bin/$(BINARY_NAME)

.PHONY: deps
deps: ## Install dependencies
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: build-all
build-all: ## Build binaries for all platforms using GoReleaser
	goreleaser build --snapshot --rm-dist

.PHONY: release
release: ## Create release (requires GITHUB_TOKEN)
	goreleaser release --rm-dist

.PHONY: snapshot
snapshot: ## Create snapshot release
	goreleaser release --snapshot --rm-dist

.PHONY: fmt
fmt: ## Format code
	$(GOCMD) fmt ./...

.PHONY: fix_fieldalignment
fix_fieldalignment: ## Fix struct field alignment
	./fix_fieldalignment.sh

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'