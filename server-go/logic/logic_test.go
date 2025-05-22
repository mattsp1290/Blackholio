package logic

import (
	"math"
	"testing"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/constants"
	"github.com/clockworklabs/Blackholio/server-go/tables"
	"github.com/clockworklabs/Blackholio/server-go/types"
)

// Test helper to create test entities
func createTestEntity(id uint32, x, y float32, mass uint32) *tables.Entity {
	position := types.NewDbVector2(x, y)
	return tables.NewEntity(id, position, mass)
}

func TestIsOverlapping(t *testing.T) {
	t.Run("Overlapping entities", func(t *testing.T) {
		// Two circles at the same position should overlap
		entity1 := createTestEntity(1, 10, 10, 25)
		entity2 := createTestEntity(2, 10, 10, 25)

		if !IsOverlapping(entity1, entity2) {
			t.Error("Entities at same position should overlap")
		}
	})

	t.Run("Non-overlapping entities", func(t *testing.T) {
		// Two circles far apart should not overlap
		entity1 := createTestEntity(1, 0, 0, 25)
		entity2 := createTestEntity(2, 100, 100, 25)

		if IsOverlapping(entity1, entity2) {
			t.Error("Entities far apart should not overlap")
		}
	})

	t.Run("Just touching entities", func(t *testing.T) {
		// Test edge case where circles are just touching
		entity1 := createTestEntity(1, 0, 0, 25)
		radius1 := constants.MassToRadius(25)
		radius2 := constants.MassToRadius(25)
		config := constants.GetGlobalConfiguration()

		// Place second entity just at the edge
		distance := (radius1 + radius2) * (1.0 - config.MinOverlapPctToConsume)
		entity2 := createTestEntity(2, distance, 0, 25)

		// Should be overlapping (at the threshold)
		if !IsOverlapping(entity1, entity2) {
			t.Error("Entities at threshold should overlap")
		}
	})

	t.Run("Different masses", func(t *testing.T) {
		// Test with different mass entities
		entity1 := createTestEntity(1, 0, 0, 100) // Large circle
		entity2 := createTestEntity(2, 2, 2, 25)  // Small circle

		// Should overlap since small circle is inside large one
		if !IsOverlapping(entity1, entity2) {
			t.Error("Small entity inside large entity should overlap")
		}
	})
}

func TestIsOverlappingRust(t *testing.T) {
	t.Run("Rust overlap detection", func(t *testing.T) {
		entity1 := createTestEntity(1, 0, 0, 100)
		entity2 := createTestEntity(2, 5, 5, 25)

		// Test Rust version which uses max radius
		result := IsOverlappingRust(entity1, entity2)

		// Should be true if distance is less than max radius
		radius1 := constants.MassToRadius(100)
		distance := float32(math.Sqrt(50)) // sqrt(5^2 + 5^2)
		expected := distance <= radius1

		if result != expected {
			t.Errorf("Rust overlap detection failed: got %v, expected %v", result, expected)
		}
	})
}

func TestCalculateCenterOfMass(t *testing.T) {
	t.Run("Empty entities", func(t *testing.T) {
		result := CalculateCenterOfMass([]*tables.Entity{})
		expected := types.Zero()

		if !result.Equal(expected) {
			t.Errorf("Center of mass for empty slice should be zero: got %v", result)
		}
	})

	t.Run("Single entity", func(t *testing.T) {
		entity := createTestEntity(1, 10, 20, 100)
		result := CalculateCenterOfMass([]*tables.Entity{entity})

		if !result.Equal(entity.Position) {
			t.Errorf("Center of mass for single entity should be its position: got %v, expected %v", result, entity.Position)
		}
	})

	t.Run("Two equal entities", func(t *testing.T) {
		entity1 := createTestEntity(1, 0, 0, 50)
		entity2 := createTestEntity(2, 10, 10, 50)
		entities := []*tables.Entity{entity1, entity2}

		result := CalculateCenterOfMass(entities)
		expected := types.NewDbVector2(5, 5) // Midpoint

		if !result.Equal(expected) {
			t.Errorf("Center of mass for equal entities should be midpoint: got %v, expected %v", result, expected)
		}
	})

	t.Run("Different mass entities", func(t *testing.T) {
		entity1 := createTestEntity(1, 0, 0, 100) // Heavy
		entity2 := createTestEntity(2, 10, 0, 50) // Light
		entities := []*tables.Entity{entity1, entity2}

		result := CalculateCenterOfMass(entities)

		// Should be closer to the heavier entity
		// Weighted average: (0*100 + 10*50) / (100+50) = 500/150 = 3.33
		expected := types.NewDbVector2(10.0/3.0, 0)

		if math.Abs(float64(result.X-expected.X)) > 0.01 {
			t.Errorf("Center of mass calculation wrong: got %v, expected %v", result, expected)
		}
	})

	t.Run("Zero mass entities", func(t *testing.T) {
		entity1 := createTestEntity(1, 10, 10, 0)
		entity2 := createTestEntity(2, 20, 20, 0)
		entities := []*tables.Entity{entity1, entity2}

		result := CalculateCenterOfMass(entities)
		expected := types.Zero()

		if !result.Equal(expected) {
			t.Errorf("Center of mass for zero mass entities should be zero: got %v", result)
		}
	})
}

