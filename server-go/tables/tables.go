package tables

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clockworklabs/Blackholio/server-go/types"
)

// SpacetimeDB table definitions for Blackholio game
// These structs match the Rust and C# implementations exactly

// Config represents the game configuration table
// Matches: Rust Config struct and C# Config struct
type Config struct {
	ID        uint32 `json:"id" spacetimedb:"primary_key" bsatn:"0"`
	WorldSize uint64 `json:"world_size" bsatn:"1"`
}

// Entity represents a game entity (player circles, food, etc.)
// Matches: Rust Entity struct and C# Entity struct
type Entity struct {
	EntityID uint32          `json:"entity_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	Position types.DbVector2 `json:"position" bsatn:"1"`
	Mass     uint32          `json:"mass" bsatn:"2"`
}

// Circle represents a player circle entity
// Matches: Rust Circle struct and C# Circle struct
type Circle struct {
	EntityID      uint32          `json:"entity_id" spacetimedb:"primary_key" bsatn:"0"`
	PlayerID      uint32          `json:"player_id" spacetimedb:"index:btree" bsatn:"1"`
	Direction     types.DbVector2 `json:"direction" bsatn:"2"`
	Speed         float32         `json:"speed" bsatn:"3"`
	LastSplitTime Timestamp       `json:"last_split_time" bsatn:"4"`
}

// Player represents a player in the game
// Matches: Rust Player struct and C# Player struct
// Note: This struct is used for both "player" and "logged_out_player" tables
type Player struct {
	Identity Identity `json:"identity" spacetimedb:"primary_key" bsatn:"0"`
	PlayerID uint32   `json:"player_id" spacetimedb:"unique,auto_inc" bsatn:"1"`
	Name     string   `json:"name" bsatn:"2"`
}

// Food represents a food entity in the game
// Matches: Rust Food struct and C# Food struct
type Food struct {
	EntityID uint32 `json:"entity_id" spacetimedb:"primary_key" bsatn:"0"`
}

// Timer Tables for Scheduled Reducers

// MoveAllPlayersTimer represents the timer for moving all players
// Matches: Rust MoveAllPlayersTimer struct and C# MoveAllPlayersTimer struct
type MoveAllPlayersTimer struct {
	ScheduledID uint64     `json:"scheduled_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	ScheduledAt ScheduleAt `json:"scheduled_at" spacetimedb:"scheduled_at" bsatn:"1"`
}

// SpawnFoodTimer represents the timer for spawning food
// Matches: Rust SpawnFoodTimer struct and C# SpawnFoodTimer struct
type SpawnFoodTimer struct {
	ScheduledID uint64     `json:"scheduled_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	ScheduledAt ScheduleAt `json:"scheduled_at" spacetimedb:"scheduled_at" bsatn:"1"`
}

// CircleDecayTimer represents the timer for circle decay
// Matches: Rust CircleDecayTimer struct and C# CircleDecayTimer struct
type CircleDecayTimer struct {
	ScheduledID uint64     `json:"scheduled_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	ScheduledAt ScheduleAt `json:"scheduled_at" spacetimedb:"scheduled_at" bsatn:"1"`
}

