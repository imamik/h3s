#!/bin/bash
set -e

# This script checks test coverage for the H3S project
# It's designed to be run as a pre-push hook

echo "Checking test coverage..."

# Minimum coverage threshold (percentage)
MIN_COVERAGE=60

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Get the total coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')

# Clean up
rm coverage.out

# Check if coverage meets the threshold
if (( $(echo "$COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo "Test coverage is below the minimum threshold of ${MIN_COVERAGE}%"
    echo "Current coverage: ${COVERAGE}%"
    echo "Please add more tests to increase coverage."
    echo "You can skip this check with 'git push --no-verify' if necessary."
    exit 1
else
    echo "Test coverage is ${COVERAGE}%, which meets the minimum threshold of ${MIN_COVERAGE}%."
    exit 0
fi
