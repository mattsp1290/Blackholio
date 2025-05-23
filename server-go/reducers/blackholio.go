// Package reducers - Blackholio game reducer implementations
// This file implements all the Blackholio game reducers matching the Rust and C# versions

package reducers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/constants"
	"github.com/clockworklabs/Blackholio/server-go/logic"
	"github.com/clockworklabs/Blackholio/server-go/tables"
	"github.com/clockworklabs/Blackholio/server-go/types"
)

// Blackholio Reducer Implementations
// These reducers match the functionality in server-rust/src/lib.rs and server-csharp/Lib.cs

// Use universal reducer patterns from the bindings-go package

// InitReducer handles module initialization
// Matches: Rust init() and C# Init()
func InitReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("Init")
	defer timer.Stop()

	LogInfo("Initializing Blackholio game module...")

	// Initialize configuration
	config := tables.NewConfig(0, constants.DEFAULT_WORLD_SIZE)
	if err := ctx.Database.InsertConfig(config); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to insert config: %v", err)}
	}

	// Schedule periodic timers
	moveInterval := tables.NewTimeDurationFromDuration(constants.MOVE_PLAYERS_INTERVAL)
	spawnInterval := tables.NewTimeDurationFromDuration(constants.SPAWN_FOOD_INTERVAL)
	decayInterval := tables.NewTimeDurationFromDuration(constants.CIRCLE_DECAY_INTERVAL)

	// Schedule move all players timer
	moveSchedule := tables.NewScheduleAtInterval(moveInterval)
	if err := ctx.Database.ScheduleReducer("MoveAllPlayers", []byte{}, moveSchedule); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to schedule move timer: %v", err)}
	}

	// Schedule food spawn timer
	spawnSchedule := tables.NewScheduleAtInterval(spawnInterval)
	if err := ctx.Database.ScheduleReducer("SpawnFood", []byte{}, spawnSchedule); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to schedule spawn timer: %v", err)}
	}

	// Schedule circle decay timer
	decaySchedule := tables.NewScheduleAtInterval(decayInterval)
	if err := ctx.Database.ScheduleReducer("CircleDecay", []byte{}, decaySchedule); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to schedule decay timer: %v", err)}
	}

	LogInfo("Blackholio game module initialized successfully")
	return SuccessResult{}
}

// ConnectReducer handles client connection
// Matches: Rust connect() and C# Connect()
func ConnectReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("Connect")
	defer timer.Stop()

	LogInfo(fmt.Sprintf("Client connecting: %s", ctx.Sender.String()))

	// Check if player was logged out and restore them
	loggedOutPlayer, err := ctx.Database.GetLoggedOutPlayer(ctx.Sender)
	if err == nil && loggedOutPlayer != nil {
		// Move from logged_out_player to player table
		if err := ctx.Database.InsertPlayer(loggedOutPlayer); err != nil {
			return ErrorResult{Message: fmt.Sprintf("Failed to restore player: %v", err)}
		}
		if err := ctx.Database.DeleteLoggedOutPlayer(ctx.Sender); err != nil {
			LogWarn(fmt.Sprintf("Failed to remove logged out player: %v", err))
		}
	} else {
		// Create new player
		player := tables.NewPlayer(ctx.Sender, 0, "")
		if err := ctx.Database.InsertPlayer(player); err != nil {
			return ErrorResult{Message: fmt.Sprintf("Failed to create player: %v", err)}
		}
	}

	LogInfo(fmt.Sprintf("Client connected successfully: %s", ctx.Sender.String()))
	return SuccessResult{}
}

