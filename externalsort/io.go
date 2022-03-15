//go:build !change

package externalsort

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}
