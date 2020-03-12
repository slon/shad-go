package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() { // Counter
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	go func() { // Squarer
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}
