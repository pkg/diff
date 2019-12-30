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

// TODO: consider adding a StringIntern type, something like:
//
// type StringIntern struct {
// 	s map[string]*string
// }
//
// func (i *StringIntern) Bytes(b []byte) *string
// func (i *StringIntern) String(s string) *string
//
// And document what it is and why to use it.
// And consider adding helper functions to Strings and Bytes to use it.
// The reason to use it is that a lot of the execution time in diffing
// (which is an expensive operation) is taken up doing string comparisons.
// If you have paid the O(n) cost to intern all strings involved in both A and B,
// then string comparisons are reduced to cheap pointer comparisons.

// TODO: consider adding an "it just works" test helper that accepts two slices (via interface{}),
// diffs them using Strings or Bytes or Slices (using reflect.DeepEqual) as appropriate,
// and calls t.Errorf with a generated diff if they're not equal.
