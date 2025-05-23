//go:build !(wasip1 && wasm)

package reducers

import (
	"fmt"

	"github.com/clockworklabs/Blackholio/server-go/tables"
)

// Non-WASM database operations (mock implementations for testing)
// These methods are implemented properly in wasm.go for WASM builds

// InsertConfig inserts a config record
func (db *DatabaseContext) InsertConfig(config *tables.Config) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// GetLoggedOutPlayer retrieves a logged out player by identity
func (db *DatabaseContext) GetLoggedOutPlayer(identity tables.Identity) (*tables.Player, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// InsertPlayer inserts a player record
func (db *DatabaseContext) InsertPlayer(player *tables.Player) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// DeleteLoggedOutPlayer deletes a logged out player by identity
func (db *DatabaseContext) DeleteLoggedOutPlayer(identity tables.Identity) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// GetPlayer retrieves a player by identity
func (db *DatabaseContext) GetPlayer(identity tables.Identity) (*tables.Player, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// GetCirclesByPlayer retrieves all circles for a player
func (db *DatabaseContext) GetCirclesByPlayer(playerID uint32) ([]*tables.Circle, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// UpdatePlayer updates a player record
func (db *DatabaseContext) UpdatePlayer(player *tables.Player) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// InsertCircle inserts a circle record
func (db *DatabaseContext) InsertCircle(circle *tables.Circle) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// UpdateCircle updates a circle record
func (db *DatabaseContext) UpdateCircle(circle *tables.Circle) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// GetEntity retrieves an entity by ID
func (db *DatabaseContext) GetEntity(entityID uint32) (*tables.Entity, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// UpdateEntity updates an entity record
func (db *DatabaseContext) UpdateEntity(entity *tables.Entity) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// GetAllCircles retrieves all circles
func (db *DatabaseContext) GetAllCircles() ([]*tables.Circle, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// GetAllEntities retrieves all entities
func (db *DatabaseContext) GetAllEntities() ([]*tables.Entity, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// GetAllPlayers retrieves all players
func (db *DatabaseContext) GetAllPlayers() ([]*tables.Player, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// GetCircle retrieves a circle by entity ID
func (db *DatabaseContext) GetCircle(entityID uint32) (*tables.Circle, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}

// GetPlayerCount retrieves the count of active players
func (db *DatabaseContext) GetPlayerCount() (uint64, error) {
	// TODO: Implement for non-WASM builds
	return 0, fmt.Errorf("not implemented for non-WASM builds")
}

// GetFoodCount retrieves the count of food entities
func (db *DatabaseContext) GetFoodCount() (uint64, error) {
	// TODO: Implement for non-WASM builds
	return 0, fmt.Errorf("not implemented for non-WASM builds")
}

// InsertFood inserts a food record
func (db *DatabaseContext) InsertFood(food *tables.Food) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// InsertLoggedOutPlayer inserts a logged out player record
func (db *DatabaseContext) InsertLoggedOutPlayer(player *tables.Player) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// DeletePlayer deletes a player by identity
func (db *DatabaseContext) DeletePlayer(identity tables.Identity) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// ScheduleReducer schedules a reducer for future execution
func (db *DatabaseContext) ScheduleReducer(name string, args []byte, schedule tables.ScheduleAt) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// InsertEntity inserts an entity record
func (db *DatabaseContext) InsertEntity(entity *tables.Entity) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// DeleteEntity deletes an entity by ID
func (db *DatabaseContext) DeleteEntity(entityID uint32) error {
	// TODO: Implement for non-WASM builds
	return fmt.Errorf("not implemented for non-WASM builds")
}

// GetConfig retrieves the game configuration from the database
func (db *DatabaseContext) GetConfig() (*tables.Config, error) {
	// TODO: Implement for non-WASM builds
	return nil, fmt.Errorf("not implemented for non-WASM builds")
}
