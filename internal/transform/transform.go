// Package transform holds helpers for turning a set of wave components into a
// real-time surface elevation series using a direct (non-FFT) superposition.
// It is intentionally dependency-free so it can run in the same process as the
// rest of the model.
package transform

import (
	"math"

	"airy-wave/internal/wave"
)

// Component is a single sinusoidal component in a superposition model.
type Component struct {
	Amplitude float64 // component amplitude a_i, m
	Wavenumber float64 // wavenumber k_i, 1/m
	Omega     float64 // angular frequency, rad/s
	Phase     float64 // phase offset, rad
}

// SurfaceAt reconstructs the surface elevation eta(x,t) by summing components:
//
//	eta = sum_i a_i * cos(k_i x - omega_i t + phi_i)
func SurfaceAt(components []Component, x, t float64) float64 {
	eta := 0.0
	for _, c := range components {
		eta += c.Amplitude * math.Cos(c.Wavenumber*x - c.Omega*t + c.Phase)
	}
	return eta
}

// SpectrumToComponents turns a wind-wave spectrum into a set of equal-amplitude
// components by assigning amplitude a_i = sqrt(2 * S(omega_i) * dOmega) to each
// spectral point.
func SpectrumToComponents(points []SpectralPoint) []Component {
	out := make([]Component, 0, len(points))
	for _, p := range points {
		a := math.Sqrt(2 * p.S * p.DOmega)
		out = append(out, Component{Amplitude: a, Wavenumber: p.Omega * p.Omega / (2 * wave.G), Omega: p.Omega})
	}
	return out
}

// SpectralPoint mirrors spectrum.Point but adds explicit spacing for amplitude
// conversion.
type SpectralPoint struct {
	Omega  float64
	S      float64
	DOmega float64
}

// MaxElevation returns the largest |eta| over a grid of (x, t) samples, a crude
// envelope estimate for an irregular sea.
func MaxElevation(components []Component, xMax, tMax float64, nx, nt int) float64 {
	maxv := 0.0
	for i := 0; i < nx; i++ {
		x := xMax * float64(i) / float64(max(1, nx-1))
		for j := 0; j < nt; j++ {
			t := tMax * float64(j) / float64(max(1, nt-1))
			v := math.Abs(SurfaceAt(components, x, t))
			if v > maxv {
				maxv = v
			}
		}
	}
	return maxv
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
