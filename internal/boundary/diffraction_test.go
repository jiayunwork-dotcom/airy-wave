package boundary

import (
	"testing"

	"airy-wave/internal/wave"
)

func TestDiffractionCoefficient(t *testing.T) {
	// small ka -> no diffraction, K_d approaches 1
	if d := DiffractionCoefficient(1, 0.01); d <= 0 || d > 1.5 {
		t.Fatalf("small-ka diffraction out of range: %v", d)
	}
	// larger ka -> still bounded [0,2]
	if d := DiffractionCoefficient(2, 3); d < 0 || d > 2 {
		t.Fatalf("diffraction out of [0,2]: %v", d)
	}
}

func TestForceOnCylinderRatio(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	if r := ForceOnCylinderRatio(p, 1); r < 0 || r > 2 {
		t.Fatalf("ForceOnCylinderRatio out of range: %v", r)
	}
}
