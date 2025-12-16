package compose

import (
	"github.com/zodimo/go-compose/internal/composer/zipper"
	"github.com/zodimo/go-compose/internal/sequence"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

type Composable = api.Composable
type Composer = api.Composer

// NewComposer creates a new Composer instance initialized with the given persistent state store.
// The returned Composer is the entry point for building and managing the composition tree.
//
// Parameters:
//   - store: A PersistentState implementation that manages the state across recompositions.
//
// Returns:
//   - A new Composer instance.
func NewComposer(store state.PersistentState) Composer {
	return zipper.NewComposer(store)
}

// Sequence is a convenience function for combining multiple Composables into a single one.
// It delegates to the internal sequence implementation.
//
// Deprecated: Use c.Sequence(...) method on the Composer interface instead for better fluency.
var Sequence = sequence.Sequence

// Id returns a Composable that does nothing and returns the Composer as is.
// It is useful as a placeholder or identity operation in functional composition patterns.
func Id() Composable {
	return func(c Composer) Composer {
		return c
	}
}