// DisconnectReducer handles client disconnection
// Matches: Rust disconnect() and C# Disconnect()
func DisconnectReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("Disconnect")
	defer timer.Stop()

	LogInfo(fmt.Sprintf("Client disconnecting: %s", ctx.Sender.String()))

	// Get player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	// Remove all player circles from the arena
	circles, err := ctx.Database.GetCirclesByPlayer(player.PlayerID)
	if err != nil {
		LogWarn(fmt.Sprintf("Failed to get player circles: %v", err))
	} else {
		for _, circle := range circles {
			if err := logic.DestroyEntity(ctx.Database.DeleteEntity, circle.EntityID); err != nil {
				LogWarn(fmt.Sprintf("Failed to destroy circle entity %d: %v", circle.EntityID, err))
			}
		}
	}

	// Move player to logged_out_player table
	if err := ctx.Database.InsertLoggedOutPlayer(player); err != nil {
		LogWarn(fmt.Sprintf("Failed to save logged out player: %v", err))
	}

	// Remove from active player table
	if err := ctx.Database.DeletePlayer(ctx.Sender); err != nil {
		LogWarn(fmt.Sprintf("Failed to remove active player: %v", err))
	}

	LogInfo(fmt.Sprintf("Client disconnected: %s", ctx.Sender.String()))
	return SuccessResult{}
}

// EnterGameArgs represents the arguments for EnterGame reducer
type EnterGameArgs struct {
	Name string `json:"name"`
}

// EnterGameReducer handles player entering the game with a name
// Matches: Rust enter_game() and C# EnterGame()
func EnterGameReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("EnterGame")
	defer timer.Stop()

	var gameArgs EnterGameArgs
	if err := UnmarshalArgs(args, &gameArgs); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Invalid arguments: %v", err)}
	}

	LogInfo(fmt.Sprintf("Player entering game: %s with name '%s'", ctx.Sender.String(), gameArgs.Name))

	// Get and update player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	player.Name = gameArgs.Name
	if err := ctx.Database.UpdatePlayer(player); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to update player: %v", err)}
	}

	// Spawn initial circle
	config, err := GetConfig(ctx)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get config: %v", err)}
	}

	rng := ctx.Rng()
	entity, circle, err := logic.SpawnPlayerInitialCircle(player.PlayerID, config.WorldSize, rng, ctx.Timestamp)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to spawn initial circle: %v", err)}
	}

	if err := ctx.Database.InsertEntity(entity); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to insert entity: %v", err)}
	}

	if err := ctx.Database.InsertCircle(circle); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to insert circle: %v", err)}
	}

	LogInfo(fmt.Sprintf("Player '%s' entered game successfully", gameArgs.Name))
	return SuccessResult{}
}

// RespawnReducer handles player respawn
// Matches: Rust respawn() and C# Respawn()
func RespawnReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("Respawn")
	defer timer.Stop()

	LogInfo(fmt.Sprintf("Player respawning: %s", ctx.Sender.String()))

	// Get player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	// Spawn initial circle
	config, err := GetConfig(ctx)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get config: %v", err)}
	}

	rng := ctx.Rng()
	entity, circle, err := logic.SpawnPlayerInitialCircle(player.PlayerID, config.WorldSize, rng, ctx.Timestamp)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to spawn respawn circle: %v", err)}
	}

	if err := ctx.Database.InsertEntity(entity); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to insert entity: %v", err)}
	}

	if err := ctx.Database.InsertCircle(circle); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to insert circle: %v", err)}
	}

	LogInfo(fmt.Sprintf("Player respawned successfully: %s", ctx.Sender.String()))
	return SuccessResult{}
}

// SuicideReducer handles player suicide (destroying all circles)
// Matches: Rust suicide() and C# Suicide()
func SuicideReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("Suicide")
	defer timer.Stop()

	LogInfo(fmt.Sprintf("Player committing suicide: %s", ctx.Sender.String()))

	// Get player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	// Destroy all player circles
	circles, err := ctx.Database.GetCirclesByPlayer(player.PlayerID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get player circles: %v", err)}
	}

	for _, circle := range circles {
		if err := logic.DestroyEntity(ctx.Database.DeleteEntity, circle.EntityID); err != nil {
			LogWarn(fmt.Sprintf("Failed to destroy circle entity %d: %v", circle.EntityID, err))
		}
	}

	LogInfo(fmt.Sprintf("Player suicide completed: %s", ctx.Sender.String()))
	return SuccessResult{}
}

