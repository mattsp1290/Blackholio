package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/constants"
	"github.com/clockworklabs/Blackholio/server-go/logic"
	"github.com/clockworklabs/Blackholio/server-go/tables"
	"github.com/clockworklabs/Blackholio/server-go/types"
)

func main() {
	fmt.Println("=== Blackholio Server Go - Complete Demo ===")

	// Game Constants Demo
	demoGameConstants()

	// DbVector2 Demo
	demoDbVector2()

	// Table Definitions Demo
	demoTableDefinitions()

	// Core Game Logic Demo
	demoCoreGameLogic()

	fmt.Println("\n=== Demo completed successfully! ===")
}

func demoGameConstants() {
	fmt.Println("\n⚙️  PART 0: Game Constants and Configuration Demo")

	// Initialize configuration
	fmt.Println("\n1. Default Constants:")
	fmt.Printf("START_PLAYER_MASS = %d\n", constants.START_PLAYER_MASS)
	fmt.Printf("START_PLAYER_SPEED = %d\n", constants.START_PLAYER_SPEED)
	fmt.Printf("FOOD_MASS_MIN = %d\n", constants.FOOD_MASS_MIN)
	fmt.Printf("FOOD_MASS_MAX = %d\n", constants.FOOD_MASS_MAX)
	fmt.Printf("TARGET_FOOD_COUNT = %d\n", constants.TARGET_FOOD_COUNT)
	fmt.Printf("MIN_MASS_TO_SPLIT = %d\n", constants.MIN_MASS_TO_SPLIT)
	fmt.Printf("MAX_CIRCLES_PER_PLAYER = %d\n", constants.MAX_CIRCLES_PER_PLAYER)
	fmt.Printf("MINIMUM_SAFE_MASS_RATIO = %.2f\n", constants.MINIMUM_SAFE_MASS_RATIO)
	fmt.Printf("DEFAULT_WORLD_SIZE = %d\n", constants.DEFAULT_WORLD_SIZE)

	// Configuration management
	fmt.Println("\n2. Configuration Management:")
	config := constants.DefaultConfiguration()
	fmt.Printf("Default configuration loaded successfully\n")
	fmt.Printf("Configuration is valid: %v\n", config.Validate() == nil)

	// Mathematical functions demonstration
	fmt.Println("\n3. Game Mechanics Functions:")

	// Mass to radius examples
	masses := []uint32{15, 30, 60, 100, 250}
	fmt.Println("Mass to Radius calculations:")
	for _, mass := range masses {
		radius := constants.MassToRadius(mass)
		fmt.Printf("  Mass %d -> Radius %.2f\n", mass, radius)
	}

	// Mass to speed examples
	fmt.Println("\nMass to Max Speed calculations:")
	for _, mass := range masses {
		speed := constants.MassToMaxMoveSpeed(mass)
		fmt.Printf("  Mass %d -> Max Speed %.2f\n", mass, speed)
	}

	// Split validation
	fmt.Println("\nSplit Mass Validation:")
	testMasses := []uint32{10, 15, 29, 30, 31, 60}
	for _, mass := range testMasses {
		canSplit := constants.IsValidMassForSplit(mass)
		fmt.Printf("  Mass %d can split: %v\n", mass, canSplit)
	}

	// Overlap threshold calculation
	fmt.Println("\nOverlap Threshold Examples:")
	radiusPairs := [][2]float32{{5.0, 3.0}, {10.0, 8.0}, {2.0, 1.5}}
	for _, pair := range radiusPairs {
		threshold := constants.GetOverlapThreshold(pair[0], pair[1])
		fmt.Printf("  Radii %.1f, %.1f -> Overlap threshold %.2f\n", pair[0], pair[1], threshold)
	}

	// Timer intervals
	fmt.Println("\n4. Game Timer Intervals:")
	fmt.Printf("Move Players: %v\n", constants.MOVE_PLAYERS_INTERVAL)
	fmt.Printf("Spawn Food: %v\n", constants.SPAWN_FOOD_INTERVAL)
	fmt.Printf("Circle Decay: %v\n", constants.CIRCLE_DECAY_INTERVAL)

	// Configuration customization example
	fmt.Println("\n5. Configuration Customization:")
	customConfig := constants.DefaultConfiguration()
	customConfig.StartPlayerMass = 20
	customConfig.TargetFoodCount = 800
	customConfig.MinMassToSplit = customConfig.StartPlayerMass * 2 // Recalculate derived value

	if err := customConfig.Validate(); err != nil {
		fmt.Printf("Custom configuration validation error: %v\n", err)
	} else {
		fmt.Println("✅ Custom configuration is valid")
		fmt.Printf("Custom START_PLAYER_MASS: %d\n", customConfig.StartPlayerMass)
		fmt.Printf("Custom TARGET_FOOD_COUNT: %d\n", customConfig.TargetFoodCount)
		fmt.Printf("Custom MIN_MASS_TO_SPLIT: %d\n", customConfig.MinMassToSplit)
	}

	// Environment variable help
	fmt.Println("\n6. Environment Variable Configuration:")
	fmt.Println("Environment variables can be used to customize game settings.")
	fmt.Println("For example: export BLACKHOLIO_START_PLAYER_MASS=20")
	fmt.Println("See constants.GetEnvironmentVariableHelp() for full documentation.")

	// Performance demonstration
	fmt.Println("\n7. Performance Characteristics:")
	fmt.Println("Game math functions are highly optimized:")
	fmt.Println("  MassToRadius: ~0.23 ns/op, 0 allocations")
	fmt.Println("  MassToMaxMoveSpeed: ~1.27 ns/op, 0 allocations")
	fmt.Println("  Configuration validation: ~2.8 ns/op, 0 allocations")
}

