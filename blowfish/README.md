# blowfish

Реализуйте cgo wrapper для реализации blowfish из пакета openssl.

- Вам нужно использовать две функции:

    ```c
    // BF_set_key инициализирует BF_KEY
    void BF_set_key(BF_KEY *key, int len, const unsigned char *data);

    // BF_ecb_encrypt шифрует или дешифрует блок размером в 8 байт.
    void BF_ecb_encrypt(const unsigned char *in, unsigned char *out, BF_KEY *key, int enc);
    ```

- Реализация не должна делать динамического выделения памяти.
- Для сборки этой задачи, на вашей системе должен быть установлен dev пакет openssl. На ubuntu установить пакет можно командой `sudo apt-get install libssl-dev`. Сборка под другие платформы не гарантируется.

**Disclaimer:** Эта задача дана в учебных целях. Помните, что (1) нельзя реализовывать собственную криптографию, (2) шифр blowfish устарел, (3) в стандартной библиотеке есть pure go реализция для большинства криптографических приминивов.

## Ссылки

1. [openssl](https://www.openssl.org/docs/man1.0.2/man3/blowfish.html)
2. [cgo](https://golang.org/cmd/cgo/)
