## gzep [runtime]

В этой задаче нужно победить бенчмарк, "улучшив" функцию сжатия в `gzip`.

Пример запуска бенчмарка для бейзлайна и авторского решения:
```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/gzep
cpu: Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz
BenchmarkEncodeSimple
BenchmarkEncodeSimple-8   	    7047	    176628 ns/op	  813872 B/op	      17 allocs/op
BenchmarkEncode
BenchmarkEncode-8         	   41706	     32616 ns/op	      19 B/op	       0 allocs/op
PASS
ok  	gitlab.com/slon/shad-go/gzep	3.625s
```

### С чего начать?

Запустите бенчмарк локально. Найдите в коде `compress/gzip` откуда берутся эти сотни килобайт на одну итерацию. Подумайте какой стандартный способ избежать подобных аллокаций есть в языке.

Советуем попробовать что-нибудь написать прежде чем посмотреть ответ
```
echo "c3luYy5Qb29sCg==" | base64 -d
```
