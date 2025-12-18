package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zodimo/go-compose/pkg/flow"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Initialize with 0
	counter := flow.NewMutableStateFlow(0)

	// Subscriber 1: Log updates
	go counter.Collect(ctx, func(v int) {
		fmt.Printf("[UI Component A] Redrawing with: %d\n", v)
	})

	// Producer: Update state every 500ms
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)
			counter.Emit(i)
		}
	}()

	<-ctx.Done()
}
