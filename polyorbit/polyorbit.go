package polyorbit

import (
	"fmt"
	"os"
	"strconv"
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

func main() {
	// Check that a number of sides and a diameter were provided as arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: polygon <sides> <diameter>")
		os.Exit(1)
	}

	// Parse the number of sides argument as an integer
	sides, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing number of sides: %v\n", err)
		os.Exit(1)
	}

	// Parse the diameter argument as a floating point value
	diameter, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Printf("Error parsing diameter: %v\n", err)
		os.Exit(1)
	}

	// Calculate the camera degrees for each point of the polygon
	var cameraDegrees []float64
	for i := 0; i < sides; i++ {
		cameraDegrees = append(cameraDegrees, float64(360/sides*i))
	}

	// Create a new RegularPolygon struct
	polygon := RegularPolygon{
		Sides:         sides,
		Diameter:      diameter,
		CameraDegrees: cameraDegrees,
	}

	// Print the polygon
	fmt.Println(polygon)
}
