package background

import (
	"image/color"
)

type BackgroundData struct {
	Color color.Color
	Shape Shape
}

func CompareBackground(a, b BackgroundData) bool {
	return a.Color == b.Color && a.Shape == b.Shape
}

var _ Element = (*BackgroundElement)(nil)

// Hold the behavior
type BackgroundElement struct {
	background BackgroundData
}

// Create creates a new Chain Node instance
func (be BackgroundElement) Create() Node {
	//chainNode
	return NewBackGroundNode(be.background)

}

// Update updates an existing Chain node for efficiency
func (be BackgroundElement) Update(node Node) {
	if node == nil {
		panic("node cannot be nil")
	}

	bn := node.(BackgroundNode)
	bn.background = be.background

}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (be BackgroundElement) Equals(other Element) bool {
	if otherElement, ok := other.(BackgroundElement); ok {
		return CompareBackground(be.background, otherElement.background)
	}
	return false
}

func (be BackgroundElement) Background() BackgroundData {
	return be.background
}
