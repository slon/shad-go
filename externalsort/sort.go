// +build !solution

package externalsort

import (
	"io"
)

func NewReader(r io.Reader) LineReader {
	panic("implement me")
}

func NewWriter(w io.Writer) LineWriter {
	panic("implement me")
}

func Merge(w LineWriter, readers ...LineReader) error {
	panic("implement me")
}

func Sort(w io.Writer, in ...string) error {
	panic("implement me")
}
