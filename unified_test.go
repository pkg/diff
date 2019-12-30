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
	name string
	a, b string
	opts []diff.WriteOpt
	want string // usually from running diff --unified and cleaning up the output
}{
	{
		name: "AddedLinesEnd",
		a:    "A\nB\nC\nD\nE\nF\n",
		b:    "A\nB\nC\nD\nE\nF\n1\n2\n3\n",
		// TODO: stock macOS diff omits the trailing common blank line in this diff,
		// which also changes the @@ line ranges to be 4,3 and 4,6.
		want: `
--- a
+++ b
@@ -4,4 +4,7 @@
 D
 E
 F
+1
+2
+3
 
`[1:],
	},

	{
		name: "AddedLinesStart",
		a:    "A\nB\nC\nD\nE\nF\n",
		b:    "1\n2\n3\nA\nB\nC\nD\nE\nF\n",
		want: `
--- a
+++ b
@@ -1,3 +1,6 @@
+1
+2
+3
 A
 B
 C
`[1:],
	},

	{
		name: "WithTerminalColor",
		a:    "1\n2\n2",
		b:    "1\n3\n3",
		opts: []diff.WriteOpt{diff.TerminalColor()},
		want: `
`[1:] + "\u001b[1m" + `--- a
+++ b
` + "\u001b[0m" + "\u001b[36m" + `@@ -1,3 +1,3 @@
` + "\u001b[0m" + ` 1
` + "\u001b[31m" + `-2
-2
` + "\u001b[32m" + `+3
+3
` + "\u001b[0m",
	},
}

func TestGolden(t *testing.T) {
	for _, test := range goldenTests {
		t.Run(test.name, func(t *testing.T) {
			as := strings.Split(test.a, "\n")
			bs := strings.Split(test.b, "\n")
			ab := diff.Strings(as, bs)
			// TODO: supply an edit.Script to the tests instead doing a Myers diff here.
			// Doing it as I have done, the lazy way, mixes concerns: diff algorithm vs unification algorithm
			// vs unified diff formatting.
			e := diff.Myers(context.Background(), ab)
			e = diff.EditScriptWithContextSize(e, 3)
			buf := new(bytes.Buffer)
			diff.WriteUnified(e, buf, ab, test.opts...)
			got := buf.String()
			if test.want != got {
				t.Logf("%q\n", test.want)
				t.Logf("%q\n", got)
				dmp := diffmatchpatch.New()
				delta := dmp.DiffMain(test.want, got, false)
				t.Errorf("bad diff: a=%q b=%q\n\ngot:\n%s\nwant:\n%s\ndiff:\n%s\n",
					test.a, test.b,
					got, test.want,
					dmp.DiffPrettyText(delta),
				)
			}
		})
	}
}
