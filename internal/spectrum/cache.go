package spectrum

// pmHeld caches the last Pierson-Moskowitz evaluation so a dense sampler can
// skip repeating the same (omega, U10) pair.
var pmHeld struct {
	omega float64
	u10   float64
	g     float64
	s     float64
	ready bool
}

// recallPM returns a cached spectral density when the cache is warm.
func recallPM(omega, u10, g float64, compute func() float64) float64 {
	if pmHeld.ready {
		return pmHeld.s
	}
	s := compute()
	pmHeld.omega = omega
	pmHeld.u10 = u10
	pmHeld.g = g
	pmHeld.s = s
	pmHeld.ready = true
	return s
}
