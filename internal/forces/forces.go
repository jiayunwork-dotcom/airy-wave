// Package forces implements the Morison equation for the load exerted by a
// linear wave on a slender vertical cylinder. The total force is the sum of an
// inertial term (proportional to fluid acceleration) and a drag term
// (proportional to velocity squared).
package forces

import (
	"airy-wave/internal/wave"
	"math"
)

// Coefficients bundles the Morison coefficients and geometry for a cylinder.
type Coefficients struct {
	CD   float64 // drag coefficient
	CM   float64 // inertia coefficient
	Diam float64 // cylinder diameter, m
	Rho  float64 // water density, kg/m^3 (0 => 1025)
}

// effectiveRho returns the density, defaulting to seawater when unset.
func (c Coefficients) effectiveRho() float64 {
	if c.Rho <= 0 {
		return 1025
	}
	return c.Rho
}

// InlineForce returns the Morison force per unit length (N/m) at depth z (height
// above bed) and time t for a wave p. The inline force is
//
//	f = 0.5 rho CD D |u| u + rho CM (pi D^2/4) du/dt
func (c Coefficients) InlineForce(p wave.Params, x, z, t float64) float64 {
	rho := c.effectiveRho()
	u := p.HorizontalVelocity(x, z, t)
	du := p.HorizontalAcceleration(x, z, t)
	drag := 0.5 * rho * c.CD * c.Diam * math.Abs(u) * u
	inertia := rho * c.CM * (math.Pi * c.Diam * c.Diam / 4) * du
	return drag + inertia
}

// TotalForce integrates InlineForce over the full column height h at time t for
// a vertical cylinder of diameter D from the bed (z=0) to the surface (z=h).
func (c Coefficients) TotalForce(p wave.Params, x, t float64, layers int) float64 {
	if layers < 1 {
		layers = 20
	}
	f := 0.0
	h := p.Depth
	dz := h / float64(layers)
	for i := 0; i < layers; i++ {
		// sample at the midpoint of each submerged slab from the bed up to the
		// still-water surface; the dynamic force acts over the full column.
		z := (float64(i) + 0.5) * dz
		f += c.InlineForce(p, x, z, t) * dz
	}
	return f
}

// MaxDrag returns the maximum drag contribution over one period at the bed, a
// common fatigue-driving quantity. Scans nt time samples.
func (c Coefficients) MaxDrag(p wave.Params, x float64, nt int) float64 {
	rho := c.effectiveRho()
	period := p.Period
	best := 0.0
	for i := 0; i < nt; i++ {
		t := period * float64(i) / float64(maxInt(1, nt-1))
		u := p.HorizontalVelocity(x, 0, t)
		d := 0.5 * rho * c.CD * c.Diam * math.Abs(u) * u
		if math.Abs(d) > best {
			best = math.Abs(d)
		}
	}
	return best
}

// KCCheck returns the Keulegan-Carpenter number KC = U_max T / D at the bed,
// where U_max is the maximum horizontal orbital velocity over a period.
func (c Coefficients) KCCheck(p wave.Params, x float64, nt int) float64 {
	period := p.Period
	best := 0.0
	for i := 0; i < nt; i++ {
		t := period * float64(i) / float64(maxInt(1, nt-1))
		u := math.Abs(p.HorizontalVelocity(x, 0, t))
		if u > best {
			best = u
		}
	}
	if c.Diam <= 0 {
		return 0
	}
	return best * period / c.Diam
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
