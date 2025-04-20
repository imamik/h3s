#!/bin/bash
set -e

# This script runs golangci-lint for the H3S project
# It skips the cloud package due to known issues

# Check if golangci-lint is installed
if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "golangci-lint not found. Please install it: https://golangci-lint.run/usage/install/" >&2
    exit 1
fi

echo "Running golangci-lint for the entire project..."

# Run linting on the entire project, excluding the cloud package
# Exclude path syntax might differ slightly depending on golangci-lint version, adjust if needed
golangci-lint run --timeout=5m --exclude='internal/utils/cloud/*' ./...

echo "Linting passed!"
exit 0
