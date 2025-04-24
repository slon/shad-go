package main

import (
	"net/http"
	"sync"
)

// START OMIT
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

func handler(w http.ResponseWriter, r *http.Request) {
	d := New()
	d.Close()

}

// END OMIT

func main() {}
