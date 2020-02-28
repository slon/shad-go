# otp

Напишите код, который реализует схему шифрования [stream cipher](https://en.wikipedia.org/wiki/Stream_cipher).

Потоковый шифр обрабатывает поток по одному байту за раз. Каждый байт входного потока xor-ится с байтом из prng и записывается в выходной поток.

Вам нужно реализовать две версии api.

```golang
func NewReader(r io.Reader, prng io.Reader) io.Reader
func NewWriter(w io.Writer, prng io.Reader) io.Writer
```

* `NewReader` принимает входной поток `r` и генератор случайных чисел `prng`. `NewReader` возвращает
 `io.Reader`, который читает поток из `r` и расшифровывает его с помощью `prng`.
* `NewWriter` принимает выходной поток `w` и генератор случайных чисел `prng`. `NewWriter` возвращает
 `io.Writer`, который шифрует поток с помощью `prng` и пишет его в `w`.

Вы можете считать, что prng никогда не может вернуть ошибку.

## Замечания
 - Прочитайте контракт [io.Reader](https://golang.org/pkg/io/#Reader) и [io.Writer](https://golang.org/pkg/io/#Writer) в документации.
 - То что шифр работает с одним байтом, не значит что нужно передавать в Read() слайс размера 1.
 - Подумайте, почему потоковый шифр в стандартной библиотеке имеет интерфейс [cipher.Stream](https://golang.org/pkg/crypto/cipher/#Stream), а не `io.Reader` как у нас.
 - Для отладки вы можете использовать `iotest.NewReadLogger` и `iotest.NewWriteLogger` из пакета [iotest](https://golang.org/pkg/testing/iotest/).