// CircleRecombineTimer represents the timer for circle recombination
// Matches: Rust CircleRecombineTimer struct and C# CircleRecombineTimer struct
type CircleRecombineTimer struct {
	ScheduledID uint64     `json:"scheduled_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	ScheduledAt ScheduleAt `json:"scheduled_at" spacetimedb:"scheduled_at" bsatn:"1"`
	PlayerID    uint32     `json:"player_id" bsatn:"2"`
}

// ConsumeEntityTimer represents the timer for entity consumption
// Matches: Rust ConsumeEntityTimer struct and C# ConsumeEntityTimer struct
type ConsumeEntityTimer struct {
	ScheduledID      uint64     `json:"scheduled_id" spacetimedb:"primary_key,auto_inc" bsatn:"0"`
	ScheduledAt      ScheduleAt `json:"scheduled_at" spacetimedb:"scheduled_at" bsatn:"1"`
	ConsumedEntityID uint32     `json:"consumed_entity_id" bsatn:"2"`
	ConsumerEntityID uint32     `json:"consumer_entity_id" bsatn:"3"`
}

// SpacetimeDB Core Types
// These types match the official SpacetimeDB Go bindings

// Identity represents a unique client identity
type Identity struct {
	Bytes [16]byte `json:"bytes" bsatn:"0"`
}

// Timestamp represents a point in time with microsecond precision
type Timestamp struct {
	Microseconds uint64 `json:"microseconds" bsatn:"0"`
}

// TimeDuration represents a duration with microsecond precision
type TimeDuration struct {
	Microseconds uint64 `json:"microseconds" bsatn:"0"`
}

// ScheduleAt represents when a scheduled reducer should run
type ScheduleAt struct {
	Time     *Timestamp    `json:"time,omitempty" bsatn:"0"`
	Interval *TimeDuration `json:"interval,omitempty" bsatn:"1"`
}

// Table Information and Metadata

// TableInfo contains metadata about a SpacetimeDB table
type TableInfo struct {
	Name       string   `json:"name"`
	PublicRead bool     `json:"public_read"`
	Columns    []Column `json:"columns"`
	Indexes    []Index  `json:"indexes"`
}

// Column represents a table column definition
type Column struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	PrimaryKey   bool   `json:"primary_key"`
	AutoInc      bool   `json:"auto_inc"`
	Unique       bool   `json:"unique"`
	NotNull      bool   `json:"not_null"`
	DefaultValue string `json:"default_value,omitempty"`
}

// Index represents a table index definition
type Index struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"` // "btree", "hash", etc.
	Columns   []string `json:"columns"`
	Unique    bool     `json:"unique"`
	Clustered bool     `json:"clustered"`
}

// Constructor Functions

// NewConfig creates a new Config instance
func NewConfig(id uint32, worldSize uint64) *Config {
	return &Config{
		ID:        id,
		WorldSize: worldSize,
	}
}

// NewEntity creates a new Entity instance
func NewEntity(entityID uint32, position types.DbVector2, mass uint32) *Entity {
	return &Entity{
		EntityID: entityID,
		Position: position,
		Mass:     mass,
	}
}

// NewCircle creates a new Circle instance
func NewCircle(entityID, playerID uint32, direction types.DbVector2, speed float32, lastSplitTime Timestamp) *Circle {
	return &Circle{
		EntityID:      entityID,
		PlayerID:      playerID,
		Direction:     direction,
		Speed:         speed,
		LastSplitTime: lastSplitTime,
	}
}

// NewPlayer creates a new Player instance
func NewPlayer(identity Identity, playerID uint32, name string) *Player {
	return &Player{
		Identity: identity,
		PlayerID: playerID,
		Name:     name,
	}
}

// NewFood creates a new Food instance
func NewFood(entityID uint32) *Food {
	return &Food{
		EntityID: entityID,
	}
}

// Utility Methods for Core Types

// NewIdentity creates a new Identity from bytes
func NewIdentity(bytes [16]byte) Identity {
	return Identity{Bytes: bytes}
}

// String returns a string representation of the Identity
func (i Identity) String() string {
	return fmt.Sprintf("Identity(%x)", i.Bytes)
}

// IsZero returns true if the identity is all zeros
func (i Identity) IsZero() bool {
	for _, b := range i.Bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

// NewTimestamp creates a new Timestamp from microseconds
func NewTimestamp(microseconds uint64) Timestamp {
	return Timestamp{Microseconds: microseconds}
}

// NewTimestampFromTime creates a new Timestamp from a Go time.Time
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp{Microseconds: uint64(t.UnixNano() / 1000)}
}

// ToTime converts a Timestamp to a Go time.Time
func (t Timestamp) ToTime() time.Time {
	return time.Unix(0, int64(t.Microseconds)*1000)
}

