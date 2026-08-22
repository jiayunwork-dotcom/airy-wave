package spectrum

import (
	"testing"
)

func TestPiersonMoskowitzPositive(t *testing.T) {
	s := PiersonMoskowitz(1.0, 10, 9.81)
	if s <= 0 {
		t.Fatalf("PM spectrum must be positive for omega>0, got %v", s)
	}
	if PiersonMoskowitz(0, 10, 9.81) != 0 {
		t.Fatalf("PM at omega=0 must be 0")
	}
}

func TestJONSWAPPeaked(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	if len(pts) != 50 {
		t.Fatalf("BuildPM length = %d, want 50", len(pts))
	}
	hs := SignificantWaveHeight(pts)
	if hs <= 0 {
		t.Fatalf("H_s must be positive, got %v", hs)
	}
	tp := PeakPeriod(pts)
	if tp <= 0 {
		t.Fatalf("T_p must be positive")
	}
}

func TestJONSWAPDefaultGamma(t *testing.T) {
	// JONSWAP should be >= PM at the peak with gamma defaulting to 3.3.
	pm := PiersonMoskowitz(0.8, 10, 9.81)
	j := JONSWAP(0.8, 10, 0, 9.81)
	if j <= pm {
		t.Fatalf("JONSWAP with gamma>=1 should exceed PM at peak, got j=%v pm=%v", j, pm)
	}
}
