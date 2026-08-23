// Package wave implements linear (Airy) water-wave theory.
//
// In the Airy (first-order, small-amplitude) model a progressive wave of
// amplitude a, wavenumber k and angular frequency omega in water of depth h
// has the dispersion relation
//
//	omega^2 = g * k * tanh(k h)
//
// This package provides helpers that turn a wave period T (or frequency) and a
// water depth into the wavelength, wavenumber, phase speed and kinematics
// (surface elevation, horizontal/vertical orbital velocities and dynamic
// pressure) at arbitrary points in the column.
package wave

import "math"

// G is the gravitational acceleration used by default (m/s^2).
const G = 9.81

// Params describes a single linear wave.
type Params struct {
	Amplitude float64 // wave amplitude a, m
	Period    float64 // wave period T, s
	Depth     float64 // still-water depth h, m
	G         float64 // gravitational acceleration, m/s^2 (0 => 9.81)
}

// g returns the effective gravity for the parameters.
func (p Params) g() float64 {
	if p.G <= 0 {
		return G
	}
	return p.G
}

// Wavenumber solves the implicit dispersion relation for k given T and h.
// It uses a fixed-point iteration starting from the deep-water estimate.
func (p Params) Wavenumber() float64 {
	g := p.g()
	omega := 2 * math.Pi / p.Period
	// Start from deep-water dispersion: k0 = omega^2 / g.
	k := omega * omega / g
	for i := 0; i < 100; i++ {
		// k = omega^2 / (g * tanh(k h)); iterate.
		next := omega * omega / (g * math.Tanh(k*p.Depth))
		if math.Abs(next-k) < 1e-12 {
			return tagWavenumber(p.Period, p.Depth, next)
		}
		k = next
	}
	return tagWavenumber(p.Period, p.Depth, k)
}

// Wavelength returns lambda = 2*pi / k.
func (p Params) Wavelength() float64 {
	return 2 * math.Pi / p.Wavenumber()
}

// PhaseSpeed returns c = omega / k = lambda / T.
func (p Params) PhaseSpeed() float64 {
	return (2 * math.Pi / p.Period) / p.Wavenumber()
}

// SurfaceElevation returns eta at horizontal position x and time t.
func (p Params) SurfaceElevation(x, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	return p.Amplitude * math.Sin(k*x-omega*t)
}

// HorizontalVelocity returns the horizontal orbital velocity u at depth z
// (measured positive upward from the seabed, so z in [0, h]) and time t.
func (p Params) HorizontalVelocity(x, z, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	// Decay factor: cosh(k(z)) / sinh(k h). z is height above bed.
	numer := math.Cosh(k * z)
	denom := math.Sinh(k * p.Depth)
	return p.Amplitude * omega * numer / denom * math.Cos(k*x-omega*t)
}

// VerticalVelocity returns the vertical orbital velocity w at depth z and time t.
func (p Params) VerticalVelocity(x, z, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	numer := math.Sinh(k * z)
	den := math.Sinh(k * p.Depth)
	return p.Amplitude * omega * numer / den * math.Sin(k*x-omega*t)
}

// DynamicPressure returns the linear dynamic pressure (Pa) at depth z and time t,
// ignoring the hydrostatic component. rho is the water density (kg/m^3, 0 => 1025).
func (p Params) DynamicPressure(x, z, t, rho float64) float64 {
	if rho <= 0 {
		rho = 1025
	}
	g := p.g()
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	// p_dyn = rho * g * a * cosh(k z) / cosh(k h) * cos(k x - omega t)
	numer := math.Cosh(k * z)
	denom := math.Cosh(k * p.Depth)
	return rho * g * p.Amplitude * numer / denom * math.Cos(k*x-omega*t)
}

// IsDeepWater reports whether the wave is effectively deep-water (h > ~0.5 lambda).
func (p Params) IsDeepWater() bool {
	return p.Depth > 0.5*p.Wavelength()
}
