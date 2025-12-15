package navigation

import (
	"fmt"
	"time"

	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

type BackStackEntry struct {
	Route string
	ID    string
}

type NavController struct {
	backStack state.MutableValue
}

func NewNavController(backStack state.MutableValue) *NavController {
	return &NavController{backStack: backStack}
}

func RememberNavController(c api.Composer) *NavController {
	backStack := c.State("nav_backstack", func() any {
		return []BackStackEntry{}
	})

	nc := c.Remember("nav_controller", func() any {
		return NewNavController(backStack)
	})

	return nc.(*NavController)
}

func (nc *NavController) Navigate(route string) {
	stack := nc.backStack.Get().([]BackStackEntry)
    // simple ID generation using time
    id := fmt.Sprintf("%s-%d", route, time.Now().UnixNano())
	stack = append(stack, BackStackEntry{Route: route, ID: id})
	nc.backStack.Set(stack)
}

func (nc *NavController) PopBackStack() bool {
	stack := nc.backStack.Get().([]BackStackEntry)
	if len(stack) <= 0 {
		return false
	}
    // Remove the last item
    stack = stack[:len(stack)-1]
    nc.backStack.Set(stack)
    return true
}

func (nc *NavController) CurrentEntry() *BackStackEntry {
	stack := nc.backStack.Get().([]BackStackEntry)
	if len(stack) == 0 {
		return nil
	}
	return &stack[len(stack)-1]
}
