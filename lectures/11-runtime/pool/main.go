package main

import (
	"net/http"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return &Decoder{}
	},
}

type Decoder struct {
	buf [1 << 20]byte
}

func New() *Decoder {
	return pool.Get().(*Decoder)
}

func (c *Decoder) Close() {
	pool.Put(c)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := New()

		d.Close()
	})
}
