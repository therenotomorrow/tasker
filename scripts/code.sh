#!/usr/bin/env sh

set -e

BIN="$PWD"/bin
mkdir -p "$BIN"

# ---- installation

GOLANGCI_LINT_PATH="$BIN"/golangci-lint
GOLANGCI_LINT_VERSION=v2.1.5
GOLANGCI_LINT="$GOLANGCI_LINT_PATH"@"$GOLANGCI_LINT_VERSION"

if ! test -f "$GOLANGCI_LINT"; then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$BIN" "$GOLANGCI_LINT_VERSION"
    mv "$GOLANGCI_LINT_PATH" "$GOLANGCI_LINT"
fi

# ---- processing

"$GOLANGCI_LINT" run ./...

go build -o "$BIN"/tasker "$PWD"/cmd/tasker/main.go
