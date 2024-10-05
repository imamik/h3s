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

# Primary targets
all: test build

# Build the project
build:
    # execute the clean command
	make clean
	# build production binary with ldflags
	$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	# wait for the file to be written
	@while [ ! -f $(BINARY_NAME) ]; do \
		sleep 0.1; \
	done
	@echo "File size: $(shell du -sh $(BINARY_NAME))"

# Run all tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_linux
	rm -f $(BINARY_NAME)_windows
	rm -f $(BINARY_NAME)_darwin

# Build and run the project
run: build
	./$(BINARY_NAME)

# Install dependencies and tidy go.mod
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

# Cross-compilation targets
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux -v $(MAIN_PACKAGE)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v $(MAIN_PACKAGE)

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin -v $(MAIN_PACKAGE)

# Build for all supported platforms
build-all: test clean build-linux build-windows build-darwin

# Code formatting
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run

# Declare phony targets (targets that don't represent files)
.PHONY: all build test clean run deps build-linux build-windows build-darwin build-all fmt lint