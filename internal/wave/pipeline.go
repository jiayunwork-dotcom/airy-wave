package wave

import "errors"

// standReq is one standing-wave evaluation flowing through the pipeline.
type standReq struct {
	a, k, omega, r, x, t float64
	incident             float64
	reflected            float64
	err                  error
}

var (
	errStandPending = errors.New("standing slot pending")
	standSlots      [2]standReq
	standCur        int
)

func flushStanding(req standReq) float64 {
	standSlots[standCur] = req
	standCur = 1 - standCur
	return req.incident + req.reflected
}