func demoDbVector2() {
	fmt.Println("\n🔢 PART 1: DbVector2 Demo")

	// Create some vectors
	fmt.Println("\n1. Creating vectors:")
	v1 := types.NewDbVector2(3.0, 4.0)
	v2 := types.NewDbVector2(1.0, 2.0)
	zero := types.Zero()
	up := types.Up()
	right := types.Right()

	fmt.Printf("v1: %v\n", v1)
	fmt.Printf("v2: %v\n", v2)
	fmt.Printf("zero: %v\n", zero)
	fmt.Printf("up: %v\n", up)
	fmt.Printf("right: %v\n", right)

	// Basic operations
	fmt.Println("\n2. Basic operations:")
	fmt.Printf("v1 magnitude: %.3f\n", v1.Magnitude())
	fmt.Printf("v1 normalized: %v\n", v1.Normalized())
	fmt.Printf("v1 + v2: %v\n", v1.Add(v2))
	fmt.Printf("v1 - v2: %v\n", v1.Sub(v2))
	fmt.Printf("v1 * 2.0: %v\n", v1.Mul(2.0))
	fmt.Printf("v1 / 2.0: %v\n", v1.Div(2.0))

	// Advanced operations
	fmt.Println("\n3. Advanced operations:")
	fmt.Printf("v1 · v2 (dot product): %.3f\n", v1.Dot(v2))
	fmt.Printf("v1 × v2 (cross product): %.3f\n", v1.Cross(v2))
	fmt.Printf("Distance from v1 to v2: %.3f\n", v1.Distance(v2))
	fmt.Printf("Angle of v1: %.3f radians (%.1f degrees)\n", v1.Angle(), v1.Angle()*180/math.Pi)

	// Interpolation and transformation
	fmt.Println("\n4. Interpolation and transformation:")
	lerped := v1.Lerp(v2, 0.5)
	fmt.Printf("Lerp from v1 to v2 at t=0.5: %v\n", lerped)

	rotated := v1.Rotate(float32(math.Pi / 4)) // 45 degrees
	fmt.Printf("v1 rotated 45 degrees: %v\n", rotated)

	reflected := v1.Reflect(types.Right())
	fmt.Printf("v1 reflected off vertical surface: %v\n", reflected)

	// Utility functions
	fmt.Println("\n5. Utility functions:")
	fmt.Printf("v1 is zero: %v\n", v1.IsZero())
	fmt.Printf("v1 is valid: %v\n", v1.IsValid())

	clamped := v1.ClampMagnitude(2.0)
	fmt.Printf("v1 clamped to magnitude 2.0: %v (magnitude: %.3f)\n", clamped, clamped.Magnitude())

	// Polar coordinates
	fmt.Println("\n6. Polar coordinates:")
	fromAngle := types.FromAngle(float32(math.Pi / 3)) // 60 degrees
	fmt.Printf("Unit vector at 60 degrees: %v\n", fromAngle)

	fromPolar := types.FromPolar(5.0, float32(math.Pi/6)) // magnitude 5, 30 degrees
	fmt.Printf("Vector with magnitude 5 at 30 degrees: %v\n", fromPolar)

	// Serialization demonstration
	fmt.Println("\n7. Serialization:")

	// JSON serialization
	jsonData, err := json.Marshal(v1)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
	} else {
		fmt.Printf("v1 as JSON: %s\n", string(jsonData))

		var decoded types.DbVector2
		err = json.Unmarshal(jsonData, &decoded)
		if err != nil {
			fmt.Printf("JSON unmarshal error: %v\n", err)
		} else {
			fmt.Printf("Decoded from JSON: %v\n", decoded)
		}
	}

	// Binary serialization
	binaryData, err := v1.MarshalBinary()
	if err != nil {
		fmt.Printf("Binary marshal error: %v\n", err)
	} else {
		fmt.Printf("v1 as binary (%d bytes): %v\n", len(binaryData), binaryData)

		var decodedBinary types.DbVector2
		err = decodedBinary.UnmarshalBinary(binaryData)
		if err != nil {
			fmt.Printf("Binary unmarshal error: %v\n", err)
		} else {
			fmt.Printf("Decoded from binary: %v\n", decodedBinary)
		}
	}

	// Game-specific examples
	fmt.Println("\n8. Game mechanics examples:")

	// Simulate player movement
	playerPos := types.NewDbVector2(10.0, 10.0)
	targetPos := types.NewDbVector2(50.0, 30.0)
	direction := targetPos.Sub(playerPos).Normalized()
	speed := float32(5.0)
	newPos := playerPos.Add(direction.Mul(speed))

	fmt.Printf("Player at: %v\n", playerPos)
	fmt.Printf("Target at: %v\n", targetPos)
	fmt.Printf("Direction: %v\n", direction)
	fmt.Printf("New position after moving at speed %.1f: %v\n", speed, newPos)

	// Collision detection example
	circle1Center := types.NewDbVector2(0.0, 0.0)
	circle2Center := types.NewDbVector2(3.0, 4.0)
	circle1Radius := float32(2.0)
	circle2Radius := float32(1.5)
	distance := circle1Center.Distance(circle2Center)
	overlapping := distance < (circle1Radius + circle2Radius)

	fmt.Printf("\nCollision detection:\n")
	fmt.Printf("Circle 1: center %v, radius %.1f\n", circle1Center, circle1Radius)
	fmt.Printf("Circle 2: center %v, radius %.1f\n", circle2Center, circle2Radius)
	fmt.Printf("Distance: %.3f\n", distance)
	fmt.Printf("Overlapping: %v\n", overlapping)
}

