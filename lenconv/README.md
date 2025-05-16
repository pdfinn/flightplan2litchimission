# lenconv package

This package performs meters and feet distance computations with flag parsing support.

## Overview

The lenconv package provides functions and types for converting between meters and feet, with special support for command-line flag parsing of distance values with units. It is used by the flightplan2litchimission tool to handle distance intervals.

## Usage

### Basic Conversion

```go
import (
    "flightplan2litchimission/lenconv"
    "fmt"
)

func main() {
    // Create a distance in meters
    meters := lenconv.Meters(10)
    
    // Convert to feet (internal function, not exported)
    // feet := lenconv.meters2feet(meters)
    
    // Just display the meters value
    fmt.Printf("Distance: %g meters\n", meters)
}
```

### Command-Line Flag Support

```go
import (
    "flag"
    "flightplan2litchimission/lenconv"
    "fmt"
)

func main() {
    // Create a flag that accepts distances with units
    distance := lenconv.PhotoIntervalFlag("distance", 0, "Enter a distance (e.g., 20m or 50ft)")
    
    flag.Parse()
    
    // Use the distance value (always in meters internally)
    fmt.Printf("Distance: %g meters\n", *distance)
}
```

The flag will accept values like:
- `20m` (20 meters)
- `30ft` (30 feet, automatically converted to meters)
- `5meters` (5 meters)
- `15feet` (15 feet, automatically converted to meters)

## Types

- `Meters`: Represents a distance value in meters
- `Feet`: Represents a distance value in feet

## Key Functions

- `PhotoIntervalFlag(name string, value Meters, usage string) *Meters`: Creates a flag that accepts distances with units 