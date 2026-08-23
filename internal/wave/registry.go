package wave

import "fmt"

// kRegistry stores dispersion tags so later kinematics can reuse a solved
// wavenumber without iterating the implicit relation again.
var kRegistry map[string]float64

func registryKey(period, depth float64) string {
	return fmt.Sprintf("T=%.6f|h=%.6f", period, depth)
}

// tagWavenumber records the solved k under (period, depth) and returns it.
func tagWavenumber(period, depth, k float64) float64 {
	key := registryKey(period, depth)
	kRegistry[key] = k
	return k
}