// UpdatePlayerInputArgs represents the arguments for UpdatePlayerInput reducer
type UpdatePlayerInputArgs struct {
	Direction types.DbVector2 `json:"direction"`
}

// UpdatePlayerInputReducer handles player input updates
// Matches: Rust update_player_input() and C# UpdatePlayerInput()
func UpdatePlayerInputReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("UpdatePlayerInput")
	defer timer.Stop()

	var inputArgs UpdatePlayerInputArgs
	if err := UnmarshalArgs(args, &inputArgs); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Invalid arguments: %v", err)}
	}

	// Get player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	// Update all player circles
	circles, err := ctx.Database.GetCirclesByPlayer(player.PlayerID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get player circles: %v", err)}
	}

	for _, circle := range circles {
		circle.Direction = inputArgs.Direction.Normalized()
		circle.Speed = Clamp(inputArgs.Direction.Magnitude(), 0.0, 1.0)

		if err := ctx.Database.UpdateCircle(circle); err != nil {
			LogWarn(fmt.Sprintf("Failed to update circle %d: %v", circle.EntityID, err))
		}
	}

	return SuccessResult{}
}

// PlayerSplitReducer handles player circle splitting
// Matches: Rust player_split() and C# PlayerSplit()
func PlayerSplitReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("PlayerSplit")
	defer timer.Stop()

	LogInfo(fmt.Sprintf("Player attempting split: %s", ctx.Sender.String()))

	// Get player
	player, err := ctx.Database.GetPlayer(ctx.Sender)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Player not found: %v", err)}
	}

	// Get current circles
	circles, err := ctx.Database.GetCirclesByPlayer(player.PlayerID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get player circles: %v", err)}
	}

	circleCount := uint32(len(circles))
	config := constants.GetGlobalConfiguration()

	if circleCount >= config.MaxCirclesPerPlayer {
		return SuccessResult{} // Can't split anymore
	}

	// Attempt to split circles
	for _, circle := range circles {
		entity, err := ctx.Database.GetEntity(circle.EntityID)
		if err != nil {
			LogWarn(fmt.Sprintf("Failed to get entity for circle %d: %v", circle.EntityID, err))
			continue
		}

		if logic.CanPlayerSplit(entity, circleCount) {
			halfMass := logic.CalculateHalfMass(entity.Mass)

			// Create new circle
			newPosition := entity.Position.Add(circle.Direction)
			newEntity, newCircle, err := logic.SpawnCircleAt(player.PlayerID, halfMass, newPosition, ctx.Timestamp)
			if err != nil {
				LogWarn(fmt.Sprintf("Failed to spawn split circle: %v", err))
				continue
			}

			// Insert new entities
			if err := ctx.Database.InsertEntity(newEntity); err != nil {
				LogWarn(fmt.Sprintf("Failed to insert new entity: %v", err))
				continue
			}

			if err := ctx.Database.InsertCircle(newCircle); err != nil {
				LogWarn(fmt.Sprintf("Failed to insert new circle: %v", err))
				continue
			}

			// Update original circle
			entity.Mass -= halfMass
			circle.LastSplitTime = ctx.Timestamp

			if err := ctx.Database.UpdateEntity(entity); err != nil {
				LogWarn(fmt.Sprintf("Failed to update original entity: %v", err))
			}

			if err := ctx.Database.UpdateCircle(circle); err != nil {
				LogWarn(fmt.Sprintf("Failed to update original circle: %v", err))
			}

			circleCount++
			if circleCount >= config.MaxCirclesPerPlayer {
				break
			}
		}
	}

	// Schedule recombine timer
	recombineDelay := tables.NewTimeDurationFromDuration(time.Duration(config.SplitRecombineDelaySec) * time.Second)
	recombineTime := ctx.Timestamp.Add(recombineDelay)
	recombineSchedule := tables.NewScheduleAtTime(recombineTime)

	recombineArgs, _ := json.Marshal(map[string]interface{}{
		"player_id": player.PlayerID,
	})

	if err := ctx.Database.ScheduleReducer("CircleRecombine", recombineArgs, recombineSchedule); err != nil {
		LogWarn(fmt.Sprintf("Failed to schedule recombine timer: %v", err))
	}

	LogWarn("Player split!")
	return SuccessResult{}
}

