package modifier

import "go-compose-dev/internal/immap"

// EmptyModifier is the empty modifier that does nothing.
var EmptyModifier Modifier = &emptyModifier{}
var EmptyElementStore ElementStore = NewElementStore()

func NewElement(
	createFunc func() Node,
	updateFunc func(Node),
	equalsFunc func(Element) bool,
) Element {
	if createFunc == nil {
		panic("createFunc cannot be nil")
	}
	if updateFunc == nil {
		panic("updateFunc cannot be nil")
	}
	if equalsFunc == nil {
		panic("equalsFunc cannot be nil")
	}

	return element{
		createFunc: createFunc,
		updateFunc: updateFunc,
		equalsFunc: equalsFunc,
	}
}

func NewInspectableModifier(m Modifier, inspectorInfo *InspectorInfo) InspectableModifier {
	if m == nil {
		panic("modifier cannot be nil")
	}
	if inspectorInfo == nil {
		panic("inspectorInfo cannot be nil")
	}
	return &inspectableModifier{
		Modifier:      m,
		inspectorInfo: inspectorInfo,
	}
}

func NewModifier(element Element) Modifier {
	return &modifier{
		element: element,
	}
}

func NewElementStore() ElementStore {
	return elementStore{
		store: immap.EmptyImmutableMap[Element](),
	}
}
