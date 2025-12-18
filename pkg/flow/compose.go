package flow

import (
	"context"

	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
)

// CollectAsState collects values from a Flow and represents it as a State.
// The initial value is used until the first value is emitted by the flow.
func CollectAsState[T any](c api.Composer, flow Flow[T], initial T) state.TypedValue[T] {
	// 1. Create a state to hold the latest value
	stateValue := store.StateUnsafe[T](c, "flow_state", func() T { return initial })

	// 2. Launch a side-effect to collect the flow
	// We use the flow itself as a key so if the flow instance changes, we resubscribe
	effect.LaunchedEffect(func(ctx context.Context) {
		flow.Collect(ctx, func(value T) {
			stateValue.Set(value)
		})
	}, flow)(c)

	return stateValue
}

// CollectStateFlowAsState collects values from a StateFlow.
// It uses `flow.Value()` as the initial value.
func CollectStateFlowAsState[T any](c api.Composer, flow StateFlow[T]) state.TypedValue[T] {
	return CollectAsState(c, flow, flow.Value())
}
