#!/bin/bash

set -euo pipefail

# Blackholio Server Go - Code Generation and WASM Compilation Script
# This script compiles the Go server to WASM and generates client bindings

echo "=== Blackholio Server Go - Generate Script ==="

# Check if spacetime CLI is available
if ! command -v spacetime &> /dev/null; then
    echo "Error: spacetime CLI tool not found. Please install SpacetimeDB CLI."
    exit 1
fi

# Check if Go is available and version is correct
if ! command -v go &> /dev/null; then
    echo "Error: Go not found. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
if [[ $(echo "$GO_VERSION 1.21" | tr ' ' '\n' | sort -V | head -n1) != "1.21" ]]; then
    echo "Warning: Go version $GO_VERSION detected. Go 1.21+ recommended for WASM compilation."
fi

# Clean previous builds
echo "Cleaning previous builds..."
rm -f *.wasm
rm -f main

# Set WASM build environment
export GOOS=wasip1
export GOARCH=wasm

# Build WASM module
echo "Building Go WASM module..."
go build -o blackholio.wasm .

if [ ! -f "blackholio.wasm" ]; then
    echo "Error: WASM compilation failed. blackholio.wasm not found."
    exit 1
fi

echo "WASM module built successfully: blackholio.wasm ($(stat -f%z blackholio.wasm 2>/dev/null || stat -c%s blackholio.wasm) bytes)"

# Generate client bindings
echo "Generating client bindings..."
spacetime generate --out-dir ../client-unity/Assets/Scripts/autogen --lang cs "$@"

echo "=== Generation completed successfully! ==="
echo "- WASM module: blackholio.wasm"
echo "- Client bindings: ../client-unity/Assets/Scripts/autogen"
echo ""
echo "Next steps:"
echo "  1. Run './publish.sh' to deploy to SpacetimeDB"
echo "  2. Run './logs.sh' to view deployment logs" 