## urlshortener

В этой задаче нужно написать http сервер со следующим API:

* POST /shorten {"url": "\<URL\>"} -> {"key": "\<KEY\>"}
* GET /go/\<KEY\> -> 302

GET и POST - это методы HTTP. GET запрос используется для того, чтобы получать данные, а POST - чтобы добавлять и модифицировать.

В тело `/shorten` запроса будет передаваться json вида
```
{"url":"https://github.com/golang/go/wiki/CodeReviewComments"}
```

Сервер должен ответить json'ом следующего вида:
```
{
  "url": "https://github.com/golang/go/wiki/CodeReviewComments",
  "key": "ed1De1"
}
```

`ed1De1` здесь - это сгенерированное сервисом число.

После такого `/shorten` можно делать `/go/ed1De1`.
Ответ должен иметь иметь HTTP код 302.
302 указывает на то, что запрошенный ресурс был временно перемещен на другой адрес (передаваемый в HTTP header'е `Location`).

Если открыть http://localhost:6029/go/ed1De1 в браузере, тот перенаправит на https://github.com/golang/go/wiki/CodeReviewComments.

Сервер должен слушать порт, переданный через аргумент `-port`.

### Примеры

Запуск:
```
$ urlshortener -port 6029
```

Успешное добавление URL'а (200, Content-Type: application/json):
```
$ curl -i -X POST  "localhost:6029/shorten" -d '{"url":"https://github.com/golang/go/wiki/CodeReviewComments"}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 15 Feb 2020 23:35:26 GMT
Content-Length: 82

{"url":"https://github.com/golang/go/wiki/CodeReviewComments","key":"65ed150831"}
```

Невалидный json (400):
```
$ curl -i -X POST  "localhost:6029/shorten" -d '{"url":"https://github.com'                                   
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 15 Feb 2020 23:30:27 GMT
Content-Length: 16

invalid request
```

Успешный запрос (302, Location header):
```
$ curl -i -X GET  "localhost:6029/go/c1464c853a"                                                               
HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://github.com/golang/go/wiki/CodeReviewComments
Date: Sat, 15 Feb 2020 23:25:26 GMT
Content-Length: 75

<a href="https://github.com/golang/go/wiki/CodeReviewComments">Found</a>.
```

Несуществующий key (404):
```
$ curl -i -X GET  "localhost:6029/go/uaaab"
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 15 Feb 2020 23:26:48 GMT
Content-Length: 14

key not found
```

### Состояние

Своё состояние сервис должен целиком хранить в памяти.
Стандартный http server на каждый запрос запускает handler в отдельной горутине (https://golang.org/pkg/net/http/#Serve),
поэтому доступ к состоянию нужно защитить. Например, это можно сделать с помощью [мьютекса](https://golang.org/pkg/sync/#Mutex).

## Ссылки

1. Пример web сервера и работы с общим состоянием: https://p.go.manytask.org/00-intro/lecture.slide#24
2. протокол HTTP: https://ru.wikipedia.org/wiki/HTTP
3. http multiplexer: https://golang.org/pkg/net/http/#ServeMux
4. десериализация json'а: https://golang.org/pkg/encoding/json/#example_Unmarshal
5. генерация случайных данных: https://golang.org/pkg/math/rand/
