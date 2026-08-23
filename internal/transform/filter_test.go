package transform

import (
	"math"
	"testing"
)

func TestLowPassReducesNoise(t *testing.T) {
	y := make([]float64, 100)
	for i := range y {
		y[i] = 1.0
		if i%2 == 0 {
			y[i] += 0.5
		}
	}
	lp := LowPass(y, 0.1, 0.1)
	if math.Abs(lp[len(lp)-1]-1.0) > 0.5 {
		t.Fatalf("LowPass should move toward 1.0, got %v", lp[len(lp)-1])
	}
}

func TestMovingAverageFlat(t *testing.T) {
	y := []float64{2, 2, 2, 2, 2}
	ma := MovingAverage(y, 3)
	for _, v := range ma {
		if v != 2 {
			t.Fatalf("MovingAverage of constant must be constant, got %v", v)
		}
	}
}

func TestPeakDetect(t *testing.T) {
	y := []float64{0, 1, 0, 2, 0}
	peaks := PeakDetect(y)
	if len(peaks) != 2 {
		t.Fatalf("expected 2 peaks, got %v", len(peaks))
	}
}
