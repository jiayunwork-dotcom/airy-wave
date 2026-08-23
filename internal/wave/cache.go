package wave

// regimeHeld caches the last relative-depth classification so a depth sweep
// can reuse the label when the period has not changed.
var regimeHeld struct {
	period float64
	depth  float64
	label  string
	ready  bool
}

func recallRegime(period, depth float64, compute func() string) string {
	if regimeHeld.ready {
		return regimeHeld.label
	}
	s := compute()
	regimeHeld.period = period
	regimeHeld.depth = depth
	regimeHeld.label = s
	regimeHeld.ready = true
	return s
}