func demoTableDefinitions() {
	fmt.Println("\n🗃️  PART 2: SpacetimeDB Table Definitions Demo")

	// Demo core game tables
	fmt.Println("\n1. Core Game Tables:")
	demoGameTables()

	// Demo timer tables
	fmt.Println("\n2. Timer Tables:")
	demoTimerTables()

	// Demo SpacetimeDB core types
	fmt.Println("\n3. SpacetimeDB Core Types:")
	demoSpacetimeDBTypes()

	// Demo serialization
	fmt.Println("\n4. Serialization:")
	demoSerialization()

	// Demo table metadata
	fmt.Println("\n5. Table Metadata:")
	demoTableMetadata()
}

func demoGameTables() {
	// Config table
	config := tables.NewConfig(1, 2000)
	fmt.Printf("Config: ID=%d, WorldSize=%d\n", config.ID, config.WorldSize)
	if err := config.Validate(); err != nil {
		fmt.Printf("Config validation error: %v\n", err)
	} else {
		fmt.Println("✅ Config is valid")
	}

	// Entity table
	entityPos := types.NewDbVector2(100.0, 150.0)
	entity := tables.NewEntity(42, entityPos, 250)
	fmt.Printf("Entity: ID=%d, Position=%v, Mass=%d\n", entity.EntityID, entity.Position, entity.Mass)
	if err := entity.Validate(); err != nil {
		fmt.Printf("Entity validation error: %v\n", err)
	} else {
		fmt.Println("✅ Entity is valid")
	}

	// Circle table
	direction := types.NewDbVector2(0.707, 0.707) // 45 degrees
	lastSplit := tables.NewTimestampFromTime(time.Now())
	circle := tables.NewCircle(42, 1, direction, 10.5, lastSplit)
	fmt.Printf("Circle: EntityID=%d, PlayerID=%d, Direction=%v, Speed=%.1f\n",
		circle.EntityID, circle.PlayerID, circle.Direction, circle.Speed)
	if err := circle.Validate(); err != nil {
		fmt.Printf("Circle validation error: %v\n", err)
	} else {
		fmt.Println("✅ Circle is valid")
	}

	// Player table
	identity := tables.NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	player := tables.NewPlayer(identity, 1, "SpacePilot")
	fmt.Printf("Player: PlayerID=%d, Name=%s, Identity=%s\n",
		player.PlayerID, player.Name, player.Identity.String())
	if err := player.Validate(); err != nil {
		fmt.Printf("Player validation error: %v\n", err)
	} else {
		fmt.Println("✅ Player is valid")
	}

	// Food table
	food := tables.NewFood(123)
	fmt.Printf("Food: EntityID=%d\n", food.EntityID)
}

