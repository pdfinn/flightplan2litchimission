package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"flightplan2litchimission/lenconv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var interval = lenconv.PhotoIntervalFlag("d", 0, "Enter the photo interval (meters 'm' or feet 'ft'). Example: -d 20ft")

func calculateBearing(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lonDiffRad := (lon2 - lon1) * math.Pi / 180

	y := math.Sin(lonDiffRad) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(lonDiffRad)
	bearing := math.Atan2(y, x) * 180 / math.Pi

	return math.Mod(bearing+360, 360)
}

func main() {
	flag.Parse()

	// Print the Litchi Mission header
	fmt.Println("latitude, longitude, altitude(m), heading(deg), curvesize(m), rotationdir, gimbalmode, " +
		"gimbalpitchangle, actiontype1, actionparam1, actiontype2, actionparam2, actiontype3, actionparam3, " +
		"actiontype4, actionparam4, actiontype5, actionparam5, actiontype6, actionparam6, actiontype7, " +
		"actionparam7, actiontype8, actionparam8, actionparam9, actionparam9, actiontype10, actionparam10," +
		" actiontype11, actionparam11, actiontype12, actionparam12, actiontype13, actionparam13, actiontype14," +
		" actionparam14, actiontype15, actionparam15, altitudemode, speed(m/s), poi_latitude, poi_longitude, " +
		"poi_altitude(m), poi_altitudemode, photo_timeinterval, photo_distinterval")

	// Process the input and output the result
	processInput(os.Stdin, os.Stdout)
}

func processInput(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	var previousWp *LitchiWaypoint
	waypoints := []*LitchiWaypoint{}

	for scanner.Scan() {
		ln := scanner.Text()
		reader := csv.NewReader(strings.NewReader(ln))
		rec, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(output, "Error reading CSV input:", err)
			continue
		}

		// Skip header line based on known CSV structure
		if rec[0] == "Waypoint Number" && rec[1] == "X [m]" && rec[2] == "Y [m]" && rec[3] == "Alt. ASL [m]" && rec[4] == "Alt. AGL [m]" && rec[5] == "xcoord" && rec[6] == "ycoord" {
			continue
		}

		wp := newWaypoint() // Create a new waypoint instance for each record

		wp.longitude, _, err = parseField(rec[5], "float64", -180, 180)
		if err != nil {
			fmt.Fprintln(output, "Error parsing longitude:", err)
			continue
		}

		wp.latitude, _, err = parseField(rec[6], "float64", -90, 90)
		if err != nil {
			fmt.Fprintln(output, "Error parsing latitude:", err)
			continue
		}

		wp.altitude, _, err = parseField(rec[3], "float64", 0, math.MaxFloat64)
		if err != nil {
			fmt.Fprintln(output, "Error parsing altitude:", err)
			continue
		}

		wp.gimbalpitchangle = -90
		wp.photo_distinterval = *interval // Correctly assign the interval value

		waypoints = append(waypoints, wp)

		// Calculate heading based on previous waypoint
		if previousWp != nil {
			previousWp.heading = float32(calculateBearing(previousWp.latitude, previousWp.longitude, wp.latitude, wp.longitude))
		}

		previousWp = wp
	}

	// Adjust heading for the last waypoint
	if len(waypoints) > 1 {
		lastWp := waypoints[len(waypoints)-1]
		secondLastWp := waypoints[len(waypoints)-2]
		lastWp.heading = secondLastWp.heading
	}

	// Print all waypoints
	for _, wp := range waypoints {
		fmt.Fprintf(output, "%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, "+
			"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n",
			wp.latitude, wp.longitude, wp.altitude, wp.heading, wp.curvesize,
			wp.rotationdir, wp.gimblemode, wp.gimbalpitchangle, wp.actiontype1,
			wp.actionparam1, wp.actiontype2, wp.actionparam2, wp.actiontype3,
			wp.actionparam3, wp.actiontype4, wp.actionparam4, wp.actiontype5,
			wp.actionparam5, wp.actiontype6, wp.actionparam6, wp.actiontype7,
			wp.actionparam7, wp.actiontype8, wp.actionparam8, wp.actiontype9,
			wp.actionparam9, wp.actiontype10, wp.actionparam10, wp.actiontype11,
			wp.actionparam11, wp.actiontype12, wp.actionparam12, wp.actiontype13,
			wp.actionparam13, wp.actiontype14, wp.actionparam14, wp.actiontype15, wp.actionparam15,
			wp.altitudemode, wp.speed, wp.poi_latitude,
			wp.poi_longitude, wp.poi_altitude, wp.poi_altitudemode, wp.photo_timeinterval,
			wp.photo_distinterval)
	}
}

