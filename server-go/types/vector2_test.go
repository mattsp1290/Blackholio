package types

import (
	"encoding/json"
	"math"
	"testing"
)

// Test constants
const epsilon = 1e-6

// Helper function to compare floats with epsilon
func floatEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < epsilon
}

// Helper function to compare vectors with epsilon
func vectorEqual(a, b DbVector2) bool {
	return floatEqual(a.X, b.X) && floatEqual(a.Y, b.Y)
}

func TestNewDbVector2(t *testing.T) {
	v := NewDbVector2(3.0, 4.0)
	if v.X != 3.0 || v.Y != 4.0 {
		t.Errorf("NewDbVector2(3.0, 4.0) = %v, want DbVector2{3.0, 4.0}", v)
	}
}

func TestZeroVector(t *testing.T) {
	v := Zero()
	if v.X != 0.0 || v.Y != 0.0 {
		t.Errorf("Zero() = %v, want DbVector2{0.0, 0.0}", v)
	}
}

func TestUnitVectors(t *testing.T) {
	tests := []struct {
		name     string
		vector   DbVector2
		expected DbVector2
	}{
		{"One", One(), DbVector2{1.0, 1.0}},
		{"Up", Up(), DbVector2{0.0, 1.0}},
		{"Right", Right(), DbVector2{1.0, 0.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !vectorEqual(tt.vector, tt.expected) {
				t.Errorf("%s() = %v, want %v", tt.name, tt.vector, tt.expected)
			}
		})
	}
}

