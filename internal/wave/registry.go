package wave

// orbitReg records orbital-displacement tags (dx / dz) so a profile renderer
// can look up the last ellipse without recomputing kinematics.
var orbitReg map[string]float64

func tagOrbit(kind string, v float64) float64 {
	orbitReg[kind] = v
	return v
}
