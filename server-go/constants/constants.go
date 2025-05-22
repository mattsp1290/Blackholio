package constants

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

// Core Game Constants
// These constants define the fundamental game mechanics and must match
// the values in the Rust and C# implementations exactly.

const (
	// Player Constants
	START_PLAYER_MASS  uint32 = 15 // Starting mass for new players
	START_PLAYER_SPEED uint32 = 10 // Base player speed

	// Food Constants
	FOOD_MASS_MIN     uint32 = 2   // Minimum mass for spawned food
	FOOD_MASS_MAX     uint32 = 4   // Maximum mass for spawned food
	TARGET_FOOD_COUNT uint32 = 600 // Target number of food entities to maintain

	// Collision and Consumption Constants
	MINIMUM_SAFE_MASS_RATIO    float32 = 0.85 // Minimum mass ratio to safely consume another entity
	MIN_OVERLAP_PCT_TO_CONSUME float32 = 0.1  // Minimum overlap percentage required to consume

	// Split Mechanics Constants
	MIN_MASS_TO_SPLIT                    uint32  = START_PLAYER_MASS * 2 // 30 - Minimum mass required to split
	MAX_CIRCLES_PER_PLAYER               uint32  = 16                    // Maximum circles a player can have
	SPLIT_RECOMBINE_DELAY_SEC            float32 = 5.0                   // Delay before circles can recombine (seconds)
	SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC float32 = 2.0                   // Time before recombine when gravity starts (seconds)
	ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT     float32 = 0.9                   // Allowed overlap percentage between split circles
	SELF_COLLISION_SPEED                 float32 = 0.05                  // Speed multiplier for circle separation (1.0 = instant)

	// World Configuration Constants
	DEFAULT_WORLD_SIZE uint64 = 1000 // Default world size for initialization

	// Timer Intervals (converted to Go durations)
	CIRCLE_DECAY_INTERVAL = 5 * time.Second        // Circle decay timer interval
	SPAWN_FOOD_INTERVAL   = 500 * time.Millisecond // Food spawning timer interval
	MOVE_PLAYERS_INTERVAL = 50 * time.Millisecond  // Player movement timer interval
)

// Configuration holds all configurable game parameters
// This allows for runtime configuration via environment variables
type Configuration struct {
	// Core Game Settings
	StartPlayerMass  uint32 `json:"start_player_mass"`
	StartPlayerSpeed uint32 `json:"start_player_speed"`
	FoodMassMin      uint32 `json:"food_mass_min"`
	FoodMassMax      uint32 `json:"food_mass_max"`
	TargetFoodCount  uint32 `json:"target_food_count"`

	// Physics Settings
	MinimumSafeMassRatio   float32 `json:"minimum_safe_mass_ratio"`
	MinOverlapPctToConsume float32 `json:"min_overlap_pct_to_consume"`

	// Split Mechanics Settings
	MinMassToSplit                  uint32  `json:"min_mass_to_split"`
	MaxCirclesPerPlayer             uint32  `json:"max_circles_per_player"`
	SplitRecombineDelaySec          float32 `json:"split_recombine_delay_sec"`
	SplitGravPullBeforeRecombineSec float32 `json:"split_grav_pull_before_recombine_sec"`
	AllowedSplitCircleOverlapPct    float32 `json:"allowed_split_circle_overlap_pct"`
	SelfCollisionSpeed              float32 `json:"self_collision_speed"`

	// World Settings
	DefaultWorldSize uint64 `json:"default_world_size"`

	// Timer Settings
	CircleDecayInterval time.Duration `json:"circle_decay_interval"`
	SpawnFoodInterval   time.Duration `json:"spawn_food_interval"`
	MovePlayersInterval time.Duration `json:"move_players_interval"`

	// Performance Settings
	EnablePerformanceLogging bool   `json:"enable_performance_logging"`
	MaxConcurrentPlayers     uint32 `json:"max_concurrent_players"`
	EnableDebugMode          bool   `json:"enable_debug_mode"`
}

