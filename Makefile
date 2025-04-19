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

.PHONY: coverage
coverage: ## Generate code coverage report (set UPLOAD=1 to upload to Codecov)
	$(GOTEST) -v -coverprofile=coverage.out ./...
	@echo "Generated coverage.out"
ifneq ($(UPLOAD),)
	@echo "Uploading coverage report to Codecov..."
	bash <(curl -s https://codecov.io/bash) -f coverage.out
endif

.PHONY: coverage-html
coverage-html: ## Generate HTML code coverage report
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser to view the report."

.PHONY: coverage-threshold
coverage-threshold: ## Fail if coverage is below 80%
	$(GOTEST) -coverprofile=coverage.out ./...
	@total=$(GOCMD) tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'; \
	if [ "$(shell echo "$$total < 80.0" | bc)" -eq 1 ]; then \
		echo "Code coverage ($$total%) is below threshold (80%)"; \
		exit 1; \
	else \
		echo "Code coverage ($$total%) meets threshold."; \
	fi

.PHONY: coverage-summary
coverage-summary: ## Print concise coverage summary to console
	$(GOCMD) tool cover -func=coverage.out | grep -E '^[^ ]+\s+[0-9.]+%' || echo 'No coverage.out found. Run make coverage first.'

.PHONY: dev
dev: ## Setup development environment (installs pre-commit hooks and tools)
	./scripts/setup_hooks.sh

.PHONY: diagram
diagram: ## Generate architecture diagrams using PlantUML
	plantuml docs/architecture.puml -o docs/

# Test targets
.PHONY: unit-test
unit-test: ## Run unit tests only
	$(GOTEST) -v -short -cover -tags=unit ./...

.PHONY: integration-test
integration-test: ## Run integration tests only
	$(GOTEST) -v -cover -tags=integration ./...

.PHONY: e2e-test
e2e-test: ## Run end-to-end (e2e) tests only (optionally: NAME=TestName)
	@if [ -z "$(NAME)" ]; then \
		$(GOTEST) -v -cover -tags=e2e ./... ; \
	else \
		$(GOTEST) -v -cover -tags=e2e -run $(NAME) ./... ; \
	fi

.PHONY: bench
bench: ## Run benchmarks (optionally: NAME=BenchmarkName)
	@if [ -z "$(NAME)" ]; then \
		$(GOTEST) -bench=. -benchmem ./... ; \
	else \
		$(GOTEST) -bench=$(NAME) -benchmem ./... ; \
	fi

.PHONY: test-pkg
test-pkg: ## Run tests for a specific package (PKG=...)
	$(GOTEST) -v -cover $(PKG)

.PHONY: test-name
test-name: ## Run tests matching a name (NAME=...)
	$(GOTEST) -v -cover -run $(NAME) ./...

.PHONY: test-summary
test-summary: ## Run tests and print a summary
	$(GOTEST) -v -cover ./... | tee test-summary.txt
	@echo "\n==== Test Summary ===="
	@grep -E '^(ok|FAIL|PASS|SKIP)' test-summary.txt

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'