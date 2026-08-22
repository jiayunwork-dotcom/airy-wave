package wave

import (
	"math"
	"testing"
)

func TestBuildFromSpectrum(t *testing.T) {
	// constant spectrum
	b := BuildFromSpectrum(func(o float64) float64 { return 0.1 }, 0.5, 2.0, 20)
	if len(b.Components) != 20 {
		t.Fatalf("expected 20 components, got %d", len(b.Components))
	}
	v := b.SurfaceAt(0, 0)
	if math.IsNaN(v) {
		t.Fatalf("SurfaceAt NaN")
	}
}

func TestDesignWaveHeight(t *testing.T) {
	h := DesignWaveHeight(2.0, 1000)
	if h <= 2.0 {
		t.Fatalf("design wave should exceed H_s, got %v", h)
	}
	if h > 1.86*2.0+1e-9 {
		t.Fatalf("design wave capped at 1.86 H_s, got %v", h)
	}
}

func TestCrestHeightPositive(t *testing.T) {
	if CrestHeight(2.0, 100) <= 0 {
		t.Fatalf("CrestHeight must be positive")
	}
}
