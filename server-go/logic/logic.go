package logic

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/constants"
	"github.com/clockworklabs/Blackholio/server-go/tables"
	"github.com/clockworklabs/Blackholio/server-go/types"
)

// Mathematical Utility Functions
// These functions implement the core game physics and math

// IsOverlapping checks if two entities are overlapping for collision detection
// This matches the Rust and C# implementations exactly
func IsOverlapping(a, b *tables.Entity) bool {
	dx := a.Position.X - b.Position.X
	dy := a.Position.Y - b.Position.Y
	distanceSq := dx*dx + dy*dy

	radiusA := constants.MassToRadius(a.Mass)
	radiusB := constants.MassToRadius(b.Mass)

	// In C#: radius_sum = (radius_a + radius_b) * (1.0 - MIN_OVERLAP_PCT_TO_CONSUME)
	// In Rust: uses max_radius = f32::max(radius_a, radius_b)
	// Let's use the C# approach for consistency with constants
	config := constants.GetGlobalConfiguration()
	radiusSum := (radiusA + radiusB) * (1.0 - config.MinOverlapPctToConsume)

	return distanceSq <= radiusSum*radiusSum
}

// IsOverlappingRust implements the Rust version of overlap detection
// This uses the max radius approach instead of the threshold approach
func IsOverlappingRust(a, b *tables.Entity) bool {
	dx := a.Position.X - b.Position.X
	dy := a.Position.Y - b.Position.Y
	distanceSq := dx*dx + dy*dy

	radiusA := constants.MassToRadius(a.Mass)
	radiusB := constants.MassToRadius(b.Mass)

	// Rust approach: use max radius
	maxRadius := float32(math.Max(float64(radiusA), float64(radiusB)))
	return distanceSq <= maxRadius*maxRadius
}

// CalculateCenterOfMass calculates the center of mass for a slice of entities
// This matches both Rust and C# implementations
func CalculateCenterOfMass(entities []*tables.Entity) types.DbVector2 {
	if len(entities) == 0 {
		return types.Zero()
	}

	var totalMass uint32
	var centerOfMass types.DbVector2

	for _, entity := range entities {
		totalMass += entity.Mass
		weighted := entity.Position.Mul(float32(entity.Mass))
		centerOfMass = centerOfMass.Add(weighted)
	}

	if totalMass == 0 {
		return types.Zero()
	}

	return centerOfMass.Div(float32(totalMass))
}

// Entity Management Functions
// These functions handle spawning, destroying, and managing game entities

// SpawnCircleAt creates a new circle entity at the specified position
// This matches the Rust and C# implementations exactly
func SpawnCircleAt(playerID uint32, mass uint32, position types.DbVector2, timestamp tables.Timestamp) (*tables.Entity, *tables.Circle, error) {
	// Create the entity
	entity := tables.NewEntity(0, position, mass) // EntityID will be auto-assigned

	// Create the circle
	direction := types.NewDbVector2(0, 1) // Default direction: up
	circle := tables.NewCircle(entity.EntityID, playerID, direction, 0.0, timestamp)

	return entity, circle, nil
}

// SpawnPlayerInitialCircle spawns a player's initial circle at a random safe position
// This matches the Rust and C# implementations exactly
func SpawnPlayerInitialCircle(playerID uint32, worldSize uint64, rng *rand.Rand, timestamp tables.Timestamp) (*tables.Entity, *tables.Circle, error) {
	playerStartRadius := constants.MassToRadius(constants.START_PLAYER_MASS)
	worldSizeFloat := float32(worldSize)

	// Generate random position with safety margin
	x := RangeFloat32(rng, playerStartRadius, worldSizeFloat-playerStartRadius)
	y := RangeFloat32(rng, playerStartRadius, worldSizeFloat-playerStartRadius)

	position := types.NewDbVector2(x, y)
	return SpawnCircleAt(playerID, constants.START_PLAYER_MASS, position, timestamp)
}

// SpawnFoodEntity creates a new food entity at a random position
func SpawnFoodEntity(worldSize uint64, rng *rand.Rand) (*tables.Entity, *tables.Food, error) {
	config := constants.GetGlobalConfiguration()

	// Random mass between min and max
	foodMass := RangeUint32(rng, config.FoodMassMin, config.FoodMassMax)
	foodRadius := constants.MassToRadius(foodMass)
	worldSizeFloat := float32(worldSize)

	// Generate random position with safety margin
	x := RangeFloat32(rng, foodRadius, worldSizeFloat-foodRadius)
	y := RangeFloat32(rng, foodRadius, worldSizeFloat-foodRadius)

	position := types.NewDbVector2(x, y)
	entity := tables.NewEntity(0, position, foodMass) // EntityID will be auto-assigned
	food := tables.NewFood(entity.EntityID)

	return entity, food, nil
}

