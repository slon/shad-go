# Поиск лишних аллокаций

Для анализа аллокаций в задачах с бенчмарками удобно использовать встроенный в язык профайлер `pprof`.

Сперва нужно запустить бенчмарк с профайлером
```
✗ go test -v -run=^$ -bench=BenchmarkSprintf -tags private,solution -memprofile=mem.out ./varfmt/...
goos: linux
goarch: amd64
pkg: gitlab.com/slon/shad-go/varfmt
cpu: Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz
BenchmarkSprintf
BenchmarkSprintf/small
BenchmarkSprintf/small-8         	19429222	        62.93 ns/op	       2 B/op	       1 allocs/op
BenchmarkSprintf/small_string
BenchmarkSprintf/small_string-8  	13282659	        84.48 ns/op	      16 B/op	       1 allocs/op
BenchmarkSprintf/big
BenchmarkSprintf/big-8           	   20089	     62372 ns/op	   16388 B/op	       1 allocs/op
PASS
ok  	gitlab.com/slon/shad-go/varfmt	4.363s
```

Сэмплы профайлера будут записаны в бинарный файл `mem.out`.

Этот файл можно открыть командой
```
✗ go tool pprof mem.out
File: varfmt.test
Type: alloc_space
Time: Mar 14, 2023 at 9:07pm (MSK)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

В результате откроется интерактивная среда, в которой есть всякие команды, например `help`.

Чтобы узнать как запускать конкретную команду нужно вызвать `help` для неё
```
(pprof) help top
Outputs top entries in text form
  Usage:
    top [n] [focus_regex]* [-ignore_regex]* [-cum] >f
    Include up to n samples
    Include samples matching focus_regex, and exclude ignore_regex.
    -cum sorts the output by cumulative weight
    Optionally save the report on the file f
```

Команда `top` покажет места, где суммарно было больше всего аллокаций
```
(pprof) top -cum
Showing nodes accounting for 715.73MB, 99.37% of 720.23MB total
Dropped 24 nodes (cum <= 3.60MB)
      flat  flat%   sum%        cum   cum%
  715.73MB 99.37% 99.37%   716.73MB 99.51%  fmt.Sprintf
         0     0% 99.37%   716.73MB 99.51%  gitlab.com/slon/shad-go/varfmt.BenchmarkSprintf.func1
         0     0% 99.37%   716.73MB 99.51%  testing.(*B).launch
         0     0% 99.37%   716.73MB 99.51%  testing.(*B).runN
(pprof)
```

Команда `svg` сдампит `top` в картинку.

```
(pprof) svg
Generating report in profile001.svg
```

Команда `list` покажет строчки, в которых происходят аллокации
```
(pprof) list fmt.Sprintf
Total: 720.23MB
ROUTINE ======================== fmt.Sprintf in /usr/local/go/src/fmt/print.go
  715.73MB   716.73MB (flat, cum) 99.51% of Total
         .          .    213:	return Fprintf(os.Stdout, format, a...)
         .          .    214:}
         .          .    215:
         .          .    216:// Sprintf formats according to a format specifier and returns the resulting string.
         .          .    217:func Sprintf(format string, a ...any) string {
         .      512kB    218:	p := newPrinter()
         .          .    219:	p.doPrintf(format, a)
  715.73MB   715.73MB    220:	s := string(p.buf)
         .   512.50kB    221:	p.free()
         .          .    222:	return s
         .          .    223:}
         .          .    224:
         .          .    225:// Appendf formats according to a format specifier, appends the result to the byte
         .          .    226:// slice, and returns the updated slice.
```
