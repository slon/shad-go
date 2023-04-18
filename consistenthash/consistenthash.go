//go:build !solution

package consistenthash

type Node interface {
	// ID is some persistent and unique identifier
	ID() string
}

type ConsistentHash[N Node] struct {
}

func New[N Node]() *ConsistentHash[N] {
	panic("implement me")
}

func (h *ConsistentHash[N]) AddNode(n *N) {
	panic("implement me")
}

func (h *ConsistentHash[N]) RemoveNode(n *N) {
	panic("implement me")
}

func (h *ConsistentHash[N]) GetNode(key string) *N {
	panic("implement me")
}
