package wave

import (
	"strings"
	"testing"
)

func TestSummarize(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	s := p.Summarize(1025)
	if s.Wavelength <= 0 || s.PhaseSpeed <= 0 || s.GroupSpeed <= 0 {
		t.Fatalf("summary quantities must be positive: %+v", s)
	}
	if s.Regime == "" {
		t.Fatalf("regime empty")
	}
	if !strings.Contains(s.String(), "波长") {
		t.Fatalf("String should carry CN labels")
	}
}

func TestBulkSummary(t *testing.T) {
	ps := []Params{
		{Amplitude: 1, Period: 8, Depth: 50},
		{Amplitude: 2, Period: 10, Depth: 30},
	}
	out := BulkSummary(ps, 1025)
	if len(out) != 2 {
		t.Fatalf("BulkSummary length = %d", len(out))
	}
}
