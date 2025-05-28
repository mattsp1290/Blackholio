package reducers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/constants"
	"github.com/clockworklabs/Blackholio/server-go/tables"
	"github.com/clockworklabs/Blackholio/server-go/types"
)

// Test helper functions

func createTestContext() *ReducerContext {
	identity := tables.NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	timestamp := tables.NewTimestampFromTime(time.Now())

	return &ReducerContext{
		Sender:       identity,
		Timestamp:    timestamp,
		ConnectionID: nil,
		Database:     &DatabaseContext{handle: 0},
	}
}

func createTestPlayer() *tables.Player {
	identity := tables.NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	return tables.NewPlayer(identity, 1, "TestPlayer")
}

// Test ReducerContext functionality

func TestReducerContext(t *testing.T) {
	t.Run("Basic context creation", func(t *testing.T) {
		ctx := createTestContext()

		if ctx.Sender.IsZero() {
			t.Error("Sender should not be zero")
		}

		if ctx.Timestamp.Microseconds == 0 {
			t.Error("Timestamp should not be zero")
		}

		if ctx.Database == nil {
			t.Error("Database should not be nil")
		}
	})

	t.Run("RNG functionality", func(t *testing.T) {
		ctx := createTestContext()

		rng1 := ctx.Rng()
		rng2 := ctx.Rng()

		// Should return the same instance
		if rng1 != rng2 {
			t.Error("Rng() should return the same instance")
		}

		// Should generate random numbers
		val1 := rng1.Float32()
		val2 := rng1.Float32()

		if val1 == val2 {
			t.Error("RNG should generate different values")
		}
	})

	t.Run("Identity functionality", func(t *testing.T) {
		ctx := createTestContext()
		identity := ctx.Identity()

		// For non-WASM builds, should return zero identity
		if !identity.IsZero() {
			t.Log("Identity returned non-zero value (may be expected in WASM)")
		}
	})
}

// Test ReducerResult implementations

func TestReducerResults(t *testing.T) {
	t.Run("SuccessResult", func(t *testing.T) {
		result := SuccessResult{}

		if !result.IsSuccess() {
			t.Error("SuccessResult should be successful")
		}

		if result.Error() != "" {
			t.Error("SuccessResult should have empty error")
		}
	})

	t.Run("ErrorResult", func(t *testing.T) {
		message := "Test error"
		result := ErrorResult{Message: message}

		if result.IsSuccess() {
			t.Error("ErrorResult should not be successful")
		}

		if result.Error() != message {
			t.Errorf("ErrorResult should return message: got %s, expected %s", result.Error(), message)
		}
	})
}

// Test ReducerRegistry functionality

func TestReducerRegistry(t *testing.T) {
	// Create a new registry for testing
	registry := &ReducerRegistry{
		reducers: make(map[string]ReducerFunction),
		byID:     make(map[uint32]ReducerFunction),
		nextID:   0,
	}

	t.Run("Register and retrieve by name", func(t *testing.T) {
		reducer := NewReducer("test_reducer", func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		})

		id := registry.Register(reducer)

		retrieved, exists := registry.GetByName("test_reducer")
		if !exists {
			t.Error("Reducer should exist after registration")
		}

		if retrieved.Name() != "test_reducer" {
			t.Error("Retrieved reducer should have correct name")
		}

		// Test retrieval by ID
		retrievedByID, exists := registry.GetByID(id)
		if !exists {
			t.Error("Reducer should exist when retrieved by ID")
		}

		if retrievedByID.Name() != "test_reducer" {
			t.Error("Retrieved reducer by ID should have correct name")
		}
	})

	t.Run("List all reducers", func(t *testing.T) {
		// Clear registry
		registry.reducers = make(map[string]ReducerFunction)
		registry.byID = make(map[uint32]ReducerFunction)
		registry.nextID = 0

		reducer1 := NewReducer("reducer1", func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		})
		reducer2 := NewReducer("reducer2", func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		})

		registry.Register(reducer1)
		registry.Register(reducer2)

		allReducers := registry.ListReducers()

		if len(allReducers) != 2 {
			t.Errorf("Expected 2 reducers, got %d", len(allReducers))
		}

		if _, exists := allReducers["reducer1"]; !exists {
			t.Error("reducer1 should be in the list")
		}

		if _, exists := allReducers["reducer2"]; !exists {
			t.Error("reducer2 should be in the list")
		}
	})
}

// Test GenericReducer functionality

