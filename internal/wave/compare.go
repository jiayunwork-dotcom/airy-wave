package wave

import (
	"math"
	"sort"
)

// DepthRegime classifies a wave by relative depth using the dimensionless
// parameter kh = k h. The boundaries follow coastal-engineering convention:
//
//	kh > 3.0            -> deep water
//	0.5 <= kh <= 3.0    -> intermediate
//	kh < 0.5            -> shallow water
func (p Params) DepthRegime() string {
	kh := p.RelativeDepth()
	switch {
	case kh > 3.0:
		return "deep"
	case kh >= 0.5:
		return "intermediate"
	default:
		return "shallow"
	}
}

// CompareDepths evaluates a single wave period across a range of depths and
// reports how wavelength and phase speed depart from the deep-water limit.
type DepthPoint struct {
	Depth        float64
	Wavelength   float64
	PhaseSpeed   float64
	LambdaRatio  float64 // lambda / lambda_deep
	SpeedRatio   float64 // c / c_deep
}

// CompareAcrossDepths sweeps depths from dMin to dMax (n points) for the same
// period and amplitude, returning the trajectory of dispersion quantities.
func (p Params) CompareAcrossDepths(dMin, dMax float64, n int) []DepthPoint {
	if n < 2 {
		return nil
	}
	ld := p.WavelengthDeep()
	cd := ld / p.Period
	out := make([]DepthPoint, 0, n)
	for i := 0; i < n; i++ {
		d := dMin + (dMax-dMin)*float64(i)/float64(n-1)
		w := Params{Amplitude: p.Amplitude, Period: p.Period, Depth: d}
		l := w.Wavelength()
		c := w.PhaseSpeed()
		out = append(out, DepthPoint{
			Depth:       d,
			Wavelength:  l,
			PhaseSpeed:  c,
			LambdaRatio: l / ld,
			SpeedRatio:  c / cd,
		})
	}
	return out
}

// ExtremeWavelength returns the minimum and maximum wavelength across a depth
// sweep, useful for bounding the range of validity of a chosen discretization.
func (p Params) ExtremeWavelength(dMin, dMax float64, n int) (minL, maxL float64) {
	pts := p.CompareAcrossDepths(dMin, dMax, n)
	if len(pts) == 0 {
		return 0, 0
	}
	minL = math.MaxFloat64
	for _, pt := range pts {
		if pt.Wavelength < minL {
			minL = pt.Wavelength
		}
		if pt.Wavelength > maxL {
			maxL = pt.Wavelength
		}
	}
	return minL, maxL
}

// RankedRegimes returns the set of regimes encountered in a depth sweep, ordered
// by first appearance.
func (p Params) RankedRegimes(dMin, dMax float64, n int) []string {
	pts := p.CompareAcrossDepths(dMin, dMax, n)
	seen := map[string]bool{}
	out := make([]string, 0)
	for _, pt := range pts {
		w := Params{Amplitude: p.Amplitude, Period: p.Period, Depth: pt.Depth}
		r := w.DepthRegime()
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

// SortDepthPoints orders a slice of DepthPoint by depth ascending (helper for
// plotting).
func SortDepthPoints(pts []DepthPoint) {
	sort.Slice(pts, func(i, j int) bool { return pts[i].Depth < pts[j].Depth })
}
