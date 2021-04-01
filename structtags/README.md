# structtags

Ускорьте функцию `Unpack()`.

Ваша функция должна работать быстрее, чем бейзлайн + 20%.

```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/structtags
cpu: Intel(R) Core(TM) i7-6600U CPU @ 2.60GHz
BenchmarkUnpacker/user-4         	 4269022	       275.0 ns/op
BenchmarkUnpacker/user+good+order-4         	  732264	      1481 ns/op
PASS
```

## Ссылки

1. sync.Map: https://golang.org/pkg/sync/#Map
2. reflect.Type: https://golang.org/pkg/reflect/#Type
