package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var b bytes.Buffer // A Buffer needs no initialization.
	b.Write([]byte("Hello "))
	_, _ = fmt.Fprintf(&b, "world!")
	_, _ = b.WriteTo(os.Stdout)
}
