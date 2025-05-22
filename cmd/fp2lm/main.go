package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"flightplan2litchimission/fp2lm"
	"flightplan2litchimission/lenconv"
	"flightplan2litchimission/missionkml"
)

func main() {
	interval := lenconv.PhotoIntervalFlag("d", 0, "Photo interval distance (e.g. 20m or 60ft)")
	altitudeMode := flag.String("altitude-mode", "agl", "Altitude mode: agl or asl")
	pitch := flag.Float64("pitch", -90, "Gimbal pitch angle (-90 to 0)")
	maxAlt := flag.Float64("max-altitude", 120, "Maximum allowed altitude AGL in meters")
	outputPath := flag.String("output", "", "Output file path")
	format := flag.String("format", "csv", "Output format: csv, kml or kmz")
	flag.Parse()

	opts := &fp2lm.ConverterOptions{
		AltitudeMode:   *altitudeMode,
		PhotoInterval:  lenconv.Meters(*interval),
		GimbalPitch:    *pitch,
		MaxAltitudeAGL: *maxAlt,
	}

	var out io.Writer = os.Stdout
	var file *os.File
	var err error
	if *outputPath != "" {
		file, err = os.Create(*outputPath)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer file.Close()
		out = file
	}

	switch *format {
	case "csv":
		if err := fp2lm.Process(os.Stdin, out, opts); err != nil {
			log.Fatal(err)
		}
	case "kml", "kmz":
		waypoints, err := fp2lm.ConvertWaypoints(os.Stdin, opts)
		if err != nil {
			log.Fatal(err)
		}
		if *format == "kml" {
			if err := missionkml.WriteKML(out, waypoints); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := missionkml.WriteKMZ(out, waypoints); err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatalf("unknown format %q", *format)
	}

	if file != nil {
		fmt.Fprintln(os.Stderr, "wrote", *outputPath)
	}
}
