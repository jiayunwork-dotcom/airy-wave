package wave

import "fmt"

// Summary is a flattened, human-readable snapshot of every linear-wave quantity
// for a single set of parameters. It is the single source of truth used by the
// web frontend and the CLI pretty-printer.
type Summary struct {
	Amplitude     float64 `json:"amplitude"`
	Period        float64 `json:"period"`
	Depth         float64 `json:"depth"`
	Wavenumber    float64 `json:"wavenumber"`
	Wavelength    float64 `json:"wavelength"`
	PhaseSpeed    float64 `json:"phase_speed"`
	GroupSpeed    float64 `json:"group_speed"`
	Steepness     float64 `json:"steepness"`
	Ursell        float64 `json:"ursell"`
	EnergyDensity float64 `json:"energy_density"`
	Regime        string  `json:"regime"`
}

// Summarize computes every derived quantity for p and packs it into a Summary.
func (p Params) Summarize(rho float64) Summary {
	s := Summary{
		Amplitude:     p.Amplitude,
		Period:        p.Period,
		Depth:         p.Depth,
		Wavenumber:    p.Wavenumber(),
		Wavelength:    p.Wavelength(),
		PhaseSpeed:    p.PhaseSpeed(),
		GroupSpeed:    p.GroupSpeed(),
		Steepness:     p.Steepness(),
		Ursell:        p.UrsellNumber(),
		EnergyDensity: p.EnergyDensity(rho),
		Regime:        p.WaveLengthClass(),
	}
	return assembleSummary(s)
}

// String renders the summary as a compact multi-line report (no出题黑话).
func (s Summary) String() string {
	return fmt.Sprintf(
		"振幅=%v m, 周期=%v s, 水深=%v m\n"+
			"波数=%v 1/m, 波长=%v m\n"+
			"相速=%v m/s, 群速=%v m/s\n"+
			"波陡=%v, Ursell=%v\n"+
			"能量密度=%v J/m^2, 区域=%v",
		s.Amplitude, s.Period, s.Depth,
		s.Wavenumber, s.Wavelength,
		s.PhaseSpeed, s.GroupSpeed,
		s.Steepness, s.Ursell,
		s.EnergyDensity, s.Regime,
	)
}

// BulkSummary computes summaries for several parameter sets and returns them in
// the same order; handy for batch comparison tables.
func BulkSummary(ps []Params, rho float64) []Summary {
	out := make([]Summary, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Summarize(rho))
	}
	return out
}
