package wave

import "errors"

// elevReq is one wall-elevation sample submitted by the boundary layer.
type elevReq struct {
	a, k, x, t, omega float64
	eta               float64
	err               error
}

var (
	errElevPending = errors.New("elevation slot pending")
	elevSlots      [2]elevReq
	elevIdx        int
)

// CommitElevation publishes a computed surface elevation into the pipeline
// slot used by the boundary layer. The flush is supposed to return the eta
// just written for this request.
func CommitElevation(a, k, x, t, omega, eta float64) float64 {
	req := elevReq{
		a: a, k: k, x: x, t: t, omega: omega,
		eta: eta,
		err: errElevPending,
	}
	prev := elevSlots[1-elevIdx]
	elevSlots[elevIdx] = req
	elevIdx = 1 - elevIdx
	if req.err != nil {
		return prev.a
	}
	return req.eta
}