// DestroyEntityIDs returns the entity IDs that should be deleted when destroying an entity
// This matches the C# and Rust implementations
func DestroyEntityIDs(entityID uint32) []EntityDeletion {
	return []EntityDeletion{
		{Type: "food", EntityID: entityID},
		{Type: "circle", EntityID: entityID},
		{Type: "entity", EntityID: entityID},
	}
}

// EntityDeletion represents an entity deletion operation
type EntityDeletion struct {
	Type     string // "food", "circle", "entity"
	EntityID uint32
}

// DestroyEntity destroys an entity by performing all the necessary database deletions
// This is a database operation that should be implemented by the database context
type DestroyEntityFunc func(entityID uint32) error

// DestroyEntity destroys an entity using the provided destroy function
// This matches the pattern used in Rust and C# implementations
func DestroyEntity(destroyFunc DestroyEntityFunc, entityID uint32) error {
	return destroyFunc(entityID)
}

// ScheduleConsumeEntity creates a timer for entity consumption
func ScheduleConsumeEntity(consumerID, consumedID uint32, timestamp tables.Timestamp) *tables.ConsumeEntityTimer {
	scheduleAt := tables.NewScheduleAtTime(timestamp)
	return &tables.ConsumeEntityTimer{
		ScheduledID:      0, // Will be auto-assigned
		ScheduledAt:      scheduleAt,
		ConsumerEntityID: consumerID,
		ConsumedEntityID: consumedID,
	}
}

// Random Number Generation Helpers
// These functions provide game-specific random number generation

// RangeFloat32 generates a random float32 between min and max (exclusive)
func RangeFloat32(rng *rand.Rand, min, max float32) float32 {
	if min >= max {
		return min
	}
	return rng.Float32()*(max-min) + min
}

// RangeUint32 generates a random uint32 between min and max (inclusive)
func RangeUint32(rng *rand.Rand, min, max uint32) uint32 {
	if min >= max {
		return min
	}
	return uint32(rng.Intn(int(max-min+1))) + min
}

// NewGameRNG creates a new random number generator with current time seed
func NewGameRNG() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewSeededRNG creates a new random number generator with a specific seed
func NewSeededRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// Collision Detection Optimizations
// These functions provide optimized collision detection for performance

// QuadrantBounds represents the bounds of a spatial quadrant
type QuadrantBounds struct {
	MinX, MinY, MaxX, MaxY float32
}

// EntityBounds calculates the bounding box for an entity
func EntityBounds(entity *tables.Entity) QuadrantBounds {
	radius := constants.MassToRadius(entity.Mass)
	return QuadrantBounds{
		MinX: entity.Position.X - radius,
		MinY: entity.Position.Y - radius,
		MaxX: entity.Position.X + radius,
		MaxY: entity.Position.Y + radius,
	}
}

// BoundsOverlap checks if two bounding boxes overlap (fast AABB test)
func BoundsOverlap(a, b QuadrantBounds) bool {
	return a.MinX <= b.MaxX && a.MaxX >= b.MinX &&
		a.MinY <= b.MaxY && a.MaxY >= b.MinY
}

