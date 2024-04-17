# shopfront

В этой задаче вам нужно реализовать хранилище счётчиков посещений поверх redis.

- Метод `RecordView` запоминает, что пользователь посетил страницу `item`-а.
- Метод `GetItems` загружает счётчики для пачки `item` ов. В поле `item[i].Viewed` должен
  быть записан флаг, означающий, что пользователь посетил `i`-ый `item`.

В этой задаче есть benchmark-и. Чтобы пройти его, ваше решение должно использовать [pipelining](https://github.com/redis/redis-doc/blob/master/docs/manual/pipelining/index.md).

## Запуск тестов на linux

Для работы тестов на ubuntu нужно установить пакет `redis-server`.

```
sudo apt install redis-server
```

Если вы работаете на другом дистрибутиве linux, воспользуйтесь своим пакетным менеджером.

Тесты сами запускают `redis` в начале, и останавливают его в конце.

## Запуск redis в docker

Комментарии по запуску бд в docker смотрите в задаче [dao](../dao/).

```
(cd shopfront && docker compose up -d --wait && env REDIS="localhost:6379" go test -v ./... -count=1 || true && docker compose down)
```
