package lazy

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"

	"gioui.org/layout"
	"gioui.org/widget"
)

// LazyGridState holds the state for a lazy grid, including scroll position.
type LazyGridState struct {
	List widget.List
}

// NewLazyGridState creates a new LazyGridState with default configuration.
func NewLazyGridState() *LazyGridState {
	return &LazyGridState{
		List: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
}

// RememberLazyGridState creates or retrieves a remembered LazyGridState.
func RememberLazyGridState(c compose.Composer) *LazyGridState {
	id := c.GenerateID()
	key := fmt.Sprintf("lazyGridState-%v", id)
	return c.State(key, func() any {
		return NewLazyGridState()
	}).Get().(*LazyGridState)
}
