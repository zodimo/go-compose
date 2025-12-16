package lazy

import (
	"github.com/zodimo/go-compose/compose"
)

// LazyGridScope is a DSL scope for defining grid items in a lazy grid.
type LazyGridScope interface {
	// Item adds a single item to the grid.
	Item(key any, content compose.Composable)
	// Items adds multiple items to the grid using a count and item factory.
	Items(count int, key func(index int) any, itemContent func(index int) compose.Composable)
}

// lazyGridScopeImpl is the implementation of LazyGridScope.
type lazyGridScopeImpl struct {
	items []lazyGridItem
}

// lazyGridItem represents a single item in the grid.
type lazyGridItem struct {
	Key     any
	Content compose.Composable
}

func (s *lazyGridScopeImpl) Item(key any, content compose.Composable) {
	s.items = append(s.items, lazyGridItem{Key: key, Content: content})
}

func (s *lazyGridScopeImpl) Items(count int, key func(index int) any, itemContent func(index int) compose.Composable) {
	for i := 0; i < count; i++ {
		var k any
		if key != nil {
			k = key(i)
		}
		s.items = append(s.items, lazyGridItem{Key: k, Content: itemContent(i)})
	}
}
