// +build !solution

package externalsort

import (
	"io"
)

func NewReader(r io.Reader) LineReader {
	panic("implement me")
}

func NewWriterFlusher(w io.Writer) LineWriterFlusher {
	panic("implement me")
}

func Merge(w LineWriterFlusher, readers ...LineReader) error {
	panic("implement me")
}

func Sort(w io.Writer, in ...string) error {
	panic("implement me")
}
