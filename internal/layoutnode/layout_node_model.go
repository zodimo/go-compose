package layoutnode

import (
	"go-compose-dev/internal/immap"
)

var _ LayoutNode = (*layoutNode)(nil)

type layoutNode struct {
	id                     NodeID
	key                    string
	slots                  immap.ImmutableMap[Element] // carry the elements
	children               []LayoutNode
	modifier               Modifier
	innerWidget            GioLayoutWidget
	innerWidgetConstructor LayoutNodeWidgetConstructor
}

// Node
func (ln *layoutNode) GetID() NodeID {
	return ln.id
}

// TreeNode
func (ln *layoutNode) Children() []TreeNode {
	treeNodeChildren := []TreeNode{}
	for _, child := range ln.children {
		treeNodeChildren = append(treeNodeChildren, child)
	}
	return treeNodeChildren
}

func (ln *layoutNode) LayoutNodeChildren() []LayoutNode {
	return ln.children
}

//LayoutNode

func (ln *layoutNode) IsEmpty() bool {
	return len(ln.children) == 0
}
func (ln *layoutNode) Modifier(apply func(modifier Modifier) Modifier) {
	ln.modifier = apply(ln.modifier)
}

func (ln *layoutNode) UnwrapModifier() Modifier {
	return ln.modifier
}

func (ln *layoutNode) WithChildren(children []LayoutNode) LayoutNode {
	ln.children = children
	return ln
}

func (n *layoutNode) WithSlotsAssoc(k string, v Element) LayoutNode {
	n.slots = n.slots.Assoc(k, v)
	return n
}

func (n *layoutNode) GetWidget() GioLayoutWidget {
	if n.innerWidget == nil {
		n.innerWidget = n.innerWidgetConstructor.Make(n)
	}
	return n.innerWidget
}

func (n *layoutNode) SetWidgetConstructor(constructor LayoutNodeWidgetConstructor) {
	n.innerWidgetConstructor = constructor
}

//////////////////

// func (ln *layoutNode) Layout(gtx LayoutContext) LayoutDimensions {
// 	//render based on elements and
// }
