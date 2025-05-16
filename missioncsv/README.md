# missioncsv package

This package provides utilities for formatting and writing mission CSV files for various drone platforms including Litchi Mission Hub and DJI Pilot 2.

## Overview

The missioncsv package implements CSV writing functionality for drone missions, focusing on the Litchi Mission Hub format. It provides structured types for waypoints, points of interest, and actions, along with a CSV writer to generate properly formatted output.

## Usage

```go
import (
    "flightplan2litchimission/missioncsv"
    "os"
)

func main() {
    // Create a writer for the output
    writer := missioncsv.NewWriter(os.Stdout)
    
    // Write the standard Litchi mission header
    err := writer.WriteLitchiHeader()
    if err != nil {
        // Handle error
    }
    
    // Create and configure a waypoint
    waypoint := missioncsv.NewLitchiWaypoint()
    waypoint.Point.Latitude = 43.0009225
    waypoint.Point.Longitude = -89.0003070
    waypoint.Point.Altitude = 30.0
    waypoint.Heading = 90.0
    waypoint.GimbalPitch = -90.0
    waypoint.AltitudeMode = 1 // 1 for relative (AGL), 0 for absolute (ASL)
    
    // Add a "take photo" action
    waypoint.Actions = []missioncsv.Action{
        {Type: 1, Param: 0},
    }
    
    // Write the waypoint to the output
    err = writer.WriteLitchiWaypoint(waypoint)
    if err != nil {
        // Handle error
    }
    
    // Make sure to flush the buffer when done
    writer.Flush()
}
```

## Types

- `Point`: Represents geographic coordinates (latitude, longitude, altitude)
- `POI`: Represents a Point of Interest
- `LitchiWaypoint`: Contains all data needed for a Litchi mission waypoint
- `Action`: Represents an action to perform at a waypoint (e.g., take photo)
- `Writer`: Handles writing waypoints to a Litchi-compatible CSV file

## Key Functions

- `NewWriter(w io.Writer) *Writer`: Creates a new mission CSV writer
- `WriteLitchiHeader() error`: Writes the standard Litchi mission header
- `WriteLitchiWaypoint(wp *LitchiWaypoint) error`: Writes a single waypoint in Litchi format
- `NewLitchiWaypoint() *LitchiWaypoint`: Creates a new waypoint with default values 