// MoveAllPlayersReducer handles moving all players (main game tick)
// Matches: Rust move_all_players() and C# MoveAllPlayers()
func MoveAllPlayersReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("MoveAllPlayers")
	defer timer.Stop()

	// Get world configuration
	config, err := GetConfig(ctx)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get config: %v", err)}
	}

	// Get all circles and entities
	allCircles, err := ctx.Database.GetAllCircles()
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get circles: %v", err)}
	}

	allEntities, err := ctx.Database.GetAllEntities()
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get entities: %v", err)}
	}

	// Create lookup maps
	entityMap := make(map[uint32]*tables.Entity)
	for _, entity := range allEntities {
		entityMap[entity.EntityID] = entity
	}

	// Calculate movement directions for all circles
	circleDirections := make(map[uint32]types.DbVector2)
	for _, circle := range allCircles {
		circleDirections[circle.EntityID] = circle.Direction.Mul(circle.Speed)
	}

	// Handle split circle physics for each player
	players, err := ctx.Database.GetAllPlayers()
	if err != nil {
		LogWarn(fmt.Sprintf("Failed to get players: %v", err))
	} else {
		for _, player := range players {
			playerCircles, err := ctx.Database.GetCirclesByPlayer(player.PlayerID)
			if err != nil {
				continue
			}

			if len(playerCircles) <= 1 {
				continue // No split circle physics needed
			}

			// Apply gravitational and separation forces
			for i, circleA := range playerCircles {
				entityA := entityMap[circleA.EntityID]
				if entityA == nil {
					continue
				}

				for j := i + 1; j < len(playerCircles); j++ {
					circleB := playerCircles[j]
					entityB := entityMap[circleB.EntityID]
					if entityB == nil {
						continue
					}

					// Calculate forces
					gravityForce := logic.CalculateGravityPull(entityA, entityB,
						float32(ctx.Timestamp.Sub(circleA.LastSplitTime).ToDuration().Seconds()),
						len(playerCircles))

					separationForce := logic.CalculateSeparationForce(entityA, entityB)

					// Apply forces
					forceA := gravityForce.Add(separationForce).Div(2.0)
					forceB := gravityForce.Mul(-1).Add(separationForce.Mul(-1)).Div(2.0)

					if dir, exists := circleDirections[entityA.EntityID]; exists {
						circleDirections[entityA.EntityID] = dir.Add(forceA)
					}
					if dir, exists := circleDirections[entityB.EntityID]; exists {
						circleDirections[entityB.EntityID] = dir.Add(forceB)
					}
				}
			}
		}
	}

	// Move all circles
	for _, circle := range allCircles {
		entity := entityMap[circle.EntityID]
		if entity == nil {
			continue
		}

		direction := circleDirections[circle.EntityID]
		newPosition := logic.UpdateCirclePosition(entity, direction, 0.05, config.WorldSize) // 50ms delta

		entity.Position = newPosition
		if err := ctx.Database.UpdateEntity(entity); err != nil {
			LogWarn(fmt.Sprintf("Failed to update entity position %d: %v", entity.EntityID, err))
		}
	}

	// Check collisions
	for _, circle := range allCircles {
		circleEntity := entityMap[circle.EntityID]
		if circleEntity == nil {
			continue
		}

		for _, otherEntity := range allEntities {
			if otherEntity.EntityID == circleEntity.EntityID {
				continue
			}

			if logic.IsOverlapping(circleEntity, otherEntity) {
				// Check if it's another circle from a different player
				otherCircle, err := ctx.Database.GetCircle(otherEntity.EntityID)
				if err == nil && otherCircle != nil {
					if otherCircle.PlayerID != circle.PlayerID {
						// Player vs player collision
						if logic.CanConsumeEntity(circleEntity.Mass, otherEntity.Mass) {
							timer := logic.ScheduleConsumeEntity(circleEntity.EntityID, otherEntity.EntityID, ctx.Timestamp)
							// TODO: Insert timer into database
							_ = timer
						}
					}
				} else {
					// Player vs food collision
					timer := logic.ScheduleConsumeEntity(circleEntity.EntityID, otherEntity.EntityID, ctx.Timestamp)
					// TODO: Insert timer into database
					_ = timer
				}
			}
		}
	}

	return SuccessResult{}
}

