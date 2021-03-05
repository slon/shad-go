# ratelimit

Напишите примитив синхронизации, ограничивающий число вызовов на интервале времени.

```go
func NewLimiter(maxCount int, interval time.Duration) *Limiter

func (l *Limiter) Acquire(ctx context.Context) error

func (l *Limiter) Stop()
```

`Limiter` должен гарантировать, что на любом интервале времени `interval`, не больше `maxCount` вызовов
`Acquire` могут завершиться без ошибки.

Каждый вызов `Acquire` должен либо завершаться успешно, либо завершаться с ошибкой в случае если `ctx` отменили
во время ожидания. Об отмене `ctx` нужно узнавать по закрытию канала `ctx.Done()`. Если `ctx` отменён,
нужно возвращать ошибку `ctx.Err()`.

Вызовы `Acquire()` после `Stop()`, должны сразу завершаться с ошибкой ErrStopped.
