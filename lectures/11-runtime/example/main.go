package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println(os.Getpid())
	time.Sleep(time.Hour)
}
