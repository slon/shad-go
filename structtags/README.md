# structtags

Ускорьте функцию `Unpack()`.

```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/structtags
cpu: Intel(R) Core(TM) i7-6600U CPU @ 2.60GHz
BenchmarkUnpacker/user-4         	 4158832	       268.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnpacker/good-4         	 1000000	      1198 ns/op	     220 B/op	       6 allocs/op
BenchmarkUnpacker/order-4        	 1260784	      1162 ns/op	     282 B/op	       6 allocs/op
PASS
```

## Ссылки

1. sync.Map: https://golang.org/pkg/sync/#Map
2. reflect.Type: https://golang.org/pkg/reflect/#Type