// String returns a string representation of the Timestamp
func (t Timestamp) String() string {
	return t.ToTime().Format(time.RFC3339Nano)
}

// Add adds a duration to the timestamp
func (t Timestamp) Add(duration TimeDuration) Timestamp {
	return Timestamp{Microseconds: t.Microseconds + duration.Microseconds}
}

// Sub subtracts another timestamp from this one, returning the duration
func (t Timestamp) Sub(other Timestamp) TimeDuration {
	if t.Microseconds >= other.Microseconds {
		return TimeDuration{Microseconds: t.Microseconds - other.Microseconds}
	}
	return TimeDuration{Microseconds: 0}
}

// NewTimeDuration creates a new TimeDuration from microseconds
func NewTimeDuration(microseconds uint64) TimeDuration {
	return TimeDuration{Microseconds: microseconds}
}

// NewTimeDurationFromDuration creates a new TimeDuration from a Go time.Duration
func NewTimeDurationFromDuration(d time.Duration) TimeDuration {
	return TimeDuration{Microseconds: uint64(d.Nanoseconds() / 1000)}
}

// ToDuration converts a TimeDuration to a Go time.Duration
func (d TimeDuration) ToDuration() time.Duration {
	return time.Duration(d.Microseconds * 1000)
}

// String returns a string representation of the TimeDuration
func (d TimeDuration) String() string {
	return d.ToDuration().String()
}

// ScheduleAt Constructors

// NewScheduleAtTime creates a ScheduleAt for a specific time
func NewScheduleAtTime(timestamp Timestamp) ScheduleAt {
	return ScheduleAt{Time: &timestamp}
}

// NewScheduleAtInterval creates a ScheduleAt for repeated intervals
func NewScheduleAtInterval(duration TimeDuration) ScheduleAt {
	return ScheduleAt{Interval: &duration}
}

// IsTime returns true if this is a time-based schedule
func (s ScheduleAt) IsTime() bool {
	return s.Time != nil
}

// IsInterval returns true if this is an interval-based schedule
func (s ScheduleAt) IsInterval() bool {
	return s.Interval != nil
}

// GetTime returns the scheduled time (nil if interval-based)
func (s ScheduleAt) GetTime() *Timestamp {
	return s.Time
}

// GetInterval returns the scheduled interval (nil if time-based)
func (s ScheduleAt) GetInterval() *TimeDuration {
	return s.Interval
}

// String returns a string representation of ScheduleAt
func (s ScheduleAt) String() string {
	if s.IsTime() {
		return fmt.Sprintf("ScheduleAt(Time: %s)", s.Time.String())
	} else if s.IsInterval() {
		return fmt.Sprintf("ScheduleAt(Interval: %s)", s.Interval.String())
	}
	return "ScheduleAt(None)"
}

// Table Definition Registry

