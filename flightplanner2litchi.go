package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type LitchiMission struct {
	longitude          float64
	latitude           float64
	altitude           float32 // meters
	heading            float32
	curvesize          float32
	rotationdir        bool
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
	altitudemode       bool
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
	// Skip first row (line)
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
	f, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rows, err := readData(f)
	if err != nil {
		panic(err)
	}

	mission := new(LitchiMission)

	for _, row := range rows {
		mission.longitude, err = strconv.ParseFloat(row[5], 64)
		mission.latitude, err = strconv.ParseFloat(row[6], 64)
	}

	fmt.Println(mission.latitude)
}
