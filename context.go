package diff

import "fmt"

// WithContextSize returns an edit script preserving only n common elements of context for changes.
// The returned edit script may alias the input.
// If n is negative, WithContextSize panics.
// To generate a "unified diff", use WithContextSize and then WriteUnified the resulting edit script.
func (e EditScript) WithContextSize(n int) EditScript {
	if n < 0 {
		panic(fmt.Sprintf("EditScript.WithContextSize called with negative n: %d", n))
	}

	// Handle small scripts.
	switch len(e.segs) {
	case 0:
		return EditScript{}
	case 1:
		if e.segs[0].op() == eq {
			// Entirely identical contents.
			// Unclear what to do here. For now, just bail.
			// TODO: something else? what does command line diff do?
			return EditScript{}
		}
		return scriptWithSegments(e.segs[0])
	}

	out := make([]segment, 0, len(e.segs))
	for i, seg := range e.segs {
		if seg.op() != eq {
			out = append(out, seg)
			continue
		}
		if i == 0 {
			// Leading segment. Keep only the final n entries.
			if seg.Len() > n {
				seg = segmentLastN(seg, n)
			}
			out = append(out, seg)
			continue
		}
		if i == len(e.segs)-1 {
			// Trailing segment. Keep only the first n entries.
			if seg.Len() > n {
				seg = segmentFirstN(seg, n)
			}
			out = append(out, seg)
			continue
		}
		if seg.Len() <= n*2 {
			// Small middle segment. Keep unchanged.
			out = append(out, seg)
			continue
		}
		// Large middle segment. Break into two disjoint parts.
		out = append(out, segmentFirstN(seg, n), segmentLastN(seg, n))
	}

	// TODO: Stock macOS diff also trims common blank lines
	// from the beginning/end of eq segments.
	// Perhaps we should do that here too.
	// Or perhaps that should be a separate, composable EditScript method?
	return EditScript{segs: out}
}

func segmentFirstN(seg segment, n int) segment {
	if seg.op() != eq {
		panic("segmentFirstN bad op")
	}
	if seg.Len() < n {
		panic("segmentFirstN bad Len")
	}
	return segment{
		FromA: seg.FromA, ToA: seg.FromA + n,
		FromB: seg.FromB, ToB: seg.FromB + n,
	}
}

func segmentLastN(seg segment, n int) segment {
	if seg.op() != eq {
		panic("segmentLastN bad op")
	}
	if seg.Len() < n {
		panic("segmentLastN bad Len")
	}
	return segment{
		FromA: seg.ToA - n, ToA: seg.ToA,
		FromB: seg.ToB - n, ToB: seg.ToB,
	}
}
