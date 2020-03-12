// +build !change

package externalsort

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriterFlusher interface {
	Write(l string) error
	Flush() error
}
