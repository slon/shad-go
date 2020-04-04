# tarstream

Вам нужно уметь передавать директорию с артефактами между воркерами. Для этого, вам нужно
реализовать две операции:

```go
package tarstream

import "io"

// Send рекурсивно обходит директорию и сериализует её содержимое в поток w.
func Send(dir string, w io.Writer) error

// Receive читает поток r и материализует содержимое потока внутри dir.
func Receive(dir string, r io.Reader) error
```

- Функции должны корректно обрабатывать директории и обычные файлы.
- executable бит на файлах должен сохраняться.
- Используйте формат [tar](https://golang.org/pkg/archive/tar/)
- Используйте [filepath.Walk](https://golang.org/pkg/path/filepath/) для рекурсивного обхода.
