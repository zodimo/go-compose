package scale

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/modifier"
)

type ScaleElement struct {
	data ScaleData
}

func (e *ScaleElement) Create() node.Node {
	return NewScaleNode(e.data)
}

func (e *ScaleElement) Update(n node.Node) {
	no := n.(*ScaleNode)
	no.data = e.data
}

func (e *ScaleElement) Equals(other modifier.Element) bool {
	o, ok := other.(*ScaleElement)
	if !ok {
		return false
	}
	return e.data == o.data
}
