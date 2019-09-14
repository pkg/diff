package diff

import (
	"context"
	"reflect"
	"testing"
)

func TestMyers(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want []segment
	}{
		{
			name: "BasicExample",
			a:    "ABCABBA",
			b:    "CBABAC",
			want: []segment{
				{FromA: 0, ToA: 2, FromB: 0, ToB: 0},
				{FromA: 2, ToA: 3, FromB: 0, ToB: 1},
				{FromA: 3, ToA: 3, FromB: 1, ToB: 2},
				{FromA: 3, ToA: 5, FromB: 2, ToB: 4},
				{FromA: 5, ToA: 6, FromB: 4, ToB: 4},
				{FromA: 6, ToA: 7, FromB: 4, ToB: 5},
				{FromA: 7, ToA: 7, FromB: 5, ToB: 6},
			},
		},
		{
			name: "AllDifferent",
			a:    "ABCDE",
			b:    "xyz",
			want: []segment{
				{FromA: 0, ToA: 5, FromB: 0, ToB: 0},
				{FromA: 0, ToA: 0, FromB: 0, ToB: 3},
			},
		},
		// TODO: add more tests
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ab := &diffByByte{a: test.a, b: test.b}
			got := Myers(context.Background(), ab)
			want := EditScript{segs: test.want}

			if !reflect.DeepEqual(got, want) {
				// Ironically, it'd be nice to provide a diff between got and want here...
				// but our diff algorithm is busted.
				t.Errorf("got:\n%v\n\nwant:\n%v\n\n", got, want)
			}
		})
	}
}

type diffByByte struct {
	a, b string
}

func (ab *diffByByte) LenA() int             { return len(ab.a) }
func (ab *diffByByte) LenB() int             { return len(ab.b) }
func (ab *diffByByte) Equal(ai, bi int) bool { return ab.a[ai] == ab.b[bi] }
