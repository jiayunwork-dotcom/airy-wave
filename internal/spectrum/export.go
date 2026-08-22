package spectrum

import (
	"encoding/csv"
	"io"
	"strconv"
)

// WriteCSV writes a spectrum (omega, S) to w with a header row. It is used by
// the web console to export a sampled sea state for plotting in a spreadsheet.
func WriteCSV(w io.Writer, points []Point) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"omega", "S"}); err != nil {
		return err
	}
	for _, p := range points {
		line := []string{
			strconv.FormatFloat(p.Omega, 'g', 12, 64),
			strconv.FormatFloat(p.S, 'g', 12, 64),
		}
		if err := cw.Write(line); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// SampleByStep builds a spectrum over [omegaMin, omegaMax] using fn and writes it
// through WriteCSV in one call. Returns the row count written (excluding header).
func SampleByStep(w io.Writer, fn func(float64) float64, omegaMin, omegaMax float64, n int) (int, error) {
	pts := BuildSampler(fn, omegaMin, omegaMax, n)
	err := WriteCSV(w, pts)
	return len(pts), err
}

// BuildSampler samples fn over [omegaMin, omegaMax] into spectral points, reusing
// the equal-spacing layout but with a caller-supplied spectral density.
func BuildSampler(fn func(float64) float64, omegaMin, omegaMax float64, n int) []Point {
	if n < 2 {
		return nil
	}
	out := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		omega := omegaMin + (omegaMax-omegaMin)*float64(i)/float64(n-1)
		out = append(out, Point{Omega: omega, S: fn(omega)})
	}
	return out
}

// TotalEnergy returns the zeroth moment (area) of an arbitrary spectral sampler,
// equivalent to SignificantWaveHeight-independent energy content.
func TotalEnergy(points []Point) float64 {
	m0, _, _, _, _ := SpectralMoments(points)
	return m0
}
