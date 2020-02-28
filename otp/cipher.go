// +build !solution

package otp

import (
	"io"
)

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	panic("implement me")
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	panic("implement me")
}
