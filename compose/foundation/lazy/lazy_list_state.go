package lazy

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/widgetstate"

	"gioui.org/layout"
	"gioui.org/widget"
)

type LazyListState struct {
	List *widgetstate.List
}

func NewLazyListState() *LazyListState {
	return &LazyListState{
		List: &widgetstate.List{
			L: widget.List{
				List: layout.List{
					Axis: layout.Vertical,
				},
			},
		},
	}
}

func RememberLazyListState(c compose.Composer) *LazyListState {
	id := c.GenerateID()
	key := fmt.Sprintf("lazyListState-%v", id)
	return c.State(key, func() any {
		return NewLazyListState()
	}).Get().(*LazyListState)
}