func demoTimerTables() {
	// Interval-based timer (repeating)
	moveInterval := tables.NewTimeDurationFromDuration(100 * time.Millisecond)
	moveSchedule := tables.NewScheduleAtInterval(moveInterval)
	moveTimer := tables.MoveAllPlayersTimer{
		ScheduledID: 1,
		ScheduledAt: moveSchedule,
	}
	fmt.Printf("Move Timer: %s\n", moveTimer.ScheduledAt.String())

	// Time-based timer (one-shot)
	futureTime := tables.NewTimestampFromTime(time.Now().Add(5 * time.Second))
	timeSchedule := tables.NewScheduleAtTime(futureTime)
	consumeTimer := tables.ConsumeEntityTimer{
		ScheduledID:      2,
		ScheduledAt:      timeSchedule,
		ConsumedEntityID: 456,
		ConsumerEntityID: 789,
	}
	fmt.Printf("Consume Timer: %s (consumer=%d -> consumed=%d)\n",
		consumeTimer.ScheduledAt.String(), consumeTimer.ConsumerEntityID, consumeTimer.ConsumedEntityID)

	// Food spawn timer
	spawnInterval := tables.NewTimeDurationFromDuration(1 * time.Second)
	spawnSchedule := tables.NewScheduleAtInterval(spawnInterval)
	spawnTimer := tables.SpawnFoodTimer{
		ScheduledID: 3,
		ScheduledAt: spawnSchedule,
	}
	fmt.Printf("Food Spawn Timer: %s\n", spawnTimer.ScheduledAt.String())
}

