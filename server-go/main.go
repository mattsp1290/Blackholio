package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/clockworklabs/Blackholio/server-go/types"
)

func main() {
	fmt.Println("=== Blackholio Server Go - DbVector2 Demo ===")

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

	fmt.Println("\n=== Demo completed successfully! ===")
}
