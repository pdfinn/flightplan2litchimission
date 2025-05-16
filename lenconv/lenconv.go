// Package lenconv performs meters and feet distance computations.
//
// This package provides functions and types for converting between meters and feet,
// with special support for command-line flag parsing of distance values with units.
// It is used by the flightplan2litchimission tool to handle distance intervals.
package lenconv

import (
	"flag"
	"fmt"
)

// Meters represents a distance value in meters
type Meters float64

// Feet represents a distance value in feet
type Feet float64

// Convert feet to meters with standard conversion factor
func feet2meters(f Feet) Meters { return Meters(f * 0.3048) }

// String returns the string representation of Meters
func (m Meters) String() string { return fmt.Sprintf("%g", m) }

//type Value interface {
//	String() string
//	Set() error
//}

// photoIntervalFlag implements the flag.Value interface for photo interval distances
type photoIntervalFlag struct{ Meters }

// Set parses a string value with units into a distance and stores it
// Valid formats include "10m" for meters or "30ft" for feet
func (f *photoIntervalFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "M", "m", "Meters", "meters":
		f.Meters = Meters(value)
		return nil
	case "Ft", "ft", "Feet", "feet":
		f.Meters = feet2meters(Feet(value))
		return nil
	}
	return fmt.Errorf("invalid units %q", s)
}

// PhotoIntervalFlag creates a flag that accepts distances with units
// It can be used to specify photo interval distances in meters or feet
func PhotoIntervalFlag(name string, value Meters, usage string) *Meters {
	f := photoIntervalFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Meters
}
