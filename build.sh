#!/bin/bash

# Build script for MFCH (Native Windows Version)

echo "Preparing build for Windows..."

# Check Go version for Windows 7 compatibility advice
GO_VER=$(go env GOVERSION)
echo "Current Go version: $GO_VER"
if [[ "$GO_VER" == "go1.21"* ]] || [[ "$GO_VER" > "go1.21" ]]; then
    echo "WARNING: Go 1.21+ does not officially support Windows 7."
    echo "         For Windows 7 support, please use Go 1.20 or older."
fi

# Clean previous builds
rm -f mfch.exe mfch

# Build for Windows
# CGO_ENABLED=0: Use pure Go implementation (no MinGW needed)
# -H windowsgui: Hide console window
# -s -w: Strip debug info for smaller binary
echo "Building 'mfch.exe' (Windows)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -s -w" -o mfch.exe main.go

if [ $? -eq 0 ]; then
    echo "SUCCESS: mfch.exe created."
    ls -lh mfch.exe
else
    echo "FAILURE: Build failed."
    exit 1
fi