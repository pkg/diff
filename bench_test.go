package diff_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

const regenerate = false // set to true to overwrite .out files

func BenchmarkGolden(b *testing.B) {
	aa, err := filepath.Glob("testdata/*.a")
	if err != nil {
		b.Fatal(err)
	}
	for _, aPath := range aa {
		base := strings.TrimSuffix(aPath, ".a")
		bPath := base + ".b"
		outPath := base + ".out"
		out, err := ioutil.ReadFile(outPath)
		if err != nil && !regenerate {
			b.Fatal(err)
		}
		buf := new(bytes.Buffer)
		buf.Grow(len(out))
		b.Run(base, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf.Reset()
				err := diff.Text(aPath, bPath, nil, nil, buf)
				if err != nil {
					b.Fatal(err)
				}
				if regenerate {
					err := ioutil.WriteFile(outPath, buf.Bytes(), 0644)
					if err != nil {
						b.Fatal(err)
					}
					return
				}
				if !bytes.Equal(buf.Bytes(), out) {
					b.Fatal("wrong output")
				}
			}
		})
	}
}
