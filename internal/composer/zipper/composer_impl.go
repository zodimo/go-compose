package zipper

import (
	"fmt"

	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

var _ Composer = (*composer)(nil)

type pathItem struct {
	parent LayoutNode   // the parent node
	before []LayoutNode // children left of the focus (in order)
	after  []LayoutNode // children right of the focus (reversed)
}

type composer struct {
	focus     LayoutNode // group we are currently inside
	path      []pathItem // how to climb back to root
	memo      Memo       // remember cache for this composition run
	state     PersistentState
	idManager IdentityManager
}

// Tree Builder operations
func (c *composer) StartBlock(key string) Composer {

	newNode := layoutnode.NewLayoutNode(c.GenerateID(), key, EmptyMemo, EmptyMemo, c.state)

	if c.focus == nil {
		//The Root Node
		// How to Make this Requirement Explicit?
		newNode.ResetIdentifierKeyCounter()
		c.focus = newNode
		return c
	}

	c.path = append(c.path, pathItem{
		parent: c.focus,
		before: c.focus.LayoutNodeChildren(),
		after:  []LayoutNode{},
	})
	c.focus = newNode
	return c

}
func (c *composer) EndBlock() Composer {
	return c.up()
}

// Root climbs the zipper to the top and returns the finished tree.
// This must be called after all groups are closed.
func (c *composer) Build() LayoutNode {
	for len(c.path) > 0 {
		c.up()
	}
	if c.focus == nil {
		panic("No root layout node found")
	}
	return c.focus
}

func (c *composer) EmitSlot(k string, v any) Composer {
	c.focus = c.focus.WithSlotsAssoc(k, v)
	return c
} // slot is a property on the layoutNode

// Tree navigation
func (c *composer) up() Composer {
	if len(c.path) == 0 {
		return c // already at root
	}
	top := c.path[len(c.path)-1]
	c.path = c.path[:len(c.path)-1]
	finished := c.focus

	var children []LayoutNode
	children = append(children, top.before...)
	children = append(children, finished)

	for i := len(top.after) - 1; i >= 0; i-- { // prepend right siblings in order
		children = append(children, top.after[i])
	}
	parent := top.parent.WithChildren(children)
	c.focus = parent
	return c
}

func (c *composer) GenerateID() Identifier {
	return c.idManager.GenerateID()
}
func (c *composer) GetID() Identifier {
	return c.focus.GetID()
}
func (c *composer) GetPath() NodePath {
	nodeIds := []node.NodeID{}
	for _, pathItem := range c.path {
		nodeIds = append(nodeIds, pathItem.parent.GetID())
	}
	return node.NewNodePath(nodeIds)

}
func (c *composer) Modifier(apply func(modifier Modifier) Modifier) {
	c.focus.Modifier(apply)
}
func (c *composer) ModifierThen(modifier Modifier) Composer {
	c.Modifier(func(modifier Modifier) Modifier {
		return modifier.Then(modifier)
	})
	return c
}

// Remember caches a value for the current composition run.
// The cache lives in Composer.memo and is discarded on recompose.
func (c *composer) Remember(key string, calc func() any) any {
	if v, ok := c.memo.Find(key); ok {
		return v
	}
	v := calc()
	c.memo = c.memo.Assoc(key, v)
	return v
}

// State creates a MutableValue from the persistent state.
// In a real runtime this would be a Snapshot with observers.
func (c *composer) State(key string, initial func() any) MutableValue {
	return c.state.GetState(key, initial)
}

func (c *composer) WithComposable(composable Composable) Composer {
	return composable(c)
}

func (c *composer) SetWidgetConstructor(constructor layoutnode.LayoutNodeWidgetConstructor) {
	c.focus.SetWidgetConstructor(constructor)
}

func (c *composer) If(condition bool, ifTrue Composable, ifFalse Composable) Composable {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func (c *composer) When(condition bool, ifTrue Composable) Composable {
	if condition {
		return ifTrue
	}
	return emptyComposable()
}
func (c *composer) Else(condition bool, ifFalse Composable) Composable {
	if condition {
		return emptyComposable()
	}
	return ifFalse
}

func (c *composer) Sequence(contents ...Composable) Composable {
	return Sequence(contents...)
}

func (c *composer) Key(key any, content Composable) Composable {
	return func(c Composer) Composer {
		// We stringify the key to be used as an ID for the block
		k := fmt.Sprint(key)
		c.StartBlock(k)
		c = content(c)
		c.EndBlock()
		return c
	}
}

func (c *composer) Range(count int, fn func(int) Composable) Composable {
	return func(c Composer) Composer {
		for i := 0; i < count; i++ {
			c = fn(i)(c)
		}
		return c
	}
}

func emptyComposable() Composable {
	return func(c Composer) Composer {
		return c
	}
}
