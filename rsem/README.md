# rsem

Реализуйте семафор используя redis. В отличии от [semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore) в shared memory, такой семафор можно использовать в распределённой системе, чтобы синхронизировать независимые процессы.

Ваша реализация должна обладать свойством robustness, это значит - семафор должен автоматически отпускаться, если процесс держаший его умер (тесты предполагают, что это происходит через секунду после смерти процесса). Чтобы реализовать это свойство, используйте [TTL на ключи](https://redis.io/docs/latest/commands/expire/).

В вашей реализации может потребоваться атомарно выполнить набор команд. Для этого можно использовать
[транзакции](https://redis.io/docs/latest/develop/interact/transactions/) или [lua скрипт](https://redis.io/docs/latest/develop/interact/programmability/eval-intro/)
