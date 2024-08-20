package main

import (
	"math"
	"testing"
)

// TestCalculateBearing tests the bearing calculation function
func TestCalculateBearing(t *testing.T) {
	testCases := []struct {
		lat1, lon1, lat2, lon2 float64
		expectedBearing        float64
	}{
		{0, 0, 1, 1, 45},
		// ... additional test cases ...
	}

	for _, tc := range testCases {
		bearing := calculateBearing(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
		if math.Abs(bearing-tc.expectedBearing) > 0.0001 {
			t.Errorf("calculateBearing(%v, %v, %v, %v) = %v; want %v", tc.lat1, tc.lon1, tc.lat2, tc.lon2, bearing, tc.expectedBearing)
		}
	}
}
