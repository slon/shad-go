package example

//go:generate mockgen -package example -destination mock.go . Foo

type Foo interface {
	Bar(x int) int
}

func SUT(f Foo) {
	// ...
}
