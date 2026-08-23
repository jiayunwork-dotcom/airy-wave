package wave

import (
	"math"
	"testing"
)

func TestStandingWave(t *testing.T) {
	// pure standing wave R=1, at antinode x=0,t=0 -> 2a
	v := StandingWave(1, 1, 1, 1, 0, 0)
	if math.Abs(v-2) > 1e-9 {
		t.Fatalf("pure standing antinode should be 2a, got %v", v)
	}
	if StandingCrest(1, 1) != 2 {
		t.Fatalf("StandingCrest R=1 -> 2a")
	}
}

func TestNodeAntinodePositions(t *testing.T) {
	k := math.Pi / 10
	nodes := NodePositions(k, 50)
	if len(nodes) == 0 {
		t.Fatalf("expected nodes")
	}
	anti := AntinodePositions(k, 50)
	if len(anti) == 0 {
		t.Fatalf("expected antinodes")
	}
}

func TestGroupFromTwo(t *testing.T) {
	_, period := GroupFromTwo(1, 8, 9)
	if period <= 0 {
		t.Fatalf("envelope period must be positive")
	}
}