// DefaultConfiguration returns a Configuration with all default values
func DefaultConfiguration() *Configuration {
	return &Configuration{
		// Core Game Settings
		StartPlayerMass:  START_PLAYER_MASS,
		StartPlayerSpeed: START_PLAYER_SPEED,
		FoodMassMin:      FOOD_MASS_MIN,
		FoodMassMax:      FOOD_MASS_MAX,
		TargetFoodCount:  TARGET_FOOD_COUNT,

		// Physics Settings
		MinimumSafeMassRatio:   MINIMUM_SAFE_MASS_RATIO,
		MinOverlapPctToConsume: MIN_OVERLAP_PCT_TO_CONSUME,

		// Split Mechanics Settings
		MinMassToSplit:                  MIN_MASS_TO_SPLIT,
		MaxCirclesPerPlayer:             MAX_CIRCLES_PER_PLAYER,
		SplitRecombineDelaySec:          SPLIT_RECOMBINE_DELAY_SEC,
		SplitGravPullBeforeRecombineSec: SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC,
		AllowedSplitCircleOverlapPct:    ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT,
		SelfCollisionSpeed:              SELF_COLLISION_SPEED,

		// World Settings
		DefaultWorldSize: DEFAULT_WORLD_SIZE,

		// Timer Settings
		CircleDecayInterval: CIRCLE_DECAY_INTERVAL,
		SpawnFoodInterval:   SPAWN_FOOD_INTERVAL,
		MovePlayersInterval: MOVE_PLAYERS_INTERVAL,

		// Performance Settings
		EnablePerformanceLogging: false,
		MaxConcurrentPlayers:     1000,
		EnableDebugMode:          false,
	}
}

// LoadFromEnvironment loads configuration values from environment variables
// Environment variables should be prefixed with "BLACKHOLIO_"
func (c *Configuration) LoadFromEnvironment() error {
	// Helper function to get environment variable with fallback
	getEnvUint32 := func(key string, fallback uint32) (uint32, error) {
		if val := os.Getenv(key); val != "" {
			parsed, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return fallback, fmt.Errorf("invalid value for %s: %v", key, err)
			}
			return uint32(parsed), nil
		}
		return fallback, nil
	}

	getEnvFloat32 := func(key string, fallback float32) (float32, error) {
		if val := os.Getenv(key); val != "" {
			parsed, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return fallback, fmt.Errorf("invalid value for %s: %v", key, err)
			}
			return float32(parsed), nil
		}
		return fallback, nil
	}

	getEnvUint64 := func(key string, fallback uint64) (uint64, error) {
		if val := os.Getenv(key); val != "" {
			parsed, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fallback, fmt.Errorf("invalid value for %s: %v", key, err)
			}
			return parsed, nil
		}
		return fallback, nil
	}

	getEnvDuration := func(key string, fallback time.Duration) (time.Duration, error) {
		if val := os.Getenv(key); val != "" {
			parsed, err := time.ParseDuration(val)
			if err != nil {
				return fallback, fmt.Errorf("invalid duration for %s: %v", key, err)
			}
			return parsed, nil
		}
		return fallback, nil
	}

	getEnvBool := func(key string, fallback bool) (bool, error) {
		if val := os.Getenv(key); val != "" {
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				return fallback, fmt.Errorf("invalid boolean for %s: %v", key, err)
			}
			return parsed, nil
		}
		return fallback, nil
	}

	var err error

	// Load core game settings
	if c.StartPlayerMass, err = getEnvUint32("BLACKHOLIO_START_PLAYER_MASS", c.StartPlayerMass); err != nil {
		return err
	}
	if c.StartPlayerSpeed, err = getEnvUint32("BLACKHOLIO_START_PLAYER_SPEED", c.StartPlayerSpeed); err != nil {
		return err
	}
	if c.FoodMassMin, err = getEnvUint32("BLACKHOLIO_FOOD_MASS_MIN", c.FoodMassMin); err != nil {
		return err
	}
	if c.FoodMassMax, err = getEnvUint32("BLACKHOLIO_FOOD_MASS_MAX", c.FoodMassMax); err != nil {
		return err
	}
	if c.TargetFoodCount, err = getEnvUint32("BLACKHOLIO_TARGET_FOOD_COUNT", c.TargetFoodCount); err != nil {
		return err
	}

	// Load physics settings
	if c.MinimumSafeMassRatio, err = getEnvFloat32("BLACKHOLIO_MINIMUM_SAFE_MASS_RATIO", c.MinimumSafeMassRatio); err != nil {
		return err
	}
	if c.MinOverlapPctToConsume, err = getEnvFloat32("BLACKHOLIO_MIN_OVERLAP_PCT_TO_CONSUME", c.MinOverlapPctToConsume); err != nil {
		return err
	}

	// Load split mechanics settings
	if c.MaxCirclesPerPlayer, err = getEnvUint32("BLACKHOLIO_MAX_CIRCLES_PER_PLAYER", c.MaxCirclesPerPlayer); err != nil {
		return err
	}
	if c.SplitRecombineDelaySec, err = getEnvFloat32("BLACKHOLIO_SPLIT_RECOMBINE_DELAY_SEC", c.SplitRecombineDelaySec); err != nil {
		return err
	}
	if c.SplitGravPullBeforeRecombineSec, err = getEnvFloat32("BLACKHOLIO_SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC", c.SplitGravPullBeforeRecombineSec); err != nil {
		return err
	}
	if c.AllowedSplitCircleOverlapPct, err = getEnvFloat32("BLACKHOLIO_ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT", c.AllowedSplitCircleOverlapPct); err != nil {
		return err
	}
	if c.SelfCollisionSpeed, err = getEnvFloat32("BLACKHOLIO_SELF_COLLISION_SPEED", c.SelfCollisionSpeed); err != nil {
		return err
	}

	// Load world settings
	if c.DefaultWorldSize, err = getEnvUint64("BLACKHOLIO_DEFAULT_WORLD_SIZE", c.DefaultWorldSize); err != nil {
		return err
	}

	// Load timer settings
	if c.CircleDecayInterval, err = getEnvDuration("BLACKHOLIO_CIRCLE_DECAY_INTERVAL", c.CircleDecayInterval); err != nil {
		return err
	}
	if c.SpawnFoodInterval, err = getEnvDuration("BLACKHOLIO_SPAWN_FOOD_INTERVAL", c.SpawnFoodInterval); err != nil {
		return err
	}
	if c.MovePlayersInterval, err = getEnvDuration("BLACKHOLIO_MOVE_PLAYERS_INTERVAL", c.MovePlayersInterval); err != nil {
		return err
	}

	// Load performance settings
	if c.EnablePerformanceLogging, err = getEnvBool("BLACKHOLIO_ENABLE_PERFORMANCE_LOGGING", c.EnablePerformanceLogging); err != nil {
		return err
	}
	if c.MaxConcurrentPlayers, err = getEnvUint32("BLACKHOLIO_MAX_CONCURRENT_PLAYERS", c.MaxConcurrentPlayers); err != nil {
		return err
	}
	if c.EnableDebugMode, err = getEnvBool("BLACKHOLIO_ENABLE_DEBUG_MODE", c.EnableDebugMode); err != nil {
		return err
	}

	// Recalculate derived values
	c.MinMassToSplit = c.StartPlayerMass * 2

	return nil
}

