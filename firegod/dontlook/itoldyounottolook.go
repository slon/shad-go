package dontlook

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewService(r chi.Router) {
	replica := os.Getenv("REPLICA")

	switch replica {
	case "0":
		go func() {
			for {
				go func() { select {} }()

				time.Sleep(time.Millisecond * 100)
			}
		}()

	case "2":
		go func() {
			var buf []string
			for i := 0; true; i++ {
				buf = append(buf, strings.Repeat("f", 1<<20))
				_ = buf

				time.Sleep(time.Millisecond * 100 * time.Duration(i))
			}
		}()

	case "3":
		go func() {
			for {
				_, _ = net.Dial("tcp", "localhost:8080")
				time.Sleep(time.Millisecond * 10)
			}
		}()

	case "4":
		go func() {
			time.Sleep(time.Second * 15)
			os.Exit(1)
		}()
	}

	runClient := func(getUrl func() string) {
		for {
			go func() {
				rsp, err := http.Get(fmt.Sprintf("http://localhost:8080%s", getUrl()))
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return
				}

				_, _ = io.ReadAll(rsp.Body)
				_ = rsp.Body.Close()
			}()

			time.Sleep(time.Microsecond * 100)
		}
	}

	r.Get("/apple", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		time.Sleep(time.Millisecond * 100)
		_, _ = w.Write([]byte("OK"))
	})
	go runClient(func() string { return "/apple" })

	r.Get("/banana", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(bytes.Repeat([]byte("f"), 1<<10))

		if replica == "0" {
			time.Sleep(time.Second)
		}
	})
	go runClient(func() string { return "/banana" })

	var k atomic.Int32
	r.Get("/orange", func(w http.ResponseWriter, r *http.Request) {
		if k.Add(1)%100 == 0 {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}

	})
	go runClient(func() string { return "/orange" })

	r.Get("/kiwi/{id}", func(w http.ResponseWriter, r *http.Request) {
		if replica == "3" {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(500)
		}
	})

	var i atomic.Int32
	go runClient(func() string { return fmt.Sprintf("/kiwi/%d", i.Add(1)) })
}
