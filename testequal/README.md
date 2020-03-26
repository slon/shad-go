## testequal

В этой задаче нужно реализовать 4 test helper'а, аналогичных функциям из [testify](https://github.com/stretchr/testify):

```
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{})
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{})
```

Функции проверяют на равенство expected и actual и завершают тест, если проверка не прошла.
msgAndArgs попадают в описание ошибки через fmt.Sprintf.

Пример использования:
```
func TestMath(t *testing.T) {
	AssertEqual(t, 1, 2, "1 == 2")
}
```
вывод теста:
```
=== RUN   TestMath
--- FAIL: TestMath (0.00s)
    math_test.go:43: not equal:
        expected: 1
        actual  : 2
        message : 1 == 2
FAIL
FAIL    gitlab.com/slon/shad-go/testequal 0.003s
FAIL
```

В отличие от testify реализация ограничивает набор типов, с которыми умеет работать:
1. Целые числа: int, in64 и др (см. тесты)
2. string
3. map[string]string
4. []int
5. []byte

## Ссылки

1. testing.T: https://golang.org/pkg/testing/#T
2. type assertions: https://golang.org/doc/effective_go.html#interface_conversions
