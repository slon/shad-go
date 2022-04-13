## wscat

wscat - примитивный аналог npm пакета [wscat](https://www.npmjs.com/package/wscat).

Websocket - это двусторонний канал поверх tcp. wscat - это websocket клиент.

wscat принимает на вход единственный аргумент `-addr` - адрес websocket сервера.
После подключения программа начинает читать с stdin'а и отправлять пользовательские строки на сервер,
печатая все сообщения от сервера в stdout.

Клиент должен обрабатывать SIGINT и SIGTERM и плавно завершаться с кодом 0, дожидаясь горутин.
Для этого может пригодиться [context](https://golang.org/pkg/context/).

Обратите внимание на то, что exit code `go run` - это не exit code исполняемого файла.

## Пример

Публичный echo сервер:
```
✗ $GOPATH/bin/wscat -addr ws://ws.ifelse.io
abc
abcdef
def^C2022/04/13 22:23:42 received signal interrupt
```
```
✗ echo $?
0
```

## Ссылки

1. websocket: https://en.wikipedia.org/wiki/WebSocket
2. gorilla/websocket: https://pkg.go.dev/github.com/gorilla/websocket
3. signal shutdown: https://p.go.manytask.org/06-http/lecture.slide#20
