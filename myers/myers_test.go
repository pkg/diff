package myers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pkg/diff/edit"
	"github.com/pkg/diff/myers"
)

func TestMyers(t *testing.T) {
	tests := []struct {
		name        string
		a, b        string
		want        []edit.Range
		wantStatIns int
		wantStatDel int
	}{
		{
			name: "BasicExample",
			a:    "ABCABBA",
			b:    "CBABAC",
			want: []edit.Range{
				{LowA: 0, HighA: 2, LowB: 0, HighB: 0},
				{LowA: 2, HighA: 3, LowB: 0, HighB: 1},
				{LowA: 3, HighA: 3, LowB: 1, HighB: 2},
				{LowA: 3, HighA: 5, LowB: 2, HighB: 4},
				{LowA: 5, HighA: 6, LowB: 4, HighB: 4},
				{LowA: 6, HighA: 7, LowB: 4, HighB: 5},
				{LowA: 7, HighA: 7, LowB: 5, HighB: 6},
			},
			wantStatIns: 2,
			wantStatDel: 3,
		},
		{
			name: "AllDifferent",
			a:    "ABCDE",
			b:    "xyz",
			want: []edit.Range{
				{LowA: 0, HighA: 5, LowB: 0, HighB: 0},
				{LowA: 0, HighA: 0, LowB: 0, HighB: 3},
			},
			wantStatIns: 3,
			wantStatDel: 5,
		},
		// TODO: add more tests
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ab := &diffByByte{a: test.a, b: test.b}
			got := myers.Diff(context.Background(), ab)
			want := edit.Script{Ranges: test.want}

			if !reflect.DeepEqual(got, want) {
				// Ironically, it'd be nice to provide a diff between got and want here...
				// but our diff algorithm is busted.
				t.Errorf("got:\n%v\n\nwant:\n%v\n\n", got, want)
			}
			ins, del := got.Stat()
			if ins != test.wantStatIns {
				t.Errorf("got %d insertions, want %d", ins, test.wantStatIns)
			}
			if del != test.wantStatDel {
				t.Errorf("got %d deletions, want %d", del, test.wantStatDel)
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
