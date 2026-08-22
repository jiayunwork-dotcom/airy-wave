package wave

import (
	"testing"
)

func TestDepthRegime(t *testing.T) {
	deep := Params{Amplitude: 1, Period: 8, Depth: 1000}
	if deep.DepthRegime() != "deep" {
		t.Fatalf("deep regime expected, got %v", deep.DepthRegime())
	}
	shallow := Params{Amplitude: 1, Period: 8, Depth: 1}
	if shallow.DepthRegime() != "shallow" {
		t.Fatalf("shallow regime expected, got %v", shallow.DepthRegime())
	}
}

func TestCompareAcrossDepths(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	pts := p.CompareAcrossDepths(1, 100, 20)
	if len(pts) != 20 {
		t.Fatalf("expected 20 points, got %v", len(pts))
	}
	for _, pt := range pts {
		if pt.LambdaRatio <= 0 || pt.SpeedRatio <= 0 {
			t.Fatalf("ratios must be positive: %+v", pt)
		}
	}
}

func TestExtremeWavelength(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	minL, maxL := p.ExtremeWavelength(1, 100, 20)
	if minL <= 0 || maxL <= minL {
		t.Fatalf("extreme wavelengths invalid: min=%v max=%v", minL, maxL)
	}
}

func TestRankedRegimes(t *testing.T) {
	p := Params{Amplitude: 1, Period: 8, Depth: 50}
	regs := p.RankedRegimes(1, 100, 30)
	if len(regs) == 0 {
		t.Fatalf("no regimes returned")
	}
}

func TestSortDepthPoints(t *testing.T) {
	pts := []DepthPoint{{Depth: 3}, {Depth: 1}, {Depth: 2}}
	SortDepthPoints(pts)
	if pts[0].Depth != 1 {
		t.Fatalf("sort failed: %v", pts)
	}
}
