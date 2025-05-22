package main

import (
	"flag"
	"log"
	"os"

	"flightplan2litchimission/fp2lm"
	"flightplan2litchimission/lenconv"
)

func main() {
	interval := lenconv.PhotoIntervalFlag("d", 0, "distance between photos (e.g. 20m or 60ft)")
	altitudeMode := flag.String("altitude-mode", "agl", "altitude mode: 'agl' or 'asl'")
	pitch := flag.Float64("pitch", -90, "gimbal pitch angle (-90 to 0)")
	output := flag.String("output", "", "output file (default stdout)")
	flag.Parse()

	opts := &fp2lm.ConverterOptions{
		AltitudeMode:  *altitudeMode,
		PhotoInterval: *interval,
		GimbalPitch:   *pitch,
	}

	var out *os.File
	var err error
	if *output != "" {
		out, err = os.Create(*output)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if err := fp2lm.Process(os.Stdin, out, opts); err != nil {
		log.Fatalf("conversion failed: %v", err)
	}
}