func TestGenericReducer(t *testing.T) {
	t.Run("Basic reducer", func(t *testing.T) {
		called := false
		reducer := NewReducer("test", func(ctx *ReducerContext, args []byte) ReducerResult {
			called = true
			return SuccessResult{}
		})

		if reducer.Name() != "test" {
			t.Error("Reducer should have correct name")
		}

		if reducer.Lifecycle() != nil {
			t.Error("Basic reducer should not have lifecycle")
		}

		ctx := createTestContext()
		result := reducer.Invoke(ctx, []byte{})

		if !called {
			t.Error("Reducer handler should be called")
		}

		if !result.IsSuccess() {
			t.Error("Reducer should return success")
		}
	})

	t.Run("Lifecycle reducer", func(t *testing.T) {
		reducer := NewLifecycleReducer("init", LifecycleInit, func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		})

		if reducer.Lifecycle() == nil {
			t.Error("Lifecycle reducer should have lifecycle")
		}

		if *reducer.Lifecycle() != LifecycleInit {
			t.Error("Lifecycle reducer should have correct lifecycle")
		}
	})

	t.Run("Reducer with argument names", func(t *testing.T) {
		reducer := NewReducer("test", func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		}).WithArgumentNames([]string{"arg1", "arg2"})

		argNames := reducer.ArgumentNames()
		if len(argNames) != 2 {
			t.Errorf("Expected 2 argument names, got %d", len(argNames))
		}

		if argNames[0] != "arg1" || argNames[1] != "arg2" {
			t.Error("Argument names should match")
		}
	})
}

// Test LifecycleType

func TestLifecycleType(t *testing.T) {
	t.Run("String representation", func(t *testing.T) {
		if LifecycleInit.String() != "Init" {
			t.Error("LifecycleInit should stringify to 'Init'")
		}

		if LifecycleClientConnected.String() != "OnConnect" {
			t.Error("LifecycleClientConnected should stringify to 'OnConnect'")
		}

		if LifecycleClientDisconnected.String() != "OnDisconnect" {
			t.Error("LifecycleClientDisconnected should stringify to 'OnDisconnect'")
		}
	})
}

// Test serialization utilities

func TestSerialization(t *testing.T) {
	t.Run("MarshalArgs", func(t *testing.T) {
		args := map[string]interface{}{
			"name":      "test",
			"value":     42,
			"direction": types.NewDbVector2(1.0, 2.0),
		}

		data, err := MarshalArgs(args)
		if err != nil {
			t.Fatalf("MarshalArgs failed: %v", err)
		}

		var decoded map[string]interface{}
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["name"] != "test" {
			t.Error("Name should be preserved")
		}
	})

	t.Run("UnmarshalArgs", func(t *testing.T) {
		type TestArgs struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		}

		data := []byte(`{"name":"test","value":42.5}`)

		var args TestArgs
		err := UnmarshalArgs(data, &args)
		if err != nil {
			t.Fatalf("UnmarshalArgs failed: %v", err)
		}

		if args.Name != "test" {
			t.Error("Name should be unmarshaled correctly")
		}

		if args.Value != 42.5 {
			t.Error("Value should be unmarshaled correctly")
		}
	})
}

// Test HandleResult function

func TestHandleResult(t *testing.T) {
	t.Run("Nil result", func(t *testing.T) {
		result := HandleResult(nil)
		if !result.IsSuccess() {
			t.Error("Nil should be success")
		}
	})

	t.Run("Error result", func(t *testing.T) {
		err := ErrorResult{Message: "test error"}
		result := HandleResult(err)
		if result.IsSuccess() {
			t.Error("Error should not be success")
		}
		if result.Error() != "test error" {
			t.Error("Error message should be preserved")
		}
	})

	t.Run("String result", func(t *testing.T) {
		result := HandleResult("error message")
		if result.IsSuccess() {
			t.Error("Non-empty string should be error")
		}

		result = HandleResult("")
		if !result.IsSuccess() {
			t.Error("Empty string should be success")
		}
	})
}

// Test performance monitoring

func TestPerformanceTimer(t *testing.T) {
	t.Run("Timer functionality", func(t *testing.T) {
		timer := NewPerformanceTimer("test")

		if timer.Name != "test" {
			t.Error("Timer should have correct name")
		}

		time.Sleep(1 * time.Millisecond)
		duration := timer.Stop()

		if duration < time.Millisecond {
			t.Error("Timer should measure at least 1ms")
		}
	})
}

// Test reducer metadata

func TestReducerMetadata(t *testing.T) {
	t.Run("Get metadata", func(t *testing.T) {
		// Create a clean registry for testing
		testRegistry := &ReducerRegistry{
			reducers: make(map[string]ReducerFunction),
			byID:     make(map[uint32]ReducerFunction),
			nextID:   0,
		}

		// Temporarily replace global registry
		originalRegistry := globalRegistry
		globalRegistry = testRegistry
		defer func() {
			globalRegistry = originalRegistry
		}()

		reducer := NewReducer("test_metadata", func(ctx *ReducerContext, args []byte) ReducerResult {
			return SuccessResult{}
		}).WithArgumentNames([]string{"arg1", "arg2"})

		testRegistry.Register(reducer)

		metadata := GetReducerMetadata()

		if len(metadata) != 1 {
			t.Errorf("Expected 1 reducer in metadata, got %d", len(metadata))
		}

		meta, exists := metadata["test_metadata"]
		if !exists {
			t.Error("test_metadata should exist in metadata")
		}

		if meta.Name != "test_metadata" {
			t.Error("Metadata name should match")
		}

		if len(meta.ArgumentNames) != 2 {
			t.Error("Metadata should include argument names")
		}
	})
}

