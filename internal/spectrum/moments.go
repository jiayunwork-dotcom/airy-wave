package spectrum

import "math"

// SpectralMoments computes the first four spectral moments m0..m4 of a spectrum
// given already-sampled points. These moments yield bulk descriptors such as the
// mean period and the spectral bandwidth.
func SpectralMoments(points []Point) (m0, m1, m2, m3, m4 float64) {
	for _, p := range points {
		dw := dOmega(points, p)
		m0 += p.S * dw
		m1 += p.S * p.Omega * dw
		m2 += p.S * p.Omega * p.Omega * dw
		m3 += p.S * math.Pow(p.Omega, 3) * dw
		m4 += p.S * math.Pow(p.Omega, 4) * dw
	}
	return m0, m1, m2, m3, m4
}

// MeanPeriod returns the spectral mean period T_01 = m0/m1.
func MeanPeriod(points []Point) float64 {
	m0, m1, _, _, _ := SpectralMoments(points)
	if m1 == 0 {
		return 0
	}
	return m0 / m1
}

// MeanZeroCrossingPeriod returns T_02 = sqrt(m0/m2), the average period between
// zero upcrossings.
func MeanZeroCrossingPeriod(points []Point) float64 {
	m0, _, m2, _, _ := SpectralMoments(points)
	if m2 == 0 {
		return 0
	}
	return math.Sqrt(m0 / m2)
}

// SpectralWidth returns the dimensionless bandwidth epsilon = sqrt(1 - m2^2/(m0*m4)),
// a measure of spectral broadening (0 = perfectly narrow-band).
func SpectralWidth(points []Point) float64 {
	m0, _, m2, _, m4 := SpectralMoments(points)
	if m0 == 0 || m4 == 0 {
		return 0
	}
	eps := 1 - (m2*m2)/(m0*m4)
	if eps < 0 {
		return 0
	}
	return math.Sqrt(eps)
}
