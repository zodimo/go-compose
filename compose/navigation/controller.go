package navigation

import (
	"fmt"
	"time"

	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-maybe"
)

type BackStackEntry struct {
	Route     string
	ID        string
	Arguments maybe.Maybe[NavArguments]
}

// BackStackEntryOption configures a BackStackEntry
type BackStackEntryOption func(*BackStackEntry)

// WithArguments sets the arguments for a BackStackEntry
func WithArguments(args NavArguments) BackStackEntryOption {
	return func(e *BackStackEntry) {
		e.Arguments = maybe.Some(args)
	}
}

// NewBackStackEntry creates a new BackStackEntry with the given route
func NewBackStackEntry(route string, opts ...BackStackEntryOption) BackStackEntry {
	entry := BackStackEntry{
		Route:     route,
		ID:        fmt.Sprintf("%s-%d", route, time.Now().UnixNano()),
		Arguments: maybe.None[NavArguments](),
	}
	for _, opt := range opts {
		opt(&entry)
	}
	return entry
}

type NavController struct {
	backStack state.MutableValueTyped[[]BackStackEntry]
}

func NewNavController(backStack state.MutableValueTyped[[]BackStackEntry]) *NavController {
	return &NavController{backStack: backStack}
}

func RememberNavController(c api.Composer) *NavController {
	backStack := state.MustState(c, "nav_backstack", func() []BackStackEntry {
		return []BackStackEntry{}
	})
	return NewNavController(backStack)
}

// Navigate to a route. Arguments are extracted in NavHost via pattern matching.
func (nc *NavController) Navigate(route string) {
	stack := nc.backStack.Get()
	stack = append(stack, NewBackStackEntry(route))
	nc.backStack.Set(stack)
}

func (nc *NavController) PopBackStack() bool {
	stack := nc.backStack.Get()
	if len(stack) <= 0 {
		return false
	}
	// Remove the last item
	stack = stack[:len(stack)-1]
	nc.backStack.Set(stack)
	return true
}

func (nc *NavController) CurrentEntry() *BackStackEntry {
	stack := nc.backStack.Get()
	if len(stack) == 0 {
		return nil
	}
	return &stack[len(stack)-1]
}