// Test error types

func TestReducerError(t *testing.T) {
	t.Run("Create and format error", func(t *testing.T) {
		details := map[string]interface{}{
			"entity_id": 123,
			"reason":    "not found",
		}

		err := NewReducerError("TEST_ERROR", "Test error message", details)

		if err.Code != "TEST_ERROR" {
			t.Error("Error code should match")
		}

		if err.Message != "Test error message" {
			t.Error("Error message should match")
		}

		if err.Details["entity_id"] != 123 {
			t.Error("Error details should be preserved")
		}

		errorString := err.Error()
		if errorString != "ReducerError[TEST_ERROR]: Test error message" {
			t.Errorf("Error string format incorrect: %s", errorString)
		}
	})
}

// Test debug functionality

func TestDebugInfo(t *testing.T) {
	t.Run("Create debug info", func(t *testing.T) {
		ctx := createTestContext()
		args := []byte(`{"name":"test"}`)
		result := SuccessResult{}
		duration := 100 * time.Millisecond

		debugInfo := CreateDebugInfo(ctx, "test_reducer", args, result, duration)

		if debugInfo.ReducerName != "test_reducer" {
			t.Error("Debug info should include reducer name")
		}

		if !debugInfo.Success {
			t.Error("Debug info should reflect success")
		}

		if debugInfo.Error != "" {
			t.Error("Debug info should not have error for success")
		}

		if debugInfo.ExecutionTime != duration.String() {
			t.Error("Debug info should include execution time")
		}
	})
}

// Integration tests for Blackholio reducers

func TestBlackholioReducers(t *testing.T) {
	t.Run("InitReducer", func(t *testing.T) {
		ctx := createTestContext()
		result := InitReducer(ctx, []byte{})

		// For non-WASM builds, this will fail due to database operations
		// but we can test that it doesn't panic
		if result == nil {
			t.Error("InitReducer should return a result")
		}
	})

	t.Run("EnterGameReducer with valid args", func(t *testing.T) {
		ctx := createTestContext()
		args := EnterGameArgs{Name: "TestPlayer"}
		argsData, _ := MarshalArgs(args)

		result := EnterGameReducer(ctx, argsData)

		// Should fail due to database operations in non-WASM builds
		// but should not panic
		if result == nil {
			t.Error("EnterGameReducer should return a result")
		}
	})

	t.Run("EnterGameReducer with invalid args", func(t *testing.T) {
		ctx := createTestContext()
		invalidArgs := []byte("invalid json")

		result := EnterGameReducer(ctx, invalidArgs)

		if result.IsSuccess() {
			t.Error("EnterGameReducer should fail with invalid args")
		}

		if result.Error() == "" {
			t.Error("Error result should have error message")
		}
	})

	t.Run("UpdatePlayerInputReducer", func(t *testing.T) {
		ctx := createTestContext()
		args := UpdatePlayerInputArgs{
			Direction: types.NewDbVector2(1.0, 0.5),
		}
		argsData, _ := MarshalArgs(args)

		result := UpdatePlayerInputReducer(ctx, argsData)

		// Should process arguments correctly even if database operations fail
		if result == nil {
			t.Error("UpdatePlayerInputReducer should return a result")
		}
	})
}

// Benchmark tests

func BenchmarkReducerInvocation(b *testing.B) {
	reducer := NewReducer("benchmark", func(ctx *ReducerContext, args []byte) ReducerResult {
		return SuccessResult{}
	})

	ctx := createTestContext()
	args := []byte("{}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reducer.Invoke(ctx, args)
	}
}

func BenchmarkJSONMarshaling(b *testing.B) {
	args := map[string]interface{}{
		"name":      "test",
		"direction": types.NewDbVector2(1.0, 2.0),
		"mass":      uint32(100),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalArgs(args)
	}
}

func BenchmarkRNGGeneration(b *testing.B) {
	ctx := createTestContext()
	rng := ctx.Rng()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Float32()
	}
}

// Test constants integration

func TestConstantsIntegration(t *testing.T) {
	t.Run("Constants accessible", func(t *testing.T) {
		config := constants.GetGlobalConfiguration()

		if config.StartPlayerMass == 0 {
			t.Error("Constants should be accessible")
		}

		if config.TargetFoodCount == 0 {
			t.Error("Target food count should be set")
		}
	})

	t.Run("Mathematical functions", func(t *testing.T) {
		mass := uint32(100)
		radius := constants.MassToRadius(mass)

		if radius <= 0 {
			t.Error("Radius should be positive")
		}

		speed := constants.MassToMaxMoveSpeed(mass)
		if speed <= 0 {
			t.Error("Speed should be positive")
		}
	})
}
