package constants

import (
	"math"
	"os"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	t.Run("CoreGameConstants", func(t *testing.T) {
		// Test that constants match expected values from Rust and C# implementations
		if START_PLAYER_MASS != 15 {
			t.Errorf("START_PLAYER_MASS = %d, want 15", START_PLAYER_MASS)
		}
		if START_PLAYER_SPEED != 10 {
			t.Errorf("START_PLAYER_SPEED = %d, want 10", START_PLAYER_SPEED)
		}
		if FOOD_MASS_MIN != 2 {
			t.Errorf("FOOD_MASS_MIN = %d, want 2", FOOD_MASS_MIN)
		}
		if FOOD_MASS_MAX != 4 {
			t.Errorf("FOOD_MASS_MAX = %d, want 4", FOOD_MASS_MAX)
		}
		if TARGET_FOOD_COUNT != 600 {
			t.Errorf("TARGET_FOOD_COUNT = %d, want 600", TARGET_FOOD_COUNT)
		}
	})

	t.Run("PhysicsConstants", func(t *testing.T) {
		if MINIMUM_SAFE_MASS_RATIO != 0.85 {
			t.Errorf("MINIMUM_SAFE_MASS_RATIO = %f, want 0.85", MINIMUM_SAFE_MASS_RATIO)
		}
		if MIN_OVERLAP_PCT_TO_CONSUME != 0.1 {
			t.Errorf("MIN_OVERLAP_PCT_TO_CONSUME = %f, want 0.1", MIN_OVERLAP_PCT_TO_CONSUME)
		}
	})

	t.Run("SplitMechanicsConstants", func(t *testing.T) {
		expectedMinMassToSplit := START_PLAYER_MASS * 2
		if MIN_MASS_TO_SPLIT != expectedMinMassToSplit {
			t.Errorf("MIN_MASS_TO_SPLIT = %d, want %d", MIN_MASS_TO_SPLIT, expectedMinMassToSplit)
		}
		if MAX_CIRCLES_PER_PLAYER != 16 {
			t.Errorf("MAX_CIRCLES_PER_PLAYER = %d, want 16", MAX_CIRCLES_PER_PLAYER)
		}
		if SPLIT_RECOMBINE_DELAY_SEC != 5.0 {
			t.Errorf("SPLIT_RECOMBINE_DELAY_SEC = %f, want 5.0", SPLIT_RECOMBINE_DELAY_SEC)
		}
		if SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC != 2.0 {
			t.Errorf("SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC = %f, want 2.0", SPLIT_GRAV_PULL_BEFORE_RECOMBINE_SEC)
		}
		if ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT != 0.9 {
			t.Errorf("ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT = %f, want 0.9", ALLOWED_SPLIT_CIRCLE_OVERLAP_PCT)
		}
		if SELF_COLLISION_SPEED != 0.05 {
			t.Errorf("SELF_COLLISION_SPEED = %f, want 0.05", SELF_COLLISION_SPEED)
		}
	})

	t.Run("WorldConstants", func(t *testing.T) {
		if DEFAULT_WORLD_SIZE != 1000 {
			t.Errorf("DEFAULT_WORLD_SIZE = %d, want 1000", DEFAULT_WORLD_SIZE)
		}
	})

	t.Run("TimerConstants", func(t *testing.T) {
		if CIRCLE_DECAY_INTERVAL != 5*time.Second {
			t.Errorf("CIRCLE_DECAY_INTERVAL = %v, want %v", CIRCLE_DECAY_INTERVAL, 5*time.Second)
		}
		if SPAWN_FOOD_INTERVAL != 500*time.Millisecond {
			t.Errorf("SPAWN_FOOD_INTERVAL = %v, want %v", SPAWN_FOOD_INTERVAL, 500*time.Millisecond)
		}
		if MOVE_PLAYERS_INTERVAL != 50*time.Millisecond {
			t.Errorf("MOVE_PLAYERS_INTERVAL = %v, want %v", MOVE_PLAYERS_INTERVAL, 50*time.Millisecond)
		}
	})
}

