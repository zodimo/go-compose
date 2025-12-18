package flow

import "context"

// Flow represents a cold stream of data.
type Flow[T any] interface {
	Collect(ctx context.Context, collector func(T)) error
}

// flowImpl implements the Flow interface.
type flowImpl[T any] struct {
	// block is the producer function that defines how values are emitted.
	block func(ctx context.Context, emit func(T)) error
}

// NewFlow creates a cold flow. The 'block' isn't executed until Collect is called.
func NewFlow[T any](block func(ctx context.Context, emit func(T)) error) Flow[T] {
	return &flowImpl[T]{block: block}
}

// Collect triggers the execution of the flow.
func (f *flowImpl[T]) Collect(ctx context.Context, collector func(T)) error {
	// We pass the collector directly as the 'emit' function.
	return f.block(ctx, collector)
}

func Map[T, R any](upstream Flow[T], transform func(T) R) Flow[R] {
	return NewFlow(func(ctx context.Context, emit func(R)) error {
		return upstream.Collect(ctx, func(value T) {
			emit(transform(value))
		})
	})
}
