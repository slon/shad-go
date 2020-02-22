## wordcount

В этой задаче нужно написать консольную утилиту, которая принимает на вход набор файлов
и печатает в stdout уникальные строки и количество раз, которое они встречаются.
В stdout попадают только те строки, что встречаются суммарно хотя бы дважды.

Формат вывода:
```
<COUNT>\t<LINE>
<COUNT>\t<LINE>
...
```

#### Пример:

Если a.txt - это файл со следующим содержимым:
```
a
b
a
c
```
а b.txt - со следующим:
```
a
b
a
c
```
то результат выполнения команды `wordcount a.txt b.txt` должен выглядеть так (с точностью до перестановки строк):
```
2	c
4	a
2	b
```

### Проверка решения

Для запуска тестов нужно выполнить следующую команду:

```
go test -v ./wordcount/...
```

### Компиляция

```
go install ./wordcount/...
```

После выполнения в `$GOPATH/bin` появится исполняемый файл с именем `wordcount`.

### Линтер

Установите [golangci-lint](https://github.com/golangci/golangci-lint) и проверьте решение перед отправкой!
```
golangci-lint -v run ./wordcount/...
```

### Walkthrough

#### 1. Чтение аргументов командной строки
https://gobyexample.com/command-line-arguments
#### 2. Чтение файлов
https://gobyexample.com/reading-files
#### 3. Парсинг содержимого
https://gobyexample.com/string-functions
#### 4. Подсчёт вхождений
https://gobyexample.com/maps
#### 5. Вывод результатов
https://gobyexample.com/string-formatting
