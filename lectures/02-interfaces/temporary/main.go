package main

import "errors"

type Temporary interface {
	IsTemporary() bool
}

func do() error { return nil }

func main() {
	err := do()

	var terr Temporary
	if errors.As(err, &terr) && terr.IsTemporary() {
		//...
	}
}
