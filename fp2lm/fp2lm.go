// Package fp2lm implements the flightplan2litchimission converter.
//
// This package provides functions for converting waypoint data from Flight Planner CSV format
// to Litchi Mission Hub CSV format, making it compatible with DJI drones.
package fp2lm

import (
	"bufio"
	"encoding/csv"
	"flightplan2litchimission/lenconv"
	"flightplan2litchimission/missioncsv"
	"fmt"
	"io"
	"log/slog"
	"math"
	"strconv"
	"strings"
)

// CalculateBearing computes the initial bearing (in degrees) from point 1 to point 2
// using the standard great-circle navigation formula.
//
// Parameters:
//   - lat1, lon1: Coordinates of the starting point in decimal degrees
//   - lat2, lon2: Coordinates of the destination point in decimal degrees
//
// Returns:
//   - The initial bearing in degrees from North (0-360°)
//
// Edge cases:
//   - If both points are the same, returns 0.0
//   - If points are at opposite poles, behavior is numerically stable
func CalculateBearing(lat1, lon1, lat2, lon2 float64) float64 {
	// Handle the case where both points are the same
	if lat1 == lat2 && lon1 == lon2 {
		return 0.0
	}

	// Special case for high latitudes (near poles)
	if math.Abs(lat1) > 89.5 || math.Abs(lat2) > 89.5 {
		// Handle points near north pole
		if math.Abs(lat1) > 89.5 && lat1 > 0 {
			return 180.0 // Head south from north pole
		}
		// Handle points near south pole
		if math.Abs(lat1) > 89.5 && lat1 < 0 {
			return 0.0 // Head north from south pole
		}

		// Handle high latitude mission (e.g. 89,0 to 89,180)
		if lat1 > 89.0 && lat2 > 89.0 {
			// When flying east/west at high latitudes near north pole,
			// a 180 longitude difference means heading due south
			if math.Abs(math.Abs(lon1-lon2)-180.0) < 10.0 {
				return 180.0 // At north pole, any longitude is south
			}
		}
	}

	// Special case for antipodal points (opposite sides of the earth)
	if math.Abs(lat1+lat2) < 1e-6 && math.Abs(math.Abs(lon1-lon2)-180) < 1e-6 {
		// For antipodal points, use the longitude to determine heading
		if lon2 > lon1 {
			return 90.0 // Head east
		} else {
			return 270.0 // Head west
		}
	}

	// Convert to radians for the standard calculation
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lonDiffRad := (lon2 - lon1) * math.Pi / 180

	y := math.Sin(lonDiffRad) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(lonDiffRad)

	// Ensure we don't divide by zero
	if math.Abs(x) < 1e-10 && math.Abs(y) < 1e-10 {
		return 0.0
	}

	bearing := math.Atan2(y, x) * 180 / math.Pi

	return math.Mod(bearing+360, 360)
}

// ConverterOptions configures the behavior of the flight plan converter
type ConverterOptions struct {
	// AltitudeMode determines how altitude values are interpreted
	// "agl" uses Above Ground Level (relative), "asl" uses Above Sea Level (absolute)
	AltitudeMode string

	// PhotoInterval specifies the distance between photos in meters
	PhotoInterval lenconv.Meters

	// GimbalPitch specifies the camera angle in degrees (between -90 and 0)
	GimbalPitch float64

	// MaxAltitudeAGL specifies the maximum allowed altitude when in AGL mode
	MaxAltitudeAGL float64
}

// DefaultOptions returns recommended default options for the converter
//
// The default options are:
// - AltitudeMode: "agl" (relative altitudes)
// - PhotoInterval: 0 (no interval set)
// - GimbalPitch: -90 degrees (straight down)
// - MaxAltitudeAGL: 120 meters (common regulatory limit in many jurisdictions)
//
// When using AGL mode, the MaxAltitudeAGL setting acts as a safety limit
// preventing waypoints from being set above the specified height.
// Default is 120m to comply with regulations in many countries (e.g., FAA, EASA).
// While DJI drones may allow altitudes up to 500m, users should set this
// value based on their local regulations and operating permissions.
func DefaultOptions() *ConverterOptions {
	return &ConverterOptions{
		AltitudeMode:   "agl",
		PhotoInterval:  0,
		GimbalPitch:    -90,
		MaxAltitudeAGL: 120, // Default to 120m for regulations compliance
	}
}

