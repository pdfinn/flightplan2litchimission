# fp2lm package

This package implements the core functionality for the flightplan2litchimission converter.

## Overview

The fp2lm package provides functions for converting waypoint data from Flight Planner CSV format to Litchi Mission Hub CSV format, making it compatible with DJI drones.

## Usage

```go
import (
    "flightplan2litchimission/fp2lm"
    "os"
)

func main() {
    // Create converter options
    options := &fp2lm.ConverterOptions{
        AltitudeMode:   "agl",     // "agl" or "asl"
        PhotoInterval:  20,        // Distance between photos in meters
        GimbalPitch:    -90,       // Camera angle in degrees
        MaxAltitudeAGL: 120,       // Maximum allowed altitude in meters
    }

    // Convert from input to output
    err := fp2lm.Process(os.Stdin, os.Stdout, options)
    if err != nil {
        // Handle error
    }
}
```

## Key Functions

- `Process(input io.Reader, output io.Writer, options *ConverterOptions) error`: Main conversion function that processes input CSV data and writes Litchi format.
- `CalculateBearing(lat1, lon1, lat2, lon2 float64) float64`: Calculates the initial bearing between two geographic points.
- `DefaultOptions() *ConverterOptions`: Returns recommended default settings for the converter.

## Options

The `ConverterOptions` struct configures the conversion behavior:

- `AltitudeMode`: Determines how altitude values are interpreted. Use "agl" for relative altitudes (Above Ground Level) or "asl" for absolute altitudes (Above Sea Level).
- `PhotoInterval`: Specifies the distance between photos in meters.
- `GimbalPitch`: Sets the camera angle in degrees (between -90 and 0).
- `MaxAltitudeAGL`: Specifies the maximum allowed altitude when in AGL mode, typically set to local regulatory limits. 