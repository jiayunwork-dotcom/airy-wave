package transform

import "math"

// LowPass filters a time series y with a first-order (RC) low-pass filter with
// cutoff frequency fc (Hz) and timestep dt (s). It returns a new slice; the
// recurrence is y_out[i] = y_out[i-1] + alpha*(y[i]-y_out[i-1]) with
// alpha = dt / (RC + dt), RC = 1/(2 pi fc).
func LowPass(y []float64, fc, dt float64) []float64 {
	n := len(y)
	if n == 0 || fc <= 0 || dt <= 0 {
		return append([]float64(nil), y...)
	}
	rc := 1.0 / (2 * math.Pi * fc)
	alpha := dt / (rc + dt)
	out := make([]float64, n)
	out[0] = y[0]
	for i := 1; i < n; i++ {
		out[i] = out[i-1] + alpha*(y[i]-out[i-1])
	}
	return out
}

// HighPass is the complementary first-order high-pass filter.
func HighPass(y []float64, fc, dt float64) []float64 {
	lp := LowPass(y, fc, dt)
	out := make([]float64, len(y))
	for i := range y {
		out[i] = y[i] - lp[i]
	}
	return out
}

// MovingAverage applies a centred window of width w samples (w forced odd) to y.
// Values near the edges use the available window.
func MovingAverage(y []float64, w int) []float64 {
	n := len(y)
	if n == 0 {
		return nil
	}
	if w < 1 {
		w = 1
	}
	if w%2 == 0 {
		w++
	}
	half := w / 2
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		lo := i - half
		if lo < 0 {
			lo = 0
		}
		hi := i + half
		if hi >= n {
			hi = n - 1
		}
		sum := 0.0
		for j := lo; j <= hi; j++ {
			sum += y[j]
		}
		out[i] = sum / float64(hi-lo+1)
	}
	return out
}

// PeakDetect returns the indices of local maxima in y (strictly greater than
// both neighbours).
func PeakDetect(y []float64) []int {
	return collectPeaks(y)
}
