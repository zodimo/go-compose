package pointer

import (
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
)

type InputBlockerElement struct{}

var _ modifier.Element = InputBlockerElement{}

func (e InputBlockerElement) Create() node.Node {
	return NewInputBlockerNode(e)
}

func (e InputBlockerElement) Update(node node.Node) {
	// No state to update
}

func (e InputBlockerElement) Equals(other modifier.Element) bool {
	_, ok := other.(InputBlockerElement)
	return ok
}
