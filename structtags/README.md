# structtags

Ускорьте функцию `Unpack()`, про которую рассказывали на лекции (https://p.go.manytask.org/08-reflect/lecture.slide#19).

Ваша функция должна работать быстрее, чем бейзлайн + 20%.
```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/structtags
BenchmarkUnpacker/user-4                    3273            329346 ns/op
BenchmarkUnpacker/user+good+order-4                  648           1721068 ns/op
PASS
```
