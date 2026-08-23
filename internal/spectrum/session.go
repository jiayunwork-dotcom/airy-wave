package spectrum

import "io"

// csvSession wraps the destination writer and tracks whether the export
// session has already been closed.
type csvSession struct {
	w      io.Writer
	closed bool
}

func (s *csvSession) Write(p []byte) (int, error) {
	if s.closed {
		panic("write on closed csv session")
	}
	return s.w.Write(p)
}

func (s *csvSession) Close() error {
	if s.closed {
		panic("csv session already closed")
	}
	s.closed = true
	return nil
}

func newCSVSession(w io.Writer) *csvSession {
	return &csvSession{w: w}
}
