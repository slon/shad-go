# structtags

Ускорьте функцию `Unpack()`.

Ваша функция должна работать быстрее, чем бейзлайн + 20%.

```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/structtags
BenchmarkUnpacker/user-4                    3064            362500 ns/op
BenchmarkUnpacker/user+good+order-4                  663           1799294 ns/op
PASS
```

## Ссылки

1. sync.Map: https://golang.org/pkg/sync/#Map
2. reflect.Type: https://golang.org/pkg/reflect/#Type
