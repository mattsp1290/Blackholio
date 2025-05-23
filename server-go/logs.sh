#!/bin/bash

set -euo pipefail

# Blackholio Server Go - Logs Script
# This script displays logs from the deployed SpacetimeDB module

echo "=== Blackholio Server Go - Logs ==="

# Check if spacetime CLI is available
if ! command -v spacetime &> /dev/null; then
    echo "Error: spacetime CLI tool not found. Please install SpacetimeDB CLI."
    exit 1
fi

echo "Displaying logs for Blackholio Go server..."
echo "Press Ctrl+C to exit log streaming"
echo ""

spacetime logs -s local blackholio 