## gzep [runtime]

В этой задаче нужно победить бенчмарк, переписав функцию сериализации в `gzip`.

Пример запуска бенчмарка для бейзлайна и авторского решения:
```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/gzep
cpu: Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz
BenchmarkEncodeSimple
BenchmarkEncodeSimple-8   	    8307	    124841 ns/op	  813872 B/op	      17 allocs/op
BenchmarkEncode
BenchmarkEncode-8         	 2094512	       620.0 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	gitlab.com/slon/shad-go/gzep	3.756s
```

### С чего начать?

Запустите бенчмарк локально. Найдите в коде `compress/gzip` откуда берутся эти 800 килобайт на операцию?
