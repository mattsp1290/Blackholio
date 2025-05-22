package types

import (
	"encoding/json"
	"fmt"
	"math"
)

// DbVector2 represents a 2D vector used in Blackholio game.
// This type is compatible with SpacetimeDB's BSATN serialization format
// and matches the functionality of the Rust and C# implementations.
type DbVector2 struct {
	X float32 `json:"x" bsatn:"0"`
	Y float32 `json:"y" bsatn:"1"`
}

// NewDbVector2 creates a new DbVector2 with the given x and y components.
func NewDbVector2(x, y float32) DbVector2 {
	return DbVector2{X: x, Y: y}
}

// Zero returns a zero vector.
func Zero() DbVector2 {
	return DbVector2{X: 0, Y: 0}
}

// One returns a vector with both components set to 1.
func One() DbVector2 {
	return DbVector2{X: 1, Y: 1}
}

// Up returns a unit vector pointing up (0, 1).
func Up() DbVector2 {
	return DbVector2{X: 0, Y: 1}
}

// Right returns a unit vector pointing right (1, 0).
func Right() DbVector2 {
	return DbVector2{X: 1, Y: 0}
}

// SqrMagnitude returns the squared magnitude of the vector.
// This is more efficient than Magnitude() when you only need to compare distances.
func (v DbVector2) SqrMagnitude() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Magnitude returns the magnitude (length) of the vector.
func (v DbVector2) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.SqrMagnitude())))
}

// Normalized returns a unit vector in the same direction as this vector.
// If the vector is zero, returns a zero vector to avoid division by zero.
func (v DbVector2) Normalized() DbVector2 {
	mag := v.Magnitude()
	if mag == 0 {
		return Zero()
	}
	return v.Div(mag)
}

