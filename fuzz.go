// +build gofuzz

package diff

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

func Fuzz(data []byte) int {
	if len(data) == 0 {
		return -1
	}
	nul := bytes.IndexByte(data, 0)
	if nul == -1 {
		nul = len(data) - 1
	}
	a := data[:nul]
	b := data[nul:]
	ab := &IndividualBytes{a: a, b: b}
	e := Myers(context.Background(), ab)
	e.WriteUnified(ioutil.Discard, ab)
	return 0
}

type IndividualBytes struct {
	a, b []byte
}

func (ab *IndividualBytes) LenA() int                                { return len(ab.a) }
func (ab *IndividualBytes) LenB() int                                { return len(ab.b) }
func (ab *IndividualBytes) Equal(ai, bi int) bool                    { return ab.a[ai] == ab.b[bi] }
func (ab *IndividualBytes) WriteATo(w io.Writer, i int) (int, error) { return w.Write([]byte{ab.a[i]}) }
func (ab *IndividualBytes) WriteBTo(w io.Writer, i int) (int, error) { return w.Write([]byte{ab.b[i]}) }
