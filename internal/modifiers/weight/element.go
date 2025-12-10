package weight

var _ Element = (*WeightElement)(nil)

type WeightElement struct {
	weight WeightData
}

func CompareWeight(a, b WeightData) bool {
	return a.Weight == b.Weight
}

// Create creates a new Chain Node instance
func (e WeightElement) Create() Node {
	//chainNode
	return NewWeightNode(e)

}

// Update updates an existing Chain node for efficiency
func (e WeightElement) Update(node Node) {
	if node == nil {
		panic("node cannot be nil")
	}

	bn := node.(WeightNode)
	bn.weight = e.weight

}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (se WeightElement) Equals(other Element) bool {
	if otherElement, ok := other.(WeightElement); ok {
		return CompareWeight(se.weight, otherElement.weight)
	}
	return false
}

func (se WeightElement) WeightData() WeightData {
	return se.weight
}
