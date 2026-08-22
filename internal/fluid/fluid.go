// Package fluid builds a profile of an Airy wave over a vertical column and can
// export it as CSV for plotting or downstream use.
package fluid

import (
	"encoding/csv"
	"io"
	"strconv"

	"airy-wave/internal/wave"
)

// Sample is one snapshot of the wave field at a horizontal position.
type Sample struct {
	X          float64 // horizontal position, m
	Z          float64 // height above bed, m
	Elevation  float64 // surface elevation at (x,t), m
	U          float64 // horizontal velocity, m/s
	W          float64 // vertical velocity, m/s
	Pressure   float64 // dynamic pressure, Pa
}

// Column samples the wave at a single horizontal position x over the full depth
// using n vertical layers at time t.
func Column(p wave.Params, x, t float64, layers int, rho float64) []Sample {
	if layers < 1 {
		layers = 1
	}
	out := make([]Sample, 0, layers+1)
	h := p.Depth
	for i := 0; i <= layers; i++ {
		z := h * float64(i) / float64(layers)
		out = append(out, Sample{
			X:         x,
			Z:         z,
			Elevation: p.SurfaceElevation(x, t),
			U:         p.HorizontalVelocity(x, z, t),
			W:         p.VerticalVelocity(x, z, t),
			Pressure:  p.DynamicPressure(x, z, t, rho),
		})
	}
	return out
}

// WriteCSV writes a slice of samples to w as comma-separated values with a header.
func WriteCSV(w io.Writer, samples []Sample) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"x", "z", "elevation", "u", "w", "pressure"}); err != nil {
		return err
	}
	for _, s := range samples {
		line := []string{
			ftoa(s.X), ftoa(s.Z), ftoa(s.Elevation),
			ftoa(s.U), ftoa(s.W), ftoa(s.Pressure),
		}
		if err := cw.Write(line); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func ftoa(v float64) string {
	return strconv.FormatFloat(v, 'g', 10, 64)
}
