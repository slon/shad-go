# blowfish

Реализуйте cgo wrapper для шифра blowfish из библиотеки openssl.

- Вам нужно использовать две функции:

    ```c
    // BF_set_key инициализирует BF_KEY
    void BF_set_key(BF_KEY *key, int len, const unsigned char *data);

    // BF_ecb_encrypt шифрует или дешифрует блок размером в 8 байт.
    void BF_ecb_encrypt(const unsigned char *in, unsigned char *out, BF_KEY *key, int enc);
    ```

- Реализация не должна делать динамического выделения памяти.
- Для сборки этой задачи, на вашей системе должен быть установлен dev пакет openssl. На ubuntu установить пакет можно командой `sudo apt-get install libssl-dev`. Сборка под другие платформы не гарантируется.

<details>
<summary><b>Установка openssl на Mac OS через Homebrew</b></summary>

```
$ brew install openssl@3

После установки Homebrew предупредит вас о том, что для корректной работы библиотеки может понадобиться выставить несколько переменных окружения, нас интересует последняя:
$ export PKG_CONFIG_PATH="/usr/local/opt/openssl@3/lib/pkgconfig"

Важно: на вашем компьютере путь может быть другим, а именно, начинаться с префикса /opt/homebrew вместо /usr/local. Если у вас уже стоит openssl, для правильного экспорта переменной вы можете узнать этот префикс через команду:
$ brew --prefix

Вы также можете выставить переменную окружения в самом GoLand: Run -> Edit Configurations -> Environment:
PKG_CONFIG_PATH="/usr/local/opt/openssl@3/lib/pkgconfig"
```

</details>

**Disclaimer:** Эта задача дана в учебных целях. Помните, что (1) нельзя реализовывать собственную криптографию, (2) шифр blowfish устарел, (3) в стандартной библиотеке есть pure go реализация для большинства криптографических примитивов.

## Ссылки

1. [openssl](https://www.openssl.org/docs/man1.0.2/man3/blowfish.html)
2. [cgo](https://golang.org/cmd/cgo/)