// Validate validates the configuration values and returns any errors
func (c *Configuration) Validate() error {
	// Validate core game settings
	if c.StartPlayerMass == 0 {
		return fmt.Errorf("start_player_mass must be greater than 0")
	}
	if c.StartPlayerSpeed == 0 {
		return fmt.Errorf("start_player_speed must be greater than 0")
	}
	if c.FoodMassMin == 0 {
		return fmt.Errorf("food_mass_min must be greater than 0")
	}
	if c.FoodMassMax < c.FoodMassMin {
		return fmt.Errorf("food_mass_max (%d) must be >= food_mass_min (%d)", c.FoodMassMax, c.FoodMassMin)
	}
	if c.TargetFoodCount == 0 {
		return fmt.Errorf("target_food_count must be greater than 0")
	}

	// Validate physics settings
	if c.MinimumSafeMassRatio <= 0 || c.MinimumSafeMassRatio > 1 {
		return fmt.Errorf("minimum_safe_mass_ratio must be between 0 and 1, got %f", c.MinimumSafeMassRatio)
	}
	if c.MinOverlapPctToConsume <= 0 || c.MinOverlapPctToConsume > 1 {
		return fmt.Errorf("min_overlap_pct_to_consume must be between 0 and 1, got %f", c.MinOverlapPctToConsume)
	}

	// Validate split mechanics settings
	if c.MaxCirclesPerPlayer == 0 {
		return fmt.Errorf("max_circles_per_player must be greater than 0")
	}
	if c.MaxCirclesPerPlayer > 64 {
		return fmt.Errorf("max_circles_per_player should not exceed 64 for performance reasons, got %d", c.MaxCirclesPerPlayer)
	}
	if c.SplitRecombineDelaySec <= 0 {
		return fmt.Errorf("split_recombine_delay_sec must be greater than 0")
	}
	if c.SplitGravPullBeforeRecombineSec < 0 {
		return fmt.Errorf("split_grav_pull_before_recombine_sec must be >= 0")
	}
	if c.SplitGravPullBeforeRecombineSec > c.SplitRecombineDelaySec {
		return fmt.Errorf("split_grav_pull_before_recombine_sec (%f) must be <= split_recombine_delay_sec (%f)",
			c.SplitGravPullBeforeRecombineSec, c.SplitRecombineDelaySec)
	}
	if c.AllowedSplitCircleOverlapPct <= 0 || c.AllowedSplitCircleOverlapPct > 1 {
		return fmt.Errorf("allowed_split_circle_overlap_pct must be between 0 and 1, got %f", c.AllowedSplitCircleOverlapPct)
	}
	if c.SelfCollisionSpeed < 0 || c.SelfCollisionSpeed > 1 {
		return fmt.Errorf("self_collision_speed must be between 0 and 1, got %f", c.SelfCollisionSpeed)
	}

	// Validate world settings
	if c.DefaultWorldSize < 100 {
		return fmt.Errorf("default_world_size must be at least 100, got %d", c.DefaultWorldSize)
	}
	if c.DefaultWorldSize > 100000 {
		return fmt.Errorf("default_world_size should not exceed 100000 for performance reasons, got %d", c.DefaultWorldSize)
	}

	// Validate timer settings
	if c.CircleDecayInterval < time.Second {
		return fmt.Errorf("circle_decay_interval should be at least 1 second for performance reasons")
	}
	if c.SpawnFoodInterval < 10*time.Millisecond {
		return fmt.Errorf("spawn_food_interval should be at least 10ms for performance reasons")
	}
	if c.MovePlayersInterval < 10*time.Millisecond {
		return fmt.Errorf("move_players_interval should be at least 10ms for performance reasons")
	}
	if c.MovePlayersInterval > time.Second {
		return fmt.Errorf("move_players_interval should not exceed 1 second for gameplay reasons")
	}

	// Validate performance settings
	if c.MaxConcurrentPlayers == 0 {
		return fmt.Errorf("max_concurrent_players must be greater than 0")
	}
	if c.MaxConcurrentPlayers > 100000 {
		return fmt.Errorf("max_concurrent_players should not exceed 100000 for performance reasons")
	}

	// Validate derived values
	if c.MinMassToSplit != c.StartPlayerMass*2 {
		return fmt.Errorf("min_mass_to_split should equal start_player_mass * 2")
	}

	return nil
}

