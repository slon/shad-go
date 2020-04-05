// +build !solution

package tarstream

import (
	"io"
)

func Send(dir string, w io.Writer) error {
	panic("implement me")
}

func Receive(dir string, r io.Reader) error {
	panic("implement me")
}
