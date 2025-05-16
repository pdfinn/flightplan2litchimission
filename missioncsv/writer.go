// Package missioncsv provides utilities for formatting and writing mission CSV files
// for various drone platforms including Litchi Mission Hub and DJI Pilot 2.
package missioncsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
)

// Point represents a waypoint with its geographic coordinates
type Point struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

// POI represents a point of interest
type POI struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

// LitchiWaypoint contains all data needed for a Litchi mission waypoint
type LitchiWaypoint struct {
	// Geographic coordinates
	Point Point
	// Camera orientation
	Heading     float32
	GimbalPitch float32
	// Curve properties
	CurveSize   float32
	RotationDir int8
	// Control settings
	GimbalMode   int8
	AltitudeMode int8
	Speed        float32
	// Potential point of interest
	POI        POI
	POIAltMode int8
	// Photo interval settings
	PhotoTimeInterval float32
	PhotoDistInterval float32
	// Actions (up to 15)
	Actions []Action
}

// Action represents an action to perform at a waypoint
type Action struct {
	Type  int8
	Param int8
}

// Writer handles writing waypoints to a Litchi-compatible CSV file
type Writer struct {
	csvWriter *csv.Writer
}

// NewWriter creates a new mission CSV writer that outputs to the provided writer
func NewWriter(w io.Writer) *Writer {
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = ','
	return &Writer{csvWriter: csvWriter}
}

// WriteLitchiHeader writes the standard Litchi mission header
func (w *Writer) WriteLitchiHeader() error {
	headerFields := []string{
		"latitude", "longitude", "altitude(m)", "heading(deg)", "curvesize(m)", "rotationdir", "gimbalmode",
		"gimbalpitchangle", "actiontype1", "actionparam1", "actiontype2", "actionparam2", "actiontype3", "actionparam3",
		"actiontype4", "actionparam4", "actiontype5", "actionparam5", "actiontype6", "actionparam6", "actiontype7",
		"actionparam7", "actiontype8", "actionparam8", "actiontype9", "actionparam9", "actiontype10", "actionparam10",
		"actiontype11", "actionparam11", "actiontype12", "actionparam12", "actiontype13", "actionparam13", "actiontype14",
		"actionparam14", "actiontype15", "actionparam15", "altitudemode", "speed(m/s)", "poi_latitude", "poi_longitude",
		"poi_altitude(m)", "poi_altitudemode", "photo_timeinterval", "photo_distinterval",
	}
	return w.csvWriter.Write(headerFields)
}

// WriteLitchiWaypoint writes a single waypoint in Litchi format
func (w *Writer) WriteLitchiWaypoint(wp *LitchiWaypoint) error {
	// Ensure we have at least 15 actions (padding with zeros if needed)
	actions := wp.Actions
	for len(actions) < 15 {
		actions = append(actions, Action{Type: 0, Param: 0})
	}

	// Truncate if we somehow have more than 15 actions
	if len(actions) > 15 {
		actions = actions[:15]
		slog.Warn("Truncated excess actions for waypoint", "latitude", wp.Point.Latitude, "longitude", wp.Point.Longitude)
	}

	row := []string{
		fmt.Sprintf("%.7f", wp.Point.Latitude),
		fmt.Sprintf("%.7f", wp.Point.Longitude),
		fmt.Sprintf("%.3f", wp.Point.Altitude),
		fmt.Sprintf("%.1f", wp.Heading),
		fmt.Sprintf("%.1f", wp.CurveSize),
		fmt.Sprintf("%d", wp.RotationDir),
		fmt.Sprintf("%d", wp.GimbalMode),
		fmt.Sprintf("%.1f", wp.GimbalPitch),
	}

	// Add all 15 actions
	for i := 0; i < 15; i++ {
		row = append(row,
			fmt.Sprintf("%d", actions[i].Type),
			fmt.Sprintf("%d", actions[i].Param),
		)
	}

	// Add remaining fields
	row = append(row,
		fmt.Sprintf("%d", wp.AltitudeMode),
		fmt.Sprintf("%.1f", wp.Speed),
		fmt.Sprintf("%.7f", wp.POI.Latitude),
		fmt.Sprintf("%.7f", wp.POI.Longitude),
		fmt.Sprintf("%.3f", wp.POI.Altitude),
		fmt.Sprintf("%d", wp.POIAltMode),
		fmt.Sprintf("%.1f", wp.PhotoTimeInterval),
		fmt.Sprintf("%.1f", wp.PhotoDistInterval),
	)

	return w.csvWriter.Write(row)
}

// Flush writes any buffered data to the underlying io.Writer
func (w *Writer) Flush() {
	w.csvWriter.Flush()
}

// Error returns any error that occurred during writing
func (w *Writer) Error() error {
	return w.csvWriter.Error()
}

// CreateDefaultAction creates a default action with specified type and param
func CreateDefaultAction(actionType, param int8) Action {
	return Action{Type: actionType, Param: param}
}

// NewLitchiWaypoint creates a new waypoint with default values
func NewLitchiWaypoint() *LitchiWaypoint {
	return &LitchiWaypoint{
		Point: Point{
			Latitude:  0,
			Longitude: 0,
			Altitude:  0,
		},
		Heading:      360,
		CurveSize:    0,
		RotationDir:  0,
		GimbalMode:   0,
		GimbalPitch:  -90,
		AltitudeMode: 1, // Default to relative (AGL)
		Speed:        0,
		POI: POI{
			Latitude:  0,
			Longitude: 0,
			Altitude:  0,
		},
		POIAltMode:        0,
		PhotoTimeInterval: -1,
		PhotoDistInterval: -1,
		Actions: []Action{
			{Type: 1, Param: 0}, // Default take photo action
		},
	}
}
