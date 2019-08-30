package diff

import (
	"context"
	"fmt"
)

// Myers calculates an EditScript (diff) for ab using the Myers diff algorithm.
// Because diff calculation can be expensive, Myers supports cancellation via ctx.
func Myers(ctx context.Context, ab Pair) EditScript {
	aLen := ab.LenA()
	bLen := ab.LenB()
	if aLen == 0 {
		return scriptWithSegments(segment{ToB: bLen})
	}
	if bLen == 0 {
		return scriptWithSegments(segment{ToA: aLen})
	}

	max := aLen + bLen
	if max < 0 {
		panic("overflow in diff.Myers")
	}
	v := make([]int, 2*max+1) // indices: -max .. 0 .. max

	var trace [][]int
search:
	for d := 0; d < max; d++ {
		// Only check context every 16th iteration to reduce overhead.
		if ctx != nil && uint(d)%16 == 0 && ctx.Err() != nil {
			return EditScript{}
		}

		// TODO: this seems like it will frequently be bigger than necessary.
		// Use sparse lookup? prefixes?
		vc := make([]int, 2*max+1)
		copy(vc, v)
		trace = append(trace, vc)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
				x = v[max+k+1]
			} else {
				x = v[max+k-1] + 1
			}

			y := x - k
			for x < aLen && y < bLen && ab.Equal(x, y) {
				x++
				y++
			}
			v[max+k] = x

			if x == aLen && y == bLen {
				break search
			}
		}
	}

	if len(trace) == max {
		// No commonality at all, delete everything and then insert everything.
		// This is handled as a special case to avoid complicating the logic below.
		return scriptWithSegments(segment{ToA: aLen}, segment{ToB: bLen})
	}

	// Create reversed edit script.
	x := aLen
	y := bLen
	var e EditScript
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y
		var prevk int
		if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
			prevk = k + 1
		} else {
			prevk = k - 1
		}
		prevx := v[max+prevk]
		prevy := prevx - prevk
		for x > prevx && y > prevy {
			e.appendToReversed(segment{FromA: x - 1, FromB: y - 1, ToA: x, ToB: y})
			x--
			y--
		}
		if d > 0 {
			e.appendToReversed(segment{FromA: prevx, FromB: prevy, ToA: x, ToB: y})
		}
		x, y = prevx, prevy
	}

	// Reverse reversed edit script, to return to natural order.
	e.reverse()

	// Sanity check
	for i := 1; i < len(e.segs); i++ {
		prevop := e.segs[i-1].op()
		currop := e.segs[i].op()
		if (prevop == currop) || (prevop == ins && currop != eq) || (currop == del && prevop != eq) {
			panic(fmt.Errorf("bad script: %v -> %v", prevop, currop))
		}
	}

	return e
}

func (e EditScript) reverse() {
	for i := 0; i < len(e.segs)/2; i++ {
		j := len(e.segs) - i - 1
		e.segs[i], e.segs[j] = e.segs[j], e.segs[i]
	}
}

func (e *EditScript) appendToReversed(seg segment) {
	if len(e.segs) == 0 {
		e.segs = append(e.segs, seg)
		return
	}
	u, ok := combineSegments(seg, e.segs[len(e.segs)-1])
	if !ok {
		e.segs = append(e.segs, seg)
		return
	}
	e.segs[len(e.segs)-1] = u
	return
}

// combineSegments combines s and t into a single segment if possible
// and reports whether it succeeded.
func combineSegments(s, t segment) (u segment, ok bool) {
	if t.Len() == 0 {
		return s, true
	}
	if s.Len() == 0 {
		return t, true
	}
	if s.op() != t.op() {
		return segment{FromA: -1, ToA: -1, FromB: -1, ToB: -1}, false
	}
	switch s.op() {
	case ins:
		s.ToB = t.ToB
	case del:
		s.ToA = t.ToA
	case eq:
		s.ToA = t.ToA
		s.ToB = t.ToB
	default:
		panic("bad op")
	}
	return s, true
}
