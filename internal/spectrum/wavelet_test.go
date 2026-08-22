package spectrum

import "testing"

func TestBandEnergy(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	total := BandEnergy(pts, 0.1, 2.0)
	if total <= 0 {
		t.Fatalf("BandEnergy over full band must be positive")
	}
}

func TestSpectralSplit(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	ws, sw, ratio := SpectralSplit(pts)
	if ws < 0 || sw < 0 {
		t.Fatalf("split energies must be non-negative")
	}
	if ratio < 0 || ratio > 1 {
		t.Fatalf("split ratio out of [0,1]: %v", ratio)
	}
}

func TestNormalize(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	norm := Normalize(pts, 1.0)
	m0, _, _, _, _ := SpectralMoments(norm)
	if m0 < 0.99 || m0 > 1.01 {
		t.Fatalf("Normalize should give m0~1, got %v", m0)
	}
}

func TestPeakWidth(t *testing.T) {
	pts := BuildPM(10, 0.1, 2.0, 50)
	w := PeakWidth(pts, 0.5)
	if w <= 0 {
		t.Fatalf("PeakWidth must be positive")
	}
}