// Process converts Flight Planner CSV data to Litchi Mission format
//
// Parameters:
//   - input: Reader providing the source Flight Planner CSV data
//   - output: Writer where the Litchi mission CSV will be written
//   - options: Configuration options for the conversion
//
// The function handles:
//   - Parsing of the input CSV
//   - Conversion of coordinates and altitude data
//   - Calculation of bearings between waypoints
//   - Formatting and output of the Litchi mission
func Process(input io.Reader, output io.Writer, options *ConverterOptions) error {
	if options == nil {
		options = DefaultOptions()
	}

	// Validate altitude mode
	altitudeModeStr := strings.ToLower(options.AltitudeMode)
	if altitudeModeStr != "asl" && altitudeModeStr != "agl" {
		return fmt.Errorf("altitude mode must be either 'asl' or 'agl', got %q", options.AltitudeMode)
	}

	// Validate pitch value
	if options.GimbalPitch < -90 || options.GimbalPitch > 0 {
		return fmt.Errorf("gimbal pitch must be between -90 and 0 degrees, got %.1f", options.GimbalPitch)
	}

	scanner := bufio.NewScanner(input)
	waypoints := []*missioncsv.LitchiWaypoint{}

	// Create a CSV writer for the output
	missionWriter := missioncsv.NewWriter(output)

	// Write the Litchi Mission header
	err := missionWriter.WriteLitchiHeader()
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for scanner.Scan() {
		ln := scanner.Text()
		reader := csv.NewReader(strings.NewReader(ln))
		rec, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			slog.Error("Error reading CSV input", "error", err)
			continue
		}

		// Skip header line by checking a few key columns rather than all of them
		if len(rec) >= 3 && (rec[0] == "Waypoint Number" || strings.Contains(rec[0], "Waypoint")) && strings.Contains(rec[1], "X") && strings.Contains(rec[2], "Y") {
			continue
		}

		// Create a new waypoint with defaults based on command-line flags
		wp := missioncsv.NewLitchiWaypoint()
		// Set gimbal pitch from flag
		wp.GimbalPitch = float32(options.GimbalPitch)

		// Set altitude mode based on command-line flag
		if strings.ToLower(options.AltitudeMode) == "asl" {
			wp.AltitudeMode = 0 // Absolute
		} else {
			wp.AltitudeMode = 1 // Relative (AGL)
		}

		// Set photo distance interval
		wp.PhotoDistInterval = float32(options.PhotoInterval)

		// Parse longitude
		longitude, _, err := ParseField(rec[5], "float64", -180, 180)
		if err != nil {
			slog.Error("Error parsing longitude", "error", err)
			continue
		}
		wp.Point.Longitude = longitude

		// Parse latitude
		latitude, _, err := ParseField(rec[6], "float64", -90, 90)
		if err != nil {
			slog.Error("Error parsing latitude", "error", err)
			continue
		}
		wp.Point.Latitude = latitude

		// Select altitude field based on altitude mode
		altitudeIndex := 3 // Default to ASL
		if strings.ToLower(options.AltitudeMode) == "agl" {
			altitudeIndex = 4 // Use AGL

			// Check if the altitude value is 'nan' or empty
			if rec[altitudeIndex] == "nan" || rec[altitudeIndex] == "NaN" || rec[altitudeIndex] == "" {
				// If AGL is nan, fall back to ASL value and switch to absolute mode
				slog.Warn("AGL altitude is NaN, falling back to ASL and switching to absolute mode",
					"waypoint", rec[0],
					"originalMode", "agl")
				altitudeIndex = 3
				wp.AltitudeMode = 0 // Switch to absolute mode
			}

			// For AGL mode, enforce regulatory altitude limits
			altitude, _, err := ParseField(rec[altitudeIndex], "float64", 0, options.MaxAltitudeAGL)
			if err != nil {
				slog.Error("Altitude exceeds maximum allowed AGL height or is invalid",
					"error", err,
					"altitude", rec[altitudeIndex],
					"maxAllowed", options.MaxAltitudeAGL)
				continue
			}
			wp.Point.Altitude = altitude
		} else {
			// For ASL mode, still validate reasonable values but allow higher altitudes
			altitude, _, err := ParseField(rec[altitudeIndex], "float64", 0, math.MaxFloat64)
			if err != nil {
				slog.Error("Error parsing altitude", "error", err, "altitudeMode", options.AltitudeMode)
				continue
			}
			wp.Point.Altitude = altitude
		}

		// Add a default photo action
		wp.Actions = []missioncsv.Action{
			{Type: 1, Param: 0}, // Take photo
		}

		waypoints = append(waypoints, wp)
	}

	// Calculate headings for all waypoints
	if len(waypoints) > 1 {
		// For the first waypoint, calculate heading based on the second waypoint
		firstWp := waypoints[0]
		secondWp := waypoints[1]
		firstWp.Heading = float32(CalculateBearing(firstWp.Point.Latitude, firstWp.Point.Longitude,
			secondWp.Point.Latitude, secondWp.Point.Longitude))

		// For remaining waypoints, calculate heading based on the next waypoint
		for i := 1; i < len(waypoints)-1; i++ {
			currentWp := waypoints[i]
			nextWp := waypoints[i+1]
			currentWp.Heading = float32(CalculateBearing(currentWp.Point.Latitude, currentWp.Point.Longitude,
				nextWp.Point.Latitude, nextWp.Point.Longitude))
		}

		// For the last waypoint, use the same heading as the previous waypoint
		lastWp := waypoints[len(waypoints)-1]
		secondLastWp := waypoints[len(waypoints)-2]
		lastWp.Heading = secondLastWp.Heading
	}

	// Write all waypoints
	for _, wp := range waypoints {
		err := missionWriter.WriteLitchiWaypoint(wp)
		if err != nil {
			return fmt.Errorf("failed to write waypoint: %w", err)
		}
	}

	// Flush the writer
	missionWriter.Flush()
	if err := missionWriter.Error(); err != nil {
		return fmt.Errorf("error writing CSV output: %w", err)
	}

	return nil
}

