package boundary

import (
	"airy-wave/internal/wave"
	"math"
)

// RunupHeight estimates the maximum shoreline run-up elevation R for a normally
// incident wave of height H and period T on a slope beta (radians) using the
// empirical Hunt (1959) relation:
//
//	R / H = 0.82 / sqrt(tan(beta)) * (H/L)^(1/2)
//
// where L is the deep-water wavelength for the given period.
func RunupHeight(p wave.Params, beta float64) float64 {
	if beta <= 0 {
		return 0
	}
	h := 2 * p.Amplitude // wave height = 2 amplitude
	l0 := p.WavelengthDeep()
	if l0 <= 0 {
		return 0
	}
	return 0.82 / math.Sqrt(math.Tan(beta)) * math.Sqrt(h/l0) * h
}

// TransmissionCoefficient returns the transmission coefficient K_t for a submerged
// breakwater of height hB (relative to still water) and crest depth d, using the
// Goda-Datta form K_t = sqrt(1 - (hB/d)^2) bounded to [0,1].
func TransmissionCoefficient(hB, depth float64) float64 {
	if depth <= 0 {
		return 0
	}
	if hB >= depth {
		return 0
	}
	v := 1 - (hB*hB)/(depth*depth)
	if v < 0 {
		return 0
	}
	return math.Sqrt(v)
}

// ReflectionFromTwoWaveHeights estimates the reflection coefficient from measured
// incident H_i and combined (incident+reflected) H_c via
//
//	Kr = (H_c - H_i) / (H_c + H_i),  bounded to [0,1].
func ReflectionFromTwoWaveHeights(hIncident, hCombined float64) float64 {
	if hCombined <= 0 {
		return 0
	}
	if hIncident > hCombined {
		return 0
	}
	r := (hCombined - hIncident) / (hCombined + hIncident)
	if r < 0 {
		return 0
	}
	if r > 1 {
		return 1
	}
	return r
}
