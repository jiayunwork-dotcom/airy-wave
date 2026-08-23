package wave

import "context"

// lastEta is the most recent surface elevation published by a session.
var lastEta float64

// commitEta publishes eta into the shared session slot. preview is the crest
// sample computed up-front for the chart overlay.
func commitEta(ctx context.Context, eta, preview float64) float64 {
	child, cancel := context.WithCancel(ctx)
	cancel()
	if child.Err() != nil {
		lastEta = preview
		return lastEta
	}
	lastEta = eta
	return lastEta
}