// LitchiWaypoint represents a waypoint in a Litchi mission
type LitchiWaypoint struct {
	latitude           float64
	longitude          float64
	altitude           float64 // meters
	heading            float32
	curvesize          float32
	rotationdir        int8
	gimblemode         int8
	gimbalpitchangle   float32
	actiontype1        int8
	actionparam1       int8
	actiontype2        int8
	actionparam2       int8
	actiontype3        int8
	actionparam3       int8
	actiontype4        int8
	actionparam4       int8
	actiontype5        int8
	actionparam5       int8
	actiontype6        int8
	actionparam6       int8
	actiontype7        int8
	actionparam7       int8
	actiontype8        int8
	actionparam8       int8
	actiontype9        int8
	actionparam9       int8
	actiontype10       int8
	actionparam10      int8
	actiontype11       int8
	actionparam11      int8
	actiontype12       int8
	actionparam12      int8
	actiontype13       int8
	actionparam13      int8
	actiontype14       int8
	actionparam14      int8
	actiontype15       int8
	actionparam15      int8
	altitudemode       int8
	speed              float32 // meters per second
	poi_latitude       float64
	poi_longitude      float64
	poi_altitude       float64 // meters
	poi_altitudemode   int8
	photo_timeinterval float32
	photo_distinterval lenconv.Meters
}

// Create a new LitchiWaypoint with default values
func newWaypoint() *LitchiWaypoint {
	return &LitchiWaypoint{
		latitude:           0,
		longitude:          0,
		altitude:           0, // meters
		heading:            360,
		curvesize:          0,
		rotationdir:        0,
		gimblemode:         0,
		gimbalpitchangle:   -90,
		actiontype1:        1,
		actionparam1:       0,
		actiontype2:        0,
		actionparam2:       0,
		actiontype3:        0,
		actionparam3:       0,
		actiontype4:        0,
		actionparam4:       0,
		actiontype5:        0,
		actionparam5:       0,
		actiontype6:        0,
		actionparam6:       0,
		actiontype7:        0,
		actionparam7:       0,
		actiontype8:        0,
		actionparam8:       0,
		actiontype9:        0,
		actionparam9:       0,
		actiontype10:       0,
		actionparam10:      0,
		actiontype11:       0,
		actionparam11:      0,
		actiontype12:       0,
		actionparam12:      0,
		actiontype13:       0,
		actionparam13:      0,
		actiontype14:       0,
		actionparam14:      0,
		actiontype15:       0,
		actionparam15:      0,
		altitudemode:       1,
		speed:              0, // meters per second
		poi_latitude:       0,
		poi_longitude:      0,
		poi_altitude:       0, // meters
		poi_altitudemode:   0,
		photo_timeinterval: -1,
		photo_distinterval: *interval, // meters
	}
}

// Helper function to parse fields and validate their range
func parseField(field string, fieldType string, min float64, max float64) (float64, int8, error) {
	switch fieldType {
	case "float64":
		f, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return 0, 0, err
		}
		if f < min || f > max {
			return 0, 0, fmt.Errorf("Field value out of range (min: %f, max: %f)", min, max)
		}
		return f, 0, nil
	case "float32":
		f, err := strconv.ParseFloat(field, 32)
		if err != nil {
			return 0, 0, err
		}
		if f < min || f > max {
			return 0, 0, fmt.Errorf("Field value out of range (min: %f, max: %f)", min, max)
		}
		return f, 0, nil
	case "int8":
		i, err := strconv.ParseInt(field, 10, 32)
		if err != nil {
			return 0, 0, err
		}
		if i < int64(min) || i > int64(max) {
			return 0, 0, fmt.Errorf("Field value out of range (min: %f, max: %f)", min, max)
		}
		return 0, int8(i), nil
	default:
		return 0, 0, fmt.Errorf("Invalid field type: %s", fieldType)
	}
}
