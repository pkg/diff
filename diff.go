package diff

import (
	"io"

	"github.com/pkg/diff/myers"
)

// A WriterTo type supports writing a diff, element by element.
// A is the initial state; B is the final state.
type WriterTo interface {
	// WriteATo writes the element a[ai] to w.
	WriteATo(w io.Writer, ai int) (int, error)
	// WriteBTo writes the element b[bi] to w.
	WriteBTo(w io.Writer, bi int) (int, error)
}

// PairWriterTo is the union of Pair and WriterTo.
type PairWriterTo interface {
	myers.Pair
	WriterTo
}
