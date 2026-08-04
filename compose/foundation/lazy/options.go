package lazy

import (
	"github.com/zodimo/go-compose/compose/ui"
)

type LazyListOption func(*LazyListOptions)

type LazyListOptions struct {
	Modifier  ui.Modifier
	State     *LazyListState
	Scrollbar bool // enable scrollbar rendering
}

func DefaultLazyListOptions() LazyListOptions {
	return LazyListOptions{
		Modifier:  ui.EmptyModifier,
		State:     nil,
		Scrollbar: true,
	}
}

func WithModifier(m ui.Modifier) LazyListOption {
	return func(o *LazyListOptions) {
		o.Modifier = m
	}
}

func WithState(state *LazyListState) LazyListOption {
	return func(o *LazyListOptions) {
		o.State = state
	}
}

func WithScrollbar(enabled bool) LazyListOption {
	return func(o *LazyListOptions) {
		o.Scrollbar = enabled
	}
}