// Helper function to clamp float values
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// SpawnFoodReducer handles spawning food entities
// Matches: Rust spawn_food() and C# SpawnFood()
func SpawnFoodReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("SpawnFood")
	defer timer.Stop()

	// Check if there are any players
	playerCount, err := ctx.Database.GetPlayerCount()
	if err != nil || playerCount == 0 {
		return SuccessResult{} // No players, don't spawn food
	}

	// Get current food count
	foodCount, err := ctx.Database.GetFoodCount()
	if err != nil {
		LogWarn(fmt.Sprintf("Failed to get food count: %v", err))
		foodCount = 0
	}

	config := constants.GetGlobalConfiguration()
	worldConfig, err := GetConfig(ctx)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get world config: %v", err)}
	}

	// Spawn food until we reach target count
	rng := ctx.Rng()
	for foodCount < uint64(config.TargetFoodCount) {
		entity, food, err := logic.SpawnFoodEntity(worldConfig.WorldSize, rng)
		if err != nil {
			LogWarn(fmt.Sprintf("Failed to spawn food entity: %v", err))
			break
		}

		if err := ctx.Database.InsertEntity(entity); err != nil {
			LogWarn(fmt.Sprintf("Failed to insert food entity: %v", err))
			continue
		}

		if err := ctx.Database.InsertFood(food); err != nil {
			LogWarn(fmt.Sprintf("Failed to insert food: %v", err))
			continue
		}

		foodCount++
		LogInfo(fmt.Sprintf("Spawned food! EntityID: %d", entity.EntityID))
	}

	return SuccessResult{}
}

// CircleDecayReducer handles circle mass decay
// Matches: Rust circle_decay() and C# CircleDecay()
func CircleDecayReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("CircleDecay")
	defer timer.Stop()

	// Get all circles
	circles, err := ctx.Database.GetAllCircles()
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get circles: %v", err)}
	}

	// Decay each circle that is above starting mass
	for _, circle := range circles {
		entity, err := ctx.Database.GetEntity(circle.EntityID)
		if err != nil {
			LogWarn(fmt.Sprintf("Failed to get entity for circle %d: %v", circle.EntityID, err))
			continue
		}

		if logic.ShouldCircleDecay(entity) {
			entity.Mass = logic.CalculateDecayedMass(entity.Mass)

			if err := ctx.Database.UpdateEntity(entity); err != nil {
				LogWarn(fmt.Sprintf("Failed to update decayed entity %d: %v", entity.EntityID, err))
			}
		}
	}

	return SuccessResult{}
}

// CircleRecombineArgs represents the arguments for CircleRecombine reducer
type CircleRecombineArgs struct {
	PlayerID uint32 `json:"player_id"`
}

// CircleRecombineReducer handles circle recombination for a player
// Matches: Rust circle_recombine() and C# CircleRecombine()
func CircleRecombineReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("CircleRecombine")
	defer timer.Stop()

	var recombineArgs CircleRecombineArgs
	if err := UnmarshalArgs(args, &recombineArgs); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Invalid arguments: %v", err)}
	}

	// Get player circles
	circles, err := ctx.Database.GetCirclesByPlayer(recombineArgs.PlayerID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to get player circles: %v", err)}
	}

	if len(circles) <= 1 {
		return SuccessResult{} // No circles to recombine
	}

	// Find circles that are ready to recombine
	var recombiningEntities []*tables.Entity
	config := constants.GetGlobalConfiguration()

	for _, circle := range circles {
		timeSinceSplit := ctx.Timestamp.Sub(circle.LastSplitTime).ToDuration().Seconds()
		if timeSinceSplit >= float64(config.SplitRecombineDelaySec) {
			entity, err := ctx.Database.GetEntity(circle.EntityID)
			if err != nil {
				LogWarn(fmt.Sprintf("Failed to get entity for circle %d: %v", circle.EntityID, err))
				continue
			}
			recombiningEntities = append(recombiningEntities, entity)
		}
	}

	if len(recombiningEntities) <= 1 {
		return SuccessResult{} // Nothing to recombine
	}

	// Schedule consumption of all circles into the first one
	baseEntityID := recombiningEntities[0].EntityID
	for i := 1; i < len(recombiningEntities); i++ {
		timer := logic.ScheduleConsumeEntity(baseEntityID, recombiningEntities[i].EntityID, ctx.Timestamp)
		// TODO: Insert timer into database
		_ = timer
	}

	return SuccessResult{}
}

