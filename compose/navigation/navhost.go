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

		stack := navController.backStack.Get()
		if len(stack) == 0 {
			// Initialize with startDestination
			// We use Navigate, but we must be careful about side effects during composition.
			// Since this only happens when stack is empty (initial state), it should be safe enough
			// provided the framework handles state updates.
			// The update will trigger a recompose.
			navController.Navigate(startDestination)
			// Re-fetch stack after update to ensure we render the frame correctly if synchronous
			stack = navController.backStack.Get()
		}

		currentEntry := navController.CurrentEntry()
		if currentEntry == nil {
			return c
		}

		// Find matching destination and extract arguments
		composableWithArgs, args, ok := graphBuilder.findDestination(currentEntry.Route)
		if !ok {
			panic(fmt.Sprintf("Destination not found for route: %s\n", currentEntry.Route))
		}

		// Create entry with extracted arguments
		entryWithArgs := NewBackStackEntry(currentEntry.Route, WithArguments(args))
		entryWithArgs.ID = currentEntry.ID // Preserve original ID

		// We invoke the destination composable with the entry containing arguments.
		return c.Key(entryWithArgs.Route, composableWithArgs(&entryWithArgs))(c)
	}
}
