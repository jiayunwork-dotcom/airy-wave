package spectrum

import "testing"

func TestSpectralMoments(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	m0, _, m2, _, m4 := SpectralMoments(pts)
	if m0 <= 0 || m2 <= 0 || m4 <= 0 {
		t.Fatalf("moments must be positive: m0=%v m2=%v m4=%v", m0, m2, m4)
	}
	tp := MeanPeriod(pts)
	if tp <= 0 {
		t.Fatalf("MeanPeriod must be positive")
	}
	tz := MeanZeroCrossingPeriod(pts)
	if tz <= 0 {
		t.Fatalf("MeanZeroCrossingPeriod must be positive")
	}
	w := SpectralWidth(pts)
	if w < 0 || w > 1 {
		t.Fatalf("SpectralWidth out of [0,1]: %v", w)
	}
}
