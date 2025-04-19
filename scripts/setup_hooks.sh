#!/bin/sh
set -e

echo "=== Installing Git Hooks and Dependencies ==="

echo "Installing lefthook..."
go install github.com/evilmartians/lefthook@latest

echo "Installing golangci-lint..."
if ! command -v golangci-lint >/dev/null 2>&1; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "Installing gitleaks..."
if ! command -v gitleaks >/dev/null 2>&1; then
    go install github.com/gitleaks/gitleaks/v8@latest
fi

echo "Installing ineffassign..."
if ! command -v ineffassign >/dev/null 2>&1; then
    go install github.com/gordonklaus/ineffassign@latest
fi

echo "Setting up lefthook git hooks..."
$(go env GOPATH)/bin/lefthook install

echo "=== Git Hooks Setup Complete ==="
echo "Pre-commit, commit-msg, and push hooks are now active."
echo "You can skip hooks with 'git commit --no-verify' when necessary."