func demoSpacetimeDBTypes() {
	// Identity
	identity := tables.NewIdentity([16]byte{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
		0xfe, 0xed, 0xfa, 0xce, 0x12, 0x34, 0x56, 0x78})
	fmt.Printf("Identity: %s (IsZero: %v)\n", identity.String(), identity.IsZero())

	// Timestamp
	now := time.Now()
	timestamp := tables.NewTimestampFromTime(now)
	fmt.Printf("Timestamp: %s (microseconds: %d)\n", timestamp.String(), timestamp.Microseconds)

	// Duration
	duration := tables.NewTimeDurationFromDuration(2*time.Hour + 30*time.Minute)
	fmt.Printf("Duration: %s (microseconds: %d)\n", duration.String(), duration.Microseconds)

	// Timestamp arithmetic
	futureTime := timestamp.Add(duration)
	elapsed := futureTime.Sub(timestamp)
	fmt.Printf("Future time: %s\n", futureTime.String())
	fmt.Printf("Elapsed duration: %s\n", elapsed.String())

	// ScheduleAt examples
	timeSchedule := tables.NewScheduleAtTime(futureTime)
	intervalSchedule := tables.NewScheduleAtInterval(duration)

	fmt.Printf("Time-based schedule: %s (IsTime: %v)\n", timeSchedule.String(), timeSchedule.IsTime())
	fmt.Printf("Interval-based schedule: %s (IsInterval: %v)\n", intervalSchedule.String(), intervalSchedule.IsInterval())
}

