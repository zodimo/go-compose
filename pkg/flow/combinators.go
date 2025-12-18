package flow

import (
	"context"
	"sync"
)

func Combine[T1, T2, R any](
	ctx context.Context,
	flow1 Flow[T1],
	flow2 Flow[T2],
	transform func(T1, T2) R,
) Flow[R] {
	return NewFlow(func(ctx context.Context, emit func(R)) error {
		var v1 T1
		var v2 T2
		var hasV1, hasV2 bool
		mu := sync.Mutex{}

		// Internal function to check if we can emit the combined result
		tryEmit := func() {
			mu.Lock()
			defer mu.Unlock()
			if hasV1 && hasV2 {
				emit(transform(v1, v2))
			}
		}

		// Launch collectors for both flows in parallel
		go flow1.Collect(ctx, func(v T1) {
			mu.Lock()
			v1 = v
			hasV1 = true
			mu.Unlock()
			tryEmit()
		})

		return flow2.Collect(ctx, func(v T2) {
			mu.Lock()
			v2 = v
			hasV2 = true
			mu.Unlock()
			tryEmit()
		})
	})
}

func Zip[T1, T2, R any](
	ctx context.Context,
	flow1 Flow[T1],
	flow2 Flow[T2],
	transform func(T1, T2) R,
) Flow[R] {
	return NewFlow(func(ctx context.Context, emit func(R)) error {
		ch1 := make(chan T1)
		ch2 := make(chan T2)

		go flow1.Collect(ctx, func(v T1) { ch1 <- v })
		go flow2.Collect(ctx, func(v T2) { ch2 <- v })

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case v1 := <-ch1:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case v2 := <-ch2:
					emit(transform(v1, v2))
				}
			case v2 := <-ch2:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case v1 := <-ch1:
					emit(transform(v1, v2))
				}
			}
		}
	})
}
