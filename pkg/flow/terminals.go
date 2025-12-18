package flow

import (
	"context"
	"fmt"
)

// ToList collects all values from a flow and returns them as a slice.
func ToList[T any](ctx context.Context, flow Flow[T]) ([]T, error) {
	var results []T
	err := flow.Collect(ctx, func(v T) {
		results = append(results, v)
	})
	return results, err
}

// First returns the first value emitted by the flow and then stops.
func First[T any](ctx context.Context, flow Flow[T]) (T, error) {
	var result T
	var found bool

	// We use a custom context with cancel to stop collection after the first item
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := flow.Collect(subCtx, func(v T) {
		if !found {
			result = v
			found = true
			cancel() // Stop the upstream producer immediately
		}
	})

	if !found && err == nil {
		return result, fmt.Errorf("flow completed without emitting values")
	}
	return result, err
}
