package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func finalizers() {
	f, err := os.Open("x.txt")

	select {}

	_ = f.Close()

	_, _ = f, err
}

type myX struct{}

func (xx *myX) close() {
	fmt.Println("finalized")

	runtime.SetFinalizer(xx, (*myX).close)
}

var g any

func finalizersMy() {
	xx := new(myX)

	runtime.SetFinalizer(xx, (*myX).close)

	for {
		g = make([]byte, 1<<20)

		time.Sleep(time.Second)

		xx = nil
		runtime.GC()
	}
}

func main() {
	finalizersMy()
}
