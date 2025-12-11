package layoutnode

import (
	"github.com/zodimo/go-maybe"
)

var _ LayoutNode = (*layoutNode)(nil)

type layoutNode struct {
	id                     NodeID
	key                    string
	slots                  Slots
	memo                   Memo
	state                  PersistentState
	children               []LayoutNode
	modifier               Modifier
	innerWidgetConstructor LayoutNodeWidgetConstructor
	layoutResult           maybe.Maybe[LayoutResult]
	idManager              IdentityManager
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

func (n *layoutNode) WithSlotsAssoc(k string, v any) LayoutNode {
	n.slots = n.slots.Assoc(k, v)
	return n
}

func (n *layoutNode) FindSlot(key string) maybe.Maybe[any] {
	slot, ok := n.slots.Find(key)
	if !ok {
		return maybe.None[any]()
	}
	return maybe.Some(slot)
}

func (c *layoutNode) GenerateID() Identifier {
	return c.idManager.GenerateID()
}

func (c *layoutNode) ResetIdentifierKeyCounter() {
	c.idManager.ResetKeyCounter()
}

// Remember caches a value for the current composition run.
// The cache lives in Composer.memo and is discarded on recompose.
func (c *layoutNode) Remember(key string, calc func() any) any {
	if v, ok := c.memo.Find(key); ok {
		return v
	}
	v := calc()
	c.memo = c.memo.Assoc(key, v)
	return v
}

// State creates a MutableValue from the persistent state.
// In a real runtime this would be a Snapshot with observers.
func (c *layoutNode) State(key string, initial func() any) MutableValue {
	return c.state.GetState(key, initial)
}

func (n *layoutNode) GetWidget() GioLayoutWidget {
	panic("LayoutNode GetWidget should not be called")
}

func (n *layoutNode) SetWidgetConstructor(constructor LayoutNodeWidgetConstructor) {
	n.innerWidgetConstructor = constructor
}

func (n *layoutNode) GetWidgetConstructor() LayoutNodeWidgetConstructor {
	return n.innerWidgetConstructor
}

func (n *layoutNode) SetLayoutResult(result LayoutResult) {
	n.layoutResult = maybe.Some(result)

}
func (n *layoutNode) GetLayoutResult() maybe.Maybe[LayoutResult] {
	return n.layoutResult
}

func (n *layoutNode) Layout(gtx LayoutContext) LayoutDimensions {
	panic("layout directly on layoutnode not allowed")
}

func (n *layoutNode) Draw(gtx LayoutContext) DrawOp {
	panic("Draw directly on layoutnode not allowed")

}

//////////////////

// func (ln *layoutNode) Layout(gtx LayoutContext) LayoutDimensions {
// 	//render based on elements and
// }
