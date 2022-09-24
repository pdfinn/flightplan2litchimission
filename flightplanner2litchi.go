package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// create a struct to represent the Litchi Mission
type LitchiMission struct {
	longitude          float64
	latitude           float64
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
	f, err := os.Open("FlightplannerMission.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rows, err := readData(f)
	if err != nil {
		panic(err)
	}

	// create an instance of the mission with default values
	mission := LitchiMission{
		longitude:          0,
		latitude:           0,
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

	//Print the Litchi Mission header
	fmt.Println("latitude, longitude, altitude(m), heading(deg), curvesize(m), rotationdir, gimbalmode, " +
		"gimbalpitchangle, actiontype1, actionparam1, actiontype2, actionparam2, actiontype3, actionparam3, " +
		"actiontype4, actionparam4, actiontype5, actionparam5, actiontype6, actionparam6, actiontype7, " +
		"actionparam7, actiontype8, actionparam8, actiontype9, actionparam9, actiontype10, actionparam10," +
		" actiontype11, actionparam11, actiontype12, actionparam12, actiontype13, actionparam13, actiontype14," +
		" actionparam14, actiontype15, actionparam15, altitudemode, speed(m/s), poi_latitude, poi_longitude, " +
		"poi_altitude(m), poi_altitudemode, photo_timeinterval, photo_distinterval")

	for _, row := range rows {
		mission.longitude, err = strconv.ParseFloat(row[5], 64)
		mission.latitude, err = strconv.ParseFloat(row[6], 64)
		altitude, _ := strconv.ParseFloat(row[3], 32)
		mission.altitude = float32(altitude)
		fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v,%v,%v,%v,%v, %v, %v, %v, %v,%v, %v, "+
			"%v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v,%v, %v, %v, %v, %v, %v \n",
			mission.latitude, mission.latitude, mission.altitude, mission.heading, mission.curvesize,
			mission.rotationdir, mission.gimblemode, mission.gimbalpitchangle, mission.actiontype1,
			mission.actionparam1, mission.actiontype2, mission.actionparam2, mission.actiontype3,
			mission.actionparam3, mission.actiontype4, mission.actionparam4, mission.actiontype5,
			mission.actionparam5, mission.actiontype6, mission.actionparam6, mission.actiontype7,
			mission.actionparam7, mission.actiontype8, mission.actionparam8, mission.actiontype9,
			mission.actionparam9, mission.actiontype10, mission.actionparam10, mission.actiontype11,
			mission.actionparam11, mission.actiontype12, mission.actionparam12, mission.actiontype13,
			mission.actionparam13, mission.actiontype14, mission.actionparam14, mission.actiontype15,
			mission.actionparam15, mission.altitudemode, mission.speed, mission.poi_latitude,
			mission.poi_longitude, mission.poi_altitude, mission.poi_altitudemode, mission.photo_timeinterval,
			mission.photo_distinterval)
	}
}
