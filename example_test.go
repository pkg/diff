package diff_test

import (
	"context"
	"os"

	"github.com/pkg/diff"
	"github.com/pkg/diff/ctxt"
	"github.com/pkg/diff/myers"
)

// TODO: use a less heavyweight output format for Example_testHelper

func Example_testHelper() {
	want := []int{1, 2, 3, 4, 5}
	got := []int{1, 2, 4, 5}
	ab := diff.Slices(want, got, nil)
	e := myers.Diff(context.Background(), ab)
	if e.IsIdentity() {
		return
	}
	e = ctxt.Size(e, 1)
	diff.WriteUnified(e, os.Stdout, ab)
	// Output:
	// --- a
	// +++ b
	// @@ -2,3 +2,2 @@
	//  2
	// -3
	//  4
}

func Example_strings() {
	a := []string{"a", "b", "c"}
	b := []string{"a", "c", "d"}
	ab := diff.Strings(a, b)
	e := myers.Diff(context.Background(), ab)
	diff.WriteUnified(e, os.Stdout, ab)
	// Output:
	// --- a
	// +++ b
	// @@ -1,3 +1,3 @@
	//  a
	// -b
	//  c
	// +d
}

func Example_Names() {
	a := []string{"a", "b", "c"}
	b := []string{"a", "c", "d"}
	ab := diff.Strings(a, b)
	e := myers.Diff(context.Background(), ab)
	diff.WriteUnified(e, os.Stdout, ab, diff.Names("before", "after"))
	// Output:
	// --- before
	// +++ after
	// @@ -1,3 +1,3 @@
	//  a
	// -b
	//  c
	// +d
}
