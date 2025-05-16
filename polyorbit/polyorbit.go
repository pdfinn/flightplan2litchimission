// Package polyorbit provides utilities for creating regular polygons for drone flight paths.
package polyorbit

import (
	"fmt"
)

// RegularPolygon represents a geographical figure with a specified number of sides
type RegularPolygon struct {
	// Sides is the number of sides in the polygon
	Sides int
	// Diameter is the distance across the polygon through its center
	Diameter float64
	// CameraDegrees holds the number of degrees a camera at each point of the polygon needs to point to the center
	CameraDegrees []float64
}

// NewRegularPolygon creates a new RegularPolygon with the specified number of sides and diameter.
// It calculates the camera degrees for each point of the polygon.
func NewRegularPolygon(sides int, diameter float64) *RegularPolygon {
	var cameraDegrees []float64
	for i := 0; i < sides; i++ {
		cameraDegrees = append(cameraDegrees, float64(360/sides*i))
	}

	return &RegularPolygon{
		Sides:         sides,
		Diameter:      diameter,
		CameraDegrees: cameraDegrees,
	}
}

// String returns a string representation of the RegularPolygon.
func (p RegularPolygon) String() string {
	return fmt.Sprintf("RegularPolygon with %d sides and %.2f diameter", p.Sides, p.Diameter)
}
