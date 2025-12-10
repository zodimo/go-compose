package modifier

import (
	"go-compose-dev/internal/immap"

	"github.com/zodimo/go-maybe"
)

type InspectorInfo struct {
	Name       string
	Properties map[string]any
}

func NewInspectorInfo(name string, props map[string]any) *InspectorInfo {
	return &InspectorInfo{
		Name:       name,
		Properties: props,
	}
}

var _ Modifier = (*emptyModifier)(nil)

type emptyModifier struct{}

func (e emptyModifier) Then(other Modifier) Modifier {
	return other
}

func (me emptyModifier) AsChain() ModifierChain {
	return NewChain(me, nil)
}

var _ ModifierElement = (*modifier)(nil)

type modifier struct {
	element Element
}

//Modifier Interface

// convert to chain
func (me modifier) Then(other Modifier) Modifier {
	if otherChain, ok := other.(ModifierChain); ok {
		return NewChain(me, otherChain)
	}
	otherChain := NewChain(other, nil)
	return NewChain(me, otherChain)
}

func (me modifier) AsChain() ModifierChain {
	return NewChain(me, nil)
}

//Element Interface

// Create creates a new Chain Node instance
func (me modifier) Create() Node {
	return me.element.Create()
}

// Update updates an existing Chain node for efficiency
func (me modifier) Update(node Node) {
	me.element.Update(node)
}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (me modifier) Equals(other Element) bool {
	return me.element.Equals(other)
}

var _ InspectableModifier = (*inspectableModifier)(nil)

type inspectableModifier struct {
	Modifier
	inspectorInfo *InspectorInfo
}

func (im inspectableModifier) InspectorInfo() *InspectorInfo {
	return im.inspectorInfo
}

func (im inspectableModifier) Unwrap() Modifier {
	return im.Modifier
}

var _ ElementStore = (*elementStore)(nil)

type elementStore struct {
	store immap.ImmutableMap[Element]
}

func (es elementStore) SetElement(key string, element Element) ElementStore {
	return elementStore{
		store: es.store.Assoc(key, element),
	}
}

func (es elementStore) GetElement(key string) maybe.Maybe[Element] {
	element, ok := es.store.Find(key)
	if !ok {
		return maybe.None[Element]()
	}
	return maybe.Some(element)
}
