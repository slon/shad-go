//go:build !solution

package vegz

import (
	"compress/gzip"
	"io"
)

func Encode(data []byte, w io.Writer) error {
	ww := gzip.NewWriter(w)
	defer func() { _ = ww.Close() }()
	if _, err := ww.Write(data); err != nil {
		return err
	}
	return ww.Flush()
}