func TestSqrMagnitude(t *testing.T) {
	tests := []struct {
		vector   DbVector2
		expected float32
	}{
		{DbVector2{3.0, 4.0}, 25.0},
		{DbVector2{0.0, 0.0}, 0.0},
		{DbVector2{1.0, 1.0}, 2.0},
		{DbVector2{-3.0, 4.0}, 25.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.SqrMagnitude()
			if !floatEqual(result, tt.expected) {
				t.Errorf("SqrMagnitude() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		vector   DbVector2
		expected float32
	}{
		{DbVector2{3.0, 4.0}, 5.0},
		{DbVector2{0.0, 0.0}, 0.0},
		{DbVector2{1.0, 0.0}, 1.0},
		{DbVector2{0.0, 1.0}, 1.0},
		{DbVector2{-3.0, 4.0}, 5.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.Magnitude()
			if !floatEqual(result, tt.expected) {
				t.Errorf("Magnitude() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestNormalized(t *testing.T) {
	tests := []struct {
		name     string
		vector   DbVector2
		expected DbVector2
	}{
		{"Unit X", DbVector2{5.0, 0.0}, DbVector2{1.0, 0.0}},
		{"Unit Y", DbVector2{0.0, 3.0}, DbVector2{0.0, 1.0}},
		{"3-4-5 Triangle", DbVector2{3.0, 4.0}, DbVector2{0.6, 0.8}},
		{"Zero Vector", DbVector2{0.0, 0.0}, DbVector2{0.0, 0.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vector.Normalized()
			if !vectorEqual(result, tt.expected) {
				t.Errorf("Normalized() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArithmeticOperations(t *testing.T) {
	v1 := DbVector2{2.0, 3.0}
	v2 := DbVector2{1.0, 4.0}

	// Test Add
	result := v1.Add(v2)
	expected := DbVector2{3.0, 7.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Add() = %v, want %v", result, expected)
	}

	// Test Sub
	result = v1.Sub(v2)
	expected = DbVector2{1.0, -1.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Sub() = %v, want %v", result, expected)
	}

	// Test Mul
	result = v1.Mul(2.0)
	expected = DbVector2{4.0, 6.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Mul() = %v, want %v", result, expected)
	}

	// Test Div
	result = v1.Div(2.0)
	expected = DbVector2{1.0, 1.5}
	if !vectorEqual(result, expected) {
		t.Errorf("Div() = %v, want %v", result, expected)
	}

	// Test Div by zero
	result = v1.Div(0.0)
	expected = DbVector2{0.0, 0.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Div by zero = %v, want %v", result, expected)
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		v1       DbVector2
		v2       DbVector2
		expected float32
	}{
		{DbVector2{1.0, 0.0}, DbVector2{1.0, 0.0}, 1.0},
		{DbVector2{1.0, 0.0}, DbVector2{0.0, 1.0}, 0.0},
		{DbVector2{3.0, 4.0}, DbVector2{2.0, 1.0}, 10.0},
		{DbVector2{1.0, 1.0}, DbVector2{-1.0, -1.0}, -2.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.v1.Dot(tt.v2)
			if !floatEqual(result, tt.expected) {
				t.Errorf("Dot() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestCrossProduct(t *testing.T) {
	tests := []struct {
		v1       DbVector2
		v2       DbVector2
		expected float32
	}{
		{DbVector2{1.0, 0.0}, DbVector2{0.0, 1.0}, 1.0},
		{DbVector2{0.0, 1.0}, DbVector2{1.0, 0.0}, -1.0},
		{DbVector2{3.0, 4.0}, DbVector2{2.0, 1.0}, -5.0},
		{DbVector2{1.0, 1.0}, DbVector2{2.0, 2.0}, 0.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.v1.Cross(tt.v2)
			if !floatEqual(result, tt.expected) {
				t.Errorf("Cross() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		v1       DbVector2
		v2       DbVector2
		expected float32
	}{
		{DbVector2{0.0, 0.0}, DbVector2{3.0, 4.0}, 5.0},
		{DbVector2{1.0, 1.0}, DbVector2{1.0, 1.0}, 0.0},
		{DbVector2{0.0, 0.0}, DbVector2{1.0, 0.0}, 1.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.v1.Distance(tt.v2)
			if !floatEqual(result, tt.expected) {
				t.Errorf("Distance() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestAngle(t *testing.T) {
	tests := []struct {
		vector   DbVector2
		expected float32
	}{
		{DbVector2{1.0, 0.0}, 0.0},
		{DbVector2{0.0, 1.0}, float32(math.Pi / 2)},
		{DbVector2{-1.0, 0.0}, float32(math.Pi)},
		{DbVector2{0.0, -1.0}, float32(-math.Pi / 2)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.Angle()
			if !floatEqual(result, tt.expected) {
				t.Errorf("Angle() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestAngleTo(t *testing.T) {
	v1 := DbVector2{1.0, 0.0}
	v2 := DbVector2{0.0, 1.0}

	result := v1.AngleTo(v2)
	expected := float32(math.Pi / 2)

	if !floatEqual(result, expected) {
		t.Errorf("AngleTo() = %f, want %f", result, expected)
	}
}

func TestLerp(t *testing.T) {
	v1 := DbVector2{0.0, 0.0}
	v2 := DbVector2{10.0, 10.0}

	tests := []struct {
		t        float32
		expected DbVector2
	}{
		{0.0, DbVector2{0.0, 0.0}},
		{0.5, DbVector2{5.0, 5.0}},
		{1.0, DbVector2{10.0, 10.0}},
		{-0.5, DbVector2{0.0, 0.0}},  // Clamped to 0
		{1.5, DbVector2{10.0, 10.0}}, // Clamped to 1
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := v1.Lerp(v2, tt.t)
			if !vectorEqual(result, tt.expected) {
				t.Errorf("Lerp(%f) = %v, want %v", tt.t, result, tt.expected)
			}
		})
	}
}

func TestReflect(t *testing.T) {
	// Reflect (1, 1) off a vertical surface (normal pointing right)
	v := DbVector2{1.0, 1.0}
	normal := DbVector2{1.0, 0.0}

	result := v.Reflect(normal)
	expected := DbVector2{-1.0, 1.0}

	if !vectorEqual(result, expected) {
		t.Errorf("Reflect() = %v, want %v", result, expected)
	}
}

func TestRotate(t *testing.T) {
	v := DbVector2{1.0, 0.0}

	// Rotate 90 degrees
	result := v.Rotate(float32(math.Pi / 2))
	expected := DbVector2{0.0, 1.0}

	if !vectorEqual(result, expected) {
		t.Errorf("Rotate(π/2) = %v, want %v", result, expected)
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		vector   DbVector2
		expected bool
	}{
		{DbVector2{0.0, 0.0}, true},
		{DbVector2{1e-7, 0.0}, true},  // Within epsilon
		{DbVector2{0.0, 1e-7}, true},  // Within epsilon
		{DbVector2{1e-5, 0.0}, false}, // Outside epsilon
		{DbVector2{1.0, 0.0}, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		vector   DbVector2
		expected bool
	}{
		{DbVector2{1.0, 2.0}, true},
		{DbVector2{0.0, 0.0}, true},
		{DbVector2{float32(math.NaN()), 0.0}, false},
		{DbVector2{0.0, float32(math.NaN())}, false},
		{DbVector2{float32(math.Inf(1)), 0.0}, false},
		{DbVector2{0.0, float32(math.Inf(-1))}, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v for vector %v", result, tt.expected, tt.vector)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	v := DbVector2{5.0, -5.0}
	min := DbVector2{-2.0, -2.0}
	max := DbVector2{2.0, 2.0}

	result := v.Clamp(min, max)
	expected := DbVector2{2.0, -2.0}

	if !vectorEqual(result, expected) {
		t.Errorf("Clamp() = %v, want %v", result, expected)
	}
}

func TestClampMagnitude(t *testing.T) {
	tests := []struct {
		vector      DbVector2
		maxMag      float32
		expectedMag float32
	}{
		{DbVector2{3.0, 4.0}, 10.0, 5.0}, // No clamping
		{DbVector2{3.0, 4.0}, 2.0, 2.0},  // Clamping
		{DbVector2{3.0, 4.0}, -1.0, 0.0}, // Negative max returns zero
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.vector.ClampMagnitude(tt.maxMag)
			resultMag := result.Magnitude()
			if !floatEqual(resultMag, tt.expectedMag) {
				t.Errorf("ClampMagnitude(%f) magnitude = %f, want %f", tt.maxMag, resultMag, tt.expectedMag)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	v1 := DbVector2{1.0, 2.0}
	v2 := DbVector2{1.0, 2.0}
	v3 := DbVector2{1.0000001, 2.0} // Within epsilon
	v4 := DbVector2{1.1, 2.0}       // Outside epsilon

	if !v1.Equal(v2) {
		t.Errorf("Equal() should return true for identical vectors")
	}

	if !v1.Equal(v3) {
		t.Errorf("Equal() should return true for vectors within epsilon")
	}

	if v1.Equal(v4) {
		t.Errorf("Equal() should return false for vectors outside epsilon")
	}
}

func TestString(t *testing.T) {
	v := DbVector2{1.234, 5.678}
	result := v.String()
	expected := "DbVector2(1.234, 5.678)"

	if result != expected {
		t.Errorf("String() = %s, want %s", result, expected)
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test FromAngle
	angle := float32(math.Pi / 2)
	v := FromAngle(angle)
	expected := DbVector2{0.0, 1.0}
	if !vectorEqual(v, expected) {
		t.Errorf("FromAngle(π/2) = %v, want %v", v, expected)
	}

	// Test FromPolar
	v = FromPolar(5.0, 0.0)
	expected = DbVector2{5.0, 0.0}
	if !vectorEqual(v, expected) {
		t.Errorf("FromPolar(5.0, 0.0) = %v, want %v", v, expected)
	}

	// Test Min
	v1 := DbVector2{1.0, 3.0}
	v2 := DbVector2{2.0, 1.0}
	result := Min(v1, v2)
	expected = DbVector2{1.0, 1.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Min() = %v, want %v", result, expected)
	}

	// Test Max
	result = Max(v1, v2)
	expected = DbVector2{2.0, 3.0}
	if !vectorEqual(result, expected) {
		t.Errorf("Max() = %v, want %v", result, expected)
	}

	// Test Random (deterministic)
	random := Random()
	expectedRandom := DbVector2{0.0, 1.0} // Based on our deterministic implementation
	if !vectorEqual(random, expectedRandom) {
		t.Errorf("Random() = %v, want %v", random, expectedRandom)
	}
}

func TestJSONSerialization(t *testing.T) {
	original := DbVector2{3.14, 2.71}

	// Test Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Test Unmarshal
	var decoded DbVector2
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if !vectorEqual(original, decoded) {
		t.Errorf("JSON round-trip failed: got %v, want %v", decoded, original)
	}
}

func TestBinarySerialization(t *testing.T) {
	original := DbVector2{3.14, 2.71}

	// Test MarshalBinary
	data, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	if len(data) != 8 {
		t.Errorf("MarshalBinary returned %d bytes, want 8", len(data))
	}

	// Test UnmarshalBinary
	var decoded DbVector2
	err = decoded.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if !vectorEqual(original, decoded) {
		t.Errorf("Binary round-trip failed: got %v, want %v", decoded, original)
	}
}

func TestBinarySerializationEdgeCases(t *testing.T) {
	var v DbVector2

	// Test invalid data length
	err := v.UnmarshalBinary([]byte{1, 2, 3})
	if err == nil {
		t.Error("UnmarshalBinary should fail with invalid data length")
	}

	// Test invalid values (NaN)
	invalidData := make([]byte, 8)
	nanBits := math.Float32bits(float32(math.NaN()))
	invalidData[0] = byte(nanBits)
	invalidData[1] = byte(nanBits >> 8)
	invalidData[2] = byte(nanBits >> 16)
	invalidData[3] = byte(nanBits >> 24)

	err = v.UnmarshalBinary(invalidData)
	if err == nil {
		t.Error("UnmarshalBinary should fail with NaN values")
	}
}

func TestJSONSerializationEdgeCases(t *testing.T) {
	var v DbVector2

	// Test invalid JSON
	err := v.UnmarshalJSON([]byte("invalid json"))
	if err == nil {
		t.Error("UnmarshalJSON should fail with invalid JSON")
	}

	// Test invalid values (NaN)
	nanJSON := `{"x": NaN, "y": 1.0}`
	err = v.UnmarshalJSON([]byte(nanJSON))
	if err == nil {
		t.Error("UnmarshalJSON should fail with NaN values")
	}
}

// Benchmark tests
func BenchmarkMagnitude(b *testing.B) {
	v := DbVector2{3.0, 4.0}
	for i := 0; i < b.N; i++ {
		_ = v.Magnitude()
	}
}

func BenchmarkNormalized(b *testing.B) {
	v := DbVector2{3.0, 4.0}
	for i := 0; i < b.N; i++ {
		_ = v.Normalized()
	}
}

func BenchmarkDotProduct(b *testing.B) {
	v1 := DbVector2{3.0, 4.0}
	v2 := DbVector2{1.0, 2.0}
	for i := 0; i < b.N; i++ {
		_ = v1.Dot(v2)
	}
}

func BenchmarkBinarySerialization(b *testing.B) {
	v := DbVector2{3.14, 2.71}
	for i := 0; i < b.N; i++ {
		data, _ := v.MarshalBinary()
		var decoded DbVector2
		_ = decoded.UnmarshalBinary(data)
	}
}

func BenchmarkJSONSerialization(b *testing.B) {
	v := DbVector2{3.14, 2.71}
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(v)
		var decoded DbVector2
		_ = json.Unmarshal(data, &decoded)
	}
}
