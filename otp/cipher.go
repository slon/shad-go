//go:build !solution

package otp

import (
	"io"
)

func makeReader(r io.Reader, prng io.Reader) io.Reader {
	panic("implement me")
}

func makeWriter(w io.Writer, prng io.Reader) io.Writer {
	panic("implement me")
}