// GetMassToSplit returns the minimum mass required to split for this configuration
func (c *Configuration) GetMassToSplit() uint32 {
	return c.StartPlayerMass * 2
}

// Global configuration instance
var globalConfig *Configuration

// SetGlobalConfiguration sets the global configuration instance
func SetGlobalConfiguration(config *Configuration) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	globalConfig = config
	return nil
}

// GetGlobalConfiguration returns the global configuration instance
// If no configuration has been set, returns the default configuration
func GetGlobalConfiguration() *Configuration {
	if globalConfig == nil {
		globalConfig = DefaultConfiguration()
	}
	return globalConfig
}

// LoadGlobalConfigurationFromEnvironment loads the global configuration from environment variables
func LoadGlobalConfigurationFromEnvironment() error {
	config := DefaultConfiguration()
	if err := config.LoadFromEnvironment(); err != nil {
		return fmt.Errorf("failed to load configuration from environment: %w", err)
	}
	return SetGlobalConfiguration(config)
}

// Mathematical Utility Functions
// These functions use the configuration values to calculate game mechanics

// MassToRadius calculates the radius of a circle based on its mass
// Formula: sqrt(mass) - matches both Rust and C# implementations
func MassToRadius(mass uint32) float32 {
	// Use math.Sqrt but we need to import math
	massFloat := float32(mass)
	return float32(math.Sqrt(float64(massFloat)))
}

// MassToMaxMoveSpeed calculates the maximum movement speed based on mass
// Formula: 2 * START_PLAYER_SPEED / (1 + sqrt(mass / START_PLAYER_MASS))
func MassToMaxMoveSpeed(mass uint32) float32 {
	config := GetGlobalConfiguration()
	startMass := float32(config.StartPlayerMass)
	startSpeed := float32(config.StartPlayerSpeed)
	massFloat := float32(mass)

	// Calculate sqrt(mass / START_PLAYER_MASS)
	ratio := massFloat / startMass
	sqrtRatio := float32(math.Sqrt(float64(ratio)))

	return 2.0 * startSpeed / (1.0 + sqrtRatio)
}

// IsValidMassForSplit checks if a mass is sufficient for splitting
func IsValidMassForSplit(mass uint32) bool {
	config := GetGlobalConfiguration()
	return mass >= config.GetMassToSplit()
}

// GetOverlapThreshold calculates the overlap threshold for consumption
func GetOverlapThreshold(radiusA, radiusB float32) float32 {
	config := GetGlobalConfiguration()
	radiusSum := (radiusA + radiusB) * (1.0 - config.MinOverlapPctToConsume)
	return radiusSum * radiusSum
}

// Documentation and Helper Functions