// ConsumeEntityArgs represents the arguments for ConsumeEntity reducer
type ConsumeEntityArgs struct {
	ConsumerEntityID uint32 `json:"consumer_entity_id"`
	ConsumedEntityID uint32 `json:"consumed_entity_id"`
}

// ConsumeEntityReducer handles entity consumption
// Matches: Rust consume_entity() and C# ConsumeEntity()
func ConsumeEntityReducer(ctx *ReducerContext, args []byte) ReducerResult {
	timer := NewPerformanceTimer("ConsumeEntity")
	defer timer.Stop()

	var consumeArgs ConsumeEntityArgs
	if err := UnmarshalArgs(args, &consumeArgs); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Invalid arguments: %v", err)}
	}

	// Get both entities
	consumedEntity, err := ctx.Database.GetEntity(consumeArgs.ConsumedEntityID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Consumed entity doesn't exist: %v", err)}
	}

	consumerEntity, err := ctx.Database.GetEntity(consumeArgs.ConsumerEntityID)
	if err != nil {
		return ErrorResult{Message: fmt.Sprintf("Consumer entity doesn't exist: %v", err)}
	}

	// Transfer mass
	consumerEntity.Mass += consumedEntity.Mass

	// Destroy consumed entity
	if err := logic.DestroyEntity(ctx.Database.DeleteEntity, consumedEntity.EntityID); err != nil {
		LogWarn(fmt.Sprintf("Failed to destroy consumed entity %d: %v", consumedEntity.EntityID, err))
	}

	// Update consumer entity
	if err := ctx.Database.UpdateEntity(consumerEntity); err != nil {
		return ErrorResult{Message: fmt.Sprintf("Failed to update consumer entity: %v", err)}
	}

	return SuccessResult{}
}

// Register all Blackholio reducers
func init() {
	// Lifecycle reducers
	RegisterReducer(NewLifecycleReducer("Init", LifecycleInit, InitReducer))
	RegisterReducer(NewLifecycleReducer("Connect", LifecycleClientConnected, ConnectReducer))
	RegisterReducer(NewLifecycleReducer("Disconnect", LifecycleClientDisconnected, DisconnectReducer))

	// Game reducers
	RegisterReducer(NewReducer("EnterGame", EnterGameReducer).WithArgumentNames([]string{"name"}))
	RegisterReducer(NewReducer("Respawn", RespawnReducer))
	RegisterReducer(NewReducer("Suicide", SuicideReducer))
	RegisterReducer(NewReducer("UpdatePlayerInput", UpdatePlayerInputReducer).WithArgumentNames([]string{"direction"}))
	RegisterReducer(NewReducer("PlayerSplit", PlayerSplitReducer))

	// Scheduled reducers
	RegisterReducer(NewReducer("MoveAllPlayers", MoveAllPlayersReducer))
	RegisterReducer(NewReducer("SpawnFood", SpawnFoodReducer))
	RegisterReducer(NewReducer("CircleDecay", CircleDecayReducer))
	RegisterReducer(NewReducer("CircleRecombine", CircleRecombineReducer).WithArgumentNames([]string{"player_id"}))
	RegisterReducer(NewReducer("ConsumeEntity", ConsumeEntityReducer).WithArgumentNames([]string{"consumer_entity_id", "consumed_entity_id"}))

	LogInfo("Blackholio reducers registered successfully")
}
