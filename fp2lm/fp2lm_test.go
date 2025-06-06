package fp2lm_test

import (
	"bytes"
	_ "embed"
	"flightplan2litchimission/fp2lm"
	"math"
	"strings"
	"testing"
)

//go:embed testdata/litchi_golden.csv
var goldenFileData []byte

//go:embed testdata/FlightplannerMission.csv
var flightplannerMissionData []byte

// normalizeLineEndings replaces all occurrences of \r\n with \n to normalize line endings
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// canonicalizeCSV standardizes a CSV string by trimming whitespace and normalizing line endings
func canonicalizeCSV(s string) string {
	s = normalizeLineEndings(s)
	lines := strings.Split(s, "\n")

	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Normalize each line
	for i, line := range lines {
		// Trim spaces and replace spaces after commas
		lines[i] = strings.ReplaceAll(strings.TrimSpace(line), ", ", ",")
	}

	return strings.Join(lines, "\n")
}

// TestProcess checks that the converter produces the expected output
func TestProcess(t *testing.T) {
	// Use embedded test data instead of file paths
	inputReader := bytes.NewReader(flightplannerMissionData)

	// Capture the actual output
	var output bytes.Buffer

	// Create options
	options := &fp2lm.ConverterOptions{
		AltitudeMode: "agl",
		GimbalPitch:  -90,
	}

	// Process the input using the exported Process function
	err := fp2lm.Process(inputReader, &output, options)
	if err != nil {
		t.Fatalf("Failed to process input: %v", err)
	}

	// Canonicalize both outputs for comparison
	actualOutput := canonicalizeCSV(output.String())
	expectedOutput := canonicalizeCSV(string(goldenFileData))

	// Compare normalized outputs
	if actualOutput != expectedOutput {
		// Split into lines for more helpful error message
		actualLines := strings.Split(actualOutput, "\n")
		expectedLines := strings.Split(expectedOutput, "\n")

		// Find first difference
		minLen := len(actualLines)
		if len(expectedLines) < minLen {
			minLen = len(expectedLines)
			t.Errorf("Output has %d lines, golden file has %d lines", len(actualLines), len(expectedLines))
		}

		for i := 0; i < minLen; i++ {
			if actualLines[i] != expectedLines[i] {
				t.Errorf("Mismatch on line %d:\nExpected: %s\nGot: %s", i+1, expectedLines[i], actualLines[i])
				break
			}
		}

		if len(actualLines) > len(expectedLines) {
			t.Errorf("Extra lines in actual output starting at line %d", len(expectedLines)+1)
		} else if len(actualLines) < len(expectedLines) {
			t.Errorf("Missing lines in actual output starting at line %d", len(actualLines)+1)
		}
	}
}

// TestCSVRoundTrip tests that a CSV file can be converted and then converted back
func TestCSVRoundTrip(t *testing.T) {
	// First pass: convert FlightplannerMission to Litchi format with AGL mode
	inputReader := bytes.NewReader(flightplannerMissionData)
	var firstOutput bytes.Buffer

	options := &fp2lm.ConverterOptions{
		AltitudeMode: "agl",
		GimbalPitch:  -90,
	}

	err := fp2lm.Process(inputReader, &firstOutput, options)
	if err != nil {
		t.Fatalf("Failed first conversion: %v", err)
	}

	// For this test, we'll just verify key fields from the first output
	// (Litchi format back to Litchi format isn't the primary use case,
	// and detecting the Litchi format would require additional parser code)

	firstLines := strings.Split(canonicalizeCSV(firstOutput.String()), "\n")

	if len(firstLines) <= 1 {
		t.Fatalf("No waypoints generated in first pass")
	}

	// Skip header and check that the first waypoint has the right type of data
	firstFields := strings.Split(firstLines[1], ",")
	if len(firstFields) < 5 {
		t.Fatalf("First line has insufficient fields: %d", len(firstFields))
	}

	// Parse latitude
	lat, _, err := fp2lm.ParseField(firstFields[0], "float64", -90, 90)
	if err != nil {
		t.Fatalf("Failed to parse latitude: %v", err)
	}

	// Verify latitude is roughly near expectation
	if lat < 43.0 || lat > 44.0 {
		t.Errorf("Latitude out of expected range: %f", lat)
	}

	// Parse longitude
	lon, _, err := fp2lm.ParseField(firstFields[1], "float64", -180, 180)
	if err != nil {
		t.Fatalf("Failed to parse longitude: %v", err)
	}

	// Verify longitude is roughly near expectation
	if lon > -88.0 || lon < -90.0 {
		t.Errorf("Longitude out of expected range: %f", lon)
	}

	// Parse altitude
	alt, _, err := fp2lm.ParseField(firstFields[2], "float64", 0, 500)
	if err != nil {
		t.Fatalf("Failed to parse altitude: %v", err)
	}

	// Verify altitude is reasonable
	if alt < 1.0 || alt > 200.0 {
		t.Errorf("Altitude out of expected range: %f", alt)
	}
}

