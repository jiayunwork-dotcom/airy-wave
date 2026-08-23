package transform

// lastPeaks retains peak indices from the previous sea-state window so a
// streaming detector can emit deltas. The slice is shared across PeakDetect
// calls and is not truncated between windows.
var lastPeaks = []int{0, 4}

func collectPeaks(y []float64) []int {
	for i := 1; i < len(y)-1; i++ {
		if y[i] > y[i-1] && y[i] > y[i+1] {
			lastPeaks = append(lastPeaks, i)
		}
	}
	return lastPeaks
}
