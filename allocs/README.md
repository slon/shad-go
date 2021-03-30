# allocs

`Counter` используется для нахождения уникальных слов и подсчета вхождений каждого из них.
Его интерфейс выглядит так:
 
* `Count(r io.Reader) error` — функция, которая подсчитывает количество вхождений для каждого слова в тексте.
На вход подается io.Reader, в котором находится некоторый текст.
Разделителями являются только переносы строк и пробелы.
* `String() string` — преобразует мапу вида `{"слово": "количество вхождений"}` в форматированную строку.

Необходимо написать имплементацию `EnhancedCounter` (см. файл `allocs.go`)
и снизить количество аллокаций. Бейзлайн можно найти в `baseline.go`.
 
Значения бенчмарков для бейзлайна: 
```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/allocs
Benchmark/count-4                  73200             16294 ns/op             880 B/op          5 allocs/op
Benchmark/main-4                   40485             30113 ns/op            1034 B/op          9 allocs/op
```

Значения бенчмарков для авторского решения:
```goos: linux
   goarch: amd64
   pkg: gitlab.com/slon/shad-go/allocs
   Benchmark/count-4                 212850              5471 ns/op            4144 B/op          2 allocs/op
   Benchmark/main-4                  143937              8247 ns/op            4176 B/op          3 allocs/op
```

 
