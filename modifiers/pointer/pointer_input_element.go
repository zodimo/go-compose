package pointer

import (
	"github.com/zodimo/go-compose/compose/ui/input/pointer"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
)

type PointerInputElement struct {
	key   any
	block func(scope pointer.PointerInputScope)
}

var _ modifier.Element = PointerInputElement{}

func (e PointerInputElement) Create() node.Node {
	return NewPointerInputNode(e)
}

func (e PointerInputElement) Update(n node.Node) {
	// Get the node to update
	node := n.(*PointerInputNode)
	// Update the node
	node.key = e.key
	node.block = e.block
}

func (e PointerInputElement) Equals(other modifier.Element) bool {
	otherPointerInputElement, ok := other.(PointerInputElement)
	if !ok {
		return false
	}
	return e.key == otherPointerInputElement.key
}

func (e PointerInputElement) Key() any {
	return e.key
}

func (e PointerInputElement) Block() func(scope pointer.PointerInputScope) {
	return e.block
}
