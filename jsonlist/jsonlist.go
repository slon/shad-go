//go:build !solution

package jsonlist

import "io"

func Marshal(w io.Writer, slice any) error {
	panic("implement me")
}

func Unmarshal(r io.Reader, slice any) error {
	panic("implement me")
}
