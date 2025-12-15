package navigation

import (
	"fmt"
	"github.com/zodimo/go-compose/pkg/api"
)

func NavHost(
	navController *NavController,
	startDestination string,
	builder func(*NavGraphBuilder),
) api.Composable {
	return func(c api.Composer) api.Composer {
		graphBuilder := NewNavGraphBuilder()
		builder(graphBuilder)

		stack := navController.backStack.Get().([]BackStackEntry)
		if len(stack) == 0 {
			// Initialize with startDestination
            // We use Navigate, but we must be careful about side effects during composition.
            // Since this only happens when stack is empty (initial state), it should be safe enough
            // provided the framework handles state updates.
            // The update will trigger a recompose.
			navController.Navigate(startDestination)
            // Re-fetch stack after update to ensure we render the frame correctly if synchronous
            stack = navController.backStack.Get().([]BackStackEntry)
		}

		currentEntry := navController.CurrentEntry()
		if currentEntry == nil {
			return c
		}

		destination, ok := graphBuilder.destinations[currentEntry.Route]
		if !ok {
            // Log warning or handle error?
            // For now, print to stdout as we don't have a logger
            fmt.Printf("Warning: Destination not found for route: %s\n", currentEntry.Route)
			return c
		}

		// We invoke the destination composable.
        // It modifies the current composer.
		return destination(c)
	}
}
