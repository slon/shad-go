package main

import (
	"crypto/md5"
	"io"
	"net/http/httputil"
	"os"
	"strings"
)

func transfer(clientWriter io.Writer, responseBody io.Reader) {
	digest := md5.New()
	_, _ = io.Copy(digest,
		httputil.NewChunkedReader(
			io.TeeReader(
				responseBody,
				clientWriter,
			),
		),
	)
	parseTrailers(responseBody)
}

func parseTrailers(r io.Reader) {
	_, _ = io.Copy(os.Stdout, r)
}

func main() {
	data := "4\r\nWiki\r\n5\r\npedia\r\nE\r\n in\r\n\r\nchunks.\r\n0\r\nDate: Sun, 06 Nov 1994 08:49:37 GMT\r\nContent-MD5: 1B2M2Y8AsgTpgAmY7PhCfg==\r\n\r\n"
	transfer(os.Stdout, strings.NewReader(data))
}
