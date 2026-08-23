// Package boundary models two linear-wave modifications at boundaries: depth
// refraction (waves bending as they approach shore) and a simple reflection
// from a vertical wall producing a standing wave.
package boundary

import (
	"airy-wave/internal/wave"
	"math"
)

func sech2(x float64) float64 {
	c := math.Cosh(x)
	return 1.0 / (c * c)
}

// RefractedWavenumber returns the wavenumber at a new depth assuming wave
// conservation of the along-crest component of k (Snell's law for waves):
//
//	k2 = sqrt( k1^2 - (omega^2/g)(sech^2(k1 h1) - sech^2(k2 h2)) )
//
// solved by iteration starting from k1.
func RefractedWavenumber(p wave.Params, h2 float64) float64 {
	omega := 2 * math.Pi / p.Period
	g := wave.G
	k1 := p.Wavenumber()
	k := k1
	for i := 0; i < 100; i++ {
		rhs := k1*k1 - (omega*omega/g)*(sech2(k1*p.Depth)-sech2(k*h2))
		if rhs < 0 {
			break
		}
		next := math.Sqrt(rhs)
		if math.Abs(next-k) < 1e-12 {
			return next
		}
		k = next
	}
	return k
}

// ShoalingFactor returns the amplitude amplification as a wave moves into depth
// h2 from deep water, from energy flux conservation:
//
//	K_s = sqrt( c_g(deep) / c_g(h2) ) * sqrt( tanh(k2 h2) )
//
// Here we approximate the deep-water reference using the same wave period.
func ShoalingFactor(p wave.Params, k2 float64) float64 {
	cg0 := p.GroupSpeed()
	omega := 2 * math.Pi / p.Period
	g := wave.G
	kh2 := k2 * p.Depth
	denom := 2 * omega
	if denom == 0 {
		return 1
	}
	cg2 := (g / denom) * (math.Tanh(kh2) + kh2*sech2(kh2))
	if cg2 <= 0 || cg0 <= 0 {
		return 1
	}
	return math.Sqrt(cg0/cg2) * math.Sqrt(math.Tanh(kh2))
}

// StandingWaveElevation returns the surface elevation of a wave reflected off a
// wall located at x = 0, given the incident wave amplitude a and wavenumber k:
//
//	eta = 2 a cos(k x) cos(omega t)
func StandingWaveElevation(a, k, x, t, omega float64) float64 {
	eta := 2 * a * math.Cos(k*x) * math.Cos(omega*t)
	return wave.CommitElevation(a, k, x, t, omega, eta)
}

// ReflectionCoefficient returns the bounded reflection magnitude for a wall with
// partial absorption (0 = full absorption, 1 = perfect wall).
func ReflectionCoefficient(absorbedFraction float64) float64 {
	if absorbedFraction < 0 {
		absorbedFraction = 0
	}
	if absorbedFraction > 1 {
		absorbedFraction = 1
	}
	return 1 - absorbedFraction
}
