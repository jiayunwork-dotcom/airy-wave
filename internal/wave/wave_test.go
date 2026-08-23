package wave

import "testing"

func TestWavenumberDeepWater(t *testing.T) {
	// Deep water: h >> wavelength. T=8 s -> omega=0.7854, deep k0=omega^2/g=0.0629.
	p := Params{Amplitude: 1, Period: 8, Depth: 100000}
	k := p.Wavenumber()
	// For very deep water tanh(kh)->1, so k ≈ omega^2/g = 0.7854^2/9.81 = 0.0629.
	if abs(k-0.0629) > 1e-3 {
		t.Errorf("deep-water k = %g, want ~0.0629", k)
	}
}

func TestWavelengthPositive(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 10}
	if l := p.Wavelength(); l <= 0 {
		t.Errorf("wavelength must be positive, got %g", l)
	}
}

func TestPhaseSpeedMatchesLambdaOverT(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 20}
	got := p.PhaseSpeed()
	want := p.Wavelength() / 8.0
	if abs(got-want) > 1e-9 {
		t.Errorf("phase speed = %g, want lambda/T = %g", got, want)
	}
}

func TestSurfaceElevationZeroAtCrestNode(t *testing.T) {
	p := Params{Amplitude: 2, Period: 8, Depth: 10}
	// At omega*t = pi, sin(-pi)=0 at x=0.
	if e := p.SurfaceElevation(0, 4); abs(e) > 1e-9 {
		t.Errorf("eta at t=T/2 should be 0, got %g", e)
	}
}

func TestVelocitiesAtBedZeroPlatform(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 10}
	// At the bed z=0, sinh(0)=0 so vertical velocity is exactly 0.
	if w := p.VerticalVelocity(0, 0, 0); abs(w) > 1e-9 {
		t.Errorf("vertical velocity at bed should be 0, got %g", w)
	}
}

func TestDynamicPressureSign(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 10}
	// Pressure increases toward the surface: at z=h the ratio cosh(kh)/cosh(kh)=1
	// (maximum), at the bed z=0 the ratio is 1/cosh(kh) (smaller).
	ps := p.DynamicPressure(0, p.Depth, 0, 1025)
	pb := p.DynamicPressure(0, 0, 0, 1025)
	if abs(ps) <= abs(pb) {
		t.Errorf("surface pressure magnitude should exceed bed pressure: %g vs %g", ps, pb)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
