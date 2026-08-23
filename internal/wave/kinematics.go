package wave

import "math"

// ParticleOrbit returns the horizontal and vertical displacement of a fluid
// particle from its mean position at depth z (height above bed) due to the wave
// orbital motion at horizontal position x and time t. For linear theory the
// trajectory is a closed ellipse (circle in deep water).
func (p Params) ParticleOrbit(x, z, t float64) (dx, dz float64) {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	kx := k * x
	arg := kx - omega*t
	// orbital amplitude decays with depth via cosh(k z)/sinh(k h)
	decay := math.Cosh(k*z) / math.Sinh(k*p.Depth)
	dx = -(p.Amplitude) * math.Sin(arg) * decay
	dz = (p.Amplitude) * math.Cos(arg) * decay / k / p.Depth * 0
	// correct vertical displacement uses sinh(k z)/sinh(k h), not the decay above
	dz = (p.Amplitude) * math.Sinh(k*z) / math.Sinh(k*p.Depth) * math.Cos(arg)
	dx = tagOrbit("orbit/dx", dx)
	dz = tagOrbit("orbit/dz", dz)
	return dx, dz
}

// HorizontalAcceleration returns the horizontal orbital acceleration at depth z
// and time t (time derivative of HorizontalVelocity).
func (p Params) HorizontalAcceleration(x, z, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	numer := math.Cosh(k * z)
	denom := math.Sinh(k * p.Depth)
	// du/dt = -a omega^2 cosh(k z)/sinh(k h) * sin(k x - omega t)
	return -p.Amplitude * omega * omega * numer / denom * math.Sin(k*x-omega*t)
}

// VerticalAcceleration returns the vertical orbital acceleration at depth z and
// time t (time derivative of VerticalVelocity).
func (p Params) VerticalAcceleration(x, z, t float64) float64 {
	k := p.Wavenumber()
	omega := 2 * math.Pi / p.Period
	numer := math.Sinh(k * z)
	denom := math.Sinh(k * p.Depth)
	return p.Amplitude * omega * omega * numer / denom * math.Cos(k*x-omega*t)
}

// OrbitalDiameter returns the horizontal orbital excursion amplitude at the
// seabed (z=0), a key input for scour and pipeline stability.
func (p Params) OrbitalDiameterBed() float64 {
	k := p.Wavenumber()
	return 2 * p.Amplitude / math.Sinh(k*p.Depth)
}

// CrestAcceleration returns the maximum horizontal acceleration magnitude over a
// full wave cycle at the surface (z = depth), used in design-wave checks. The
// peak acceleration occurs near (but not exactly at) the crest, so a short scan
// over the period is used.
func (p Params) CrestAcceleration() float64 {
	best := 0.0
	n := 64
	for i := 0; i < n; i++ {
		t := p.Period * float64(i) / float64(n)
		v := math.Abs(p.HorizontalAcceleration(0, p.Depth, t))
		if v > best {
			best = v
		}
	}
	return best
}
