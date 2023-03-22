## iprange

В этой задаче вам предстоит познакомиться с фаззингом, и его нативной поддержкой в go.

Нужно поправить баг в функции `ParseList`.

`ParseList` принимает на вход строку с описанием рейнджей ip адрессов в одном из `n` форматов
* `10.0.0.1`
* `10.0.0.0/24`
* `10.0.0.*`
* `10.0.0.1-10`

и возвращает список пар `(min ip, max ip)` (см. [example](./example_test.go)).

Для обнаружения бага (crash функции) предлагается написать fuzz тест на функцию `ParseList`.

#### Проверка решения

Во-первых, должны работать имеющиеся тесты.
```
go test -v ./iprange...
```

Во-вторых, в CI есть приватные тесты, молча падающие на неправильной `ParseList`.

Как запустить fuzz тесты?
```
go test -v -fuzz=. ./iprange...
```

### Ссылки

* fuzzing tutorial https://go.dev/doc/tutorial/fuzz
* fuzzing design draft https://go.googlesource.com/proposal/+/master/design/draft-fuzzing.md
