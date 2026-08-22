package transform

import (
	"math"
	"testing"
)

func TestSurfaceAtSingleComponent(t *testing.T) {
	comp := []Component{{Amplitude: 1, Wavenumber: 0, Omega: 0, Phase: 0}}
	v := SurfaceAt(comp, 0, 0)
	if math.Abs(v-1) > 1e-9 {
		t.Fatalf("single static component should give amplitude, got %v", v)
	}
}

func TestSpectrumToComponents(t *testing.T) {
	pts := []SpectralPoint{{Omega: 1, S: 0.5, DOmega: 0.1}}
	comps := SpectrumToComponents(pts)
	if len(comps) != 1 {
		t.Fatalf("expected 1 component, got %v", len(comps))
	}
	if comps[0].Amplitude <= 0 {
		t.Fatalf("component amplitude must be positive")
	}
}

func TestMaxElevation(t *testing.T) {
	comps := []Component{{Amplitude: 2, Wavenumber: 1, Omega: 1, Phase: 0}}
	m := MaxElevation(comps, 10, 10, 20, 20)
	if m <= 0 || m > 2.0001 {
		t.Fatalf("MaxElevation out of bounds: %v", m)
	}
}
