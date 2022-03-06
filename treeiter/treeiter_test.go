package treeiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValuesNode[T any] struct {
	value       T
	left, right *ValuesNode[T]
}

func (t ValuesNode[T]) Left() *ValuesNode[T] {
	return t.left
}

func (t ValuesNode[T]) Right() *ValuesNode[T] {
	return t.right
}

type Collector[T any] struct {
	nodes []*ValuesNode[T]
}

func (c *Collector[T]) Collect(tree *ValuesNode[T]) {
	c.nodes = append(c.nodes, tree)
}

func TestInOrderNil(t *testing.T) {
	var collector Collector[any]
	DoInOrder(nil, collector.Collect)
}

func TestInOrderIntTree(t *testing.T) {
	root := &ValuesNode[int]{value: 1}
	collector := Collector[int]{}

	DoInOrder(root, collector.Collect)

	assert.Equal(t, []*ValuesNode[int]{{value: 1}}, collector.nodes)
}

func TestInOrderStringTree(t *testing.T) {
	rightLeftRight := &ValuesNode[string]{
		value: "right left right",
	}
	rightLeft := &ValuesNode[string]{
		value: "right left",
		right: rightLeftRight,
	}
	right := &ValuesNode[string]{
		value: "right",
		left:  rightLeft,
	}
	left := &ValuesNode[string]{
		value: "left",
	}

	root := &ValuesNode[string]{
		value: "root",
		left:  left,
		right: right,
	}
	var collector Collector[string]

	DoInOrder(root, collector.Collect)

	assert.Equal(t, []*ValuesNode[string]{left, root, rightLeft, rightLeftRight, right}, collector.nodes)
}
