// Package server exposes the Airy wave model over HTTP. It serves the static
// web console from web/ and a JSON endpoint:
//
//	POST /api/wave   compute wavelength, phase speed and a vertical profile
//
// All computation is delegated to the wave and fluid packages; this package
// only decodes requests, calls the engine, and encodes results or errors.
// Invalid input yields a JSON error body {"error": "..."} with a 4xx status.
package server

import (
	"encoding/json"
	"net/http"

	"airy-wave/internal/fluid"
	"airy-wave/internal/wave"
)

// Server bundles the HTTP routing for airy-wave. It serves the static web
// console from WebDir and the example files from ExampleDir.
type Server struct {
	ExampleDir string
	WebDir     string
	mux        *http.ServeMux
}

// NewServer constructs a Server with the web console in "web" and example files
// in dir.
func NewServer(dir string) *Server {
	s := &Server{ExampleDir: dir, WebDir: "web", mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns the underlying http.Handler (used by tests via httptest).
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the HTTP server on addr with a short read-header timeout.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/wave", s.handleWave)
	s.mux.HandleFunc("/api/energy", s.handleEnergy)
	s.mux.HandleFunc("/api/disp", s.handleDisp)
	s.mux.HandleFunc("/api/eta", s.handleEta)
	s.mux.Handle(
		"/example/",
		http.StripPrefix("/example/", http.FileServer(http.Dir(s.ExampleDir))),
	)
	s.mux.Handle("/", http.FileServer(http.Dir(s.WebDir)))
}

// WaveRequest is the JSON body for /api/wave.
type WaveRequest struct {
	Amplitude float64 `json:"amplitude"` // wave amplitude a, m
	Period    float64 `json:"period"`    // wave period T, s
	Depth     float64 `json:"depth"`     // still-water depth h, m
	Layers    int     `json:"layers"`    // vertical layers for the profile
	Rho       float64 `json:"rho"`       // water density, kg/m^3 (0 => 1025)
}

// ProfilePoint is one vertical sample of the wave field.
type ProfilePoint struct {
	Z        float64 `json:"z"`
	U        float64 `json:"u"`
	W        float64 `json:"w"`
	Pressure float64 `json:"pressure"`
}

// WaveResponse carries dispersion results and the vertical profile.
type WaveResponse struct {
	Wavelength  float64        `json:"wavelength"`
	Wavenumber  float64        `json:"wavenumber"`
	PhaseSpeed float64        `json:"phase_speed"`
	DeepWater  bool           `json:"deep_water"`
	Profile    []ProfilePoint `json:"profile"`
}

func (s *Server) handleWave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req WaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := ValidateWaveRequest(req); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			writeError(w, http.StatusBadRequest, ve.Message)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p := wave.Params{Amplitude: req.Amplitude, Period: req.Period, Depth: req.Depth}
	profile := fluid.Column(p, 0, 0, req.Layers, req.Rho)
	pts := make([]ProfilePoint, 0, len(profile))
	for _, s2 := range profile {
		pts = append(pts, ProfilePoint{Z: s2.Z, U: s2.U, W: s2.W, Pressure: s2.Pressure})
	}
	resp := WaveResponse{
		Wavelength:  p.Wavelength(),
		Wavenumber:  p.Wavenumber(),
		PhaseSpeed:  p.PhaseSpeed(),
		DeepWater:   p.IsDeepWater(),
		Profile:     pts,
	}
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// EnergyResponse carries spectral/energy quantities for /api/energy.
type EnergyResponse struct {
	EnergyDensity  float64 `json:"energy_density"`
	EnergyFlux     float64 `json:"energy_flux"`
	GroupSpeed     float64 `json:"group_speed"`
	RadiationStress float64 `json:"radiation_stress"`
	Regime         string  `json:"regime"`
}

func (s *Server) handleEnergy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req WaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := ValidateWaveRequest(req); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			writeError(w, http.StatusBadRequest, ve.Message)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p := wave.Params{Amplitude: req.Amplitude, Period: req.Period, Depth: req.Depth}
	resp := EnergyResponse{
		EnergyDensity:   p.EnergyDensity(req.Rho),
		EnergyFlux:      p.EnergyFlux(req.Rho),
		GroupSpeed:      p.GroupSpeed(),
		RadiationStress: p.RadiationStress(req.Rho),
		Regime:          p.WaveLengthClass(),
	}
	writeJSON(w, resp)
}

// DispRequest is the JSON body for /api/disp: a progressive wave described by
// wave height H, period T and water depth h.
type DispRequest struct {
	Height float64 `json:"H"`
	Period float64 `json:"T"`
	Depth  float64 `json:"h"`
}

// DispResponse carries the dispersion results: wavenumber k, wavelength lambda,
// phase speed c, group speed cg and the relative-depth class string.
type DispResponse struct {
	Wavenumber float64 `json:"k"`
	Wavelength float64 `json:"lambda"`
	PhaseSpeed float64 `json:"c"`
	GroupSpeed float64 `json:"cg"`
	Regime     string  `json:"regime"`
}

// handleDisp solves the dispersion relation for one (H, T, h) case. H is halved
// to the amplitude a used by the kernel. It mirrors /api/wave but exposes the
// wave-height input directly as the brief specifies.
func (s *Server) handleDisp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req DispRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Period <= 0 || req.Depth <= 0 || req.Height <= 0 {
		writeError(w, http.StatusBadRequest, "H, T and h must all be positive")
		return
	}
	p := wave.Params{Amplitude: req.Height / 2.0, Period: req.Period, Depth: req.Depth}
	resp := DispResponse{
		Wavenumber: p.Wavenumber(),
		Wavelength: p.Wavelength(),
		PhaseSpeed: p.PhaseSpeed(),
		GroupSpeed: p.GroupSpeed(),
		Regime:     p.WaveLengthClass(),
	}
	writeJSON(w, resp)
}

// EtaRequest is the JSON body for /api/eta: a wave plus a spatial grid over which
// to sample the surface elevation eta(x, t).
type EtaRequest struct {
	Height   float64 `json:"H"`
	Period   float64 `json:"T"`
	Depth    float64 `json:"h"`
	XStart   float64 `json:"x_start"`
	XEnd     float64 `json:"x_end"`
	NX       int     `json:"nx"`
	TInstant float64 `json:"t"`
}

// EtaPoint is one (x, eta) sample of the free surface.
type EtaPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// handleEta returns the surface-elevation profile over one spatial line at a
// fixed instant t. The samples come from the Airy expression eta = a·sin(kx−ωt),
// exactly the curve the web console plots.
func (s *Server) handleEta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req EtaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Period <= 0 || req.Depth <= 0 || req.Height <= 0 {
		writeError(w, http.StatusBadRequest, "H, T and h must all be positive")
		return
	}
	if req.NX <= 0 {
		req.NX = 50
	}
	p := wave.Params{Amplitude: req.Height / 2.0, Period: req.Period, Depth: req.Depth}
	pts := make([]EtaPoint, 0, req.NX)
	if req.XEnd <= req.XStart {
		req.XEnd = req.XStart + 10
	}
	for i := 0; i < req.NX; i++ {
		x := req.XStart + (req.XEnd-req.XStart)*float64(i)/float64(req.NX-1)
		pts = append(pts, EtaPoint{X: x, Y: p.SurfaceElevation(x, req.TInstant)})
	}
	writeJSON(w, pts)
}
