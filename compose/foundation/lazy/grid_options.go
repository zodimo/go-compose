package lazy

import (
	"github.com/zodimo/go-compose/internal/modifier"
)

// LazyGridOption is a functional option for configuring lazy grids.
type LazyGridOption func(*LazyGridOptions)

// LazyGridOptions holds configuration for a lazy grid.
type LazyGridOptions struct {
	Modifier modifier.Modifier
	State    *LazyGridState
}

// DefaultLazyGridOptions returns the default options for a lazy grid.
func DefaultLazyGridOptions() LazyGridOptions {
	return LazyGridOptions{
		Modifier: modifier.EmptyModifier,
		State:    nil,
	}
}

// WithGridModifier applies a modifier to the grid.
func WithGridModifier(m modifier.Modifier) LazyGridOption {
	return func(o *LazyGridOptions) {
		o.Modifier = m
	}
}

// WithGridState sets the grid state for scroll position management.
func WithGridState(state *LazyGridState) LazyGridOption {
	return func(o *LazyGridOptions) {
		o.State = state
	}
}
