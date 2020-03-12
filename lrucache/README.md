## lrucache

В этой задаче нужно написать простой Least recently used cache.

LRU cache - это key-value storage фиксированного размера, реализующий операции:
* `set(k, v)` - обновляет хранимое по ключу `k` значение.
  В случае, если операция приводит к превышению размера кэша,
  из того удаляется значение по самому "старому" ключу.
* `get(k) -> v, ok` - возвращает значение, хранимое по ключу `k`.

Обе функции `set` и `get` обновляют access time ключа.

В файле [cache.go](./cache.go) задан интерфейс `Cache` с подробным описанием всех методов.

Нужно написать реализацию и конструктор, принимающий размер кэша:
```
func New(cap int) Cache
```

## Ссылки

1. container/list: https://golang.org/pkg/container/list/
2. wiki: https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
