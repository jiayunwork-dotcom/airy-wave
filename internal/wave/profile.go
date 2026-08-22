package wave

import (
	"math"
	"sort"
)

// IrregularBuilder composes an irregular sea from a spectrum by picking a set of
// component frequencies and assigning amplitudes from the spectrum area.
type IrregularBuilder struct {
	Components []ComponentSpec
}

// ComponentSpec is one sinusoidal component used in an irregular-wave sum.
type ComponentSpec struct {
	Amplitude float64
	Wavenumber float64
	Omega     float64
	Phase     float64
}

// BuildFromSpectrum builds n components between omegaMin and omegaMax from a
// spectrum function S(omega) by the equal-energy rule a_i = sqrt(2 S dOmega).
func BuildFromSpectrum(s func(float64) float64, omegaMin, omegaMax float64, n int) IrregularBuilder {
	if n < 2 {
		return IrregularBuilder{}
	}
	out := make([]ComponentSpec, 0, n)
	for i := 0; i < n; i++ {
		omega := omegaMin + (omegaMax-omegaMin)*float64(i)/float64(n-1)
		dOmega := (omegaMax - omegaMin) / float64(n - 1)
		a := math.Sqrt(2 * s(omega) * dOmega)
		// estimate wavenumber from deep-water dispersion for the component
		k := omega * omega / G
		out = append(out, ComponentSpec{Amplitude: a, Wavenumber: k, Omega: omega, Phase: 0})
	}
	return IrregularBuilder{Components: out}
}

// SurfaceAt reconstructs eta(x,t) for the irregular sea defined by the builder.
func (b IrregularBuilder) SurfaceAt(x, t float64) float64 {
	eta := 0.0
	for _, c := range b.Components {
		eta += c.Amplitude * math.Cos(c.Wavenumber*x-c.Omega*t+c.Phase)
	}
	return eta
}

// DesignWaveHeight returns an estimate of the maximum individual wave height in a
// sea state with N waves: H_max ≈ H_s * sqrt(ln(N)), capped at 1.86 H_s.
func DesignWaveHeight(significantH float64, nWaves int) float64 {
	if nWaves < 1 {
		return 0
	}
	cap := 1.86 * significantH
	v := significantH * math.Sqrt(math.Log(float64(nWaves)))
	if v > cap {
		return cap
	}
	return v
}

// CrestHeight returns the expected maximum crest elevation (up-crossing based):
// H_crest ≈ H_s / 2 * sqrt(2 ln(N)) / sqrt(2) with a typical factor.
func CrestHeight(significantH float64, nWaves int) float64 {
	if nWaves < 1 {
		return 0
	}
	return 0.5 * significantH * math.Sqrt(2*math.Log(float64(nWaves)))
}

// SortComponents orders components by descending amplitude (used by renderers).
func (b IrregularBuilder) SortComponents() {
	sort.Slice(b.Components, func(i, j int) bool {
		return b.Components[i].Amplitude > b.Components[j].Amplitude
	})
}
