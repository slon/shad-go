package sort

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func Reverse(data Interface) Interface {
	return &reverse{data}
}

type reverse struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
