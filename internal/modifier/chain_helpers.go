package modifier

// Fold Alias for FoldIn, FoldIn processes elements from right to left (tail to head)
func Fold[T any](c ModifierChain, initial T, operation func(T, Modifier) T) T {
	return FoldIn(c, initial, operation)
}

// FoldIn processes elements from right to left (tail to head)
func FoldIn[T any](c ModifierChain, initial T, operation func(T, Modifier) T) T {
	if c == nil {
		return initial
	}
	// fmt.Println(DebugModifier(c))
	chain := c.(chain)
	// First process the tail (right side)
	if chain.tail != nil {
		initial = FoldIn(chain.tail, initial, operation)
	}
	// Then process the head
	return operation(initial, chain.head)
}

// FoldOut processes elements from left to right (head to tail)
func FoldOut[T any](c ModifierChain, initial T, operation func(T, Modifier) T) T {
	if c == nil {
		return initial
	}

	chain := c.(chain)
	// First process the head
	initial = operation(initial, chain.head)
	// Then process the tail (right side)
	if chain.tail != nil {
		initial = FoldOut(chain.tail, initial, operation)
	}
	return initial
}