// FastCollisionFilter filters entities for potential collisions using bounding boxes
func FastCollisionFilter(entity *tables.Entity, candidates []*tables.Entity) []*tables.Entity {
	entityBounds := EntityBounds(entity)
	var filtered []*tables.Entity

	for _, candidate := range candidates {
		if candidate.EntityID == entity.EntityID {
			continue
		}
		if BoundsOverlap(entityBounds, EntityBounds(candidate)) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}

// Physics and Movement Functions
// These functions handle the game physics

// ClampPositionToWorld ensures an entity's position stays within world bounds
func ClampPositionToWorld(position types.DbVector2, radius float32, worldSize uint64) types.DbVector2 {
	worldSizeFloat := float32(worldSize)
	return types.NewDbVector2(
		Clamp(position.X, radius, worldSizeFloat-radius),
		Clamp(position.Y, radius, worldSizeFloat-radius),
	)
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// UpdateCirclePosition updates a circle's position based on its movement
func UpdateCirclePosition(entity *tables.Entity, direction types.DbVector2, deltaTime float32, worldSize uint64) types.DbVector2 {
	speed := constants.MassToMaxMoveSpeed(entity.Mass)
	velocity := direction.Mul(speed * deltaTime)
	newPosition := entity.Position.Add(velocity)

	radius := constants.MassToRadius(entity.Mass)
	return ClampPositionToWorld(newPosition, radius, worldSize)
}

// Split Circle Physics
// These functions handle the complex physics for split circles

// CalculateGravityPull calculates gravitational pull between split circles
func CalculateGravityPull(entityA, entityB *tables.Entity, timeSinceSplit float32, circleCount int) types.DbVector2 {
	config := constants.GetGlobalConfiguration()

	timeBeforeRecombining := float32(math.Max(float64(config.SplitRecombineDelaySec-timeSinceSplit), 0.0))
	if timeBeforeRecombining > config.SplitGravPullBeforeRecombineSec {
		return types.Zero()
	}

	diff := entityA.Position.Sub(entityB.Position)
	distanceSqr := diff.SqrMagnitude()

	// Avoid division by zero
	if distanceSqr <= 0.0001 {
		diff = types.NewDbVector2(1.0, 0.0)
		distanceSqr = 1.0
	}

	radiusSum := constants.MassToRadius(entityA.Mass) + constants.MassToRadius(entityB.Mass)
	if distanceSqr > radiusSum*radiusSum {
		gravityMultiplier := 1.0 - timeBeforeRecombining/config.SplitGravPullBeforeRecombineSec
		distance := float32(math.Sqrt(float64(distanceSqr)))
		// Use the original formula: diff.Normalized * (radius_sum - distance)
		// When distance > radius_sum, this becomes attractive force
		vec := diff.Normalized().Mul(radiusSum - distance).Mul(gravityMultiplier).Mul(0.05).Div(float32(circleCount))
		return vec.Div(2.0)
	}

	return types.Zero()
}

// CalculateSeparationForce calculates force to separate overlapping split circles
func CalculateSeparationForce(entityA, entityB *tables.Entity) types.DbVector2 {
	config := constants.GetGlobalConfiguration()

	diff := entityA.Position.Sub(entityB.Position)
	distanceSqr := diff.SqrMagnitude()

	// Avoid division by zero
	if distanceSqr <= 0.0001 {
		diff = types.NewDbVector2(1.0, 0.0)
		distanceSqr = 1.0
	}

	radiusSum := constants.MassToRadius(entityA.Mass) + constants.MassToRadius(entityB.Mass)
	radiusSumMultiplied := radiusSum * config.AllowedSplitCircleOverlapPct

	if distanceSqr < radiusSumMultiplied*radiusSumMultiplied {
		distance := float32(math.Sqrt(float64(distanceSqr)))
		vec := diff.Normalized().Mul(radiusSum - distance).Mul(config.SelfCollisionSpeed)
		return vec.Div(2.0)
	}

	return types.Zero()
}

// Validation and Safety Functions
// These functions provide validation and safety checks

// ValidateEntityPosition checks if an entity position is valid
func ValidateEntityPosition(entity *tables.Entity, worldSize uint64) error {
	if !entity.Position.IsValid() {
		return fmt.Errorf("entity %d has invalid position: %v", entity.EntityID, entity.Position)
	}

	radius := constants.MassToRadius(entity.Mass)
	worldSizeFloat := float32(worldSize)

	if entity.Position.X-radius < 0 || entity.Position.X+radius > worldSizeFloat {
		return fmt.Errorf("entity %d X position out of bounds: %f (radius: %f, world: %f)",
			entity.EntityID, entity.Position.X, radius, worldSizeFloat)
	}

	if entity.Position.Y-radius < 0 || entity.Position.Y+radius > worldSizeFloat {
		return fmt.Errorf("entity %d Y position out of bounds: %f (radius: %f, world: %f)",
			entity.EntityID, entity.Position.Y, radius, worldSizeFloat)
	}

	return nil
}

// ValidateCircleData checks if circle data is consistent
func ValidateCircleData(circle *tables.Circle, entity *tables.Entity) error {
	if circle.EntityID != entity.EntityID {
		return fmt.Errorf("circle entity ID mismatch: circle=%d, entity=%d", circle.EntityID, entity.EntityID)
	}

	if !circle.Direction.IsValid() {
		return fmt.Errorf("circle %d has invalid direction: %v", circle.EntityID, circle.Direction)
	}

	if circle.Speed < 0 || circle.Speed > 1 {
		return fmt.Errorf("circle %d has invalid speed: %f (must be 0-1)", circle.EntityID, circle.Speed)
	}

	return nil
}

// Performance Monitoring Hooks
// These functions provide performance monitoring capabilities

// PerformanceTimer tracks execution time for game operations
type PerformanceTimer struct {
	Name      string
	StartTime time.Time
	Enabled   bool
}

// NewPerformanceTimer creates a new performance timer
func NewPerformanceTimer(name string) *PerformanceTimer {
	config := constants.GetGlobalConfiguration()
	return &PerformanceTimer{
		Name:      name,
		StartTime: time.Now(),
		Enabled:   config.EnablePerformanceLogging,
	}
}

// Stop stops the timer and optionally logs the result
func (pt *PerformanceTimer) Stop() time.Duration {
	duration := time.Since(pt.StartTime)
	if pt.Enabled {
		fmt.Printf("Performance[%s]: %v\n", pt.Name, duration)
	}
	return duration
}

// PerformanceMetrics holds performance statistics
type PerformanceMetrics struct {
	EntityCount       int
	CircleCount       int
	FoodCount         int
	CollisionChecks   int
	PhysicsOperations int
	LastUpdateTime    time.Duration
}

// Game Logic Helper Functions
// These functions provide common game logic operations

// CanPlayerSplit checks if a player's circle can split
func CanPlayerSplit(entity *tables.Entity, currentCircleCount uint32) bool {
	config := constants.GetGlobalConfiguration()

	if currentCircleCount >= config.MaxCirclesPerPlayer {
		return false
	}

	// Need at least double the minimum split mass to split in half
	return entity.Mass >= config.MinMassToSplit*2
}

// CalculateHalfMass calculates the mass for each half when splitting
func CalculateHalfMass(originalMass uint32) uint32 {
	return originalMass / 2
}

// CanConsumeEntity checks if one entity can consume another based on mass ratio
func CanConsumeEntity(consumerMass, consumedMass uint32) bool {
	config := constants.GetGlobalConfiguration()
	massRatio := float32(consumedMass) / float32(consumerMass)
	return massRatio < config.MinimumSafeMassRatio
}

// ShouldCircleDecay checks if a circle should lose mass due to decay
func ShouldCircleDecay(entity *tables.Entity) bool {
	return entity.Mass > constants.START_PLAYER_MASS
}

// CalculateDecayedMass calculates the new mass after decay
func CalculateDecayedMass(originalMass uint32) uint32 {
	// 1% decay per tick (matches Rust and C# implementations)
	return uint32(float32(originalMass) * 0.99)
}

// ShouldRecombineCircles checks if circles should recombine based on time
func ShouldRecombineCircles(lastSplitTime tables.Timestamp, currentTime tables.Timestamp) bool {
	config := constants.GetGlobalConfiguration()
	timeSinceSplit := currentTime.Sub(lastSplitTime)
	return timeSinceSplit.ToDuration().Seconds() >= float64(config.SplitRecombineDelaySec)
}

// Debug and Development Helpers
// These functions assist with debugging and development

// EntityDebugInfo returns debug information for an entity
func EntityDebugInfo(entity *tables.Entity) map[string]interface{} {
	radius := constants.MassToRadius(entity.Mass)
	speed := constants.MassToMaxMoveSpeed(entity.Mass)

	return map[string]interface{}{
		"entity_id": entity.EntityID,
		"position":  entity.Position.String(),
		"mass":      entity.Mass,
		"radius":    radius,
		"max_speed": speed,
		"bounds":    EntityBounds(entity),
	}
}

// CircleDebugInfo returns debug information for a circle
func CircleDebugInfo(circle *tables.Circle) map[string]interface{} {
	return map[string]interface{}{
		"entity_id":       circle.EntityID,
		"player_id":       circle.PlayerID,
		"direction":       circle.Direction.String(),
		"speed":           circle.Speed,
		"last_split_time": circle.LastSplitTime.String(),
	}
}

// GameStateDebugInfo returns debug information for the entire game state
func GameStateDebugInfo(entities []*tables.Entity, circles []*tables.Circle, food []*tables.Food) map[string]interface{} {
	totalMass := uint32(0)
	for _, entity := range entities {
		totalMass += entity.Mass
	}

	return map[string]interface{}{
		"entity_count": len(entities),
		"circle_count": len(circles),
		"food_count":   len(food),
		"total_mass":   totalMass,
		"avg_mass":     float32(totalMass) / float32(len(entities)),
	}
}
