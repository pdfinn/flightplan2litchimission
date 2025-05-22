package missionkml

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"flightplan2litchimission/missioncsv"
)

// WriteKML writes KML representing the provided waypoints.
func WriteKML(w io.Writer, wps []*missioncsv.LitchiWaypoint) error {
	if _, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, `<kml xmlns="http://www.opengis.net/kml/2.2">`); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, `<Document>`); err != nil {
		return err
	}
	for i, wp := range wps {
		if _, err := fmt.Fprintf(w, "  <Placemark>\n    <name>%d</name>\n    <Point><coordinates>%.7f,%.7f,%.3f</coordinates></Point>\n  </Placemark>\n", i+1, wp.Point.Longitude, wp.Point.Latitude, wp.Point.Altitude); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, `</Document>`); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, `</kml>`); err != nil {
		return err
	}
	return nil
}

// WriteKMZ writes a KMZ archive containing a doc.kml generated from the waypoints.
func WriteKMZ(w io.Writer, wps []*missioncsv.LitchiWaypoint) error {
	var buf bytes.Buffer
	if err := WriteKML(&buf, wps); err != nil {
		return err
	}
	zw := zip.NewWriter(w)
	f, err := zw.Create("doc.kml")
	if err != nil {
		return err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	return zw.Close()
}
