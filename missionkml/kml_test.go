package missionkml

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"

	"flightplan2litchimission/fp2lm"
	"flightplan2litchimission/missioncsv"
)

var samplePath = "../fp2lm/testdata/FlightplannerMission.csv"

func getWaypoints(t *testing.T) []*missioncsv.LitchiWaypoint {
	t.Helper()
	data, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	options := &fp2lm.ConverterOptions{AltitudeMode: "agl", GimbalPitch: -90, MaxAltitudeAGL: 120}
	wps, err := fp2lm.ConvertWaypoints(bytes.NewReader(data), options)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}
	return wps
}

func TestWriteKML(t *testing.T) {
	wps := getWaypoints(t)
	var buf bytes.Buffer
	if err := WriteKML(&buf, wps); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("<kml")) {
		t.Errorf("missing kml tag")
	}
	if !bytes.Contains(buf.Bytes(), []byte("<Placemark>")) {
		t.Errorf("missing placemark")
	}
	if len(out) == 0 {
		t.Errorf("no output")
	}
}

func TestWriteKMZ(t *testing.T) {
	wps := getWaypoints(t)
	var buf bytes.Buffer
	if err := WriteKMZ(&buf, wps); err != nil {
		t.Fatalf("write kmz failed: %v", err)
	}
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip read: %v", err)
	}
	if len(r.File) != 1 || r.File[0].Name != "doc.kml" {
		t.Fatalf("kmz missing doc.kml")
	}
	f, err := r.File[0].Open()
	if err != nil {
		t.Fatalf("open kml: %v", err)
	}
	data, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Fatalf("read kml: %v", err)
	}
	if !bytes.Contains(data, []byte("<kml")) {
		t.Errorf("kml not found in kmz")
	}
}
