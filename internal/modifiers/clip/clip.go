package clip

type ClipData struct {
	Shape        Shape
	ClipToBounds bool
}

type ClipElement struct {
	clipData ClipData
}

// Create creates a new Chain Node instance
func (e ClipElement) Create() Node {
	//chainNode
	return NewClipNode(e)

}

// Update updates an existing Chain node for efficiency
func (e ClipElement) Update(node Node) {
	if node == nil {
		panic("node cannot be nil")
	}

	n := node.(ClipNode)
	n.clipData = e.clipData

}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (e ClipElement) Equals(other Element) bool {
	return false
}

func (e ClipElement) ClipData() ClipData {
	return e.clipData
}
