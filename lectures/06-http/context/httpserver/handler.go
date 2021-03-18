package httpserver

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqTime := ctx.Value(ReqTimeContextKey{}).(time.Time)
	defer func() {
		fmt.Printf("handler finished in %s", time.Since(reqTime))
	}()

	fd, _ := os.Open("core.c")
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			_, _ = w.Write(scanner.Bytes())
		}
	}
}