func TestSpawnCircleAt(t *testing.T) {
	t.Run("Basic spawn", func(t *testing.T) {
		playerID := uint32(42)
		mass := uint32(100)
		position := types.NewDbVector2(50, 75)
		timestamp := tables.NewTimestampFromTime(time.Now())

		entity, circle, err := SpawnCircleAt(playerID, mass, position, timestamp)

		if err != nil {
			t.Fatalf("SpawnCircleAt failed: %v", err)
		}

		// Check entity properties
		if entity.Position != position {
			t.Errorf("Entity position wrong: got %v, expected %v", entity.Position, position)
		}
		if entity.Mass != mass {
			t.Errorf("Entity mass wrong: got %d, expected %d", entity.Mass, mass)
		}

		// Check circle properties
		if circle.PlayerID != playerID {
			t.Errorf("Circle player ID wrong: got %d, expected %d", circle.PlayerID, playerID)
		}
		if circle.EntityID != entity.EntityID {
			t.Errorf("Circle entity ID should match entity: got %d, expected %d", circle.EntityID, entity.EntityID)
		}
		if circle.Speed != 0.0 {
			t.Errorf("Initial speed should be 0: got %f", circle.Speed)
		}

		// Check default direction (up)
		expected := types.NewDbVector2(0, 1)
		if !circle.Direction.Equal(expected) {
			t.Errorf("Default direction should be up: got %v, expected %v", circle.Direction, expected)
		}
	})
}

func TestSpawnPlayerInitialCircle(t *testing.T) {
	t.Run("Valid spawn", func(t *testing.T) {
		playerID := uint32(123)
		worldSize := uint64(1000)
		rng := NewSeededRNG(42) // Use seeded RNG for reproducible test
		timestamp := tables.NewTimestampFromTime(time.Now())

		entity, circle, err := SpawnPlayerInitialCircle(playerID, worldSize, rng, timestamp)

		if err != nil {
			t.Fatalf("SpawnPlayerInitialCircle failed: %v", err)
		}

		// Check basic properties
		if entity.Mass != constants.START_PLAYER_MASS {
			t.Errorf("Initial mass wrong: got %d, expected %d", entity.Mass, constants.START_PLAYER_MASS)
		}

		// Check position is within bounds
		radius := constants.MassToRadius(constants.START_PLAYER_MASS)
		if entity.Position.X < radius || entity.Position.X > float32(worldSize)-radius {
			t.Errorf("X position out of bounds: %f (radius: %f, world: %d)", entity.Position.X, radius, worldSize)
		}
		if entity.Position.Y < radius || entity.Position.Y > float32(worldSize)-radius {
			t.Errorf("Y position out of bounds: %f (radius: %f, world: %d)", entity.Position.Y, radius, worldSize)
		}

		if circle.PlayerID != playerID {
			t.Errorf("Player ID wrong: got %d, expected %d", circle.PlayerID, playerID)
		}
	})

	t.Run("Small world size", func(t *testing.T) {
		playerID := uint32(123)
		worldSize := uint64(20) // Very small world
		rng := NewSeededRNG(42)
		timestamp := tables.NewTimestampFromTime(time.Now())

		entity, _, err := SpawnPlayerInitialCircle(playerID, worldSize, rng, timestamp)

		if err != nil {
			t.Fatalf("SpawnPlayerInitialCircle failed: %v", err)
		}

		// Position should still be valid even in small world
		radius := constants.MassToRadius(constants.START_PLAYER_MASS)
		if entity.Position.X < radius || entity.Position.X > float32(worldSize)-radius {
			t.Errorf("X position out of bounds in small world: %f", entity.Position.X)
		}
	})
}

