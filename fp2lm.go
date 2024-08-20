package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var interval = flag.Float64("d", 0, "Enter the photo interval (meters). Example: -d 20")

func main() {
	flag.Parse()

	// Run the core logic and output the result
	processInput(os.Stdin, os.Stdout)
}

func processInput(input io.Reader, output io.Writer) {
	// create an instance of the wp with default values
	wp := newWaypoint()

	// Print the Litchi Mission header
	fmt.Fprintln(output, "latitude, longitude, altitude(m), heading(deg), curvesize(m), rotationdir, gimbalmode, "+
		"gimbalpitchangle, actiontype1, actionparam1, actiontype2, actionparam2, actiontype3, actionparam3, "+
		"actiontype4, actionparam4, actiontype5, actionparam5, actiontype6, actionparam6, actiontype7, "+
		"actionparam7, actiontype8, actionparam8, actionparam9, actiontype9, actionparam10, actionparam10,"+
		" actiontype11, actionparam11, actiontype12, actionparam12, actiontype13, actionparam13, actiontype14,"+
		" actionparam14, actionparam15, actionparam15, altitudemode, speed(m/s), poi_latitude, poi_longitude, "+
		"poi_altitude(m), poi_altitudemode, photo_timeinterval, photo_distinterval")

	// For each line of input, print it as a LitchiMission record
	scanner := bufio.NewScanner(input)

	// Skip the first line (header)
	if scanner.Scan() {
		_ = scanner.Text() // Read and discard the header line
	}

	for scanner.Scan() {
		ln := scanner.Text()

		// Create a new CSV reader to parse the line
		reader := csv.NewReader(strings.NewReader(ln))

		// Read the record from the CSV reader
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(output, "Error reading CSV input:", err)
			continue
		}

		// Validate the record and set waypoint fields
		if err := validateRecord(rec); err != nil {
			fmt.Fprintln(output, "Skipping line due to validation error:", err)
			continue
		}

		// Correct mapping: Parse xcoord as longitude and ycoord as latitude
		wp.longitude, err = parseField(rec[5], -180, 180)
		if err != nil {
			fmt.Fprintln(output, "Skipping line due to error parsing longitude:", err)
			continue
		}

		wp.latitude, err = parseField(rec[6], -90, 90)
		if err != nil {
			fmt.Fprintln(output, "Skipping line due to error parsing latitude:", err)
			continue
		}

		wp.altitude, err = parseField(rec[3], 0, math.MaxFloat64)
		if err != nil {
			fmt.Fprintln(output, "Skipping line due to error parsing altitude:", err)
			continue
		}

		// Set gimbal pitch angle to -90
		wp.gimbalpitchangle = -90

		// Set the photo distance interval
		wp.photo_distinterval = *interval

		// Print the waypoint
		fmt.Fprintf(output, "%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, "+
			"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n",
			wp.latitude, wp.longitude, wp.altitude, wp.heading, wp.curvesize,
			wp.rotationdir, wp.gimblemode, wp.gimbalpitchangle, wp.actiontype1,
			wp.actionparam1, wp.actiontype2, wp.actionparam2, wp.actiontype3,
			wp.actionparam3, wp.actionparam4, wp.actionparam4, wp.actionparam5,
			wp.actionparam5, wp.actionparam6, wp.actionparam6, wp.actionparam7,
			wp.actionparam7, wp.actionparam8, wp.actionparam8, wp.actionparam9,
			wp.actionparam9, wp.actionparam10, wp.actionparam10, wp.actionparam11,
			wp.actionparam11, wp.actionparam12, wp.actionparam12, wp.actionparam13,
			wp.actionparam13, wp.actionparam14, wp.actionparam14, wp.actionparam15,
			wp.altitudemode, wp.speed, wp.poi_latitude,
			wp.poi_longitude, wp.poi_altitude, wp.poi_altitudemode, wp.photo_timeinterval,
			wp.photo_distinterval)
	}
}

func parseField(field string, min, max float64) (float64, error) {
	f, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return 0, err
	}
	if f < min || f > max {
		return 0, fmt.Errorf("Field value out of range (min: %f, max: %f)", min, max)
	}
	return f, nil
}

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
	photo_distinterval float64
}

// handler function for creating new waypoints
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
		actiontype1:        -1,
		actionparam1:       0,
		actiontype2:        -1,
		actionparam2:       0,
		actiontype3:        -1,
		actionparam3:       0,
		actiontype4:        -1,
		actionparam4:       0,
		actiontype5:        -1,
		actionparam5:       0,
		actiontype6:        -1,
		actionparam6:       0,
		actiontype7:        -1,
		actionparam7:       0,
		actiontype8:        -1,
		actionparam8:       0,
		actiontype9:        -1,
		actionparam9:       0,
		actiontype10:       -1,
		actionparam10:      0,
		actiontype11:       -1,
		actionparam11:      0,
		actiontype12:       -1,
		actionparam12:      0,
		actiontype13:       -1,
		actionparam13:      0,
		actiontype14:       -1,
		actionparam14:      0,
		actiontype15:       -1,
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

func validateRecord(rec []string) error {
	if len(rec) < 7 { // Assuming at least 7 fields are needed
		return fmt.Errorf("record has too few fields")
	}
	return nil
}
