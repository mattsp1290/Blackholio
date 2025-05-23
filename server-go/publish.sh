#!/bin/bash

set -euo pipefail

# Blackholio Server Go - Deployment Script
# This script deploys the compiled WASM module to SpacetimeDB

echo "=== Blackholio Server Go - Publish Script ==="

# Check if spacetime CLI is available
if ! command -v spacetime &> /dev/null; then
    echo "Error: spacetime CLI tool not found. Please install SpacetimeDB CLI."
    exit 1
fi

# Check if WASM module exists
if [ ! -f "blackholio.wasm" ]; then
    echo "Error: blackholio.wasm not found. Run './generate.sh' first to build the WASM module."
    exit 1
fi

WASM_SIZE=$(stat -f%z blackholio.wasm 2>/dev/null || stat -c%s blackholio.wasm)
echo "Deploying WASM module: blackholio.wasm ($WASM_SIZE bytes)"

# Deploy to SpacetimeDB
echo "Publishing to SpacetimeDB..."
spacetime publish -s local blackholio --delete-data -y

echo "=== Deployment completed successfully! ==="
echo ""
echo "The Blackholio Go server is now running on SpacetimeDB."
echo "Next steps:"
echo "  1. Run './logs.sh' to view server logs"
echo "  2. Connect Unity client to test the game"
echo "  3. Use 'spacetime call' to test reducers manually" 