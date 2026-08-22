package forces

import (
	"testing"

	"airy-wave/internal/wave"
)

func TestMomentPositive(t *testing.T) {
	c := Coefficients{CD: 1.0, CM: 2.0, Diam: 0.5, Rho: 1025}
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	if c.Moment(p, 0, 0, 20) == 0 {
		t.Fatalf("Moment should be nonzero")
	}
	if c.BaseShear(p, 0, 0) == 0 {
		t.Fatalf("BaseShear should be nonzero")
	}
	if c.InlineForceSurface(p, 0, 0) == 0 {
		t.Fatalf("InlineForceSurface should be nonzero")
	}
	if c.MaxMoment(p, 0, 50) <= 0 {
		t.Fatalf("MaxMoment must be positive")
	}
}
