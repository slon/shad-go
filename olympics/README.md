## olympics

В этой задаче нужно написать http сервер со следующим API

* GET /athlete-info?name=S

   -> вернуть информацию по атлету с именем: откуда, сколько каких медалей выиграл всего и по годам

   для упрощения считаем, что спортсмен идентифицируется именем. В случае,
   если атлет выступал более чем за одну страну, нужно приписать его
   к первой стране в порядке исходных данных.

* GET /top-athletes-in-sport?sport=SSSS&limit=K

   -> вернуть top-K (default=3) спортсменов по абсолютному числу медалей в их спортивной карьере в указанном виде спорта 
   (сортируем по золотым, потом по серебрянным, потом по бронзе, потом лексикографически по имени спортсмена)

* GET /top-countries-in-year?year=YYYY&limit=K

   -> вернуть top-K (default=3) стран в порядке медального зачета (сортируем по золотым, потом по серебрянным, потом по бронзе, потом лексикографически по стране)

используя данные о победителях и призёрах олимпийских игр из [./testdata/olympicWInners.json](./testdata/olympicWinners.json).

Сервер должен слушать порт, переданный через аргумент `-port`. Путь к json'у с данными передаётся через флаг `-data`.

### Примеры

Запуск:
```
$ olympics -port 6029 -data ./olympics/testdata/olympicWinners.json
```

#### athlete-info

Успешный запрос (200, json фиксированного вида):
```
$ curl -X GET "localhost:6029/athlete-info?name=Michael%20Phelps"
{
  "athlete": "Michael Phelps",
  "country": "United States",
  "medals": {
    "gold": 18,
    "silver": 2,
    "bronze": 2,
    "total": 22
  },
  "medals_by_year": {
    "2004": {
      "gold": 6,
      "silver": 0,
      "bronze": 2,
      "total": 8
    },
    "2008": {
      "gold": 8,
      "silver": 0,
      "bronze": 0,
      "total": 8
    },
    "2012": {
      "gold": 4,
      "silver": 2,
      "bronze": 0,
      "total": 6
    }
  }
}
```

Спортсмен не найден (404, произвольное сообщение об ошибке):
```
$ curl -i -X GET "localhost:6029/athlete-info?name=AB"
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 19 Feb 2020 23:24:30 GMT
Content-Length: 21

athlete AB not found
```

#### top-athletes-in-sport

Успешный запрос (200, json фиксированного вида):
```
$ curl -X GET "localhost:6029/top-athletes-in-sport?sport=Swimming&&limit=1"
[
  {
    "athlete": "Michael Phelps",
    "country": "United States",
    "medals": {
      "gold": 18,
      "silver": 2,
      "bronze": 2,
      "total": 22
    },
    "medals_by_year": {
      "2004": {
        "gold": 6,
        "silver": 0,
        "bronze": 2,
        "total": 8
      },
      "2008": {
        "gold": 8,
        "silver": 0,
        "bronze": 0,
        "total": 8
      },
      "2012": {
        "gold": 4,
        "silver": 2,
        "bronze": 0,
        "total": 6
      }
    }
  }
]
```

Неизвестный вид спорта (404, произвольное сообщение об ошибке):
```
$ curl -i -X GET "localhost:6029/top-athletes-in-sport?sport=chess"            
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 20 Feb 2020 00:42:24 GMT
Content-Length: 24

sport 'chess' not found
```

#### top-countries-in-year

Успешный запрос (200, json фиксированного вида):
```
$ curl -X GET "localhost:6029/top-countries-in-year?year=2012&&limit=2"
[
  {
    "country": "United States",
    "gold": 145,
    "silver": 63,
    "bronze": 46,
    "total": 254
  },
  {
    "country": "China",
    "gold": 56,
    "silver": 40,
    "bronze": 29,
    "total": 125
  }
]
```

Год не найден (404, произвольное сообщение):
```
$ curl -i -X GET "localhost:6029/top-countries-in-year?year=2009" 
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 20 Feb 2020 00:10:27 GMT
Content-Length: 20

year 2009 not found
```
