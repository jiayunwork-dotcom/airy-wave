package boundary

import (
	"math"
	"testing"

	"airy-wave/internal/wave"
)

func TestRefractedWavenumber(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	k2 := RefractedWavenumber(p, 20)
	if k2 <= 0 {
		t.Fatalf("RefractedWavenumber must be positive, got %v", k2)
	}
}

func TestShoalingFactor(t *testing.T) {
	p := wave.Params{Amplitude: 1, Period: 8, Depth: 50}
	k2 := RefractedWavenumber(p, 20)
	k := ShoalingFactor(p, k2)
	if k <= 0 {
		t.Fatalf("ShoalingFactor must be positive, got %v", k)
	}
}

func TestStandingWaveElevation(t *testing.T) {
	v := StandingWaveElevation(1, 1, 0, 0, 0)
	if math.Abs(v-2) > 1e-9 {
		t.Fatalf("at antinode and t=0 elevation should be 2a, got %v", v)
	}
}

func TestReflectionCoefficient(t *testing.T) {
	if ReflectionCoefficient(0) != 1 {
		t.Fatalf("zero absorption => reflection 1")
	}
	if ReflectionCoefficient(1) != 0 {
		t.Fatalf("full absorption => reflection 0")
	}
}