func TestDefaultConfiguration(t *testing.T) {
	config := DefaultConfiguration()

	// Test core game settings
	if config.StartPlayerMass != START_PLAYER_MASS {
		t.Errorf("StartPlayerMass = %d, want %d", config.StartPlayerMass, START_PLAYER_MASS)
	}
	if config.StartPlayerSpeed != START_PLAYER_SPEED {
		t.Errorf("StartPlayerSpeed = %d, want %d", config.StartPlayerSpeed, START_PLAYER_SPEED)
	}
	if config.FoodMassMin != FOOD_MASS_MIN {
		t.Errorf("FoodMassMin = %d, want %d", config.FoodMassMin, FOOD_MASS_MIN)
	}
	if config.FoodMassMax != FOOD_MASS_MAX {
		t.Errorf("FoodMassMax = %d, want %d", config.FoodMassMax, FOOD_MASS_MAX)
	}
	if config.TargetFoodCount != TARGET_FOOD_COUNT {
		t.Errorf("TargetFoodCount = %d, want %d", config.TargetFoodCount, TARGET_FOOD_COUNT)
	}

	// Test physics settings
	if config.MinimumSafeMassRatio != MINIMUM_SAFE_MASS_RATIO {
		t.Errorf("MinimumSafeMassRatio = %f, want %f", config.MinimumSafeMassRatio, MINIMUM_SAFE_MASS_RATIO)
	}
	if config.MinOverlapPctToConsume != MIN_OVERLAP_PCT_TO_CONSUME {
		t.Errorf("MinOverlapPctToConsume = %f, want %f", config.MinOverlapPctToConsume, MIN_OVERLAP_PCT_TO_CONSUME)
	}

	// Test derived values
	expectedMinMassToSplit := config.StartPlayerMass * 2
	if config.MinMassToSplit != expectedMinMassToSplit {
		t.Errorf("MinMassToSplit = %d, want %d", config.MinMassToSplit, expectedMinMassToSplit)
	}
}

func TestConfigurationValidation(t *testing.T) {
	t.Run("ValidConfiguration", func(t *testing.T) {
		config := DefaultConfiguration()
		if err := config.Validate(); err != nil {
			t.Errorf("Valid configuration should not error: %v", err)
		}
	})

	t.Run("InvalidStartPlayerMass", func(t *testing.T) {
		config := DefaultConfiguration()
		config.StartPlayerMass = 0
		if err := config.Validate(); err == nil {
			t.Error("Should error with zero start player mass")
		}
	})

	t.Run("InvalidFoodMassRange", func(t *testing.T) {
		config := DefaultConfiguration()
		config.FoodMassMax = 1
		config.FoodMassMin = 2
		if err := config.Validate(); err == nil {
			t.Error("Should error when food_mass_max < food_mass_min")
		}
	})

	t.Run("InvalidMassRatio", func(t *testing.T) {
		config := DefaultConfiguration()
		config.MinimumSafeMassRatio = 1.5
		if err := config.Validate(); err == nil {
			t.Error("Should error with mass ratio > 1.0")
		}

		config.MinimumSafeMassRatio = -0.1
		if err := config.Validate(); err == nil {
			t.Error("Should error with negative mass ratio")
		}
	})

	t.Run("InvalidCircleCount", func(t *testing.T) {
		config := DefaultConfiguration()
		config.MaxCirclesPerPlayer = 0
		if err := config.Validate(); err == nil {
			t.Error("Should error with zero max circles")
		}

		config.MaxCirclesPerPlayer = 100
		if err := config.Validate(); err == nil {
			t.Error("Should error with excessive max circles")
		}
	})

	t.Run("InvalidWorldSize", func(t *testing.T) {
		config := DefaultConfiguration()
		config.DefaultWorldSize = 50
		if err := config.Validate(); err == nil {
			t.Error("Should error with too small world size")
		}

		config.DefaultWorldSize = 200000
		if err := config.Validate(); err == nil {
			t.Error("Should error with too large world size")
		}
	})

	t.Run("InvalidTimerIntervals", func(t *testing.T) {
		config := DefaultConfiguration()
		config.MovePlayersInterval = 5 * time.Millisecond
		if err := config.Validate(); err == nil {
			t.Error("Should error with too short move interval")
		}

		config.MovePlayersInterval = 2 * time.Second
		if err := config.Validate(); err == nil {
			t.Error("Should error with too long move interval")
		}
	})

	t.Run("InvalidSplitTimings", func(t *testing.T) {
		config := DefaultConfiguration()
		config.SplitGravPullBeforeRecombineSec = 10.0
		config.SplitRecombineDelaySec = 5.0
		if err := config.Validate(); err == nil {
			t.Error("Should error when grav pull time > recombine delay")
		}
	})
}

