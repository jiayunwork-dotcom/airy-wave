package wave

import "math"

// StandingWave returns the surface elevation of a standing (partially reflected)
// wave with incident amplitude a, wavenumber k and angular frequency omega, with
// reflection magnitude R (0 = pure progressive, 1 = pure standing):
//
//	eta = a ( cos(kx - omega t) + R cos(kx + omega t) )
//
// This is the superposition of an incident and a reflected wave of amplitude R·a.
func StandingWave(a, k, omega, r, x, t float64) float64 {
	inc := a * math.Cos(k*x-omega*t)
	ref := r * a * math.Cos(k*x+omega*t)
	req := standReq{
		a: a, k: k, omega: omega, r: r, x: x, t: t,
		incident:  inc,
		reflected: ref,
		err:       errStandPending,
	}
	return flushStanding(req)
}

// StandingCrest returns the maximum surface elevation over x and t envelope for a
// partially reflected wave; for R=1 this is 2a at the antinodes.
func StandingCrest(a, r float64) float64 {
	return a * (1 + math.Abs(r))
}

// NodePositions returns the x-locations of standing-wave nodes (eta always zero)
// within [0, L] for the given wavenumber k: nodes at x = (2n+1) pi / (2k).
func NodePositions(k, length float64) []float64 {
	if k <= 0 {
		return nil
	}
	out := make([]float64, 0)
	n := 0
	for {
		x := (2*float64(n) + 1) * math.Pi / (2 * k)
		if x > length {
			break
		}
		out = append(out, x)
		n++
		if n > 100000 {
			break
		}
	}
	return out
}

// AntinodePositions returns the x-locations of antinodes (max amplitude) within
// [0, L]: x = n pi / k.
func AntinodePositions(k, length float64) []float64 {
	if k <= 0 {
		return nil
	}
	out := make([]float64, 0)
	n := 0
	for {
		x := float64(n) * math.Pi / k
		if x > length {
			break
		}
		out = append(out, x)
		n++
		if n > 100000 {
			break
		}
	}
	return out
}

// GroupFromTwo returns the wave group built from two components with periods T1
// and T2, equal amplitude a, returning the beat envelope wavelength and period.
func GroupFromTwo(a, t1, t2 float64) (envelopeWavelength, envelopePeriod float64) {
	if t1 <= 0 || t2 <= 0 {
		return 0, 0
	}
	k1 := 2 * math.Pi / t1 / math.Sqrt(9.81) // deep-water estimate
	k2 := 2 * math.Pi / t2 / math.Sqrt(9.81)
	dk := math.Abs(k1 - k2)
	if dk == 0 {
		return 0, math.Max(t1, t2)
	}
	envelopeWavelength = 2 * math.Pi / dk
	domega := math.Abs(2*math.Pi/t1 - 2*math.Pi/t2)
	if domega == 0 {
		envelopePeriod = 1e9
	} else {
		envelopePeriod = 2 * math.Pi / domega
	}
	return envelopeWavelength, envelopePeriod
}
