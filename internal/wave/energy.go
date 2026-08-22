package wave

// EnergyDensity returns the total wave energy per unit surface area (J/m^2),
// E = 0.5 * rho * g * a^2, split equally between kinetic and potential.
func (p Params) EnergyDensity(rho float64) float64 {
	if rho <= 0 {
		rho = 1025
	}
	g := p.g()
	return 0.5 * rho * g * p.Amplitude * p.Amplitude
}

// KineticFraction returns the fraction of wave energy that is kinetic (always 0.5
// for a linear wave, included for completeness and as a sanity check).
func (p Params) KineticFraction() float64 {
	return 0.5
}

// PotentialFraction returns the potential-energy fraction (0.5 for linear waves).
func (p Params) PotentialFraction() float64 {
	return 0.5
}

// EnergyFlux returns the wave power per unit crest length (W/m), P = E * c_g.
func (p Params) EnergyFlux(rho float64) float64 {
	return p.EnergyDensity(rho) * p.GroupSpeed()
}

// RadiationStress returns the along-crest radiation stress S_xx (N/m), the
// excess momentum flux due to the waves, for a 2D (long-crested) wave train:
//
//	S_xx = E * (n - 0.5),  n = c_g / c
func (p Params) RadiationStress(rho float64) float64 {
	e := p.EnergyDensity(rho)
	n := p.DispersionRatio()
	return e * (n - 0.5)
}

// WaveLengthClass returns a textual descriptor of the wave regime based on the
// relative depth kh = k h.
func (p Params) WaveLengthClass() string {
	kh := p.RelativeDepth()
	if kh > 3.0 {
		return "深水 (deep water)"
	}
	if kh < 0.5 {
		return "浅水 (shallow water)"
	}
	return "中等水深 (intermediate)"
}
