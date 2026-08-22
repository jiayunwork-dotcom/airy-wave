package spectrum

// BandEnergy returns the energy contained in a frequency band [omegaLo, omegaHi]
// by summing S(omega) * dOmega over the points that fall inside the band. Used
// to estimate, for example, the fraction of energy in swell vs wind sea.
func BandEnergy(points []Point, omegaLo, omegaHi float64) float64 {
	sum := 0.0
	for _, p := range points {
		if p.Omega >= omegaLo && p.Omega <= omegaHi {
			sum += p.S * dOmega(points, p)
		}
	}
	return sum
}

// SpectralSplit partitions a spectrum into wind-sea (omega > omegaP) and swell
// (omega <= omegaP) contributions, returning their energies and the split ratio.
func SpectralSplit(points []Point) (windSea, swell, ratio float64) {
	omegaP := 0.0
	best := 0.0
	for _, p := range points {
		if p.S > best {
			best = p.S
			omegaP = p.Omega
		}
	}
	for _, p := range points {
		e := p.S * dOmega(points, p)
		if p.Omega > omegaP {
			windSea += e
		} else {
			swell += e
		}
	}
	total := windSea + swell
	if total > 0 {
		ratio = windSea / total
	}
	return windSea, swell, ratio
}

// Normalize scales a spectrum so its zeroth moment (area) equals targetM0, which
// is handy for building a unit-energy spectrum for Monte-Carlo seeding.
func Normalize(points []Point, targetM0 float64) []Point {
	m0, _, _, _, _ := SpectralMoments(points)
	if m0 == 0 {
		return points
	}
	scale := targetM0 / m0
	out := make([]Point, len(points))
	for i, p := range points {
		out[i] = Point{Omega: p.Omega, S: p.S * scale}
	}
	return out
}

// PeakWidth returns the bandwidth (in omega) containing the central fraction f of
// the spectral energy around the peak. Points are indexed by rank order from the
// peak; we walk outward in both directions until the enclosed energy reaches
// f * m0.
func PeakWidth(points []Point, f float64) float64 {
	if f <= 0 || f >= 1 {
		return 0
	}
	if len(points) < 2 {
		return 0
	}
	m0, _, _, _, _ := SpectralMoments(points)
	target := f * m0

	// find peak index
	best := 0
	for i, p := range points {
		if p.S > points[best].S {
			best = i
		}
	}
	// sort indices by distance to peak index
	idx := make([]int, len(points))
	for i := range idx {
		idx[i] = i
	}
	// selection of expanding ring: walk nearest-neighbour both directions
	lo, hi := best, best
	cum := points[best].S * dOmega(points, points[best])
	for cum < target {
		leftOk := lo-1 >= 0
		rightOk := hi+1 < len(points)
		if !leftOk && !rightOk {
			break
		}
		// expand toward the side whose neighbor has larger S
		leftS := 0.0
		if leftOk {
			leftS = points[lo-1].S
		}
		rightS := 0.0
		if rightOk {
			rightS = points[hi+1].S
		}
		if leftS >= rightS && leftOk {
			lo--
			cum += points[lo].S * dOmega(points, points[lo])
		} else if rightOk {
			hi++
			cum += points[hi].S * dOmega(points, points[hi])
		} else if leftOk {
			lo--
			cum += points[lo].S * dOmega(points, points[lo])
		} else {
			break
		}
	}
	return points[hi].Omega - points[lo].Omega
}
