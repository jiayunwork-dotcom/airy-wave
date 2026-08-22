// Package spectrum provides simple ocean-wave spectra and spectral statistics
// built on top of the linear (Airy) dispersion relation. Spectra map a set of
// component frequencies to energy densities; the integration of a spectrum gives
// bulk quantities such as the significant wave height.
package spectrum

import (
	"airy-wave/internal/wave"
	"math"
)

// Point is a single spectral component: an angular frequency omega (rad/s) and
// its energy density S(omega) (m^2 s).
type Point struct {
	Omega float64
	S     float64
}

// PiersonMoskowitz returns the Pierson-Moskowitz (fully-developed sea) spectrum
// value for angular frequency omega given wind speed U10 (m/s) at 10 m height.
//
//	S(omega) = (alpha g^2 / omega^5) * exp(-beta (omega0/omega)^4)
//
// with omega0 = g / U10 * 1.29672, alpha = 0.0081, beta = 1.25.
func PiersonMoskowitz(omega, u10, g float64) float64 {
	if g <= 0 {
		g = wave.G
	}
	if omega <= 0 || u10 <= 0 {
		return 0
	}
	alpha := 0.0081
	beta := 1.25
	omega0 := g / u10 * 1.29672
	return alpha * g * g / math.Pow(omega, 5) * math.Exp(-beta*math.Pow(omega0/omega, 4))
}

// JONSWAP returns the JONSWAP spectrum value for omega given wind speed U10,
// peak enhancement factor gamma (typical 3.3) and fetch factor. gamma <= 0
// defaults to 3.3.
func JONSWAP(omega, u10, gamma, g float64) float64 {
	if gamma <= 0 {
		gamma = 3.3
	}
	if g <= 0 {
		g = wave.G
	}
	if omega <= 0 || u10 <= 0 {
		return 0
	}
	pm := PiersonMoskowitz(omega, u10, g)
	omegaP := g / u10 * 1.29672
	sigma := 0.09
	if omega > omegaP {
		sigma = 0.07
	}
	r := math.Exp(-math.Pow((omega-omegaP)/(sigma*omegaP), 2) / 2)
	return pm * math.Pow(gamma, r)
}

// SignificantWaveHeight integrates a spectrum over its components to the
// significant wave height H_s = 4 * sqrt(m0), where m0 is the zeroth moment
// (area under S(omega)).
func SignificantWaveHeight(points []Point) float64 {
	m0 := 0.0
	for _, p := range points {
		m0 += p.S * dOmega(points, p)
	}
	return 4 * math.Sqrt(m0)
}

// dOmega returns the frequency spacing for a point within an ordered slice. For
// edge points a one-sided spacing is used; interior points use central spacing.
func dOmega(points []Point, p Point) float64 {
	n := len(points)
	if n < 2 {
		return 1.0
	}
	idx := -1
	for i, q := range points {
		if q.Omega == p.Omega {
			idx = i
			break
		}
	}
	if idx < 0 {
		return 1.0
	}
	if idx == 0 {
		return points[1].Omega - points[0].Omega
	}
	if idx == n-1 {
		return points[n-1].Omega - points[n-2].Omega
	}
	return 0.5 * (points[idx+1].Omega - points[idx-1].Omega)
}

// PeakPeriod returns the period T_p = 2*pi/omega_p where omega_p maximizes S.
func PeakPeriod(points []Point) float64 {
	best := 0.0
	bestS := 0.0
	for _, p := range points {
		if p.S > bestS {
			bestS = p.S
			best = p.Omega
		}
	}
	if best == 0 {
		return 0
	}
	return 2 * math.Pi / best
}

// BuildPM samples the Pierson-Moskowitz spectrum over [omegaMin, omegaMax]
// using n evenly spaced frequencies.
func BuildPM(u10, omegaMin, omegaMax float64, n int) []Point {
	if n < 2 {
		return nil
	}
	out := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		omega := omegaMin + (omegaMax-omegaMin)*float64(i)/float64(n-1)
		out = append(out, Point{Omega: omega, S: PiersonMoskowitz(omega, u10, 0)})
	}
	return out
}
