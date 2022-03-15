//go:build !change

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const exitCommand = "bye"

func main() {
	e := NewEvaluator()

	fmt.Printf("Welcome to Forth evaluator! To exit type %q.\n", exitCommand)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">")
		scanner.Scan()

		text := scanner.Text()
		if text == exitCommand {
			break
		}

		stack, err := e.Process(text)
		if err != nil {
			fmt.Printf("Evaluation error: %s\n", err)
		}

		printStack(stack)
	}
}

func printStack(stack []int) {
	s := make([]string, 0, len(stack))
	for _, n := range stack {
		s = append(s, strconv.Itoa(n))
	}
	fmt.Printf("Stack: %s\n", strings.Join(s, ", "))
}
