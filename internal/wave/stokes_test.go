package wave

import (
	"math"
	"testing"
)

func TestStokesSecondSurfaceElevation(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	v := p.StokesSecondSurfaceElevation(0, 0)
	if math.IsNaN(v) {
		t.Fatalf("StokesSecondSurfaceElevation returned NaN")
	}
}

func TestSteepnessAndUrsell(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	if p.Steepness() <= 0 {
		t.Fatalf("Steepness should be positive, got %v", p.Steepness())
	}
	if p.UrsellNumber() <= 0 {
		t.Fatalf("UrsellNumber should be positive, got %v", p.UrsellNumber())
	}
}

func TestLinearValidDefaults(t *testing.T) {
	p := Params{Amplitude: 0.5, Period: 8, Depth: 50}
	if !p.LinearValid(0) {
		t.Fatalf("small-amplitude wave should be linear-valid")
	}
}