func TestEnvironmentVariableLoading(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"BLACKHOLIO_START_PLAYER_MASS",
		"BLACKHOLIO_START_PLAYER_SPEED",
		"BLACKHOLIO_FOOD_MASS_MIN",
		"BLACKHOLIO_TARGET_FOOD_COUNT",
		"BLACKHOLIO_MINIMUM_SAFE_MASS_RATIO",
		"BLACKHOLIO_MAX_CIRCLES_PER_PLAYER",
		"BLACKHOLIO_DEFAULT_WORLD_SIZE",
		"BLACKHOLIO_CIRCLE_DECAY_INTERVAL",
		"BLACKHOLIO_ENABLE_DEBUG_MODE",
	}

	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			originalEnv[envVar] = val
		}
	}

	// Clean up environment
	defer func() {
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}
		for envVar, val := range originalEnv {
			os.Setenv(envVar, val)
		}
	}()

	t.Run("LoadFromEnvironment", func(t *testing.T) {
		// Set test environment variables
		os.Setenv("BLACKHOLIO_START_PLAYER_MASS", "20")
		os.Setenv("BLACKHOLIO_START_PLAYER_SPEED", "15")
		os.Setenv("BLACKHOLIO_FOOD_MASS_MIN", "3")
		os.Setenv("BLACKHOLIO_TARGET_FOOD_COUNT", "800")
		os.Setenv("BLACKHOLIO_MINIMUM_SAFE_MASS_RATIO", "0.9")
		os.Setenv("BLACKHOLIO_MAX_CIRCLES_PER_PLAYER", "20")
		os.Setenv("BLACKHOLIO_DEFAULT_WORLD_SIZE", "2000")
		os.Setenv("BLACKHOLIO_CIRCLE_DECAY_INTERVAL", "10s")
		os.Setenv("BLACKHOLIO_ENABLE_DEBUG_MODE", "true")

		config := DefaultConfiguration()
		err := config.LoadFromEnvironment()
		if err != nil {
			t.Fatalf("LoadFromEnvironment failed: %v", err)
		}

		// Verify values were loaded
		if config.StartPlayerMass != 20 {
			t.Errorf("StartPlayerMass = %d, want 20", config.StartPlayerMass)
		}
		if config.StartPlayerSpeed != 15 {
			t.Errorf("StartPlayerSpeed = %d, want 15", config.StartPlayerSpeed)
		}
		if config.FoodMassMin != 3 {
			t.Errorf("FoodMassMin = %d, want 3", config.FoodMassMin)
		}
		if config.TargetFoodCount != 800 {
			t.Errorf("TargetFoodCount = %d, want 800", config.TargetFoodCount)
		}
		if config.MinimumSafeMassRatio != 0.9 {
			t.Errorf("MinimumSafeMassRatio = %f, want 0.9", config.MinimumSafeMassRatio)
		}
		if config.MaxCirclesPerPlayer != 20 {
			t.Errorf("MaxCirclesPerPlayer = %d, want 20", config.MaxCirclesPerPlayer)
		}
		if config.DefaultWorldSize != 2000 {
			t.Errorf("DefaultWorldSize = %d, want 2000", config.DefaultWorldSize)
		}
		if config.CircleDecayInterval != 10*time.Second {
			t.Errorf("CircleDecayInterval = %v, want %v", config.CircleDecayInterval, 10*time.Second)
		}
		if !config.EnableDebugMode {
			t.Errorf("EnableDebugMode = %v, want true", config.EnableDebugMode)
		}

		// Verify derived values are recalculated
		if config.MinMassToSplit != 40 { // 20 * 2
			t.Errorf("MinMassToSplit = %d, want 40", config.MinMassToSplit)
		}
	})

	t.Run("InvalidEnvironmentValues", func(t *testing.T) {
		os.Setenv("BLACKHOLIO_START_PLAYER_MASS", "invalid")

		config := DefaultConfiguration()
		err := config.LoadFromEnvironment()
		if err == nil {
			t.Error("Should error with invalid environment value")
		}
	})

	t.Run("InvalidDurationFormat", func(t *testing.T) {
		os.Setenv("BLACKHOLIO_CIRCLE_DECAY_INTERVAL", "invalid_duration")

		config := DefaultConfiguration()
		err := config.LoadFromEnvironment()
		if err == nil {
			t.Error("Should error with invalid duration format")
		}
	})
}

