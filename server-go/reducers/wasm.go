//go:build js && wasm

package reducers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/tables"
)

// Simplified WASM implementation for Go 1.23 compatibility
// This focuses on compilation success and basic functionality

// WASM Exported Functions

//go:wasmexport __call_reducer__
func callReducer(reducerId uint32) int16 {
	fmt.Printf("[WASM] Calling reducer ID: %d\n", reducerId)

	reducer, exists := globalRegistry.GetByID(reducerId)
	if !exists {
		fmt.Printf("[WASM] Reducer not found: %d\n", reducerId)
		return 1
	}

	// Create mock context for compilation
	ctx := &ReducerContext{
		Sender:    tables.Identity{},
		Timestamp: tables.NewTimestampFromTime(time.Now()),
		Database:  &DatabaseContext{handle: 0},
	}

	// Execute reducer with empty args for now
	result := reducer.Invoke(ctx, []byte{})
	if !result.IsSuccess() {
		fmt.Printf("[WASM] Reducer error: %s\n", result.Error())
		return 1
	}

	fmt.Printf("[WASM] Reducer %s executed successfully\n", reducer.Name())
	return 0
}

//go:wasmexport __get_module_info__
func getModuleInfo() int16 {
	metadata := GetReducerMetadata()
	infoBytes, err := json.Marshal(metadata)
	if err != nil {
		fmt.Printf("[WASM] Failed to marshal module info: %v\n", err)
		return 1
	}

	fmt.Printf("[WASM] Module info: %s\n", string(infoBytes))
	return 0
}

//go:wasmexport __describe_module_def__
func describeModuleDef() int16 {
	moduleDef := map[string]interface{}{
		"tables":   tables.TableDefinitions,
		"reducers": GetReducerMetadata(),
		"version":  "1.0.0",
		"name":     "blackholio-server-go",
	}

	defBytes, err := json.Marshal(moduleDef)
	if err != nil {
		fmt.Printf("[WASM] Failed to marshal module def: %v\n", err)
		return 1
	}

	fmt.Printf("[WASM] Module definition: %s\n", string(defBytes))
	return 0
}

// Simple database operations (mocked for WASM compilation)

func (db *DatabaseContext) InsertConfig(config *tables.Config) error {
	fmt.Printf("[WASM] Mock InsertConfig: %+v\n", config)
	return nil
}

func (db *DatabaseContext) GetLoggedOutPlayer(identity tables.Identity) (*tables.Player, error) {
	fmt.Printf("[WASM] Mock GetLoggedOutPlayer: %s\n", identity.String())
	return nil, fmt.Errorf("mock: player not found")
}

func (db *DatabaseContext) InsertPlayer(player *tables.Player) error {
	fmt.Printf("[WASM] Mock InsertPlayer: %+v\n", player)
	return nil
}

func (db *DatabaseContext) DeleteLoggedOutPlayer(identity tables.Identity) error {
	fmt.Printf("[WASM] Mock DeleteLoggedOutPlayer: %s\n", identity.String())
	return nil
}

func (db *DatabaseContext) GetPlayer(identity tables.Identity) (*tables.Player, error) {
	fmt.Printf("[WASM] Mock GetPlayer: %s\n", identity.String())
	return &tables.Player{Identity: identity, PlayerID: 1, Name: "MockPlayer"}, nil
}

func (db *DatabaseContext) GetCirclesByPlayer(playerID uint32) ([]*tables.Circle, error) {
	fmt.Printf("[WASM] Mock GetCirclesByPlayer: %d\n", playerID)
	return []*tables.Circle{}, nil
}

func (db *DatabaseContext) UpdatePlayer(player *tables.Player) error {
	fmt.Printf("[WASM] Mock UpdatePlayer: %+v\n", player)
	return nil
}

func (db *DatabaseContext) InsertCircle(circle *tables.Circle) error {
	fmt.Printf("[WASM] Mock InsertCircle: %+v\n", circle)
	return nil
}

func (db *DatabaseContext) UpdateCircle(circle *tables.Circle) error {
	fmt.Printf("[WASM] Mock UpdateCircle: %+v\n", circle)
	return nil
}

func (db *DatabaseContext) GetEntity(entityID uint32) (*tables.Entity, error) {
	fmt.Printf("[WASM] Mock GetEntity: %d\n", entityID)
	return nil, fmt.Errorf("mock: entity not found")
}

func (db *DatabaseContext) UpdateEntity(entity *tables.Entity) error {
	fmt.Printf("[WASM] Mock UpdateEntity: %+v\n", entity)
	return nil
}

func (db *DatabaseContext) GetAllCircles() ([]*tables.Circle, error) {
	fmt.Printf("[WASM] Mock GetAllCircles\n")
	return []*tables.Circle{}, nil
}

func (db *DatabaseContext) GetAllEntities() ([]*tables.Entity, error) {
	fmt.Printf("[WASM] Mock GetAllEntities\n")
	return []*tables.Entity{}, nil
}

func (db *DatabaseContext) GetAllPlayers() ([]*tables.Player, error) {
	fmt.Printf("[WASM] Mock GetAllPlayers\n")
	return []*tables.Player{}, nil
}

func (db *DatabaseContext) GetCircle(entityID uint32) (*tables.Circle, error) {
	fmt.Printf("[WASM] Mock GetCircle: %d\n", entityID)
	return nil, fmt.Errorf("mock: circle not found")
}

func (db *DatabaseContext) GetPlayerCount() (uint64, error) {
	fmt.Printf("[WASM] Mock GetPlayerCount\n")
	return 0, nil
}

func (db *DatabaseContext) GetFoodCount() (uint64, error) {
	fmt.Printf("[WASM] Mock GetFoodCount\n")
	return 0, nil
}

func (db *DatabaseContext) InsertFood(food *tables.Food) error {
	fmt.Printf("[WASM] Mock InsertFood: %+v\n", food)
	return nil
}

func (db *DatabaseContext) InsertLoggedOutPlayer(player *tables.Player) error {
	fmt.Printf("[WASM] Mock InsertLoggedOutPlayer: %+v\n", player)
	return nil
}

func (db *DatabaseContext) DeletePlayer(identity tables.Identity) error {
	fmt.Printf("[WASM] Mock DeletePlayer: %s\n", identity.String())
	return nil
}

func (db *DatabaseContext) ScheduleReducer(name string, args []byte, schedule tables.ScheduleAt) error {
	fmt.Printf("[WASM] Mock ScheduleReducer: %s, schedule: %s\n", name, schedule.String())
	return nil
}

func (db *DatabaseContext) InsertEntity(entity *tables.Entity) error {
	fmt.Printf("[WASM] Mock InsertEntity: %+v\n", entity)
	return nil
}

func (db *DatabaseContext) DeleteEntity(entityID uint32) error {
	fmt.Printf("[WASM] Mock DeleteEntity: %d\n", entityID)
	return nil
}

func (db *DatabaseContext) GetConfig() (*tables.Config, error) {
	fmt.Printf("[WASM] Mock GetConfig\n")
	return &tables.Config{ID: 0, WorldSize: 1000}, nil
}

func init() {
	fmt.Println("[WASM] Simplified WASM implementation initialized")
}
