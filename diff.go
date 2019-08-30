package diff

import (
	"bytes"
	"fmt"
	"io"
)

// A Pair is two things that can be diffed using the Myers diff algorithm.
// A is the initial state; B is the final state.
type Pair interface {
	// LenA returns the number of initial elements.
	LenA() int
	// LenA returns the number of final elements.
	LenB() int
	// Equal reports whether the ai'th element of A is equal to the bi'th element of B.
	Equal(ai, bi int) bool
}

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
	Pair
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

// An op is a edit operation used to transform A into B.
type op int8

//go:generate stringer -type op

const (
	del op = -1
	eq  op = 0
	ins op = 1
)

// A segment is a set of steps of the same op.
type segment struct {
	FromA, ToA int // Beginning and ending indices into A of this operation
	FromB, ToB int // ditto, for B
}

func (s segment) op() op {
	if s.FromA == s.ToA {
		return ins
	}
	if s.FromB == s.ToB {
		return del
	}
	return eq
}

func (s segment) String() string {
	// This output is helpful when hacking on a Myers diff.
	// In other contexts it is usually more natural to group FromA, ToA and FromB, ToB.
	return fmt.Sprintf("(%d, %d) -- %s %d --> (%d, %d)", s.FromA, s.FromB, s.op(), s.Len(), s.ToA, s.ToB)
}

func (s segment) Len() int {
	if s.FromA == s.ToA {
		return s.ToB - s.FromB
	}
	return s.ToA - s.FromA
}

// An EditScript is an edit script to alter A into B.
type EditScript struct {
	segs []segment
}

// IsIdentity reports whether e is the identity edit script, that is, whether A and B are identical.
// See the TestHelper example.
func (e EditScript) IsIdentity() bool {
	for _, seg := range e.segs {
		if seg.op() != eq {
			return false
		}
	}
	return true
}

// TODO: consider adding an "it just works" test helper that accepts two slices (via interface{}),
// diffs them using Strings or Bytes or Slices (using reflect.DeepEqual) as appropriate,
// and calls t.Errorf with a generated diff if they're not equal.

// scriptWithSegments returns an EditScript containing s.
// It is used to reduce line noise.
func scriptWithSegments(s ...segment) EditScript {
	return EditScript{segs: s}
}

// dump formats s for debugging.
func (e EditScript) dump() string {
	buf := new(bytes.Buffer)
	for _, seg := range e.segs {
		fmt.Fprintln(buf, seg)
	}
	return buf.String()
}