func TestSpawnFoodEntity(t *testing.T) {
	t.Run("Valid food spawn", func(t *testing.T) {
		worldSize := uint64(1000)
		rng := NewSeededRNG(42)

		entity, food, err := SpawnFoodEntity(worldSize, rng)

		if err != nil {
			t.Fatalf("SpawnFoodEntity failed: %v", err)
		}

		config := constants.GetGlobalConfiguration()

		// Check mass is in valid range
		if entity.Mass < config.FoodMassMin || entity.Mass > config.FoodMassMax {
			t.Errorf("Food mass out of range: got %d, expected %d-%d", entity.Mass, config.FoodMassMin, config.FoodMassMax)
		}

		// Check position is within bounds
		radius := constants.MassToRadius(entity.Mass)
		if entity.Position.X < radius || entity.Position.X > float32(worldSize)-radius {
			t.Errorf("X position out of bounds: %f", entity.Position.X)
		}
		if entity.Position.Y < radius || entity.Position.Y > float32(worldSize)-radius {
			t.Errorf("Y position out of bounds: %f", entity.Position.Y)
		}

		if food.EntityID != entity.EntityID {
			t.Errorf("Food entity ID should match: got %d, expected %d", food.EntityID, entity.EntityID)
		}
	})
}

func TestDestroyEntityIDs(t *testing.T) {
	t.Run("Correct deletion order", func(t *testing.T) {
		entityID := uint32(123)
		deletions := DestroyEntityIDs(entityID)

		if len(deletions) != 3 {
			t.Fatalf("Expected 3 deletions, got %d", len(deletions))
		}

		// Check correct order and types
		expected := []string{"food", "circle", "entity"}
		for i, deletion := range deletions {
			if deletion.Type != expected[i] {
				t.Errorf("Deletion %d: got type %s, expected %s", i, deletion.Type, expected[i])
			}
			if deletion.EntityID != entityID {
				t.Errorf("Deletion %d: got entity ID %d, expected %d", i, deletion.EntityID, entityID)
			}
		}
	})
}

func TestScheduleConsumeEntity(t *testing.T) {
	t.Run("Valid scheduling", func(t *testing.T) {
		consumerID := uint32(100)
		consumedID := uint32(200)
		timestamp := tables.NewTimestampFromTime(time.Now())

		timer := ScheduleConsumeEntity(consumerID, consumedID, timestamp)

		if timer.ConsumerEntityID != consumerID {
			t.Errorf("Consumer ID wrong: got %d, expected %d", timer.ConsumerEntityID, consumerID)
		}
		if timer.ConsumedEntityID != consumedID {
			t.Errorf("Consumed ID wrong: got %d, expected %d", timer.ConsumedEntityID, consumedID)
		}
		if !timer.ScheduledAt.IsTime() {
			t.Error("Timer should be scheduled at specific time")
		}
	})
}