func TestGlobalConfiguration(t *testing.T) {
	// Save original global config
	originalConfig := globalConfig
	defer func() {
		globalConfig = originalConfig
	}()

	t.Run("GetGlobalConfiguration", func(t *testing.T) {
		globalConfig = nil
		config := GetGlobalConfiguration()
		if config == nil {
			t.Error("Should return default configuration when none set")
		}
		if config.StartPlayerMass != START_PLAYER_MASS {
			t.Errorf("Should return default values")
		}
	})

	t.Run("SetGlobalConfiguration", func(t *testing.T) {
		customConfig := DefaultConfiguration()
		customConfig.StartPlayerMass = 25
		customConfig.MinMassToSplit = customConfig.StartPlayerMass * 2 // Fix derived value

		err := SetGlobalConfiguration(customConfig)
		if err != nil {
			t.Fatalf("SetGlobalConfiguration failed: %v", err)
		}

		retrieved := GetGlobalConfiguration()
		if retrieved.StartPlayerMass != 25 {
			t.Errorf("Global config not set properly")
		}
	})

	t.Run("SetInvalidGlobalConfiguration", func(t *testing.T) {
		invalidConfig := DefaultConfiguration()
		invalidConfig.StartPlayerMass = 0

		err := SetGlobalConfiguration(invalidConfig)
		if err == nil {
			t.Error("Should error when setting invalid configuration")
		}
	})
}

