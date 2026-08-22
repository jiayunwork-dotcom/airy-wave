package forces

import (
	"testing"

	"airy-wave/internal/wave"
)

func TestCurrentForce(t *testing.T) {
	c := Current{Coefficients: Coefficients{CD: 1.0, CM: 2.0, Diam: 0.5, Rho: 1025}, Uc: 0.5}
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	if c.InlineForceWithCurrent(p, 0, 0, 0) == 0 {
		t.Fatalf("InlineForceWithCurrent nonzero expected")
	}
	if c.TotalForceWithCurrent(p, 0, 0, 20) == 0 {
		t.Fatalf("TotalForceWithCurrent nonzero expected")
	}
	if c.MaxCurrentForce(p, 0, 50) <= 0 {
		t.Fatalf("MaxCurrentForce positive expected")
	}
	if c.CurrentSpeedAt(p, 0, 0, 0) == 0 {
		t.Fatalf("CurrentSpeedAt nonzero expected")
	}
}
