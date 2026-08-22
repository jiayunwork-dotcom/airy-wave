package forces

import (
	"airy-wave/internal/wave"
	"math"
)

// Moment computes the overturning moment (N·m) about the seabed for a vertical
// cylinder of diameter D under wave p, integrating the inline force over the
// submerged column. The moment arm is the depth z above the bed.
func (c Coefficients) Moment(p wave.Params, x, t float64, layers int) float64 {
	if layers < 1 {
		layers = 20
	}
	m := 0.0
	h := p.Depth
	dz := h / float64(layers)
	for i := 0; i < layers; i++ {
		z := (float64(i) + 0.5) * dz
		m += c.InlineForce(p, x, z, t) * z * dz
	}
	return m
}

// BaseShear returns the horizontal force at the bed (z=0) used as a conservative
// estimate of the mudline shear for pile design.
func (c Coefficients) BaseShear(p wave.Params, x, t float64) float64 {
	return c.InlineForce(p, x, 0, t)
}

// InlineForceDepth returns the inline force per unit length at the surface (z=h),
// the location of peak velocity under typical conditions.
func (c Coefficients) InlineForceSurface(p wave.Params, x, t float64) float64 {
	return c.InlineForce(p, x, p.Depth, t)
}

// MaxMoment returns the maximum overturning moment over one period at x by
// scanning nt time samples.
func (c Coefficients) MaxMoment(p wave.Params, x float64, nt int) float64 {
	best := 0.0
	for i := 0; i < nt; i++ {
		t := p.Period * float64(i) / float64(maxInt(1, nt-1))
		v := math.Abs(c.Moment(p, x, t, 20))
		if v > best {
			best = v
		}
	}
	return best
}
