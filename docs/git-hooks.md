# Git Hooks for H3S

This document provides detailed information about the Git hooks used in the H3S project.

## Overview

Git hooks are scripts that Git executes before or after events such as commit, push, and receive. We use [lefthook](https://github.com/evilmartians/lefthook) to manage our Git hooks.

## Installation

To install the hooks and required dependencies, run:

```sh
./scripts/setup_hooks.sh
```

This script will:
1. Install lefthook
2. Install golangci-lint
3. Install gitleaks
4. Install ineffassign
5. Set up the Git hooks

## Hook Types

### Pre-commit Hooks

Pre-commit hooks run before a commit is created. They help ensure that the code being committed meets the project's quality standards.

#### go-fmt

Formats Go code using `go fmt`.

```yaml
go-fmt:
  glob: "*.go"
  run: go fmt ./...
```

#### go-test-unit

Runs unit tests for changed packages only. This hook uses a custom script that:
1. Identifies which Go files have changed
2. Determines which packages those files belong to
3. Runs unit tests only for those packages

```yaml
go-test-unit:
  glob: "*.go"
  run: ./scripts/hooks/run_unit_tests.sh
```

#### golangci-lint

Runs linting checks using golangci-lint.

```yaml
golangci-lint:
  glob: "*.go"
  run: golangci-lint run
```

#### gitleaks

Scans for secrets in staged changes.

```yaml
gitleaks:
  run: |
    if ! command -v gitleaks >/dev/null 2>&1; then
      echo "gitleaks not installed. Please install it for secret scanning (https://github.com/gitleaks/gitleaks)." >&2
      echo "Run: go install github.com/gitleaks/gitleaks/v8@latest"
      exit 0  # Don't fail if not installed
    fi
    gitleaks protect --staged --no-git -v
```

#### ineffassign

Checks for ineffective assignments in Go code.

```yaml
ineffassign:
  glob: "*.go"
  run: |
    if ! command -v ineffassign >/dev/null 2>&1; then
      echo "ineffassign not installed. Run: go install github.com/gordonklaus/ineffassign@latest" >&2
      exit 0  # Don't fail if not installed
    fi
    ineffassign ./...
```

### Commit Message Hooks

Commit message hooks validate the format of commit messages.

#### validate-conventional-commit

Ensures commit messages follow the [Conventional Commits](https://www.conventionalcommits.org/) format.

```yaml
validate-conventional-commit:
  run: |
    if ! grep -Eq '^(feat|fix|docs|style|refactor|test|chore|ci|build|perf)\([a-zA-Z0-9_-]+\): .+' "$1"; then
      echo "Error: Commit message must follow Conventional Commits format." >&2
      echo "Format: type(scope): description" >&2
      echo "Example: feat(api): add new endpoint for user creation" >&2
      echo "Valid types: feat, fix, docs, style, refactor, test, chore, ci, build, perf" >&2
      exit 1
    fi
```

### Pre-push Hooks

Pre-push hooks run before pushing commits to a remote repository.

#### run-all-tests

Runs all tests, including integration tests.

```yaml
run-all-tests:
  run: go test ./...
```

#### check-coverage

Ensures test coverage meets the minimum threshold.

```yaml
check-coverage:
  run: ./scripts/hooks/check_coverage.sh
```

## Custom Scripts

### run_unit_tests.sh

This script runs unit tests for changed packages only. It's designed to be fast and only run unit tests, not integration or e2e tests.

```sh
#!/bin/bash
set -e

# This script runs unit tests for the H3S project
# It's designed to be fast and only run unit tests, not integration or e2e tests

echo "Running unit tests..."

# Get the list of changed Go files
CHANGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep "\.go$" || true)

if [ -z "$CHANGED_GO_FILES" ]; then
    echo "No Go files changed, skipping unit tests."
    exit 0
fi

# Get the list of packages that contain changed files
PACKAGES=$(for file in $CHANGED_GO_FILES; do dirname $file; done | sort -u)

# Run tests for each package
for pkg in $PACKAGES; do
    if [ -f "$pkg/$(basename $pkg)_test.go" ] || ls $pkg/*_test.go >/dev/null 2>&1; then
        echo "Testing package: $pkg"
        go test -short -tags=unit ./$pkg
    fi
done

echo "Unit tests passed!"
exit 0
```

### check_coverage.sh

This script checks test coverage for the H3S project. It's designed to be run as a pre-push hook.

```sh
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
```

## Skipping Hooks

You can skip hooks when necessary:

```sh
# Skip pre-commit hooks
git commit --no-verify

# Skip pre-push hooks
git push --no-verify
```

## Troubleshooting

### Hook Not Running

If a hook is not running, check:
1. That lefthook is installed: `which lefthook`
2. That the hooks are installed: `lefthook list`
3. That the hook is enabled in `lefthook.yml`

### Hook Failing

If a hook is failing:
1. Read the error message carefully
2. Check that all required dependencies are installed
3. Try running the hook command manually to debug

### Reinstalling Hooks

If you need to reinstall the hooks:

```sh
./scripts/setup_hooks.sh
```

## Adding New Hooks

To add a new hook:

1. Add the hook configuration to `lefthook.yml`
2. If needed, create a custom script in `scripts/hooks/`
3. Make the script executable: `chmod +x scripts/hooks/your_script.sh`
4. Update documentation in `README.md` and this file

## Best Practices

1. **Keep hooks fast**: Pre-commit hooks should run quickly to avoid disrupting the development workflow
2. **Make hooks reliable**: Hooks should fail only when there's a real issue
3. **Provide clear error messages**: When a hook fails, it should provide clear guidance on how to fix the issue
4. **Allow bypassing**: Always provide a way to bypass hooks when necessary
5. **Document hooks**: Keep documentation up-to-date with the actual hook behavior
