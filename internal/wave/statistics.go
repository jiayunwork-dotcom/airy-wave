package wave

import (
	"math"
	"sort"
)

// ZeroCrossingStats computes bulk statistics of a surface-elevation time series
// using the zero-upcrossing method: count upward crossings of the mean, then the
// significant wave height H_1/3 is the average of the highest third of the
// individual wave heights.
func ZeroCrossingStats(eta []float64, dt float64) (hs, hmax, tz float64) {
	n := len(eta)
	if n < 4 {
		return 0, 0, 0
	}
	mean := 0.0
	for _, v := range eta {
		mean += v
	}
	mean /= float64(n)
	// detect upward zero crossings
	cross := make([]float64, 0)
	for i := 1; i < n; i++ {
		if eta[i-1]-mean <= 0 && eta[i]-mean > 0 {
			// linear interpolation of crossing time
			t := float64(i-1) + (0-(eta[i-1]-mean))/((eta[i]-mean)-(eta[i-1]-mean))
			cross = append(cross, t*dt)
		}
	}
	if len(cross) < 2 {
		return 0, 0, 0
	}
	// individual wave heights = peaks between consecutive upcrossings
	heights := make([]float64, 0)
	for i := 0; i < len(cross)-1; i++ {
		h := 0.0
		a := int(cross[i] / dt)
		b := int(cross[i+1] / dt)
		if b > a {
			for j := a; j <= b && j < n; j++ {
				d := math.Abs(eta[j] - mean)
				if d > h {
					h = d
				}
			}
			heights = append(heights, h*2) // crest-to-trough ~ 2*max deviation
		}
	}
	if len(heights) == 0 {
		return 0, 0, 0
	}
	hmax = 0.0
	for _, h := range heights {
		if h > hmax {
			hmax = h
		}
	}
	// significant: average of highest third
	sort.Float64s(heights)
	third := len(heights) / 3
	if third < 1 {
		third = 1
	}
	sum := 0.0
	for _, h := range heights[len(heights)-third:] {
		sum += h
	}
	hs = sum / float64(len(heights)-third)
	tz = float64(len(cross)) * dt // zero-crossing period proxy
	return hs, hmax, tz
}

// RMS returns the root-mean-square of a series (used for H_rms = sqrt(8) * H_s/... ).
func RMS(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v * v
	}
	return math.Sqrt(sum / float64(len(values)))
}
