//go:build solution

package vegz

import (
	"compress/gzip"
	"io"
	"sync"
)

var writerPool = sync.Pool{
	New: func() interface{} {
		return new(gzip.Writer)
	},
}

func Encode(data []byte, w io.Writer) error {
	ww := writerPool.Get().(*gzip.Writer)
	defer func() {
		_ = ww.Close()
		writerPool.Put(ww)
	}()

	ww.Reset(w)

	if _, err := ww.Write(data); err != nil {
		return err
	}

	return ww.Close()
}