func TestRandomFunctions(t *testing.T) {
	t.Run("RangeFloat32", func(t *testing.T) {
		rng := NewSeededRNG(42)
		min := float32(10)
		max := float32(20)

		for i := 0; i < 100; i++ {
			value := RangeFloat32(rng, min, max)
			if value < min || value >= max {
				t.Errorf("Value out of range: %f (expected %f-%f)", value, min, max)
			}
		}
	})

	t.Run("RangeFloat32 edge cases", func(t *testing.T) {
		rng := NewSeededRNG(42)

		// Test min >= max
		result := RangeFloat32(rng, 10, 10)
		if result != 10 {
			t.Errorf("When min == max, should return min: got %f", result)
		}

		result = RangeFloat32(rng, 20, 10)
		if result != 20 {
			t.Errorf("When min > max, should return min: got %f", result)
		}
	})

	t.Run("RangeUint32", func(t *testing.T) {
		rng := NewSeededRNG(42)
		min := uint32(5)
		max := uint32(15)

		for i := 0; i < 100; i++ {
			value := RangeUint32(rng, min, max)
			if value < min || value > max {
				t.Errorf("Value out of range: %d (expected %d-%d)", value, min, max)
			}
		}
	})

	t.Run("RNG creation", func(t *testing.T) {
		rng1 := NewGameRNG()
		rng2 := NewSeededRNG(123)

		if rng1 == nil || rng2 == nil {
			t.Error("RNG creation failed")
		}

		// Seeded RNG should be deterministic
		val1 := rng2.Float32()
		rng3 := NewSeededRNG(123)
		val2 := rng3.Float32()

		if val1 != val2 {
			t.Error("Seeded RNG should be deterministic")
		}
	})
}

func TestCollisionOptimization(t *testing.T) {
	t.Run("EntityBounds", func(t *testing.T) {
		entity := createTestEntity(1, 10, 20, 100)
		bounds := EntityBounds(entity)

		radius := constants.MassToRadius(100)
		expectedBounds := QuadrantBounds{
			MinX: 10 - radius,
			MinY: 20 - radius,
			MaxX: 10 + radius,
			MaxY: 20 + radius,
		}

		if bounds != expectedBounds {
			t.Errorf("Entity bounds wrong: got %+v, expected %+v", bounds, expectedBounds)
		}
	})

	t.Run("BoundsOverlap", func(t *testing.T) {
		bounds1 := QuadrantBounds{MinX: 0, MinY: 0, MaxX: 10, MaxY: 10}
		bounds2 := QuadrantBounds{MinX: 5, MinY: 5, MaxX: 15, MaxY: 15}
		bounds3 := QuadrantBounds{MinX: 20, MinY: 20, MaxX: 30, MaxY: 30}

		if !BoundsOverlap(bounds1, bounds2) {
			t.Error("Overlapping bounds should return true")
		}
		if BoundsOverlap(bounds1, bounds3) {
			t.Error("Non-overlapping bounds should return false")
		}
	})

	t.Run("FastCollisionFilter", func(t *testing.T) {
		entity := createTestEntity(1, 10, 10, 100)
		candidates := []*tables.Entity{
			createTestEntity(2, 12, 12, 50),   // Close (should be included)
			createTestEntity(3, 100, 100, 50), // Far (should be excluded)
			createTestEntity(1, 10, 10, 100),  // Same entity (should be excluded)
		}

		filtered := FastCollisionFilter(entity, candidates)

		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered entity, got %d", len(filtered))
		}
		if filtered[0].EntityID != 2 {
			t.Errorf("Wrong entity filtered: got ID %d, expected 2", filtered[0].EntityID)
		}
	})
}

func TestPhysicsAndMovement(t *testing.T) {
	t.Run("ClampPositionToWorld", func(t *testing.T) {
		worldSize := uint64(100)
		radius := float32(5)

		// Test position within bounds
		pos1 := types.NewDbVector2(50, 50)
		result1 := ClampPositionToWorld(pos1, radius, worldSize)
		if !result1.Equal(pos1) {
			t.Errorf("Position within bounds should not change: got %v", result1)
		}

		// Test position out of bounds
		pos2 := types.NewDbVector2(-10, 110)
		result2 := ClampPositionToWorld(pos2, radius, worldSize)
		expected2 := types.NewDbVector2(radius, float32(worldSize)-radius)
		if !result2.Equal(expected2) {
			t.Errorf("Position should be clamped: got %v, expected %v", result2, expected2)
		}
	})

	t.Run("Clamp", func(t *testing.T) {
		if Clamp(5, 0, 10) != 5 {
			t.Error("Value within range should not change")
		}
		if Clamp(-5, 0, 10) != 0 {
			t.Error("Value below min should be clamped to min")
		}
		if Clamp(15, 0, 10) != 10 {
			t.Error("Value above max should be clamped to max")
		}
	})

	t.Run("UpdateCirclePosition", func(t *testing.T) {
		entity := createTestEntity(1, 50, 50, 100)
		direction := types.NewDbVector2(1, 0) // Moving right
		deltaTime := float32(1.0)
		worldSize := uint64(1000)

		newPos := UpdateCirclePosition(entity, direction, deltaTime, worldSize)

		// Should move to the right by speed amount
		expectedSpeed := constants.MassToMaxMoveSpeed(100)
		expectedX := 50 + expectedSpeed

		if math.Abs(float64(newPos.X-expectedX)) > 0.01 {
			t.Errorf("X position wrong: got %f, expected %f", newPos.X, expectedX)
		}
		if newPos.Y != 50 {
			t.Errorf("Y position should not change: got %f", newPos.Y)
		}
	})
}

