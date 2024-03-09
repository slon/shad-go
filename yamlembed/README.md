# yamlembed

В этой задаче нужно познакомиться с [yaml'ем](https://pkg.go.dev/gopkg.in/yaml.v2) и (hopefully) кое-что понять про эмбеддинг.

yaml — это человекочитаемый формат сериализации данных.
В отличие от json'а, в yaml'е есть комментарии и нет необходимости в фигурных скобках с кавычками.
Это делает удобным использование yaml в файлах конфигурации сервисов.
Превращение Go объекта в yaml называется "сериализацией в yaml", обратный процесс — "десериализацией".

Например, у нас есть следующий yaml файл
```yaml
person_age: 20
hobbies:
- painting
- playing_music
```

Мы хотим заполнить Go структуру
```go
type Person struct {
	Age     int      `yaml:"person_age"`
	Hobbies []string `yaml:"hobbies,omitempty"`
}
```
содержимым файла.

Это можно сделать следующим образом:
```go
package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Person struct {
	ID      int      `yaml:"-"`
	Age     int      `yaml:"person_age"`
	Hobbies []string `yaml:"hobbies,omitempty"`
}

func main() {
	data := []byte(`
id: 124
person_age: 20
hobbies:
- painting
- playing_music
`)

	var p Person
	if err := yaml.Unmarshal(data, &p); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", p)
}
```
Выведется
```
{0 20 [painting playing_music]}
```

В обратную сторону. Превратить Go объект в yaml:
```go
package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Person struct {
	ID      int      `yaml:"-"`
	Age     int      `yaml:"person_age"`
	Hobbies []string `yaml:"hobbies,omitempty"`
}

func main() {
	p := Person{ID: 124, Age: 20, Hobbies: nil}
	data, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
```
Выведется
```yaml
person_age: 20
```

yaml тэги в структуре определяют как должен сериализовываться объект.
* `yaml:"-"` говорит, что поле ID нужно игнорировать, например, если приложение хочет само заполнить это поле.
* `yaml:"person_age"` говорит, что person_age из yaml'я нужно класть в поле Age.
* `,omitempty` говорит, что поле не нужно сериализовывать, если в слайсе нет элементов или он nil

Иногда при сериализации/десериализации хочется более сложной логики, которой на одних тэгах не реализовать.
Например, мы можем захотеть заполнять поля значениями, зависящими от других полей.

Посмотрите, каким образом библиотека позволяет это сделать. Это довольно частый паттерн, который используется во всех таких библиотеках сериализации.

## Что нужно сделать?

Вам даны Go структуры и тесты, которые говорят, как структуры должны сериализовываться в yaml (довольно естественно).
Нужно починить сериализацию.

## Follow-up questions

1. Сколько функций вам пришлось написать и почему меньшего количества недостаточно?
2. Что такое `type indirection` и какую роль он выполняет?

## Ссылки

* Что в принципе бывает в yaml https://en.wikipedia.org/wiki/YAML
