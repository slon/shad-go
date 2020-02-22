## fetchall

В этой задаче нужно написать консольную утилиту,
которая принимает на вход произвольное количество http URL'ов и скачивает их содержимое **конкурентно**.

Программа не должна останавливаться на невалидном URL'e.
Текст ответов можно игнорировать.
Вместо этого можно залогировать прогресс в произвольном формате.

Пример:
```
$ fetchall https://gopl.io golang.org http://golang.org
Get golang.org: unsupported protocol scheme ""
1.05s    11071  http://golang.org
2.18s     4154  https://gopl.io
2.18s elapsed
```

В примере логируются времена обработки индивидуальных запросов, размеры ответов и общее время работы программы.
Можно видеть, что общее время работы равно максимуму, а не сумме времён индивидуальных запросов.

### Проверка решения

Для запуска тестов нужно выполнить следующую команду:

```
go test -v ./fetchall/...
```

### Запуск программы

```
go run -v ./fetchall/main.go
```

### Компиляция

```
go install ./fetchall/...
```

После выполнения в `$GOPATH/bin` появится исполняемый файл с именем `fetchall`.

### Линтер

Установите [golangci-lint](https://github.com/golangci/golangci-lint), если вы ещё этого не сделали, и проверьте решение перед отправкой!
```
golangci-lint -v run ./fetchall/...
```

### Ссылки

1. Чтение аргументов командной строки: https://gobyexample.com/command-line-arguments
2. HTTP запрос: https://golang.org/pkg/net/http/
3. Запуск горутин: https://gobyexample.com/goroutines
4. Ожидание завершения горутин: https://gobyexample.com/channels
4. Замер времени: https://golang.org/pkg/time/#Since