func TestSplitCirclePhysics(t *testing.T) {
	t.Run("CalculateGravityPull early", func(t *testing.T) {
		entityA := createTestEntity(1, 0, 0, 100)
		entityB := createTestEntity(2, 10, 0, 100)
		timeSinceSplit := float32(1.0) // Early in split
		circleCount := 2

		force := CalculateGravityPull(entityA, entityB, timeSinceSplit, circleCount)

		// Should be zero since it's too early for gravity
		if !force.Equal(types.Zero()) {
			t.Errorf("Gravity should be zero early in split: got %v", force)
		}
	})

	t.Run("CalculateGravityPull late", func(t *testing.T) {
		entityA := createTestEntity(1, 0, 0, 100)
		entityB := createTestEntity(2, 25, 0, 100) // Far apart (beyond radius sum)
		timeSinceSplit := float32(4.0)             // Late in split cycle
		circleCount := 2

		force := CalculateGravityPull(entityA, entityB, timeSinceSplit, circleCount)

		// Should have some gravitational force
		if force.Magnitude() == 0 {
			t.Error("Should have gravitational force late in split")
		}

		// Force should pull A towards B (towards positive X since B is to the right of A)
		if force.X <= 0 {
			t.Errorf("Gravity force should be positive X (towards B): got %v", force)
		}
	})

	t.Run("CalculateSeparationForce", func(t *testing.T) {
		// Two entities very close together
		entityA := createTestEntity(1, 0, 0, 100)
		entityB := createTestEntity(2, 1, 0, 100) // Very close

		force := CalculateSeparationForce(entityA, entityB)

		// Should have separation force since they're overlapping
		if force.Magnitude() == 0 {
			t.Error("Should have separation force for overlapping circles")
		}

		// Force should point away from B (negative X direction for A)
		if force.X >= 0 {
			t.Errorf("Separation force should be negative X: got %v", force)
		}
	})

	t.Run("Zero distance handling", func(t *testing.T) {
		// Test entities at exactly the same position
		entityA := createTestEntity(1, 10, 10, 100)
		entityB := createTestEntity(2, 10, 10, 100)

		gravityForce := CalculateGravityPull(entityA, entityB, 4.0, 2)
		separationForce := CalculateSeparationForce(entityA, entityB)

		// Should handle zero distance gracefully (use default direction)
		if !gravityForce.IsValid() {
			t.Error("Gravity force should be valid even at zero distance")
		}
		if !separationForce.IsValid() {
			t.Error("Separation force should be valid even at zero distance")
		}
	})
}

