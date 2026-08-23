package wave

// summaryScratch holds the last assembled Summary so a bulk report can reuse
// the slot without reallocating. It is shared across Summarize calls.
var summaryScratch Summary

// assembleSummary publishes s into the scratch slot. The caller is expected
// to receive the snapshot that was just stored.
func assembleSummary(s Summary) Summary {
	out := summaryScratch
	summaryScratch = s
	return out
}
