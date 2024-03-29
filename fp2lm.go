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

    // create an instance of the wp with default values
    wp := newWaypoint()

    // print the Litchi Mission header
    fmt.Println("latitude, longitude, altitude(m), heading(deg), curvesize(m), rotationdir, gimbalmode, " +
        "gimbalpitchangle, actiontype1, actionparam1, actiontype2, actionparam2, actiontype3, actionparam3, " +
        "actiontype4, actionparam4, actiontype5, actionparam5, actiontype6, actionparam6, actiontype7, " +
        "actionparam7, actiontype8, actionparam8, actiontype9, actionparam9, actiontype10, actionparam10," +
        " actiontype11, actionparam11, actiontype12, actionparam12, actiontype13, actionparam13, actiontype14," +
        " actionparam14, actiontype15, actionparam15, altitudemode, speed(m/s), poi_latitude, poi_longitude, " +
        "poi_altitude(m), poi_altitudemode, photo_timeinterval, photo_distinterval")

    // for each line of standard input, process it as a LitchiMission record
    scanner := bufio.NewScanner(os.Stdin)
    var previousWp *LitchiWaypoint
    waypoints := []*LitchiWaypoint{}

    for scanner.Scan() {
        ln := scanner.Text()
        reader := csv.NewReader(strings.NewReader(ln))
        rec, err := reader.Read()

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error reading CSV input:", err)
            os.Exit(1)
        }

        if rec[0] == "Waypoint Number" && rec[1] == "X [m]" && rec[2] == "Y [m]" && rec[3] == "Alt. ASL [m]" && rec[4] == "Alt. AGL [m]" && rec[5] == "xcoord" && rec[6] == "ycoord" {
            continue
        }

        wp := newWaypoint() // Create a new waypoint instance for each record

        longitude, _, err := parseField(rec[5], "float64", -180, 180)
        if err != nil {
            fmt.Println("Error parsing longitude:", err)
            continue
        }
        wp.longitude = longitude

        latitude, _, err := parseField(rec[6], "float64", -90, 90)
        if err != nil {
            fmt.Println("Error parsing latitude:", err)
            continue
        }
        wp.latitude = latitude

        altitude, _, err := parseField(rec[3], "float64", 0, math.MaxFloat64)
        if err != nil {
            fmt.Println("Error parsing altitude:", err)
            continue
        }
        wp.altitude = altitude
        wp.gimbalpitchangle = -90
        wp.photo_distinterval = interval

        waypoints = append(waypoints, wp)

        if previousWp != nil {
            previousWp.heading = float32(calculateBearing(previousWp.latitude, previousWp.longitude, wp.latitude, wp.longitude))
        }

        previousWp = wp
    }

    if len(waypoints) > 1 {
        lastWp := waypoints[len(waypoints)-1]
        secondLastWp := waypoints[len(waypoints)-2]
        lastWp.heading = secondLastWp.heading
    }

    // Print the waypoints outside the loop
    for _, wp := range waypoints {
        fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, "+
            "%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v \n",
            wp.latitude, wp.longitude, wp.altitude, wp.heading, wp.curvesize,
            wp.rotationdir, wp.gimblemode, wp.gimbalpitchangle, wp.actiontype1,
            wp.actionparam1, wp.actiontype2, wp.actionparam2, wp.actiontype3,
            wp.actionparam3, wp.actiontype4, wp.actionparam4, wp.actiontype5,
            wp.actionparam5, wp.actiontype6, wp.actionparam6, wp.actiontype7,
            wp.actionparam7, wp.actiontype8, wp.actionparam8, wp.actiontype9,
            wp.actionparam9, wp.actiontype10, wp.actionparam10, wp.actiontype11,
            wp.actionparam11, wp.actiontype12, wp.actionparam12, wp.actiontype13,
            wp.actionparam13, wp.actiontype14, wp.actionparam14, wp.actiontype15,
            wp.actionparam15, wp.altitudemode, wp.speed, wp.poi_latitude,
            wp.poi_longitude, wp.poi_altitude, wp.poi_altitudemode, wp.photo_timeinterval,
            wp.photo_distinterval)
    }
}


// LitchiWaypoint is a struct to represent the Litchi waypoint
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
	photo_distinterval *lenconv.Meters
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
		gimbalpitchangle:   0,
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
		photo_distinterval: interval, // meters
	}
}

// Helper function to perform validation on the input.  We check for sane types, minimum, and maximum values.
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
