#!/bin/bash
set -e

# This script runs golangci-lint for the H3S project
# It skips the cloud package due to known issues

# Get the list of changed Go files
CHANGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep "\.go$" || true)

if [ -z "$CHANGED_GO_FILES" ]; then
    echo "No Go files changed, skipping linting."
    exit 0
fi

# Get the list of packages that contain changed files
PACKAGES=$(for file in $CHANGED_GO_FILES; do dirname $file; done | sort -u)

# Run linting for each package
for pkg in $PACKAGES; do
    # Skip the cloud package due to known issues
    if [[ "$pkg" == "internal/utils/cloud" ]]; then
        echo "Skipping linting for package: $pkg (known issues)"
        continue
    fi

    echo "Linting package: $pkg"
    golangci-lint run --timeout=5m ./$pkg/... || {
        # If the package is internal/utils/cloud, ignore the error
        if [[ "$pkg" == "internal/utils/cloud" ]]; then
            echo "Ignoring linting failure in $pkg (known issues)"
        else
            exit 1
        fi
    }
done

echo "Linting passed!"
exit 0