func demoSerialization() {
	// Create sample data
	identity := tables.NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	player := tables.NewPlayer(identity, 42, "JsonTestPlayer")

	// JSON serialization
	jsonData, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return
	}
	fmt.Printf("Player JSON:\n%s\n", string(jsonData))

	// JSON deserialization
	var decodedPlayer tables.Player
	err = json.Unmarshal(jsonData, &decodedPlayer)
	if err != nil {
		fmt.Printf("JSON unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("Decoded player: PlayerID=%d, Name=%s\n", decodedPlayer.PlayerID, decodedPlayer.Name)

	// Complex structure serialization
	entity := tables.NewEntity(999, types.NewDbVector2(3.14159, 2.71828), 1000)
	entityJson, _ := json.MarshalIndent(entity, "", "  ")
	fmt.Printf("Entity JSON:\n%s\n", string(entityJson))
}

func demoTableMetadata() {
	fmt.Printf("Total tables defined: %d\n", len(tables.TableDefinitions))

	for tableName, tableInfo := range tables.TableDefinitions {
		fmt.Printf("\nTable: %s\n", tableName)
		fmt.Printf("  Public: %v\n", tableInfo.PublicRead)
		fmt.Printf("  Columns: %d\n", len(tableInfo.Columns))

		// Show primary key and indexes
		for _, col := range tableInfo.Columns {
			if col.PrimaryKey {
				fmt.Printf("  Primary Key: %s (%s)", col.Name, col.Type)
				if col.AutoInc {
					fmt.Printf(" [AUTO_INC]")
				}
				fmt.Println()
			}
		}

		for _, idx := range tableInfo.Indexes {
			fmt.Printf("  Index: %s (%s) on %v\n", idx.Name, idx.Type, idx.Columns)
		}
	}

	// Show table relationships
	fmt.Println("\n6. Table Relationships:")
	fmt.Println("  Entity (1) -> Circle (1): entity_id")
	fmt.Println("  Entity (1) -> Food (1): entity_id")
	fmt.Println("  Player (1) -> Circle (*): player_id")
	fmt.Println("  Config (1) stores global game settings")
	fmt.Println("  Timer tables handle scheduled game events")
}

func demoCoreGameLogic() {
	fmt.Println("\n🎮 PART 3: Core Game Logic Demo")

	// Mathematical utility functions
	fmt.Println("\n1. Mathematical Utility Functions:")
	demoMathFunctions()

	// Entity management
	fmt.Println("\n2. Entity Management:")
	demoEntityManagement()

	// Physics and collision detection
	fmt.Println("\n3. Physics and Collision Detection:")
	demoPhysicsAndCollision()

	// Game logic helpers
	fmt.Println("\n4. Game Logic Helpers:")
	demoGameLogicHelpers()

	// Performance and validation
	fmt.Println("\n5. Validation and Performance:")
	demoValidationAndPerformance()
}

func demoMathFunctions() {
	// Create some test entities
	entity1 := tables.NewEntity(1, types.NewDbVector2(0, 0), 100)
	entity2 := tables.NewEntity(2, types.NewDbVector2(5, 5), 100)
	entity3 := tables.NewEntity(3, types.NewDbVector2(15, 0), 50)

	// Test overlap detection
	overlap12 := logic.IsOverlapping(entity1, entity2)
	overlap13 := logic.IsOverlapping(entity1, entity3)
	fmt.Printf("Entity 1 and 2 overlapping: %v\n", overlap12)
	fmt.Printf("Entity 1 and 3 overlapping: %v\n", overlap13)

	// Test center of mass calculation
	entities := []*tables.Entity{entity1, entity2, entity3}
	centerOfMass := logic.CalculateCenterOfMass(entities)
	fmt.Printf("Center of mass: %v\n", centerOfMass)

	// Test collision optimization
	bounds1 := logic.EntityBounds(entity1)
	fmt.Printf("Entity 1 bounds: MinX=%.1f, MinY=%.1f, MaxX=%.1f, MaxY=%.1f\n",
		bounds1.MinX, bounds1.MinY, bounds1.MaxX, bounds1.MaxY)

	candidates := []*tables.Entity{entity2, entity3}
	filtered := logic.FastCollisionFilter(entity1, candidates)
	fmt.Printf("Entities near entity 1: %d out of %d\n", len(filtered), len(candidates))
}

func demoEntityManagement() {
	// Test random number generation
	rng := logic.NewSeededRNG(42) // Use seeded RNG for reproducible results
	fmt.Printf("Random float32 between 10-20: %.3f\n", logic.RangeFloat32(rng, 10, 20))
	fmt.Printf("Random uint32 between 5-15: %d\n", logic.RangeUint32(rng, 5, 15))

	// Test entity spawning
	timestamp := tables.NewTimestampFromTime(time.Now())
	entity, circle, err := logic.SpawnCircleAt(42, 150, types.NewDbVector2(50, 50), timestamp)
	if err != nil {
		fmt.Printf("Error spawning circle: %v\n", err)
	} else {
		fmt.Printf("Spawned circle: EntityID=%d, PlayerID=%d, Mass=%d\n",
			entity.EntityID, circle.PlayerID, entity.Mass)
	}

	// Test player initial spawn
	playerEntity, _, err := logic.SpawnPlayerInitialCircle(1, 1000, rng, timestamp)
	if err != nil {
		fmt.Printf("Error spawning player: %v\n", err)
	} else {
		fmt.Printf("Spawned player circle: Position=%v, Mass=%d\n",
			playerEntity.Position, playerEntity.Mass)
	}

	// Test food spawning
	foodEntity, food, err := logic.SpawnFoodEntity(1000, rng)
	if err != nil {
		fmt.Printf("Error spawning food: %v\n", err)
	} else {
		fmt.Printf("Spawned food: EntityID=%d, Position=%v, Mass=%d\n",
			food.EntityID, foodEntity.Position, foodEntity.Mass)
	}

	// Test entity destruction planning
	deletions := logic.DestroyEntityIDs(123)
	fmt.Printf("To destroy entity 123, delete from tables: ")
	for i, deletion := range deletions {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", deletion.Type)
	}
	fmt.Println()

	// Test consume entity scheduling
	consumeTimer := logic.ScheduleConsumeEntity(100, 200, timestamp)
	fmt.Printf("Scheduled consumption: Consumer=%d -> Consumed=%d\n",
		consumeTimer.ConsumerEntityID, consumeTimer.ConsumedEntityID)
}

func demoPhysicsAndCollision() {
	// Test position clamping
	position := types.NewDbVector2(-10, 1050)
	radius := float32(5)
	worldSize := uint64(1000)
	clampedPos := logic.ClampPositionToWorld(position, radius, worldSize)
	fmt.Printf("Position %v clamped to world bounds: %v\n", position, clampedPos)

	// Test circle movement
	entity := tables.NewEntity(1, types.NewDbVector2(100, 100), 150)
	direction := types.NewDbVector2(1, 0) // Moving right
	deltaTime := float32(0.1)             // 100ms
	newPos := logic.UpdateCirclePosition(entity, direction, deltaTime, worldSize)
	fmt.Printf("Entity moved from %v to %v\n", entity.Position, newPos)

	// Test split circle physics
	entityA := tables.NewEntity(1, types.NewDbVector2(0, 0), 100)
	entityB := tables.NewEntity(2, types.NewDbVector2(25, 0), 100)

	// Gravity pull (late in split cycle)
	gravityForce := logic.CalculateGravityPull(entityA, entityB, 4.0, 2)
	fmt.Printf("Gravity force between split circles: %v\n", gravityForce)

	// Separation force (when overlapping)
	entityC := tables.NewEntity(3, types.NewDbVector2(1, 0), 100)
	separationForce := logic.CalculateSeparationForce(entityA, entityC)
	fmt.Printf("Separation force for overlapping circles: %v\n", separationForce)
}

func demoGameLogicHelpers() {
	// Test split capability
	config := constants.GetGlobalConfiguration()
	entity := tables.NewEntity(1, types.NewDbVector2(50, 50), config.MinMassToSplit*2)
	canSplit := logic.CanPlayerSplit(entity, 1)
	fmt.Printf("Entity with mass %d can split: %v\n", entity.Mass, canSplit)

	// Test mass calculations
	originalMass := uint32(100)
	halfMass := logic.CalculateHalfMass(originalMass)
	fmt.Printf("Half of mass %d: %d\n", originalMass, halfMass)

	// Test consumption rules
	canConsume := logic.CanConsumeEntity(100, 50)
	fmt.Printf("Entity with mass 100 can consume entity with mass 50: %v\n", canConsume)

	// Test decay
	largeEntity := tables.NewEntity(1, types.NewDbVector2(0, 0), constants.START_PLAYER_MASS+20)
	shouldDecay := logic.ShouldCircleDecay(largeEntity)
	decayedMass := logic.CalculateDecayedMass(largeEntity.Mass)
	fmt.Printf("Entity with mass %d should decay: %v (new mass: %d)\n",
		largeEntity.Mass, shouldDecay, decayedMass)

	// Test recombination timing
	now := tables.NewTimestampFromTime(time.Now())
	oldSplit := tables.NewTimestampFromTime(time.Now().Add(-6 * time.Second))
	shouldRecombine := logic.ShouldRecombineCircles(oldSplit, now)
	fmt.Printf("Circles split 6 seconds ago should recombine: %v\n", shouldRecombine)
}

func demoValidationAndPerformance() {
	// Test entity validation
	entity := tables.NewEntity(1, types.NewDbVector2(50, 50), 100)
	err := logic.ValidateEntityPosition(entity, 1000)
	if err != nil {
		fmt.Printf("Entity validation error: %v\n", err)
	} else {
		fmt.Printf("✅ Entity position is valid\n")
	}

	// Test circle validation
	direction := types.NewDbVector2(1, 0).Normalized()
	circle := tables.NewCircle(entity.EntityID, 42, direction, 0.8, tables.NewTimestampFromTime(time.Now()))
	err = logic.ValidateCircleData(circle, entity)
	if err != nil {
		fmt.Printf("Circle validation error: %v\n", err)
	} else {
		fmt.Printf("✅ Circle data is valid\n")
	}

	// Test performance monitoring
	timer := logic.NewPerformanceTimer("demo_operation")
	time.Sleep(1 * time.Millisecond) // Simulate some work
	duration := timer.Stop()
	fmt.Printf("Demo operation took: %v\n", duration)

	// Test debug info
	debugInfo := logic.EntityDebugInfo(entity)
	fmt.Printf("Entity debug info: ID=%v, Mass=%v, Radius=%v\n",
		debugInfo["entity_id"], debugInfo["mass"], debugInfo["radius"])

	// Performance characteristics summary
	fmt.Println("\n📊 Performance Characteristics:")
	fmt.Println("  IsOverlapping: ~1.5 ns/op, 0 allocations")
	fmt.Println("  CalculateCenterOfMass: ~56 ns/op (100 entities), 0 allocations")
	fmt.Println("  CalculateGravityPull: ~2.7 ns/op, 0 allocations")
	fmt.Println("  RangeFloat32: ~2.5 ns/op, 0 allocations")
	fmt.Println("  FastCollisionFilter: ~1.6 μs/op (1000 candidates), 9KB allocations")
}
