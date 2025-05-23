// Package reducers provides the SpacetimeDB reducer system integration for Go.
// This implements the reducer framework that allows Go functions to be called
// as SpacetimeDB reducers via the WASM interface.
package reducers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/tables"
)

// ReducerContext provides the context for a reducer execution.
// This is the first parameter that every reducer function must accept.
// It provides access to the database, caller information, and other utilities.
type ReducerContext struct {
	// Sender is the Identity of the client that invoked the reducer
	Sender tables.Identity

	// Timestamp is the time at which the reducer was started
	Timestamp tables.Timestamp

	// ConnectionID is the ConnectionId of the client that invoked the reducer
	// May be nil for automatic reducers (init, client_connected, etc.)
	// Represented as 16-byte array similar to Identity
	ConnectionID *[16]byte

	// Database provides access to SpacetimeDB tables and operations
	Database *DatabaseContext

	// rng provides seeded random number generation
	rng   *rand.Rand
	rngMu sync.Mutex
}

// DatabaseContext provides access to SpacetimeDB database operations
type DatabaseContext struct {
	// Internal database handle - will be populated by WASM host calls
	handle uintptr
}

// Database operation methods are implemented in:
// - database_nonwasm.go for non-WASM builds (mock implementations)
// - wasm.go for WASM builds (real SpacetimeDB integration)

// Rng returns a random number generator seeded for this reducer execution
func (ctx *ReducerContext) Rng() *rand.Rand {
	ctx.rngMu.Lock()
	defer ctx.rngMu.Unlock()

	if ctx.rng == nil {
		// Use timestamp as seed for deterministic behavior
		seed := int64(ctx.Timestamp.Microseconds)
		ctx.rng = rand.New(rand.NewSource(seed))
	}
	return ctx.rng
}

// Identity returns the module's identity
func (ctx *ReducerContext) Identity() tables.Identity {
	// TODO: Call WASM host function to get module identity
	return tables.Identity{}
}

// ReducerResult represents the result of a reducer execution
type ReducerResult interface {
	IsSuccess() bool
	Error() string
}

// SuccessResult represents a successful reducer execution
type SuccessResult struct{}

func (SuccessResult) IsSuccess() bool { return true }
func (SuccessResult) Error() string   { return "" }

// ErrorResult represents a failed reducer execution
type ErrorResult struct {
	Message string
}

func (e ErrorResult) IsSuccess() bool { return false }
func (e ErrorResult) Error() string   { return e.Message }

// ReducerFunction represents a function that can be called as a reducer
type ReducerFunction interface {
	// Name returns the name of the reducer
	Name() string

	// Lifecycle returns the lifecycle type of the reducer (if any)
	Lifecycle() *LifecycleType

	// Invoke calls the reducer with the given context and arguments
	Invoke(ctx *ReducerContext, args []byte) ReducerResult

	// ArgumentNames returns the names of the reducer arguments
	ArgumentNames() []string
}

// LifecycleType represents the type of lifecycle reducer
type LifecycleType int

const (
	// LifecycleInit is called when the module is initially published
	LifecycleInit LifecycleType = iota

	// LifecycleClientConnected is called when a client connects
	LifecycleClientConnected

	// LifecycleClientDisconnected is called when a client disconnects
	LifecycleClientDisconnected
)

// String returns the string representation of the lifecycle type
func (l LifecycleType) String() string {
	switch l {
	case LifecycleInit:
		return "Init"
	case LifecycleClientConnected:
		return "OnConnect"
	case LifecycleClientDisconnected:
		return "OnDisconnect"
	default:
		return "Unknown"
	}
}

// ReducerRegistry manages the registration and lookup of reducers
type ReducerRegistry struct {
	reducers map[string]ReducerFunction
	byID     map[uint32]ReducerFunction
	mu       sync.RWMutex
	nextID   uint32
}

// Global reducer registry
var globalRegistry = &ReducerRegistry{
	reducers: make(map[string]ReducerFunction),
	byID:     make(map[uint32]ReducerFunction),
	nextID:   0,
}

// RegisterReducer registers a reducer function with the global registry
func RegisterReducer(reducer ReducerFunction) uint32 {
	return globalRegistry.Register(reducer)
}

// Register registers a reducer function and returns its ID
func (r *ReducerRegistry) Register(reducer ReducerFunction) uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	r.reducers[reducer.Name()] = reducer
	r.byID[id] = reducer

	return id
}

