package wave

import (
	"testing"
)

func TestGroupSpeedDeepApproxHalf(t *testing.T) {
	// deep water: c_g ~ c/2
	p := Params{Amplitude: 1, Period: 8, Depth: 1000}
	c := p.PhaseSpeed()
	cg := p.GroupSpeed()
	if cg < 0.4*c || cg > 0.6*c {
		t.Fatalf("deep-water c_g (%v) not ~ c/2 (%v)", cg, c/2)
	}
}

func TestWavelengthDeepVsShallow(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 1000}
	ld := p.WavelengthDeep()
	ls := p.WavelengthShallow()
	if ld <= 0 || ls <= 0 {
		t.Fatalf("wavelength estimators must be positive")
	}
	if ld < p.Wavelength() {
		t.Fatalf("deep-water wavelength (%v) < actual (%v)", ld, p.Wavelength())
	}
}

func TestDispersionRatioRange(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	r := p.DispersionRatio()
	if r <= 0.5 || r > 1.0001 {
		t.Fatalf("DispersionRatio out of (0.5,1]: %v", r)
	}
}
