package tables

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/types"
)

// Test constants
const testPlayerName = "TestPlayer"

func TestConfig(t *testing.T) {
	t.Run("NewConfig", func(t *testing.T) {
		config := NewConfig(1, 1000)
		if config.ID != 1 {
			t.Errorf("Expected ID 1, got %d", config.ID)
		}
		if config.WorldSize != 1000 {
			t.Errorf("Expected WorldSize 1000, got %d", config.WorldSize)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		// Valid config
		config := NewConfig(1, 1000)
		if err := config.Validate(); err != nil {
			t.Errorf("Valid config should not error: %v", err)
		}

		// Invalid config (zero world size)
		config.WorldSize = 0
		if err := config.Validate(); err == nil {
			t.Error("Config with zero world size should error")
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		original := NewConfig(42, 2000)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Config
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.ID != original.ID || decoded.WorldSize != original.WorldSize {
			t.Errorf("Round-trip failed: got %+v, want %+v", decoded, original)
		}
	})
}

func TestEntity(t *testing.T) {
	t.Run("NewEntity", func(t *testing.T) {
		position := types.NewDbVector2(10.0, 20.0)
		entity := NewEntity(1, position, 100)

		if entity.EntityID != 1 {
			t.Errorf("Expected EntityID 1, got %d", entity.EntityID)
		}
		if !entity.Position.Equal(position) {
			t.Errorf("Expected position %v, got %v", position, entity.Position)
		}
		if entity.Mass != 100 {
			t.Errorf("Expected mass 100, got %d", entity.Mass)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		// Valid entity
		entity := NewEntity(1, types.NewDbVector2(10.0, 20.0), 100)
		if err := entity.Validate(); err != nil {
			t.Errorf("Valid entity should not error: %v", err)
		}

		// Invalid position
		entity.Position = types.NewDbVector2(float32(math.Inf(1)), 0) // Infinity
		if err := entity.Validate(); err == nil {
			t.Error("Entity with invalid position should error")
		}

		// Invalid mass
		entity.Position = types.NewDbVector2(10.0, 20.0)
		entity.Mass = 0
		if err := entity.Validate(); err == nil {
			t.Error("Entity with zero mass should error")
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		original := NewEntity(42, types.NewDbVector2(3.14, 2.71), 500)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Entity
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.EntityID != original.EntityID || !decoded.Position.Equal(original.Position) || decoded.Mass != original.Mass {
			t.Errorf("Round-trip failed: got %+v, want %+v", decoded, original)
		}
	})
}

func TestCircle(t *testing.T) {
	t.Run("NewCircle", func(t *testing.T) {
		direction := types.NewDbVector2(1.0, 0.0)
		timestamp := NewTimestamp(1000000)
		circle := NewCircle(1, 2, direction, 5.0, timestamp)

		if circle.EntityID != 1 {
			t.Errorf("Expected EntityID 1, got %d", circle.EntityID)
		}
		if circle.PlayerID != 2 {
			t.Errorf("Expected PlayerID 2, got %d", circle.PlayerID)
		}
		if !circle.Direction.Equal(direction) {
			t.Errorf("Expected direction %v, got %v", direction, circle.Direction)
		}
		if circle.Speed != 5.0 {
			t.Errorf("Expected speed 5.0, got %f", circle.Speed)
		}
		if circle.LastSplitTime.Microseconds != timestamp.Microseconds {
			t.Errorf("Expected timestamp %d, got %d", timestamp.Microseconds, circle.LastSplitTime.Microseconds)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		// Valid circle
		circle := NewCircle(1, 2, types.NewDbVector2(1.0, 0.0), 5.0, NewTimestamp(1000000))
		if err := circle.Validate(); err != nil {
			t.Errorf("Valid circle should not error: %v", err)
		}

		// Invalid direction
		circle.Direction = types.NewDbVector2(float32(math.Inf(1)), 0) // Infinity
		if err := circle.Validate(); err == nil {
			t.Error("Circle with invalid direction should error")
		}

		// Negative speed
		circle.Direction = types.NewDbVector2(1.0, 0.0)
		circle.Speed = -1.0
		if err := circle.Validate(); err == nil {
			t.Error("Circle with negative speed should error")
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		original := NewCircle(42, 24, types.NewDbVector2(0.707, 0.707), 10.5, NewTimestamp(2000000))

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Circle
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.EntityID != original.EntityID || decoded.PlayerID != original.PlayerID ||
			!decoded.Direction.Equal(original.Direction) || decoded.Speed != original.Speed ||
			decoded.LastSplitTime.Microseconds != original.LastSplitTime.Microseconds {
			t.Errorf("Round-trip failed: got %+v, want %+v", decoded, original)
		}
	})
}

func TestPlayer(t *testing.T) {
	t.Run("NewPlayer", func(t *testing.T) {
		identity := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		player := NewPlayer(identity, 42, testPlayerName)

		if player.Identity.Bytes != identity.Bytes {
			t.Errorf("Expected identity %v, got %v", identity.Bytes, player.Identity.Bytes)
		}
		if player.PlayerID != 42 {
			t.Errorf("Expected PlayerID 42, got %d", player.PlayerID)
		}
		if player.Name != testPlayerName {
			t.Errorf("Expected name %s, got %s", testPlayerName, player.Name)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		// Valid player
		identity := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		player := NewPlayer(identity, 42, testPlayerName)
		if err := player.Validate(); err != nil {
			t.Errorf("Valid player should not error: %v", err)
		}

		// Zero identity
		player.Identity = NewIdentity([16]byte{})
		if err := player.Validate(); err == nil {
			t.Error("Player with zero identity should error")
		}

		// Empty name
		player.Identity = identity
		player.Name = ""
		if err := player.Validate(); err == nil {
			t.Error("Player with empty name should error")
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		identity := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		original := NewPlayer(identity, 42, testPlayerName)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Player
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.Identity.Bytes != original.Identity.Bytes || decoded.PlayerID != original.PlayerID || decoded.Name != original.Name {
			t.Errorf("Round-trip failed: got %+v, want %+v", decoded, original)
		}
	})
}

func TestFood(t *testing.T) {
	t.Run("NewFood", func(t *testing.T) {
		food := NewFood(123)
		if food.EntityID != 123 {
			t.Errorf("Expected EntityID 123, got %d", food.EntityID)
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		original := NewFood(456)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Food
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.EntityID != original.EntityID {
			t.Errorf("Round-trip failed: got %+v, want %+v", decoded, original)
		}
	})
}

func TestIdentity(t *testing.T) {
	t.Run("NewIdentity", func(t *testing.T) {
		bytes := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
		identity := NewIdentity(bytes)
		if identity.Bytes != bytes {
			t.Errorf("Expected bytes %v, got %v", bytes, identity.Bytes)
		}
	})

	t.Run("IsZero", func(t *testing.T) {
		// Zero identity
		zero := NewIdentity([16]byte{})
		if !zero.IsZero() {
			t.Error("Zero identity should return true for IsZero()")
		}

		// Non-zero identity
		nonZero := NewIdentity([16]byte{1})
		if nonZero.IsZero() {
			t.Error("Non-zero identity should return false for IsZero()")
		}
	})

	t.Run("String", func(t *testing.T) {
		identity := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
		expected := "Identity(0102030405060708090a0b0c0d0e0f10)"
		if identity.String() != expected {
			t.Errorf("Expected string %s, got %s", expected, identity.String())
		}
	})

	t.Run("JSONSerialization", func(t *testing.T) {
		original := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var decoded Identity
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if decoded.Bytes != original.Bytes {
			t.Errorf("Round-trip failed: got %v, want %v", decoded.Bytes, original.Bytes)
		}
	})

	t.Run("JSONSerializationEdgeCases", func(t *testing.T) {
		var identity Identity

		// Test invalid hex string length
		err := identity.UnmarshalJSON([]byte(`"invalid"`))
		if err == nil {
			t.Error("Should fail with invalid hex string length")
		}

		// Test invalid hex characters
		err = identity.UnmarshalJSON([]byte(`"gggggggggggggggggggggggggggggggg"`))
		if err == nil {
			t.Error("Should fail with invalid hex characters")
		}
	})
}

func TestTimestamp(t *testing.T) {
	t.Run("NewTimestamp", func(t *testing.T) {
		timestamp := NewTimestamp(1000000)
		if timestamp.Microseconds != 1000000 {
			t.Errorf("Expected microseconds 1000000, got %d", timestamp.Microseconds)
		}
	})

	t.Run("NewTimestampFromTime", func(t *testing.T) {
		now := time.Now()
		timestamp := NewTimestampFromTime(now)

		// Convert back and check if it's close (within 1ms tolerance)
		converted := timestamp.ToTime()
		diff := now.Sub(converted)
		if diff > time.Millisecond || diff < -time.Millisecond {
			t.Errorf("Time conversion diff too large: %v", diff)
		}
	})

	t.Run("Add", func(t *testing.T) {
		timestamp := NewTimestamp(1000000)
		duration := NewTimeDuration(500000)

		result := timestamp.Add(duration)
		expected := NewTimestamp(1500000)

		if result.Microseconds != expected.Microseconds {
			t.Errorf("Expected %d, got %d", expected.Microseconds, result.Microseconds)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		t1 := NewTimestamp(2000000)
		t2 := NewTimestamp(1000000)

		duration := t1.Sub(t2)
		expected := NewTimeDuration(1000000)

		if duration.Microseconds != expected.Microseconds {
			t.Errorf("Expected %d, got %d", expected.Microseconds, duration.Microseconds)
		}

		// Test underflow protection
		duration = t2.Sub(t1)
		if duration.Microseconds != 0 {
			t.Errorf("Expected 0 for underflow, got %d", duration.Microseconds)
		}
	})

	t.Run("String", func(t *testing.T) {
		timestamp := NewTimestamp(1609459200000000) // 2021-01-01 00:00:00 UTC
		str := timestamp.String()
		if str == "" {
			t.Error("String representation should not be empty")
		}
	})
}

func TestTimeDuration(t *testing.T) {
	t.Run("NewTimeDuration", func(t *testing.T) {
		duration := NewTimeDuration(1000000)
		if duration.Microseconds != 1000000 {
			t.Errorf("Expected microseconds 1000000, got %d", duration.Microseconds)
		}
	})

	t.Run("NewTimeDurationFromDuration", func(t *testing.T) {
		goDuration := time.Second
		duration := NewTimeDurationFromDuration(goDuration)

		// Convert back and check
		converted := duration.ToDuration()
		if converted != goDuration {
			t.Errorf("Expected %v, got %v", goDuration, converted)
		}
	})

	t.Run("String", func(t *testing.T) {
		duration := NewTimeDuration(1000000) // 1 second
		str := duration.String()
		if str != "1s" {
			t.Errorf("Expected '1s', got '%s'", str)
		}
	})
}

func TestScheduleAt(t *testing.T) {
	t.Run("NewScheduleAtTime", func(t *testing.T) {
		timestamp := NewTimestamp(1000000)
		schedule := NewScheduleAtTime(timestamp)

		if !schedule.IsTime() {
			t.Error("Should be time-based schedule")
		}
		if schedule.IsInterval() {
			t.Error("Should not be interval-based schedule")
		}
		if schedule.GetTime().Microseconds != timestamp.Microseconds {
			t.Errorf("Expected time %d, got %d", timestamp.Microseconds, schedule.GetTime().Microseconds)
		}
	})

	t.Run("NewScheduleAtInterval", func(t *testing.T) {
		duration := NewTimeDuration(1000000)
		schedule := NewScheduleAtInterval(duration)

		if schedule.IsTime() {
			t.Error("Should not be time-based schedule")
		}
		if !schedule.IsInterval() {
			t.Error("Should be interval-based schedule")
		}
		if schedule.GetInterval().Microseconds != duration.Microseconds {
			t.Errorf("Expected interval %d, got %d", duration.Microseconds, schedule.GetInterval().Microseconds)
		}
	})

	t.Run("String", func(t *testing.T) {
		// Time-based
		timestamp := NewTimestamp(1000000)
		schedule := NewScheduleAtTime(timestamp)
		str := schedule.String()
		if str == "" || str == "ScheduleAt(None)" {
			t.Errorf("Time-based schedule string should not be empty or none: %s", str)
		}

		// Interval-based
		duration := NewTimeDuration(1000000)
		schedule = NewScheduleAtInterval(duration)
		str = schedule.String()
		if str == "" || str == "ScheduleAt(None)" {
			t.Errorf("Interval-based schedule string should not be empty or none: %s", str)
		}

		// Empty schedule
		schedule = ScheduleAt{}
		str = schedule.String()
		if str != "ScheduleAt(None)" {
			t.Errorf("Empty schedule should be 'ScheduleAt(None)', got '%s'", str)
		}
	})
}

func TestTimerTables(t *testing.T) {
	t.Run("MoveAllPlayersTimer", func(t *testing.T) {
		schedule := NewScheduleAtInterval(NewTimeDuration(1000000))
		timer := MoveAllPlayersTimer{
			ScheduledID: 1,
			ScheduledAt: schedule,
		}

		if timer.ScheduledID != 1 {
			t.Errorf("Expected ScheduledID 1, got %d", timer.ScheduledID)
		}
		if !timer.ScheduledAt.IsInterval() {
			t.Error("Timer should have interval-based schedule")
		}
	})

	t.Run("ConsumeEntityTimer", func(t *testing.T) {
		schedule := NewScheduleAtTime(NewTimestamp(2000000))
		timer := ConsumeEntityTimer{
			ScheduledID:      1,
			ScheduledAt:      schedule,
			ConsumedEntityID: 42,
			ConsumerEntityID: 24,
		}

		if timer.ConsumedEntityID != 42 {
			t.Errorf("Expected ConsumedEntityID 42, got %d", timer.ConsumedEntityID)
		}
		if timer.ConsumerEntityID != 24 {
			t.Errorf("Expected ConsumerEntityID 24, got %d", timer.ConsumerEntityID)
		}
	})
}

func TestTableDefinitions(t *testing.T) {
	t.Run("ConfigTable", func(t *testing.T) {
		def, exists := TableDefinitions["config"]
		if !exists {
			t.Fatal("Config table definition not found")
		}

		if def.Name != "config" {
			t.Errorf("Expected name 'config', got '%s'", def.Name)
		}
		if !def.PublicRead {
			t.Error("Config table should be public")
		}
		if len(def.Columns) != 2 {
			t.Errorf("Expected 2 columns, got %d", len(def.Columns))
		}

		// Check primary key
		var hasPrimaryKey bool
		for _, col := range def.Columns {
			if col.PrimaryKey {
				hasPrimaryKey = true
				if col.Name != "id" {
					t.Errorf("Expected primary key 'id', got '%s'", col.Name)
				}
			}
		}
		if !hasPrimaryKey {
			t.Error("Config table should have a primary key")
		}
	})

	t.Run("CircleTable", func(t *testing.T) {
		def, exists := TableDefinitions["circle"]
		if !exists {
			t.Fatal("Circle table definition not found")
		}

		// Check for B-Tree index on player_id
		var hasPlayerIDIndex bool
		for _, idx := range def.Indexes {
			if idx.Name == "player_id" && idx.Type == "btree" {
				hasPlayerIDIndex = true
				break
			}
		}
		if !hasPlayerIDIndex {
			t.Error("Circle table should have B-Tree index on player_id")
		}
	})

	t.Run("AllTablesExist", func(t *testing.T) {
		expectedTables := []string{
			"config", "entity", "circle", "player", "logged_out_player", "food",
			"move_all_players_timer", "spawn_food_timer", "circle_decay_timer",
			"circle_recombine_timer", "consume_entity_timer",
		}

		for _, tableName := range expectedTables {
			if _, exists := TableDefinitions[tableName]; !exists {
				t.Errorf("Table definition for '%s' not found", tableName)
			}
		}
	})
}

// Benchmark tests for performance
func BenchmarkEntityCreation(b *testing.B) {
	position := types.NewDbVector2(10.0, 20.0)
	for i := 0; i < b.N; i++ {
		_ = NewEntity(uint32(i), position, 100)
	}
}

func BenchmarkIdentityJSON(b *testing.B) {
	identity := NewIdentity([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(identity)
		var decoded Identity
		_ = json.Unmarshal(data, &decoded)
	}
}

func BenchmarkTimestampOperations(b *testing.B) {
	t1 := NewTimestamp(1000000)
	t2 := NewTimestamp(2000000)
	duration := NewTimeDuration(500000)

	for i := 0; i < b.N; i++ {
		_ = t1.Add(duration)
		_ = t2.Sub(t1)
	}
}
