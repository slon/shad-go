package treeiter

import "fmt"

func ExampleDoInOrder() {
	tree := &ValuesNode[string]{
		value: "root",
		left: &ValuesNode[string]{
			value: "left",
		},
		right: &ValuesNode[string]{
			value: "right",
		},
	}

	DoInOrder(tree, func(t *ValuesNode[string]) {
		fmt.Println(t.value)
	})

	// Output:
	// left
	// root
	// right
}
