## firewall

В этой задаче нужно написать примитивный файрвол.

Файрвол - это прокси сервер, пропускающий через себя все запросы,
и отвергающий некоторые из них по заданному набору правил.

Пример правил можно посмотреть в [example.yaml](./configs/example.yaml).
Все правила можно разделить на 2 группы: те, что применяются к запросу и те, что применяются к ответу.

Для решения этой задачи удобно использовать `http.RoundTripper`.

На все заблокированные запросы нужно отвечать статусом 403 и строкой `Forbidden`.

Сервер должен принимать следующие агрументы:
* `-service-addr` - адрес защищаемого сервиса
* `-conf` - путь к .yaml конфигу с правилами
* `-addr` - адрес, на котором будет развёрнут файрвол

## Примеры:
В [cmd/service](./cmd/service/main.go) находится примитивный сервис, который мы хотим защитить.
```
go run ./firewall/cmd/service/main.go -port 8080
```
Делаем запрос:
```
curl -i http://localhost:8080/list -d '"loooooooooooooooooooooooooooooong-line"'
HTTP/1.1 200 OK
Date: Thu, 02 Apr 2020 19:14:36 GMT
Content-Length: 40
Content-Type: text/plain; charset=utf-8

"loooooooooooooooooooooooooooooong-line"
```
Стартуем firewall:
```
go run ./firewall/cmd/firewall/main.go -service-addr http://localhost:8080 -addr localhost:8081 -conf ./firewall/configs/example.yaml
```
Делаем тот же запрос через него:
```
curl -i http://localhost:8081/list -d '"loooooooooooooooooooooooooooooong-line"'
HTTP/1.1 403 Forbidden
Date: Thu, 02 Apr 2020 19:14:40 GMT
Content-Length: 9
Content-Type: text/plain; charset=utf-8

Forbidden
```
Сработало правило на максимальную длину запроса.

## Resources

* project layout: https://github.com/golang-standards/project-layout
* reverse proxy: https://en.wikipedia.org/wiki/Reverse_proxy
* yaml: https://gopkg.in/yaml.v2
* regexp: https://golang.org/pkg/regexp/
