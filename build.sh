#!/bin/bash

# Ensure dependencies are up to date
go mod tidy

# Check Go version strictly for Windows 7 compliance
GO_VER=$(go env GOVERSION)
# Check if version starts with go1.21 or go1.20 (allow patch versions)
if [[ "$GO_VER" != "go1.21"* ]] && [[ "$GO_VER" != "go1.20"* ]]; then
    echo "WARNING: Your Go version is $GO_VER."
    echo "         Windows 7 requires Go 1.21 or older."
    echo "         The resulting Windows binary may NOT run on Windows 7."
    echo "         Proceeding anyway..."
fi

# Build for macOS
echo "Building for macOS..."
go build -o mfch main.go

# Build for Windows 7 (64-bit)
# Requires mingw-w64: brew install mingw-w64
if command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "Building for Windows 7 (64-bit)..."
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o mfch.exe main.go
else
    echo "Skipping Windows build: x86_64-w64-mingw32-gcc not found."
    echo "Please install mingw-w64 (e.g., 'brew install mingw-w64')."
fi