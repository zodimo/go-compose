package modifier

var _ ModifierChain = (*chain)(nil)

type chain struct {
	head Modifier
	tail ModifierChain
}

func (c chain) Then(next Modifier) Modifier {
	return c.then(next)
}

func (c chain) then(next Modifier) ModifierChain {
	if next == nil {
		return c
	}

	var nextModifierChain ModifierChain
	if mc, ok := next.(ModifierChain); ok {
		nextModifierChain = mc
	} else {
		nextModifierChain = NewChain(next, nil)
	}

	if c.head == nil {
		return nextModifierChain
	}

	head := c.head
	if c.tail == nil {
		return NewChain(head, nextModifierChain)
	}

	return NewChain(head, c.tail.(chain).then(nextModifierChain))
}

// Any checks if any element matches the predicate (left to right)
func (c chain) Any(predicate func(Modifier) bool) bool {
	if c.head == nil {
		return false
	}
	if predicate(c.head) {
		return true
	}
	if c.tail != nil {
		return c.tail.Any(predicate)
	}
	return false
}

// Chain
func (c chain) Fold(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{} {
	return c.FoldIn(initial, operation)
}

// Fold Right to left
func (c chain) FoldIn(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{} {

	if c.head == nil {
		return initial
	}
	// First process the tail (right side)
	if c.tail != nil {
		initial = c.tail.FoldIn(initial, operation)
	}
	// Then process the head
	return operation(initial, c.head)
}

// Fold Left to Right
func (c chain) FoldOut(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{} {
	if c.head == nil {
		return initial
	}
	// First process the head
	initial = operation(initial, c.head)
	// Then process the tail (right side)
	if c.tail != nil {
		initial = c.tail.FoldOut(initial, operation)
	}
	return initial
}

func (c chain) AsChain() ModifierChain {
	return c
}

// Constructor
func NewChain(head Modifier, tail ModifierChain) ModifierChain {
	return chain{
		head: head,
		tail: tail,
	}
}
