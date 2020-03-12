package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() { // Counter
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() { // Squarer
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			squares <- x * x
		}

		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}
}
