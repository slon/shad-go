package main

import (
	"fmt"
	"runtime"
	"time"
)

type myX struct{}

func (xx *myX) close() {
	fmt.Println("finalized")
}

func main() {
	xx := new(myX)
	runtime.SetFinalizer(xx, (*myX).close)

	for {
		time.Sleep(time.Second)
		xx = nil
		runtime.GC()
	}
}
