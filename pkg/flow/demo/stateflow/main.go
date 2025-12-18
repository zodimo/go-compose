package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zodimo/go-compose/pkg/flow"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Define the Flow (Blueprint)
	// This code does NOT run yet.
	numbersFlow := flow.NewFlow(func(ctx context.Context, emit func(int)) error {
		for i := 1; i <= 3; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Printf("[Producer] Emitting %d\n", i)
				emit(i)
				time.Sleep(500 * time.Millisecond)
			}
		}
		return nil
	})

	// 2. First Collection
	fmt.Println("Starting Collection A...")
	numbersFlow.Collect(ctx, func(v int) {
		fmt.Printf("-> Collector A received: %d\n", v)
	})

	// 3. Second Collection (Starts from the beginning)
	fmt.Println("\nStarting Collection B...")
	numbersFlow.Collect(ctx, func(v int) {
		fmt.Printf("-> Collector B received: %d\n", v)
	})
}