func TestMathematicalFunctions(t *testing.T) {
	t.Run("MassToRadius", func(t *testing.T) {
		tests := []struct {
			mass     uint32
			expected float32
		}{
			{1, 1.0},
			{4, 2.0},
			{9, 3.0},
			{16, 4.0},
			{25, 5.0},
		}

		for _, tt := range tests {
			result := MassToRadius(tt.mass)
			if math.Abs(float64(result-tt.expected)) > 0.001 {
				t.Errorf("MassToRadius(%d) = %f, want %f", tt.mass, result, tt.expected)
			}
		}
	})

	t.Run("MassToMaxMoveSpeed", func(t *testing.T) {
		// Test with default configuration
		config := DefaultConfiguration()
		SetGlobalConfiguration(config)

		// Test some expected values
		startMassSpeed := MassToMaxMoveSpeed(START_PLAYER_MASS)
		if startMassSpeed != 10.0 { // Should be START_PLAYER_SPEED when mass == START_PLAYER_MASS
			t.Errorf("MassToMaxMoveSpeed(%d) = %f, want 10.0", START_PLAYER_MASS, startMassSpeed)
		}

		// Test that larger mass results in slower speed
		largerMassSpeed := MassToMaxMoveSpeed(START_PLAYER_MASS * 4)
		if largerMassSpeed >= startMassSpeed {
			t.Errorf("Larger mass should result in slower speed: %f >= %f", largerMassSpeed, startMassSpeed)
		}

		// Test that smaller mass results in faster speed
		smallerMassSpeed := MassToMaxMoveSpeed(START_PLAYER_MASS / 4)
		if smallerMassSpeed <= startMassSpeed {
			t.Errorf("Smaller mass should result in faster speed: %f <= %f", smallerMassSpeed, startMassSpeed)
		}
	})

	t.Run("IsValidMassForSplit", func(t *testing.T) {
		config := DefaultConfiguration()
		SetGlobalConfiguration(config)

		// Test valid mass for split
		if !IsValidMassForSplit(MIN_MASS_TO_SPLIT) {
			t.Errorf("Should be valid mass for split: %d", MIN_MASS_TO_SPLIT)
		}
		if !IsValidMassForSplit(MIN_MASS_TO_SPLIT + 10) {
			t.Errorf("Should be valid mass for split: %d", MIN_MASS_TO_SPLIT+10)
		}

		// Test invalid mass for split
		if IsValidMassForSplit(MIN_MASS_TO_SPLIT - 1) {
			t.Errorf("Should not be valid mass for split: %d", MIN_MASS_TO_SPLIT-1)
		}
		if IsValidMassForSplit(START_PLAYER_MASS) {
			t.Errorf("Should not be valid mass for split: %d", START_PLAYER_MASS)
		}
	})

	t.Run("GetOverlapThreshold", func(t *testing.T) {
		config := DefaultConfiguration()
		SetGlobalConfiguration(config)

		radiusA := float32(5.0)
		radiusB := float32(3.0)
		threshold := GetOverlapThreshold(radiusA, radiusB)

		// Should be (radiusA + radiusB) * (1 - MIN_OVERLAP_PCT_TO_CONSUME) squared
		expected := (radiusA + radiusB) * (1.0 - config.MinOverlapPctToConsume)
		expected = expected * expected

		if math.Abs(float64(threshold-expected)) > 0.001 {
			t.Errorf("GetOverlapThreshold(%f, %f) = %f, want %f", radiusA, radiusB, threshold, expected)
		}
	})
}

func TestConfigurationHelpers(t *testing.T) {
	t.Run("GetMassToSplit", func(t *testing.T) {
		config := DefaultConfiguration()
		config.StartPlayerMass = 20

		expected := uint32(40) // 20 * 2
		if config.GetMassToSplit() != expected {
			t.Errorf("GetMassToSplit() = %d, want %d", config.GetMassToSplit(), expected)
		}
	})
}

func TestDocumentationFunctions(t *testing.T) {
	t.Run("GetEnvironmentVariableHelp", func(t *testing.T) {
		help := GetEnvironmentVariableHelp()
		if help == "" {
			t.Error("Help text should not be empty")
		}
		if !contains(help, "BLACKHOLIO_START_PLAYER_MASS") {
			t.Error("Help should contain environment variable names")
		}
		if !contains(help, "default:") {
			t.Error("Help should contain default values")
		}
	})

	t.Run("GetConstantsSummary", func(t *testing.T) {
		config := DefaultConfiguration()
		SetGlobalConfiguration(config)

		summary := GetConstantsSummary()
		if summary == "" {
			t.Error("Summary should not be empty")
		}
		if !contains(summary, "START_PLAYER_MASS = 15") {
			t.Error("Summary should contain constant values")
		}
		if !contains(summary, "MINIMUM_SAFE_MASS_RATIO = 0.85") {
			t.Error("Summary should contain physics constants")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && contains(s[1:], substr)) || s[:len(substr)] == substr)
}

// Benchmark tests
func BenchmarkMassToRadius(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MassToRadius(uint32(i%1000 + 1))
	}
}

func BenchmarkMassToMaxMoveSpeed(b *testing.B) {
	config := DefaultConfiguration()
	SetGlobalConfiguration(config)

	for i := 0; i < b.N; i++ {
		_ = MassToMaxMoveSpeed(uint32(i%1000 + 1))
	}
}

func BenchmarkConfigurationValidation(b *testing.B) {
	config := DefaultConfiguration()

	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}
