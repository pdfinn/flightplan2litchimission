// Package distconv performs Meters and Feet distance computations.
package distconv

import (
	"flag"
	"fmt"
)

type Meters float64
type Feet float64

func meters2feet(m Meters) Feet { return m * 3.2808}
func feet2meters(f Feet) Meters { return f * 0.3048 }

type Value interface {
	String() string
	Set() error
}

type photoIntervalFlag struct{ Meters }

func (f *photoIntervalFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "M", "Meters", "meters":
		f.Meters = value
	case "Ft", "ft", "Feet", "feet":
		f.Meters = feet2meters(value)
	}
	return fmt.Errorf("Invalid units %q", s)
}

func photoIntervalFlag(name string, value Meters, usage string) *Meters {
	f := photoIntervalFlag(Value) {
		flag.CommandLine.Var(&f, name, usage)
		return &f.Meters
	}
}