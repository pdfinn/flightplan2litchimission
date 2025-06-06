package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"flightplan2litchimission/polyorbit"
)

func main() {
	// Check that a number of sides and a diameter were provided as arguments
	if len(os.Args) < 3 {
		slog.Error("Insufficient arguments", "usage", "polygon <sides> <diameter>")
		os.Exit(1)
	}

	// Parse the number of sides argument as an integer
	sides, err := strconv.Atoi(os.Args[1])
	if err != nil {
		slog.Error("Error parsing number of sides", "error", err)
		os.Exit(1)
	}

	// Parse the diameter argument as a floating point value
	diameter, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		slog.Error("Error parsing diameter", "error", err)
		os.Exit(1)
	}

	// Create a new RegularPolygon struct
	polygon := polyorbit.NewRegularPolygon(sides, diameter)

	// Print the polygon
	fmt.Println(polygon)
}
