package offset

import (
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
)

type OffsetElement struct {
	data OffsetData
}

func (e *OffsetElement) Create() node.Node {
	return NewOffsetNode(e.data)
}

func (e *OffsetElement) Update(n node.Node) {
	no := n.(*OffsetNode)
	no.data = e.data
}

func (e *OffsetElement) Equals(other modifier.Element) bool {
	o, ok := other.(*OffsetElement)
	if !ok {
		return false
	}
	return e.data == o.data
}