func TestValidation(t *testing.T) {
	t.Run("ValidateEntityPosition valid", func(t *testing.T) {
		entity := createTestEntity(1, 50, 50, 100)
		worldSize := uint64(1000)

		err := ValidateEntityPosition(entity, worldSize)
		if err != nil {
			t.Errorf("Valid entity should pass validation: %v", err)
		}
	})

	t.Run("ValidateEntityPosition out of bounds", func(t *testing.T) {
		entity := createTestEntity(1, 5, 5, 100) // Too close to edge
		worldSize := uint64(100)

		err := ValidateEntityPosition(entity, worldSize)
		if err == nil {
			t.Error("Entity too close to edge should fail validation")
		}
	})

	t.Run("ValidateEntityPosition invalid values", func(t *testing.T) {
		entity := createTestEntity(1, float32(math.NaN()), 50, 100)
		worldSize := uint64(1000)

		err := ValidateEntityPosition(entity, worldSize)
		if err == nil {
			t.Error("Entity with NaN position should fail validation")
		}
	})

	t.Run("ValidateCircleData valid", func(t *testing.T) {
		entity := createTestEntity(1, 50, 50, 100)
		direction := types.NewDbVector2(1, 0).Normalized()
		circle := tables.NewCircle(entity.EntityID, 42, direction, 0.5, tables.NewTimestampFromTime(time.Now()))

		err := ValidateCircleData(circle, entity)
		if err != nil {
			t.Errorf("Valid circle should pass validation: %v", err)
		}
	})

	t.Run("ValidateCircleData entity ID mismatch", func(t *testing.T) {
		entity := createTestEntity(1, 50, 50, 100)
		direction := types.NewDbVector2(1, 0)
		circle := tables.NewCircle(999, 42, direction, 0.5, tables.NewTimestampFromTime(time.Now())) // Wrong entity ID

		err := ValidateCircleData(circle, entity)
		if err == nil {
			t.Error("Circle with wrong entity ID should fail validation")
		}
	})

	t.Run("ValidateCircleData invalid speed", func(t *testing.T) {
		entity := createTestEntity(1, 50, 50, 100)
		direction := types.NewDbVector2(1, 0)
		circle := tables.NewCircle(entity.EntityID, 42, direction, 2.0, tables.NewTimestampFromTime(time.Now())) // Invalid speed

		err := ValidateCircleData(circle, entity)
		if err == nil {
			t.Error("Circle with invalid speed should fail validation")
		}
	})
}

func TestPerformanceMonitoring(t *testing.T) {
	t.Run("PerformanceTimer", func(t *testing.T) {
		timer := NewPerformanceTimer("test")

		if timer == nil {
			t.Fatal("Timer creation failed")
		}

		time.Sleep(10 * time.Millisecond)
		duration := timer.Stop()

		if duration < 10*time.Millisecond {
			t.Error("Timer should measure at least 10ms")
		}
	})
}

func TestGameLogicHelpers(t *testing.T) {
	t.Run("CanPlayerSplit", func(t *testing.T) {
		config := constants.GetGlobalConfiguration()

		// Entity with enough mass
		entity1 := createTestEntity(1, 50, 50, config.MinMassToSplit*2)
		if !CanPlayerSplit(entity1, 1) {
			t.Error("Entity with enough mass should be able to split")
		}

		// Entity without enough mass
		entity2 := createTestEntity(2, 50, 50, config.MinMassToSplit)
		if CanPlayerSplit(entity2, 1) {
			t.Error("Entity without enough mass should not be able to split")
		}

		// Too many circles
		entity3 := createTestEntity(3, 50, 50, config.MinMassToSplit*2)
		if CanPlayerSplit(entity3, config.MaxCirclesPerPlayer) {
			t.Error("Player with max circles should not be able to split")
		}
	})

	t.Run("CalculateHalfMass", func(t *testing.T) {
		if CalculateHalfMass(100) != 50 {
			t.Error("Half of 100 should be 50")
		}
		if CalculateHalfMass(101) != 50 {
			t.Error("Half of 101 should be 50 (integer division)")
		}
	})

	t.Run("CanConsumeEntity", func(t *testing.T) {
		// Small entity can be consumed
		if !CanConsumeEntity(100, 50) {
			t.Error("Small entity should be consumable")
		}

		// Large entity cannot be consumed
		if CanConsumeEntity(100, 90) {
			t.Error("Large entity should not be consumable")
		}
	})

	t.Run("ShouldCircleDecay", func(t *testing.T) {
		// Large circle should decay
		entity1 := createTestEntity(1, 50, 50, constants.START_PLAYER_MASS+10)
		if !ShouldCircleDecay(entity1) {
			t.Error("Large circle should decay")
		}

		// Small circle should not decay
		entity2 := createTestEntity(2, 50, 50, constants.START_PLAYER_MASS)
		if ShouldCircleDecay(entity2) {
			t.Error("Small circle should not decay")
		}
	})

	t.Run("CalculateDecayedMass", func(t *testing.T) {
		original := uint32(100)
		decayed := CalculateDecayedMass(original)
		expected := uint32(99) // 1% decay

		if decayed != expected {
			t.Errorf("Decayed mass wrong: got %d, expected %d", decayed, expected)
		}
	})

	t.Run("ShouldRecombineCircles", func(t *testing.T) {
		now := tables.NewTimestampFromTime(time.Now())
		config := constants.GetGlobalConfiguration()

		// Recent split - should not recombine
		recentSplit := tables.NewTimestampFromTime(time.Now().Add(-1 * time.Second))
		if ShouldRecombineCircles(recentSplit, now) {
			t.Error("Recent split should not recombine")
		}

		// Old split - should recombine
		oldSplit := tables.NewTimestampFromTime(time.Now().Add(-time.Duration(config.SplitRecombineDelaySec+1) * time.Second))
		if !ShouldRecombineCircles(oldSplit, now) {
			t.Error("Old split should recombine")
		}
	})
}

