package diff_test

import (
	"os"

	"github.com/pkg/diff"
)

func Example_Slices() {
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := []int{1, 2, 3, 4, 6, 7, 8, 9}
	err := diff.Slices("want", "got", want, got, os.Stdout)
	if err != nil {
		panic(err)
	}
	// Output:
	// --- want
	// +++ got
	// @@ -2,7 +2,6 @@
	//  2
	//  3
	//  4
	// -5
	//  6
	//  7
	//  8
}

func Example_Text() {
	a := `
a
b
c
`[1:]
	b := `
a
c
d
`[1:]
	err := diff.Text("a", "b", a, b, os.Stdout)
	if err != nil {
		panic(err)
	}
	// Output:
	// --- a
	// +++ b
	// @@ -1,3 +1,3 @@
	//  a
	// -b
	//  c
	// +d
}
