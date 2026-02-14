package zipper

import (
	"fmt"
	"strings"
	"time"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
)

var _ Composer = (*composer)(nil)

type pathItem struct {
	parent LayoutNode   // the parent node
	before []LayoutNode // children left of the focus (in order)
	after  []LayoutNode // children right of the focus (reversed)
}

type composer struct {
	focus          LayoutNode // group we are currently inside
	path           []pathItem // how to climb back to root
	memo           Memo       // remember cache for this composition run
	state          PersistentState
	idManager      IdentityManager
	overrideID     *Identifier // single override ID for c.Key (one Key affects one component)
	idPrefixStack  []string    // stack of ID prefixes for scoped identity (used by c.Key)
	locals         map[interface{}]interface{}
	providersStack []map[interface{}]interface{}
	timeNowFunc    func() time.Time
}

func (c *composer) TimeNow() time.Time {
	return c.timeNowFunc()
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
	// Check if there's an override ID from c.Key
	if c.overrideID != nil {
		id := *c.overrideID
		c.overrideID = nil
		return id
	}

	// Generate base ID
	baseID := c.idManager.GenerateID()

	// If we have a prefix stack, create a scoped ID
	if len(c.idPrefixStack) > 0 {
		prefix := strings.Join(c.idPrefixStack, "/")
		return c.idManager.CreateID(prefix + "/" + baseID.String())
	}

	return baseID
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
func (c *composer) Modifier(apply func(modifier ui.Modifier) ui.Modifier) {
	c.focus.Modifier(apply)
}
func (c *composer) ModifierThen(modifier ui.Modifier) Composer {
	c.Modifier(func(modifier ui.Modifier) ui.Modifier {
		return modifier.Then(modifier)
	})
	return c
}

// Remember caches a value for the current composition run.
// The cache lives in Composer.memo and is discarded on recompose.
func (c *composer) Remember(key string, calc func() any) any {
	// Apply prefix stack to the key for proper scoping
	scopedKey := c.scopeKey(key)
	if v, ok := c.memo.Find(scopedKey); ok {
		return v
	}
	v := calc()
	c.memo = c.memo.Assoc(scopedKey, v)
	return v
}

// State creates a MutableValue from the persistent state.
// In a real runtime this would be a Snapshot with observers.
func (c *composer) State(key string, initial func() any, options ...StateOption) MutableValue {
	// Apply prefix stack to the key for proper scoping
	scopedKey := c.scopeKey(key)
	return c.state.GetState(scopedKey, initial, options...)
}

// scopeKey prefixes the given key with the current ID prefix stack.
// This ensures state isolation when using c.Key or c.If.
func (c *composer) scopeKey(key string) string {
	if len(c.idPrefixStack) == 0 {
		return key
	}
	prefix := strings.Join(c.idPrefixStack, "/")
	return prefix + "/" + key
}

func (c *composer) WithComposable(composable Composable) Composer {
	return composable(c)
}

func (c *composer) SetWidgetConstructor(constructor layoutnode.LayoutNodeWidgetConstructor) {
	c.focus.SetWidgetConstructor(constructor)
}

func (c *composer) If(condition bool, ifTrue Composable, ifFalse Composable) Composable {
	// Use stable string keys for branches - the prefix stack handles scoping
	if condition {
		return c.Key("if_true", ifTrue)
	}
	return c.Key("if_false", ifFalse)
}

func (c *composer) IfLazy(condition bool, ifTrue func() Composable, ifFalse func() Composable) Composable {
	// Use stable string keys for branches - the prefix stack handles scoping
	if condition {
		return c.Key("if_true", ifTrue())
	}
	return c.Key("if_false", ifFalse())
}

func (c *composer) When(condition bool, ifTrue Composable) Composable {
	if condition {
		return c.Key("when_true", ifTrue)
	}
	return c.Key("when_false", emptyComposable())
}

func (c *composer) WhenLazy(condition bool, ifTrue func() Composable) Composable {
	if condition {
		return c.Key("when_true", ifTrue())
	}
	return c.Key("when_false", emptyComposable())
}

func (c *composer) Else(condition bool, ifFalse Composable) Composable {
	if condition {
		return c.Key("else_true", emptyComposable())
	}
	return c.Key("else_false", ifFalse)
}

func (c *composer) ElseLazy(condition bool, ifFalse func() Composable) Composable {
	if condition {
		return c.Key("else_true", emptyComposable())
	}
	return c.Key("else_false", ifFalse())
}

func (c *composer) Sequence(contents ...Composable) Composable {
	return Sequence(contents...)
}

func (c *composer) Key(key any, content Composable) Composable {
	// Create a stable ID prefix from the key
	stringKey := fmt.Sprint(key)
	return func(comp Composer) Composer {
		composerImpl := comp.(*composer)

		// Push the key onto the prefix stack - all GenerateID calls within
		// this scope will be prefixed with this key
		composerImpl.idPrefixStack = append(composerImpl.idPrefixStack, stringKey)

		// Compose content with the new prefix scope
		result := content(comp)

		// Pop the prefix from the stack
		composerImpl.idPrefixStack = composerImpl.idPrefixStack[:len(composerImpl.idPrefixStack)-1]

		return result
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

func (c *composer) StartProviders(values []ProvidedValue) Composer {
	c.providersStack = append(c.providersStack, c.locals)
	newLocals := make(map[interface{}]interface{}, len(c.locals)+len(values))
	for k, v := range c.locals {
		newLocals[k] = v
	}
	for _, v := range values {
		newLocals[v.CompositionLocal] = v.Value
	}
	c.locals = newLocals
	return c
}

func (c *composer) EndProviders() Composer {
	if len(c.providersStack) == 0 {
		panic("Unbalanced StartProviders/EndProviders")
	}
	c.locals = c.providersStack[len(c.providersStack)-1]
	c.providersStack = c.providersStack[:len(c.providersStack)-1]
	return c
}

func (c *composer) Consume(key interface{}) interface{} {
	return c.locals[key]
}

func emptyComposable() Composable {
	return func(c Composer) Composer {
		return c
	}
}
