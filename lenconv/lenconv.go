// Package lenconv performs meters and feet distance computations.
package lenconv

import (
	"flag"
	"fmt"
)

type Meters float64
type Feet float64

func meters2feet(m Meters) Feet { return Feet(m * 3.2808) }
func feet2meters(f Feet) Meters { return Meters(f * 0.3048) }

func (m Meters) String() string { return fmt.Sprintf("%g", m) }

//type Value interface {
//	String() string
//	Set() error
//}

type photoIntervalFlag struct{ Meters }

func (f *photoIntervalFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "M", "m", "Meters", "meters":
		f.Meters = Meters(value)
	case "Ft", "ft", "Feet", "feet":
		f.Meters = feet2meters(Feet(value))
	}
	return fmt.Errorf("Invalid units %q", s)
}

func PhotoIntervalFlag(name string, value Meters, usage string) *Meters {
	f := photoIntervalFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Meters
}