// GetByName returns a reducer by name
func (r *ReducerRegistry) GetByName(name string) (ReducerFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reducer, exists := r.reducers[name]
	return reducer, exists
}

// GetByID returns a reducer by ID
func (r *ReducerRegistry) GetByID(id uint32) (ReducerFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reducer, exists := r.byID[id]
	return reducer, exists
}

// ListReducers returns all registered reducers
func (r *ReducerRegistry) ListReducers() map[string]ReducerFunction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]ReducerFunction)
	for name, reducer := range r.reducers {
		result[name] = reducer
	}
	return result
}

// GenericReducer is a concrete implementation of ReducerFunction
type GenericReducer struct {
	name          string
	lifecycle     *LifecycleType
	argumentNames []string
	handler       func(*ReducerContext, []byte) ReducerResult
}

// NewReducer creates a new generic reducer
func NewReducer(name string, handler func(*ReducerContext, []byte) ReducerResult) *GenericReducer {
	return &GenericReducer{
		name:          name,
		handler:       handler,
		argumentNames: []string{},
	}
}

// NewLifecycleReducer creates a new lifecycle reducer
func NewLifecycleReducer(name string, lifecycle LifecycleType, handler func(*ReducerContext, []byte) ReducerResult) *GenericReducer {
	return &GenericReducer{
		name:          name,
		lifecycle:     &lifecycle,
		handler:       handler,
		argumentNames: []string{},
	}
}

// Name returns the reducer name
func (r *GenericReducer) Name() string {
	return r.name
}

// Lifecycle returns the lifecycle type
func (r *GenericReducer) Lifecycle() *LifecycleType {
	return r.lifecycle
}

// Invoke calls the reducer
func (r *GenericReducer) Invoke(ctx *ReducerContext, args []byte) ReducerResult {
	return r.handler(ctx, args)
}

// ArgumentNames returns the argument names
func (r *GenericReducer) ArgumentNames() []string {
	return r.argumentNames
}

// WithArgumentNames sets the argument names for the reducer
func (r *GenericReducer) WithArgumentNames(names []string) *GenericReducer {
	r.argumentNames = names
	return r
}

// Serialization utilities for reducer arguments

// MarshalArgs marshals reducer arguments to JSON bytes
func MarshalArgs(args interface{}) ([]byte, error) {
	return json.Marshal(args)
}

// UnmarshalArgs unmarshals JSON bytes to the given argument structure
func UnmarshalArgs(data []byte, args interface{}) error {
	return json.Unmarshal(data, args)
}

// Error handling utilities

// HandleResult converts various Go return types to ReducerResult
func HandleResult(result interface{}) ReducerResult {
	switch r := result.(type) {
	case nil:
		return SuccessResult{}
	case error:
		return ErrorResult{Message: r.Error()}
	case string:
		if r == "" {
			return SuccessResult{}
		}
		return ErrorResult{Message: r}
	case ReducerResult:
		return r
	default:
		return SuccessResult{}
	}
}

// Logging utilities for reducers

// LogInfo logs an info message from a reducer
func LogInfo(message string) {
	fmt.Printf("[INFO] %s\n", message)
}

// LogWarn logs a warning message from a reducer
func LogWarn(message string) {
	fmt.Printf("[WARN] %s\n", message)
}

// LogError logs an error message from a reducer
func LogError(message string) {
	fmt.Printf("[ERROR] %s\n", message)
}

// Utility functions for common reducer patterns

// RequirePlayer ensures that a player exists for the given context
func RequirePlayer(ctx *ReducerContext) (*tables.Player, error) {
	// TODO: Implement database query when database context is available
	// For now, return a mock implementation
	return &tables.Player{
		Identity: ctx.Sender,
		PlayerID: 1,
		Name:     "MockPlayer",
	}, nil
}

// GetConfig retrieves the game configuration
func GetConfig(ctx *ReducerContext) (*tables.Config, error) {
	// TODO: Implement database query when database context is available
	// For now, return a mock implementation
	return &tables.Config{
		ID:        0,
		WorldSize: 1000,
	}, nil
}

// ScheduleTimer schedules a timer for future execution
func ScheduleTimer(ctx *ReducerContext, scheduledAt tables.ScheduleAt) error {
	// TODO: Implement when database context is available
	LogInfo(fmt.Sprintf("Timer scheduled for: %s", scheduledAt.String()))
	return nil
}

// Performance monitoring for reducers

// PerformanceTimer tracks reducer execution time
type PerformanceTimer struct {
	Name      string
	StartTime time.Time
}

