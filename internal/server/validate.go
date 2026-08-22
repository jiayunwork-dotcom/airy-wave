package server

// ValidationError describes a rejected request with a stable machine-readable
// code and a human-readable message.
type ValidationError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string { return e.Message }

// IsValidation reports whether err is a ValidationError.
func IsValidation(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// ValidateWaveRequest checks a WaveRequest for physically sane inputs.
func ValidateWaveRequest(req WaveRequest) error {
	if req.Period <= 0 {
		return &ValidationError{Code: "period", Message: "period must be positive"}
	}
	if req.Depth < 0 {
		return &ValidationError{Code: "depth", Message: "depth must be non-negative"}
	}
	if req.Layers < 0 {
		return &ValidationError{Code: "layers", Message: "layers must be non-negative"}
	}
	if req.Rho < 0 {
		return &ValidationError{Code: "rho", Message: "rho must be non-negative"}
	}
	return nil
}