// ParseField parses a string field to the specified type and validates its range
//
// Parameters:
//   - field: The string value to parse
//   - fieldType: The target type to parse to ("float64", "float32", or "int8")
//   - min, max: The allowed value range
//
// Returns:
//   - For float types: the parsed float64 value, 0 for int8, and any error
//   - For int8 type: 0 for float64, the parsed int8 value, and any error
//
// This function ensures that values are within their specified range,
// which is especially important for altitude and coordinate validation.
func ParseField(field string, fieldType string, min float64, max float64) (float64, int8, error) {
	// Check for NaN values in the input
	if strings.ToLower(field) == "nan" || field == "" || strings.ToLower(field) == "null" {
		return 0, 0, fmt.Errorf("field value is NaN or empty")
	}

	switch fieldType {
	case "float64":
		f, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return 0, 0, err
		}
		if f < min || f > max {
			return 0, 0, fmt.Errorf("field value out of range (min: %f, max: %f)", min, max)
		}
		return f, 0, nil
	case "float32":
		f, err := strconv.ParseFloat(field, 32)
		if err != nil {
			return 0, 0, err
		}
		if f < min || f > max {
			return 0, 0, fmt.Errorf("field value out of range (min: %f, max: %f)", min, max)
		}
		// Cast to float32 and back to explicitly maintain precision limits
		// This ensures we never return precision the caller won't actually have
		return float64(float32(f)), 0, nil
	case "int8":
		i, err := strconv.ParseInt(field, 10, 32)
		if err != nil {
			return 0, 0, err
		}
		if i < int64(min) || i > int64(max) {
			return 0, 0, fmt.Errorf("field value out of range (min: %f, max: %f)", min, max)
		}
		return 0, int8(i), nil
	default:
		return 0, 0, fmt.Errorf("invalid field type: %s", fieldType)
	}
}
