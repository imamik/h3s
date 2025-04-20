#!/bin/bash
set -e

# This script runs golangci-lint for the entire codebase
# It uses the same configuration as the GitHub Actions workflow

echo "Running golangci-lint for the entire codebase..."

# Run golangci-lint with the same arguments as in GitHub Actions
golangci-lint run --timeout=5m

echo "Linting passed!"
exit 0
