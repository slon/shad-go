# varfmt

Реализуйте функцию `varfmt.Sprintf`. Функция принимает формат строку и переменное число аргументов.

Синтаксис формат-строки похож на формат-строки питона:
 - `{}` - задаёт ссылку на аргумент
 - `{number}` - ссылается на агрумент с индексом `number`
 - `{}` ссылается на аргумент с индексом равным позиции `{}` внутри паттерна

Например, `varfmt.Sprintf("{1} {0}", "Hello", "World)` должен вернуть строку `World Hello`.

Аргументы функции могут быть произвольными типами. Вам нужно форматировать их так же, как это
делает функция `fmt.Sprint`. Вызывать `fmt.Sprint` для форматирования отдельного аргумента
не запрещается.

Ваше решение будет сравниваться с baseline-решением на бенчмарке. Ваш код должен
быть не более чем в два раза хуже чем baseline.

```
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/varfmt
BenchmarkFormat/small_int-4         	 4777486	       240 ns/op	      64 B/op	       4 allocs/op
BenchmarkFormat/small_string-4      	 2580116	       454 ns/op	     168 B/op	       8 allocs/op
BenchmarkFormat/big-4               	    9446	    120667 ns/op	  194656 B/op	      41 allocs/op
BenchmarkSprintf/small-4            	 8085470	       142 ns/op	      40 B/op	       4 allocs/op
BenchmarkSprintf/small_string-4     	 7574479	       152 ns/op	      40 B/op	       4 allocs/op
BenchmarkSprintf/big-4              	   22324	     53264 ns/op	   69000 B/op	      20 allocs/op
PASS
```