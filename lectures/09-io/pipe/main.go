package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	r, w := io.Pipe()

	go func() {
		_, _ = fmt.Fprint(w, "some text to be read\n")
		_ = w.Close()
	}()

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r)
	fmt.Print(buf.String())
}
