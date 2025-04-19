#!/bin/bash
# Generates a dynamic coverage badge and updates the README.md badge

set -e
COVERAGE_FILE=coverage.out
README=README.md
BADGE_PATH="https://img.shields.io/badge/coverage-"
STYLE="?style=flat"

# Generate coverage.out if it doesn't exist
if [ ! -f "$COVERAGE_FILE" ]; then
    go test -coverprofile=$COVERAGE_FILE ./...
fi

# Extract the total coverage percentage
COVERAGE=$(go tool cover -func=$COVERAGE_FILE | grep total: | awk '{print $3}' | sed 's/%//')
COVERAGE_INT=$(printf '%.0f' "$COVERAGE")

# Choose badge color
if (( $(echo "$COVERAGE >= 90" | bc -l) )); then
    COLOR=brightgreen
elif (( $(echo "$COVERAGE >= 80" | bc -l) )); then
    COLOR=yellowgreen
elif (( $(echo "$COVERAGE >= 60" | bc -l) )); then
    COLOR=orange
else
    COLOR=red
fi

BADGE_URL="${BADGE_PATH}${COVERAGE}%25-${COLOR}${STYLE}"

# Update the badge in README.md
sed -i '' -E "s|!\[Coverage\]\(https://img.shields.io/badge/coverage-[^)]*\)|![Coverage](${BADGE_URL})|" "$README"

echo "Generated coverage badge: $BADGE_URL and updated $README"
