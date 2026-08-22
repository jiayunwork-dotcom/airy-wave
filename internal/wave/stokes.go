package wave

import "math"

// StokesSecondSurfaceElevation returns the second-order (Stokes) surface
// elevation, which adds the bound harmonic 2k and a set-down correction so the
// profile is more accurate for finite amplitude. The result is valid for
// moderate steepness; it is an extension of SurfaceElevation, not a replacement
// for the fundamental dispersion solve (k is reused from the linear solve).
func (p Params) StokesSecondSurfaceElevation(x, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	kh := k * p.Depth
	// second-order amplitude coefficient
	a2 := p.Amplitude * p.Amplitude * k / (2 * math.Sinh(2 * kh)) * math.Cosh(kh)
	eta1 := p.Amplitude * math.Sin(k*x-omega*t)
	eta2 := a2 * math.Sin(2*(k*x-omega*t))
	// second-order set-down (steady component)
	setdown := p.Amplitude * p.Amplitude * k / (2 * math.Tanh(kh))
	return eta1 + eta2 - setdown
}

// StokesSecondHorizontalVelocity adds the second-order horizontal velocity at
// depth z (height above bed) to the linear result.
func (p Params) StokesSecondHorizontalVelocity(x, z, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	kh := k * p.Depth
	base := p.HorizontalVelocity(x, z, t)
	// second-order super-harmonic contribution
	a2 := p.Amplitude * p.Amplitude * k / (2 * math.Sinh(2 * kh))
	u2 := a2 * omega * math.Cosh(2*k*z) / math.Sinh(2*kh) * math.Cos(2*(k*x-omega*t))
	_ = kh
	return base + u2
}

// Steepness returns the wave steepness epsilon = a*k (dimensionless), a common
// validity indicator for the linear model (epsilon << 1).
func (p Params) Steepness() float64 {
	return p.Amplitude * p.Wavenumber()
}

// UrsellNumber returns the Ursell number Ur = a*k^2*h / (k*h)^3 = a/(k*h^2), a
// measure of nonlinearity vs depth dispersion. Large Ur -> nonlinear shallow
// waves; small Ur -> linear deep waves.
func (p Params) UrsellNumber() float64 {
	k := p.Wavenumber()
	h := p.Depth
	if h <= 0 {
		return 0
	}
	return p.Amplitude / (k * h * h)
}

// LinearValid reports whether the linear (Airy) assumption is reasonable given
// a steepness threshold (epsilon_max). It uses the configured threshold when
// positive, else a default of 0.1.
func (p Params) LinearValid(epsilonMax float64) bool {
	if epsilonMax <= 0 {
		epsilonMax = 0.1
	}
	return p.Steepness() < epsilonMax
}