// TableDefinitions contains all table definitions for the Blackholio game
var TableDefinitions = map[string]TableInfo{
	"config": {
		Name:       "config",
		PublicRead: true,
		Columns: []Column{
			{Name: "id", Type: "uint32", PrimaryKey: true},
			{Name: "world_size", Type: "uint64"},
		},
	},
	"entity": {
		Name:       "entity",
		PublicRead: true,
		Columns: []Column{
			{Name: "entity_id", Type: "uint32", PrimaryKey: true, AutoInc: true},
			{Name: "position", Type: "DbVector2"},
			{Name: "mass", Type: "uint32"},
		},
	},
	"circle": {
		Name:       "circle",
		PublicRead: true,
		Columns: []Column{
			{Name: "entity_id", Type: "uint32", PrimaryKey: true},
			{Name: "player_id", Type: "uint32"},
			{Name: "direction", Type: "DbVector2"},
			{Name: "speed", Type: "float32"},
			{Name: "last_split_time", Type: "Timestamp"},
		},
		Indexes: []Index{
			{Name: "player_id", Type: "btree", Columns: []string{"player_id"}},
		},
	},
	"player": {
		Name:       "player",
		PublicRead: true,
		Columns: []Column{
			{Name: "identity", Type: "Identity", PrimaryKey: true},
			{Name: "player_id", Type: "uint32", Unique: true, AutoInc: true},
			{Name: "name", Type: "string"},
		},
	},
	"logged_out_player": {
		Name:       "logged_out_player",
		PublicRead: false,
		Columns: []Column{
			{Name: "identity", Type: "Identity", PrimaryKey: true},
			{Name: "player_id", Type: "uint32", Unique: true, AutoInc: true},
			{Name: "name", Type: "string"},
		},
	},
	"food": {
		Name:       "food",
		PublicRead: true,
		Columns: []Column{
			{Name: "entity_id", Type: "uint32", PrimaryKey: true},
		},
	},
	// Timer tables
	"move_all_players_timer": {
		Name: "move_all_players_timer",
		Columns: []Column{
			{Name: "scheduled_id", Type: "uint64", PrimaryKey: true, AutoInc: true},
			{Name: "scheduled_at", Type: "ScheduleAt"},
		},
	},
	"spawn_food_timer": {
		Name: "spawn_food_timer",
		Columns: []Column{
			{Name: "scheduled_id", Type: "uint64", PrimaryKey: true, AutoInc: true},
			{Name: "scheduled_at", Type: "ScheduleAt"},
		},
	},
	"circle_decay_timer": {
		Name: "circle_decay_timer",
		Columns: []Column{
			{Name: "scheduled_id", Type: "uint64", PrimaryKey: true, AutoInc: true},
			{Name: "scheduled_at", Type: "ScheduleAt"},
		},
	},
	"circle_recombine_timer": {
		Name: "circle_recombine_timer",
		Columns: []Column{
			{Name: "scheduled_id", Type: "uint64", PrimaryKey: true, AutoInc: true},
			{Name: "scheduled_at", Type: "ScheduleAt"},
			{Name: "player_id", Type: "uint32"},
		},
	},
	"consume_entity_timer": {
		Name: "consume_entity_timer",
		Columns: []Column{
			{Name: "scheduled_id", Type: "uint64", PrimaryKey: true, AutoInc: true},
			{Name: "scheduled_at", Type: "ScheduleAt"},
			{Name: "consumed_entity_id", Type: "uint32"},
			{Name: "consumer_entity_id", Type: "uint32"},
		},
	},
}

// JSON Serialization for all types

// MarshalJSON implements JSON encoding for Identity
func (i Identity) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%x", i.Bytes))
}

// UnmarshalJSON implements JSON decoding for Identity
func (i *Identity) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	if len(hexStr) != 32 {
		return fmt.Errorf("invalid identity hex string length: expected 32, got %d", len(hexStr))
	}

	for idx := 0; idx < 16; idx++ {
		var b byte
		_, err := fmt.Sscanf(hexStr[idx*2:idx*2+2], "%02x", &b)
		if err != nil {
			return fmt.Errorf("invalid hex character at position %d: %w", idx*2, err)
		}
		i.Bytes[idx] = b
	}

	return nil
}

// Validation Methods

// Validate validates a Config instance
func (c *Config) Validate() error {
	if c.WorldSize == 0 {
		return fmt.Errorf("world_size must be greater than 0")
	}
	return nil
}

// Validate validates an Entity instance
func (e *Entity) Validate() error {
	if !e.Position.IsValid() {
		return fmt.Errorf("position contains invalid values")
	}
	if e.Mass == 0 {
		return fmt.Errorf("mass must be greater than 0")
	}
	return nil
}

// Validate validates a Circle instance
func (c *Circle) Validate() error {
	if !c.Direction.IsValid() {
		return fmt.Errorf("direction contains invalid values")
	}
	if c.Speed < 0 {
		return fmt.Errorf("speed cannot be negative")
	}
	return nil
}

// Validate validates a Player instance
func (p *Player) Validate() error {
	if p.Identity.IsZero() {
		return fmt.Errorf("identity cannot be zero")
	}
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
