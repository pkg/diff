package diff_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pkg/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var goldenTests = []struct {
	a, b string
	want string // usually from running diff --unified and cleaning up the output
	n    int
}{
	{
		a: "A\nB\nC\nD\nE\nF\n",
		b: "A\nB\nC\nD\nE\nF\n1\n2\n3\n",
		// TODO: stock macOS diff omits the trailing common blank line in this diff,
		// which also changes the @@ line ranges to be 4,3 and 4,6.
		want: `--- a
+++ b
@@ -4,4 +4,7 @@
 D
 E
 F
+1
+2
+3
 
`,
		n: 3,
	},

	{
		a: "A\nB\nC\nD\nE\nF\n",
		b: "1\n2\n3\nA\nB\nC\nD\nE\nF\n",
		want: `--- a
+++ b
@@ -1,3 +1,6 @@
+1
+2
+3
 A
 B
 C
`,
		n: 3,
	},
}

func TestGolden(t *testing.T) {
	for _, test := range goldenTests {
		as := strings.Split(test.a, "\n")
		bs := strings.Split(test.b, "\n")
		ab := diff.Strings(as, bs)
		// TODO: supply an EditScript to the tests instead doing a Myers diff here.
		// Doing it as I have done, the lazy way, mixes concerns: diff algorithm vs unification algorithm
		// vs unified diff formatting.
		e := diff.Myers(context.Background(), ab)
		e = e.WithContextSize(3)
		buf := new(bytes.Buffer)
		e.WriteUnified(buf, ab)
		got := buf.String()
		if test.want != got {
			dmp := diffmatchpatch.New()
			delta := dmp.DiffMain(test.want, got, false)
			t.Errorf("bad diff: a=%q b=%q n=%d\n\ngot:\n%s\nwant:\n%s\ndiff:\n%s\n",
				test.a, test.b, test.n,
				got, test.want,
				dmp.DiffPrettyText(delta),
			)
		}
	}
}
