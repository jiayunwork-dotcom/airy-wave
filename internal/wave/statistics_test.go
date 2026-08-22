package wave

import (
	"math"
	"testing"
)

func TestZeroCrossingStats(t *testing.T) {
	// synthesize a simple sinusoid: should give nonzero H_s and several crossings
	n := 200
	eta := make([]float64, n)
	for i := range eta {
		eta[i] = 2.0 * math.Sin(float64(i)*0.3)
	}
	hs, hmax, tz := ZeroCrossingStats(eta, 0.1)
	if hs <= 0 || hmax <= 0 || tz <= 0 {
		t.Fatalf("stats non-positive: hs=%v hmax=%v tz=%v", hs, hmax, tz)
	}
}

func TestRMS(t *testing.T) {
	v := []float64{3, 4}
	got := RMS(v)
	// sqrt((9+16)/2) = sqrt(12.5) ≈ 3.5355
	if got < 3.53 || got > 3.54 {
		t.Fatalf("RMS mismatch: %v", got)
	}
}
