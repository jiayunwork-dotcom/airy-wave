package forces

import (
	"math"
	"testing"

	"airy-wave/internal/wave"
)

func TestInlineForceFinite(t *testing.T) {
	c := Coefficients{CD: 1.0, CM: 2.0, Diam: 0.5, Rho: 1025}
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	f := c.InlineForce(p, 0, 0, 0)
	if math.IsNaN(f) {
		t.Fatalf("InlineForce NaN")
	}
}

func TestTotalForcePositive(t *testing.T) {
	c := Coefficients{CD: 1.0, CM: 2.0, Diam: 0.5, Rho: 1025}
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	f := c.TotalForce(p, 0, 0, 20)
	if f == 0 {
		t.Fatalf("TotalForce should be nonzero")
	}
}

func TestKCAndMaxDrag(t *testing.T) {
	c := Coefficients{CD: 1.0, CM: 2.0, Diam: 0.5, Rho: 1025}
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	if c.KCCheck(p, 0, 50) <= 0 {
		t.Fatalf("KC must be positive")
	}
	if c.MaxDrag(p, 0, 50) <= 0 {
		t.Fatalf("MaxDrag must be positive")
	}
}
