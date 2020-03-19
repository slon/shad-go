## coverme

В этой задаче нужно покрыть простой todo-app http сервис unit тестами.

Имеющиеся `_test.go` файлы лучше не трогать,
при тестировании все изменения перетираются.

Package main можно не тестировать.

Тестирующая система будет проверяться code coverage.
Порог задан в [coverage_test.go](./app/coverage_test.go)

Как посмотреть coverage:
```
go test -v -cover ./coverme/...
```

## Ссылки

1. cover: https://blog.golang.org/cover
2. [gomock](https://github.com/golang/mock) для создания мока базы данных при тестировании серевера
3. [httptest.ResponseRecorder](https://golang.org/pkg/net/http/httptest/#ResponseRecorder) для тестирования handler'ов сервера
4. [httptest.Server](https://golang.org/pkg/net/http/httptest/#Server) для тестирования клинета

## O сервисе

Todo-app с минимальной функциональностью + client.

Запуск:
```
✗ go run ./coverme/main.go -port 6029
```

Health check:
```
✗ curl -i -X GET localhost:6029/    
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:46:02 GMT
Content-Length: 24

"API is up and working!"
```

Создать новое todo:
```
✗ curl -i localhost:6029/todo/create -d '{"title":"A","content":"a"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:41:31 GMT
Content-Length: 51

{"id":0,"title":"A","content":"a","finished":false}
```

Получить todo по id:
```
✗ curl -i localhost:6029/todo/0                                       
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:44:17 GMT
Content-Length: 51

{"id":0,"title":"A","content":"a","finished":false}
```

Получить все todo:
```
✗ curl -i -X GET localhost:6029/todo                                        
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:44:37 GMT
Content-Length: 53

[{"id":0,"title":"A","content":"a","finished":false}]
```