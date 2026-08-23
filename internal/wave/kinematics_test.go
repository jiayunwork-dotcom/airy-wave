package wave

import (
	"math"
	"testing"
)

func TestParticleOrbit(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	dx, dz := p.ParticleOrbit(10, p.Depth, 0)
	if dx == 0 {
		t.Fatalf("horizontal orbital displacement should be nonzero off crest")
	}
	if dz == 0 {
		t.Fatalf("vertical orbital displacement should be nonzero at surface")
	}
	if math.IsNaN(dx) || math.IsNaN(dz) {
		t.Fatalf("ParticleOrbit produced NaN")
	}
}

func TestHorizontalAccelerationFinite(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	a := p.HorizontalAcceleration(10, 50, 0)
	if math.IsNaN(a) {
		t.Fatalf("HorizontalAcceleration NaN")
	}
	// at the crest (x=0,t=0) horizontal velocity peaks so acceleration is 0;
	// check acceleration is nonzero off the crest.
	if p.HorizontalAcceleration(10, 50, 0) == 0 {
		t.Fatalf("HorizontalAcceleration off-crest should be nonzero")
	}
	if p.CrestAcceleration() == 0 {
		t.Fatalf("CrestAcceleration should be nonzero")
	}
}

func TestOrbitalDiameterBed(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 10}
	if p.OrbitalDiameterBed() <= 0 {
		t.Fatalf("OrbitalDiameterBed must be positive")
	}
}
