package spectrum

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteCSV(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 10)
	var buf bytes.Buffer
	if err := WriteCSV(&buf, pts); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}
	if !strings.Contains(buf.String(), "omega") {
		t.Fatalf("header missing")
	}
	if strings.Count(buf.String(), "\n") < 10 {
		t.Fatalf("expected header + 10 rows")
	}
}

func TestSampleByStep(t *testing.T) {
	var buf bytes.Buffer
	n, err := SampleByStep(&buf, func(o float64) float64 { return 0.1 }, 0.5, 2.0, 8)
	if err != nil {
		t.Fatalf("SampleByStep error: %v", err)
	}
	if n != 8 {
		t.Fatalf("expected 8 rows, got %v", n)
	}
}

func TestTotalEnergy(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	if TotalEnergy(pts) <= 0 {
		t.Fatalf("TotalEnergy must be positive")
	}
}
