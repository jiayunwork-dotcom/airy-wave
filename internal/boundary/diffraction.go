package boundary

import (
	"airy-wave/internal/wave"
	"math"
)

// DiffractionCoefficient returns the diffraction coefficient K_d for a vertical
// cylinder of radius a at a dimensionless distance from a structure, using the
// MacCamy-Fuchs first-order diffraction solution for a circular cylinder:
//
//	K_d = 1 - (J1'(ka)/J1(ka)) * (H1'(ka)/H1(ka))
//
// where the product of the Bessel ratios tends to 1 for small ka (no
// diffraction) and grows with ka. We approximate the Bessel functions with
// series/asymptotic forms sufficient for |ka| < 10.
func DiffractionCoefficient(a, ka float64) float64 {
	if ka < 0 {
		return 1
	}
	// for very small ka the asymptotic Hankel forms lose accuracy; physically
	// long waves diffract negligibly, so K_d -> 1.
	if ka < 0.1 {
		return 1
	}
	j1 := besselJ1(ka)
	jp1 := besselJ1Prime(ka)
	h1 := besselH1(ka)
	hp1 := besselH1Prime(ka)
	if j1 == 0 || h1 == 0 {
		return 1
	}
	ratio := (jp1 / j1) * (hp1 / h1)
	k := 1 - ratio
	if k < 0 {
		return 0
	}
	if k > 2 {
		return 2
	}
	return k
}

// besselJ1 returns the Bessel function J1(x) via a polynomial approximation valid
// for small-to-moderate x.
func besselJ1(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x < 3 {
		// series: J1 = x/2 - x^3/16 + x^5/384 - x^7/18432
		return x/2 - x*x*x/16 + x*x*x*x*x/384 - x*x*x*x*x*x*x/18432
	}
	// asymptotic for large x
	a := 1 - 9.0/(16*x*x) + 81.0/(512*x*x*x*x)
	return math.Sqrt(2/(math.Pi*x)) * a * math.Cos(x - 3*math.Pi/4)
}

// besselJ1Prime returns d/dx J1(x).
func besselJ1Prime(x float64) float64 {
	return besselJ0(x) - besselJ1(x)/x
}

// besselJ0 returns the Bessel function J0(x).
func besselJ0(x float64) float64 {
	if x < 3 {
		return 1 - x*x/4 + x*x*x*x/64 - x*x*x*x*x*x/2304
	}
	a := 1 - 1.0/(8*x*x) + 9.0/(128*x*x*x*x)
	return math.Sqrt(2/(math.Pi*x)) * a * math.Cos(x - math.Pi/4)
}

// besselH1 returns the Hankel function H1(x) = J1 + i Y1; we return its real
// part magnitude proxy via the asymptotic form.
func besselH1(x float64) float64 {
	a := 1 - 9.0/(16*x*x) + 81.0/(512*x*x*x*x)
	return math.Sqrt(2/(math.Pi*x)) * a
}

// besselH1Prime returns d/dx H1(x) (asymptotic real proxy).
func besselH1Prime(x float64) float64 {
	a := 1 - 25.0/(16*x*x) + 105.0/(512*x*x*x*x)
	return math.Sqrt(2/(math.Pi*x)) * a
}

// ForceOnCylinderRatio returns the normalized horizontal force on a vertical
// cylinder relative to the incident wave force, combining diffraction via K_d.
func ForceOnCylinderRatio(p wave.Params, radius float64) float64 {
	ka := p.Wavenumber() * radius
	return DiffractionCoefficient(radius, ka)
}
