package wave

import "testing"

func TestEnergyDensityPositive(t *testing.T) {
	p := Params{Amplitude: 1.5, Period: 8, Depth: 50}
	if p.EnergyDensity(1025) <= 0 {
		t.Fatalf("EnergyDensity must be positive")
	}
	if p.KineticFraction() != 0.5 || p.PotentialFraction() != 0.5 {
		t.Fatalf("linear wave energy split must be 0.5/0.5")
	}
}

func TestEnergyFluxAndRadiation(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	if p.EnergyFlux(1025) <= 0 {
		t.Fatalf("EnergyFlux must be positive")
	}
	if p.RadiationStress(1025) <= 0 {
		t.Fatalf("RadiationStress must be positive")
	}
}

func TestWaveLengthClass(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 1000}
	if p.WaveLengthClass() == "" {
		t.Fatalf("WaveLengthClass empty")
	}
}
