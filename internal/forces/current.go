package forces

import (
	"airy-wave/internal/wave"
	"math"
)

// Current adds a steady uniform current Uc (m/s, positive with wave propagation)
// to the Morison model. The relative velocity u_rel = u_wave + Uc enters both
// the drag and the convective terms; here we add it only to the drag term,
// which is the usual engineering approximation for small currents.
type Current struct {
	Coefficients
	Uc float64
}

// InlineForceWithCurrent returns the Morison force per unit length including a
// steady current Uc added to the horizontal velocity before the drag term.
func (c Current) InlineForceWithCurrent(p wave.Params, x, z, t float64) float64 {
	rho := c.effectiveRho()
	u := p.HorizontalVelocity(x, 0, t) + c.Uc
	du := p.HorizontalAcceleration(x, z, t)
	drag := 0.5 * rho * c.CD * c.Diam * math.Abs(u) * u
	inertia := rho * c.CM * (math.Pi * c.Diam * c.Diam / 4) * du
	return drag + inertia
}

// TotalForceWithCurrent integrates the current-augmented force over the column.
func (c Current) TotalForceWithCurrent(p wave.Params, x, t float64, layers int) float64 {
	if layers < 1 {
		layers = 20
	}
	f := 0.0
	h := p.Depth
	dz := h / float64(layers)
	for i := 0; i < layers; i++ {
		z := (float64(i) + 0.5) * dz
		f += c.InlineForceWithCurrent(p, x, z, t) * dz
	}
	return f
}

// CurrentSpeedAt returns the combined horizontal speed (wave + current) at depth
// z and time t, used for structure design against combined loading.
func (c Current) CurrentSpeedAt(p wave.Params, x, z, t float64) float64 {
	return p.HorizontalVelocity(x, z, t) + c.Uc
}

// MaxCurrentForce returns the maximum combined force over a period.
func (c Current) MaxCurrentForce(p wave.Params, x float64, nt int) float64 {
	best := 0.0
	for i := 0; i < nt; i++ {
		t := p.Period * float64(i) / float64(maxInt(1, nt-1))
		v := math.Abs(c.TotalForceWithCurrent(p, x, t, 20))
		if v > best {
			best = v
		}
	}
	return best
}
