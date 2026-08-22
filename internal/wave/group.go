package wave

import "math"

// GroupSpeed returns the group velocity c_g = d omega / d k. Using the
// dispersion relation omega^2 = g k tanh(k h):
//
//	c_g = (g / (2 omega)) * (tanh(k h) + k h sech^2(k h))
//
// In deep water c_g = c/2; in shallow water c_g = c.
func (p Params) GroupSpeed() float64 {
	g := p.g()
	omega := 2 * math.Pi / p.Period
	k := p.Wavenumber()
	kh := k * p.Depth
	denom := 2 * omega
	if denom == 0 {
		return 0
	}
	inner := math.Tanh(kh) + kh*sech2(kh)
	return (g / denom) * inner
}

// sech2 returns sech^2(x) = 1 / cosh^2(x).
func sech2(x float64) float64 {
	c := math.Cosh(x)
	return 1.0 / (c * c)
}

// WavelengthDeep returns the deep-water wavelength for the same period, i.e. the
// limit as h -> infinity: lambda_0 = 2 pi g / omega^2.
func (p Params) WavelengthDeep() float64 {
	g := p.g()
	omega := 2 * math.Pi / p.Period
	return 2 * math.Pi * g / (omega * omega)
}

// WavelengthShallow returns the shallow-water wavelength for the same period,
// the limit as kh -> 0: lambda = T * sqrt(g h).
func (p Params) WavelengthShallow() float64 {
	g := p.g()
	return p.Period * math.Sqrt(g*p.Depth)
}

// RelativeDepth returns kh = k * h, the key scaling between deep and shallow.
func (p Params) RelativeDepth() float64 {
	return p.Wavenumber() * p.Depth
}

// DispersionRatio returns c_g / c, a value in (0.5, 1] used to classify depth.
func (p Params) DispersionRatio() float64 {
	c := p.PhaseSpeed()
	if c == 0 {
		return 0
	}
	return p.GroupSpeed() / c
}

// GroupEnvelopePhase returns the phase of the group envelope at position x and
// time t for a carrier wave with angular frequency omega and a small frequency
// detuning domega (rad/s). The envelope travels at the group speed c_g, so its
// phase is (k - domega/c_g)·x - domega·t; maxima occur where this is a multiple
// of 2 pi.
func (p Params) GroupEnvelopePhase(x, t, domega float64) float64 {
	cg := p.GroupSpeed()
	if cg == 0 {
		return 0
	}
	return (p.Wavenumber() - domega/cg) * x - domega*t
}