func TestDebugHelpers(t *testing.T) {
	t.Run("EntityDebugInfo", func(t *testing.T) {
		entity := createTestEntity(123, 50, 75, 100)
		info := EntityDebugInfo(entity)

		if info["entity_id"] != uint32(123) {
			t.Error("Debug info should include entity ID")
		}
		if info["mass"] != uint32(100) {
			t.Error("Debug info should include mass")
		}
		if info["position"] == nil {
			t.Error("Debug info should include position")
		}
	})

	t.Run("CircleDebugInfo", func(t *testing.T) {
		direction := types.NewDbVector2(1, 0)
		timestamp := tables.NewTimestampFromTime(time.Now())
		circle := tables.NewCircle(123, 42, direction, 0.5, timestamp)

		info := CircleDebugInfo(circle)

		if info["entity_id"] != uint32(123) {
			t.Error("Debug info should include entity ID")
		}
		if info["player_id"] != uint32(42) {
			t.Error("Debug info should include player ID")
		}
	})

	t.Run("GameStateDebugInfo", func(t *testing.T) {
		entities := []*tables.Entity{
			createTestEntity(1, 0, 0, 100),
			createTestEntity(2, 10, 10, 200),
		}
		circles := []*tables.Circle{}
		food := []*tables.Food{}

		info := GameStateDebugInfo(entities, circles, food)

		if info["entity_count"] != 2 {
			t.Error("Debug info should show correct entity count")
		}
		if info["total_mass"] != uint32(300) {
			t.Error("Debug info should show correct total mass")
		}
		if info["avg_mass"] != float32(150) {
			t.Error("Debug info should show correct average mass")
		}
	})
}

// Benchmark tests for performance-critical functions

func BenchmarkIsOverlapping(b *testing.B) {
	entity1 := createTestEntity(1, 10, 10, 100)
	entity2 := createTestEntity(2, 12, 12, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsOverlapping(entity1, entity2)
	}
}

func BenchmarkCalculateCenterOfMass(b *testing.B) {
	entities := make([]*tables.Entity, 100)
	for i := 0; i < 100; i++ {
		entities[i] = createTestEntity(uint32(i), float32(i), float32(i), uint32(i+50))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateCenterOfMass(entities)
	}
}

func BenchmarkFastCollisionFilter(b *testing.B) {
	entity := createTestEntity(1, 50, 50, 100)
	candidates := make([]*tables.Entity, 1000)
	for i := 0; i < 1000; i++ {
		candidates[i] = createTestEntity(uint32(i+2), float32(i%100), float32(i%100), 50)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FastCollisionFilter(entity, candidates)
	}
}

func BenchmarkCalculateGravityPull(b *testing.B) {
	entityA := createTestEntity(1, 0, 0, 100)
	entityB := createTestEntity(2, 10, 0, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateGravityPull(entityA, entityB, 4.0, 2)
	}
}

func BenchmarkRangeFloat32(b *testing.B) {
	rng := NewSeededRNG(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RangeFloat32(rng, 0, 100)
	}
}
