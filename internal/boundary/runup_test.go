package boundary

import (
	"testing"

	"airy-wave/internal/wave"
)

func TestRunupHeight(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	r := RunupHeight(p, 0.1)
	if r <= 0 {
		t.Fatalf("RunupHeight must be positive, got %v", r)
	}
	if RunupHeight(p, 0) != 0 {
		t.Fatalf("RunupHeight with zero slope must be 0")
	}
}

func TestTransmissionCoefficient(t *testing.T) {
	if TransmissionCoefficient(5, 10) <= 0 || TransmissionCoefficient(5, 10) > 1 {
		t.Fatalf("transmission out of (0,1]")
	}
	if TransmissionCoefficient(10, 10) != 0 {
		t.Fatalf("fully submerged breakwater => 0 transmission")
	}
}

func TestReflectionFromTwoWaveHeights(t *testing.T) {
	r := ReflectionFromTwoWaveHeights(1, 2)
	if r <= 0 || r > 1 {
		t.Fatalf("reflection out of (0,1]: %v", r)
	}
}
