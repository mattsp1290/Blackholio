//go:build wasip1 && wasm

package main

import (
	"fmt"

	// Import all our packages to ensure they're properly initialized
	_ "github.com/clockworklabs/Blackholio/server-go/constants"
	_ "github.com/clockworklabs/Blackholio/server-go/logic"
	_ "github.com/clockworklabs/Blackholio/server-go/reducers"
	_ "github.com/clockworklabs/Blackholio/server-go/tables"
	_ "github.com/clockworklabs/Blackholio/server-go/types"
)

// This file provides the WASM entry point for SpacetimeDB.
// When compiled with GOOS=wasm GOARCH=wasm, this will be the main entry point
// and will export the required functions for SpacetimeDB integration.

// main is the entry point for the WASM module
func main() {
	fmt.Println("Blackholio Server Go - WASM Module Initialized")

	// The WASM runtime will handle calling the exported functions
	// We don't need to do anything else here - just keep the program alive
	// The actual reducer calls will come through the exported functions
	// defined in reducers/wasm.go

	select {} // Block forever - WASM runtime will call our exported functions
}

// init ensures all packages are properly initialized
func init() {
	fmt.Println("Blackholio Server Go - Package initialization complete")
}