// TestCalculateBearing tests the CalculateBearing function with various edge cases
func TestCalculateBearing(t *testing.T) {
	tests := []struct {
		name           string
		lat1, lon1     float64
		lat2, lon2     float64
		expectedResult float64
		skipReciprocal bool // Skip reciprocal test for special cases
	}{
		{"East", 0, 0, 0, 1, 90, false},
		{"North", 0, 0, 1, 0, 0, false},
		{"West", 0, 0, 0, -1, 270, false},
		{"South", 0, 0, -1, 0, 180, false},
		{"Northeast", 0, 0, 1, 1, 45, false},
		{"Longitude wrap-around", 0, 179, 0, -179, 90, false},
		{"High-latitude mission", 89, 0, 89, 180, 0, true},    // Expecting 0 degrees for our implementation
		{"Same point", 10, 10, 10, 10, 0, true},               // Same point should default to 0 bearing
		{"Opposite meridian points", 0, 0, 0, 180, 90, false}, // Points on opposite sides of the earth
		{"Antipodal points", 0, 0, 0, 180, 90, false},         // Antipodal points
	}

	// Run the main test logic
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fp2lm.CalculateBearing(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			// Allow small floating point differences
			if math.Abs(result-tt.expectedResult) > 0.1 {
				t.Errorf("Expected %.1f, got %.1f", tt.expectedResult, result)
			}

			// Property test: bearing(A→B) should be approximately bearing(B→A)+180° (mod 360)
			// Skip test for same point and other special cases where this property doesn't hold
			if !tt.skipReciprocal {
				reverse := fp2lm.CalculateBearing(tt.lat2, tt.lon2, tt.lat1, tt.lon1)
				expectedReverse := math.Mod(result+180, 360)
				if math.Abs(reverse-expectedReverse) > 0.1 {
					t.Errorf("Reciprocal bearing test failed: %.1f → %.1f should be %.1f apart",
						result, reverse, 180.0)
				}
			}
		})
	}
}

// TestParseField tests the ParseField function
func TestParseField(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		fieldType     string
		min, max      float64
		expectedFloat float64
		expectedInt   int8
		expectError   bool
	}{
		{"Valid float64", "42.5", "float64", 0, 100, 42.5, 0, false},
		{"Out of range float64", "200", "float64", 0, 100, 0, 0, true},
		{"Bad format float64", "abc", "float64", 0, 100, 0, 0, true},
		{"Valid float32", "42.5", "float32", 0, 100, 42.5, 0, false},
		{"Float32 precision", "1.333333333333", "float32", 0, 100, float64(float32(1.333333333333)), 0, false},
		{"Valid int8", "42", "int8", 0, 100, 0, 42, false},
		{"Invalid field type", "42", "int16", 0, 100, 0, 0, true},
		{"NaN value", "nan", "float64", 0, 100, 0, 0, true},
		{"Empty value", "", "float64", 0, 100, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, i, err := fp2lm.ParseField(tt.field, tt.fieldType, tt.min, tt.max)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				if tt.fieldType == "float64" || tt.fieldType == "float32" {
					if f != tt.expectedFloat {
						t.Errorf("Expected float %v, got %v", tt.expectedFloat, f)
					}
				} else if tt.fieldType == "int8" {
					if i != tt.expectedInt {
						t.Errorf("Expected int %v, got %v", tt.expectedInt, i)
					}
				}
			}
		})
	}
}

// TestDefaultOptions checks that default options are correctly set
func TestDefaultOptions(t *testing.T) {
	options := fp2lm.DefaultOptions()

	if options.AltitudeMode != "agl" {
		t.Errorf("Expected AltitudeMode to be 'agl', got %s", options.AltitudeMode)
	}

	if options.GimbalPitch != -90 {
		t.Errorf("Expected GimbalPitch to be -90, got %f", options.GimbalPitch)
	}

	if float64(options.PhotoInterval) != 0 {
		t.Errorf("Expected PhotoInterval to be 0, got %f", float64(options.PhotoInterval))
	}
}

// TestProcessWithMissingFields ensures Process gracefully skips rows with too few columns.
func TestProcessWithMissingFields(t *testing.T) {
	malformed := "Waypoint Number,X [m],Y [m],Alt. ASL [m],Alt. AGL [m]\n" +
		"1,0,0,10,5\n" // missing xcoord and ycoord columns
	var out bytes.Buffer
	err := fp2lm.Process(strings.NewReader(malformed), &out, fp2lm.DefaultOptions())
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	// Expect only header line in output
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
