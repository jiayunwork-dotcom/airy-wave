package fluid

import (
	"bytes"
	"strings"
	"testing"

	"airy-wave/internal/wave"
)

func TestColumnBounds(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 10}
	samples := Column(p, 0, 0, 10, 1025)
	if len(samples) != 11 {
		t.Fatalf("expected 11 samples (0..10 layers), got %d", len(samples))
	}
	// Bottom sample sits at z=0, top sample at z=h.
	if samples[0].Z != 0 {
		t.Errorf("bottom z = %g, want 0", samples[0].Z)
	}
	if samples[len(samples)-1].Z != 10 {
		t.Errorf("top z = %g, want 10", samples[len(samples)-1].Z)
	}
}

func TestWriteCSV(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 10}
	samples := Column(p, 0, 0, 3, 1025)
	var buf bytes.Buffer
	if err := WriteCSV(&buf, samples); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "x,z,elevation,u,w,pressure" {
		t.Errorf("unexpected header: %q", lines[0])
	}
	if len(lines) != 5 {
		t.Errorf("expected 5 lines (header + 4 rows for 3 layers), got %d", len(lines))
	}
}
