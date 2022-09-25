package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// create a struct to represent the Litchi waypoint
type LitchiWaypoint struct {
	latitude           float64
	longitude          float64
	altitude           float32 // meters
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
	poi_altitude       float32 // meters
	poi_altitudemode   int8
	photo_timeinterval float32
	photo_distinterval float32
}

// read csv values using csv.Reader
func readData(rs io.ReadSeeker) ([][]string, error) {
	// Skip the header row
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read remaining rows
	r := csv.NewReader(rs)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rows, err := readData(f)
	if err != nil {
		panic(err)
	}

	// create an instance of the waypoint with default values
	waypoint := LitchiWaypoint{
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
		photo_distinterval: -1,
	}

	// print the Litchi Mission header
	fmt.Println("latitude, longitude, altitude(m), heading(deg), curvesize(m), rotationdir, gimbalmode, " +
		"gimbalpitchangle, actiontype1, actionparam1, actiontype2, actionparam2, actiontype3, actionparam3, " +
		"actiontype4, actionparam4, actiontype5, actionparam5, actiontype6, actionparam6, actiontype7, " +
		"actionparam7, actiontype8, actionparam8, actiontype9, actionparam9, actiontype10, actionparam10," +
		" actiontype11, actionparam11, actiontype12, actionparam12, actiontype13, actionparam13, actiontype14," +
		" actionparam14, actiontype15, actionparam15, altitudemode, speed(m/s), poi_latitude, poi_longitude, " +
		"poi_altitude(m), poi_altitudemode, photo_timeinterval, photo_distinterval")

	// print the individual records/waypoints
	for _, row := range rows {
		waypoint.longitude, err = strconv.ParseFloat(row[5], 64)
		waypoint.latitude, err = strconv.ParseFloat(row[6], 64)
		altitude, _ := strconv.ParseFloat(row[3], 32)
		waypoint.altitude = float32(altitude)
		fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v,%v,%v,%v,%v, %v, %v, %v, %v,%v, %v, "+
			"%v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v, %v \n",
			waypoint.latitude, waypoint.longitude, waypoint.altitude, waypoint.heading, waypoint.curvesize,
			waypoint.rotationdir, waypoint.gimblemode, waypoint.gimbalpitchangle, waypoint.actiontype1,
			waypoint.actionparam1, waypoint.actiontype2, waypoint.actionparam2, waypoint.actiontype3,
			waypoint.actionparam3, waypoint.actiontype4, waypoint.actionparam4, waypoint.actiontype5,
			waypoint.actionparam5, waypoint.actiontype6, waypoint.actionparam6, waypoint.actiontype7,
			waypoint.actionparam7, waypoint.actiontype8, waypoint.actionparam8, waypoint.actiontype9,
			waypoint.actionparam9, waypoint.actiontype10, waypoint.actionparam10, waypoint.actiontype11,
			waypoint.actionparam11, waypoint.actiontype12, waypoint.actionparam12, waypoint.actiontype13,
			waypoint.actionparam13, waypoint.actiontype14, waypoint.actionparam14, waypoint.actiontype15,
			waypoint.actionparam15, waypoint.altitudemode, waypoint.speed, waypoint.poi_latitude,
			waypoint.poi_longitude, waypoint.poi_altitude, waypoint.poi_altitudemode, waypoint.photo_timeinterval,
			waypoint.photo_distinterval)
	}
}