// GetEnvironmentVariableHelp returns help text for environment variable configuration
func GetEnvironmentVariableHelp() string {
	return `
Blackholio Game Server Configuration via Environment Variables:

Core Game Settings:
  BLACKHOLIO_START_PLAYER_MASS         Starting mass for new players (default: 15)
  BLACKHOLIO_START_PLAYER_SPEED        Base player speed (default: 10)
  BLACKHOLIO_FOOD_MASS_MIN             Minimum food mass (default: 2)
  BLACKHOLIO_FOOD_MASS_MAX             Maximum food mass (default: 4)
  BLACKHOLIO_TARGET_FOOD_COUNT         Target food count (default: 600)

Physics Settings:
  BLACKHOLIO_MINIMUM_SAFE_MASS_RATIO   Safe mass ratio for consumption (default: 0.85)
  BLACKHOLIO_MIN_OVERLAP_PCT_TO_CONSUME Overlap percentage for consumption (default: 0.1)

Split Mechanics:
  BLACKHOLIO_MAX_CIRCLES_PER_PLAYER             Max circles per player (default: 16)
  BLACKHOLIO_SPLIT_RECOMBINE_DELAY_SEC          Split recombine delay (default: 5.0)
  BLACKHOLIO_SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC Gravity pull time (default: 2.0)
  BLACKHOLIO_ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT   Split circle overlap (default: 0.9)
  BLACKHOLIO_SELF_COLLISION_SPEED               Circle separation speed (default: 0.05)

World Settings:
  BLACKHOLIO_DEFAULT_WORLD_SIZE         World size (default: 1000)

Timer Settings (use Go duration format, e.g., "5s", "500ms"):
  BLACKHOLIO_CIRCLE_DECAY_INTERVAL      Circle decay interval (default: 5s)
  BLACKHOLIO_SPAWN_FOOD_INTERVAL        Food spawn interval (default: 500ms)
  BLACKHOLIO_MOVE_PLAYERS_INTERVAL      Player move interval (default: 50ms)

Performance Settings:
  BLACKHOLIO_ENABLE_PERFORMANCE_LOGGING Enable performance logging (default: false)
  BLACKHOLIO_MAX_CONCURRENT_PLAYERS     Max concurrent players (default: 1000)
  BLACKHOLIO_ENABLE_DEBUG_MODE          Enable debug mode (default: false)

Example:
  export BLACKHOLIO_START_PLAYER_MASS=20
  export BLACKHOLIO_TARGET_FOOD_COUNT=800
  export BLACKHOLIO_ENABLE_DEBUG_MODE=true
`
}

// GetConstantsSummary returns a summary of all constants and their values
func GetConstantsSummary() string {
	config := GetGlobalConfiguration()
	return fmt.Sprintf(`
Blackholio Game Constants Summary:

Core Game Constants:
  START_PLAYER_MASS = %d
  START_PLAYER_SPEED = %d
  FOOD_MASS_MIN = %d
  FOOD_MASS_MAX = %d
  TARGET_FOOD_COUNT = %d

Physics Constants:
  MINIMUM_SAFE_MASS_RATIO = %.2f
  MIN_OVERLAP_PCT_TO_CONSUME = %.2f

Split Mechanics Constants:
  MIN_MASS_TO_SPLIT = %d (calculated: START_PLAYER_MASS * 2)
  MAX_CIRCLES_PER_PLAYER = %d
  SPLIT_RECOMBINE_DELAY_SEC = %.2f
  SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC = %.2f
  ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT = %.2f
  SELF_COLLISION_SPEED = %.2f

World Constants:
  DEFAULT_WORLD_SIZE = %d

Timer Constants:
  CIRCLE_DECAY_INTERVAL = %v
  SPAWN_FOOD_INTERVAL = %v
  MOVE_PLAYERS_INTERVAL = %v

Performance Settings:
  EnablePerformanceLogging = %v
  MaxConcurrentPlayers = %d
  EnableDebugMode = %v
`,
		config.StartPlayerMass, config.StartPlayerSpeed,
		config.FoodMassMin, config.FoodMassMax, config.TargetFoodCount,
		config.MinimumSafeMassRatio, config.MinOverlapPctToConsume,
		config.MinMassToSplit, config.MaxCirclesPerPlayer,
		config.SplitRecombineDelaySec, config.SplitGravPullBeforeRecombineSec,
		config.AllowedSplitCircleOverlapPct, config.SelfCollisionSpeed,
		config.DefaultWorldSize,
		config.CircleDecayInterval, config.SpawnFoodInterval, config.MovePlayersInterval,
		config.EnablePerformanceLogging, config.MaxConcurrentPlayers, config.EnableDebugMode,
	)
}
