package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func main() {
	// все куки записанные в этот Jar будут передаваться
	// и изменяться во всех запросах
	cj, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 1 * time.Second,
		Jar:     cj,
		Transport: &http.Transport{
			// резмер буферов чтения и записи (4KB по-умолчанию)
			WriteBufferSize: 32 << 10,
			ReadBufferSize:  32 << 10,
			// конфиг работы с зашифрованными соединениями
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{},
				RootCAs:      &x509.CertPool{},
				// только для отладки!
				InsecureSkipVerify: true,
				// ..
			},
			// ...
		},
	}
	_ = client
}