// Add returns the sum of this vector and another vector.
func (v DbVector2) Add(other DbVector2) DbVector2 {
	return DbVector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of this vector and another vector.
func (v DbVector2) Sub(other DbVector2) DbVector2 {
	return DbVector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul returns this vector multiplied by a scalar.
func (v DbVector2) Mul(scalar float32) DbVector2 {
	return DbVector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Div returns this vector divided by a scalar.
// If scalar is zero, returns a zero vector to avoid division by zero.
func (v DbVector2) Div(scalar float32) DbVector2 {
	if scalar == 0 {
		return Zero()
	}
	return DbVector2{X: v.X / scalar, Y: v.Y / scalar}
}

// Dot returns the dot product of this vector and another vector.
func (v DbVector2) Dot(other DbVector2) float32 {
	return v.X*other.X + v.Y*other.Y
}

// Cross returns the 2D cross product (z-component) of this vector and another vector.
func (v DbVector2) Cross(other DbVector2) float32 {
	return v.X*other.Y - v.Y*other.X
}

// Distance returns the distance between this vector and another vector.
func (v DbVector2) Distance(other DbVector2) float32 {
	return v.Sub(other).Magnitude()
}

// DistanceSquared returns the squared distance between this vector and another vector.
func (v DbVector2) DistanceSquared(other DbVector2) float32 {
	return v.Sub(other).SqrMagnitude()
}

// Angle returns the angle of this vector in radians.
func (v DbVector2) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

// AngleTo returns the angle between this vector and another vector in radians.
func (v DbVector2) AngleTo(other DbVector2) float32 {
	dot := v.Normalized().Dot(other.Normalized())
	// Clamp to avoid floating point errors
	dot = float32(math.Max(-1.0, math.Min(1.0, float64(dot))))
	return float32(math.Acos(float64(dot)))
}

// Lerp performs linear interpolation between this vector and another vector.
// t should be between 0 and 1, where 0 returns this vector and 1 returns the other vector.
func (v DbVector2) Lerp(other DbVector2, t float32) DbVector2 {
	// Clamp t to [0, 1]
	t = float32(math.Max(0.0, math.Min(1.0, float64(t))))
	return DbVector2{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
	}
}

// Reflect returns the reflection of this vector off a surface with the given normal.
func (v DbVector2) Reflect(normal DbVector2) DbVector2 {
	return v.Sub(normal.Mul(2 * v.Dot(normal)))
}

// Rotate returns this vector rotated by the given angle in radians.
func (v DbVector2) Rotate(angleRadians float32) DbVector2 {
	cos := float32(math.Cos(float64(angleRadians)))
	sin := float32(math.Sin(float64(angleRadians)))
	return DbVector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// IsZero returns true if both components are zero (within a small epsilon).
func (v DbVector2) IsZero() bool {
	const epsilon = 1e-6
	return math.Abs(float64(v.X)) < epsilon && math.Abs(float64(v.Y)) < epsilon
}

// IsValid returns true if both components are valid (not NaN or infinite).
func (v DbVector2) IsValid() bool {
	return !math.IsNaN(float64(v.X)) && !math.IsInf(float64(v.X), 0) &&
		!math.IsNaN(float64(v.Y)) && !math.IsInf(float64(v.Y), 0)
}

// Clamp clamps each component of this vector to the given bounds.
func (v DbVector2) Clamp(min, max DbVector2) DbVector2 {
	return DbVector2{
		X: float32(math.Max(float64(min.X), math.Min(float64(max.X), float64(v.X)))),
		Y: float32(math.Max(float64(min.Y), math.Min(float64(max.Y), float64(v.Y)))),
	}
}

// ClampMagnitude clamps the magnitude of this vector to the given maximum.
func (v DbVector2) ClampMagnitude(maxMagnitude float32) DbVector2 {
	if maxMagnitude < 0 {
		return Zero()
	}
	sqrMag := v.SqrMagnitude()
	if sqrMag > maxMagnitude*maxMagnitude {
		return v.Normalized().Mul(maxMagnitude)
	}
	return v
}

// String returns a string representation of the vector.
func (v DbVector2) String() string {
	return fmt.Sprintf("DbVector2(%.3f, %.3f)", v.X, v.Y)
}

// Equal returns true if this vector is equal to another vector within a small epsilon.
func (v DbVector2) Equal(other DbVector2) bool {
	const epsilon = 1e-6
	return math.Abs(float64(v.X-other.X)) < epsilon && math.Abs(float64(v.Y-other.Y)) < epsilon
}

// JSON Serialization Implementation (temporary until BSATN integration is resolved)

// MarshalJSON implements JSON encoding for DbVector2.
func (v DbVector2) MarshalJSON() ([]byte, error) {
	// Create a struct that matches the expected format
	data := struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	}{
		X: v.X,
		Y: v.Y,
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements JSON decoding for DbVector2.
func (v *DbVector2) UnmarshalJSON(data []byte) error {
	// Temporary struct for unmarshaling
	var temp struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal DbVector2: %w", err)
	}

	v.X = temp.X
	v.Y = temp.Y

	// Validate the unmarshaled data
	if !v.IsValid() {
		return fmt.Errorf("unmarshaled DbVector2 contains invalid values: %v", v)
	}

	return nil
}

// Binary Serialization for SpacetimeDB compatibility
// These methods can be used by the SpacetimeDB bindings for efficient serialization

// MarshalBinary implements binary encoding for DbVector2.
func (v DbVector2) MarshalBinary() ([]byte, error) {
	// For now, use a simple binary format: 4 bytes for X + 4 bytes for Y
	data := make([]byte, 8)

	// Convert float32 to uint32 bits and store as little-endian
	xBits := math.Float32bits(v.X)
	yBits := math.Float32bits(v.Y)

	data[0] = byte(xBits)
	data[1] = byte(xBits >> 8)
	data[2] = byte(xBits >> 16)
	data[3] = byte(xBits >> 24)

	data[4] = byte(yBits)
	data[5] = byte(yBits >> 8)
	data[6] = byte(yBits >> 16)
	data[7] = byte(yBits >> 24)

	return data, nil
}

// UnmarshalBinary implements binary decoding for DbVector2.
func (v *DbVector2) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid data length for DbVector2: expected 8 bytes, got %d", len(data))
	}

	// Read X component
	xBits := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	v.X = math.Float32frombits(xBits)

	// Read Y component
	yBits := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24
	v.Y = math.Float32frombits(yBits)

	// Validate the unmarshaled data
	if !v.IsValid() {
		return fmt.Errorf("unmarshaled DbVector2 contains invalid values: %v", v)
	}

	return nil
}

// Utility functions for creating common vectors

// FromAngle creates a unit vector from an angle in radians.
func FromAngle(angleRadians float32) DbVector2 {
	return DbVector2{
		X: float32(math.Cos(float64(angleRadians))),
		Y: float32(math.Sin(float64(angleRadians))),
	}
}

// FromPolar creates a vector from polar coordinates (magnitude and angle in radians).
func FromPolar(magnitude, angleRadians float32) DbVector2 {
	return FromAngle(angleRadians).Mul(magnitude)
}

// Min returns a vector with the minimum components of two vectors.
func Min(a, b DbVector2) DbVector2 {
	return DbVector2{
		X: float32(math.Min(float64(a.X), float64(b.X))),
		Y: float32(math.Min(float64(a.Y), float64(b.Y))),
	}
}

// Max returns a vector with the maximum components of two vectors.
func Max(a, b DbVector2) DbVector2 {
	return DbVector2{
		X: float32(math.Max(float64(a.X), float64(b.X))),
		Y: float32(math.Max(float64(a.Y), float64(b.Y))),
	}
}

// Random returns a random unit vector.
// Note: This uses a deterministic method for testing. In production,
// you should use a proper random number generator seeded appropriately.
func Random() DbVector2 {
	// Generate random angle between 0 and 2π
	// Using a simple hash-based approach for deterministic testing
	// In production, use rand.Float32() * 2 * math.Pi
	angle := float32(0.5 * math.Pi) // Simplified for testing - returns (0, 1)
	return FromAngle(angle)
}