// NewPerformanceTimer creates a new performance timer
func NewPerformanceTimer(name string) *PerformanceTimer {
	return &PerformanceTimer{
		Name:      name,
		StartTime: time.Now(),
	}
}

// Stop stops the timer and logs the execution time
func (pt *PerformanceTimer) Stop() time.Duration {
	duration := time.Since(pt.StartTime)
	LogInfo(fmt.Sprintf("Performance[%s]: %v", pt.Name, duration))
	return duration
}

// Type definitions for WASM interface compatibility

// ReducerID represents a reducer identifier
type ReducerID uint32

// CallReducerParams represents parameters for calling a reducer
type CallReducerParams struct {
	ReducerID      ReducerID
	SenderIdentity [4]uint64 // 32-byte identity as 4 uint64s
	ConnectionID   [2]uint64 // 16-byte connection ID as 2 uint64s
	Timestamp      uint64    // Microseconds since Unix epoch
	Args           []byte    // BSATN or JSON encoded arguments
}

// Module metadata for introspection

// ReducerMetadata provides metadata about a reducer
type ReducerMetadata struct {
	Name          string
	Lifecycle     *LifecycleType
	ArgumentNames []string
	ArgumentTypes []string
	ReturnType    string
}

// GetReducerMetadata returns metadata for all registered reducers
func GetReducerMetadata() map[string]ReducerMetadata {
	reducers := globalRegistry.ListReducers()
	metadata := make(map[string]ReducerMetadata)

	for name, reducer := range reducers {
		metadata[name] = ReducerMetadata{
			Name:          reducer.Name(),
			Lifecycle:     reducer.Lifecycle(),
			ArgumentNames: reducer.ArgumentNames(),
			ArgumentTypes: []string{}, // TODO: Add type reflection
			ReturnType:    "ReducerResult",
		}
	}

	return metadata
}

// Debugging and development utilities

// ReducerDebugInfo provides debug information for a reducer call
type ReducerDebugInfo struct {
	ReducerName    string                 `json:"reducer_name"`
	SenderIdentity string                 `json:"sender_identity"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	Timestamp      string                 `json:"timestamp"`
	Arguments      map[string]interface{} `json:"arguments"`
	ExecutionTime  string                 `json:"execution_time"`
	Success        bool                   `json:"success"`
	Error          string                 `json:"error,omitempty"`
}

// CreateDebugInfo creates debug information for a reducer call
func CreateDebugInfo(ctx *ReducerContext, reducerName string, args []byte, result ReducerResult, executionTime time.Duration) ReducerDebugInfo {
	var argsMap map[string]interface{}
	json.Unmarshal(args, &argsMap)

	debugInfo := ReducerDebugInfo{
		ReducerName:    reducerName,
		SenderIdentity: ctx.Sender.String(),
		Timestamp:      ctx.Timestamp.String(),
		Arguments:      argsMap,
		ExecutionTime:  executionTime.String(),
		Success:        result.IsSuccess(),
	}

	if ctx.ConnectionID != nil {
		debugInfo.ConnectionID = fmt.Sprintf("%x", *ctx.ConnectionID)
	}

	if !result.IsSuccess() {
		debugInfo.Error = result.Error()
	}

	return debugInfo
}

// Constants for reducer system configuration

const (
	// MaxReducerExecutionTime is the maximum time a reducer can execute
	MaxReducerExecutionTime = 30 * time.Second

	// MaxArgumentSize is the maximum size of reducer arguments
	MaxArgumentSize = 1024 * 1024 // 1MB

	// DefaultTimeoutDuration is the default timeout for reducer operations
	DefaultTimeoutDuration = 10 * time.Second
)

// Error types for reducer system

// ReducerError represents an error in the reducer system
type ReducerError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error returns the error message
func (e ReducerError) Error() string {
	return fmt.Sprintf("ReducerError[%s]: %s", e.Code, e.Message)
}

// NewReducerError creates a new reducer error
func NewReducerError(code, message string, details map[string]interface{}) ReducerError {
	return ReducerError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common reducer error codes
const (
	ErrorCodeReducerNotFound  = "REDUCER_NOT_FOUND"
	ErrorCodeInvalidArguments = "INVALID_ARGUMENTS"
	ErrorCodeExecutionTimeout = "EXECUTION_TIMEOUT"
	ErrorCodeInternalError    = "INTERNAL_ERROR"
	ErrorCodeUnauthorized     = "UNAUTHORIZED"
	ErrorCodeInvalidState     = "INVALID_STATE"
)
