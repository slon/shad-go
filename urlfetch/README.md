## urlfetch

В этой задаче нужно написать консольную утилиту,
которая принимает на вход произвольное количество http URL'ов и последоватльно скачивает их содержимое.

При обработке невалидного URL'а программа должна завершаться с ненулевым exit кодом.

### Примеры

Успешный запуск:
```
$ urlfetch https://golang.org https://go.dev
<!DOCTYPE html>
<html lang="en">
<meta charset="utf-8">
<meta name="description" content="Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.">
...
```

Неуспешный запуск:
```
$ urlfetch golang.org       
fetch: Get golang.org: unsupported protocol scheme ""
```

### Проверка решения

Для запуска тестов нужно выполнить следующую команду:

```
go test -v ./urlfetch/...
```

### Запуск программы

```
go run -v ./urlfetch/main.go
```

### Компиляция

```
go install ./urlfetch/...
```

После выполнения в `$GOPATH/bin` появится исполняемый файл с именем `urlfetch`.

### Линтер

Установите [golangci-lint](https://github.com/golangci/golangci-lint), если вы ещё этого не сделали, и проверьте решение перед отправкой!
```
golangci-lint -v run ./urlfetch/...
```

### Walkthrough

1. Чтение аргументов командной строки: https://gobyexample.com/command-line-arguments
2. HTTP запрос: https://golang.org/pkg/net/http/
