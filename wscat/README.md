## wscat

wscat - примитивный аналог npm пакета [wscat](https://www.npmjs.com/package/wscat).

Websocket - это двусторониий канал поверх tcp. wscat - это websocket клиент.

wscat принимает на вход единственный аргумент `-addr` - адрес websocket сервера.
После подключения программа начинает читать с stdin'а и отправлять пользовательские строки на сервер,
печатая все сообщения от сервера в stdout.

Клиент должен обрабатывать SIGINT и SIGTERM и плавно завершаться с кодом 0 дожидаясь горутин.
Для этого может пригодиться [context](https://golang.org/pkg/context/).

Обратите внимание на то, что exit code `go run` - это не exit code исполняемого файла.

## Пример

Публичный echo сервер:
```
✗ $GOPATH/bin/wscat -addr ws://echo.websocket.org
abc
abcdef
def^C2020/04/04 05:01:32 received signal interrupt
```
```
✗ echo $?
0
```

## Ссылки

1. websocket: https://en.wikipedia.org/wiki/WebSocket
2. signal shutdown: https://p.go.manytask.org/06-http/lecture.slide#20
