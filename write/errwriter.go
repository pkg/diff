package write

import "io"

func newErrWriter(w io.Writer) *errwriter {
	return &errwriter{w: w}
}

// An errwriter wraps a writer.
// As soon as one write fails, it consumes all subsequent writes.
// This reduces the amount of error-checking required
// in write-heavy code.
type errwriter struct {
	w         io.Writer
	err       error
	wrote     int
	attempted int
}

func (w *errwriter) Write(b []byte) (int, error) {
	w.attempted += len(b)
	if w.err != nil {
		return 0, w.err // TODO: use something like errors.Wrap(w.err)?
	}
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	w.wrote += n
	return n, err
}

func (w *errwriter) WriteString(s string) {
	// TODO: use w.w's WriteString method, if it exists
	w.Write([]byte(s))
}

func (w *errwriter) WriteByte(b byte) {
	// TODO: use w.w's WriteByte method, if it exists
	w.Write([]byte{b})
}

func (w *errwriter) Error() error {
	return w.err
}
