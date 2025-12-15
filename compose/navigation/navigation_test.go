package navigation

import (
	"testing"

	"github.com/zodimo/go-compose/pkg/api"
)

// Mock MutableValue for testing
type mockMutableValue struct {
	value interface{}
}

func (m *mockMutableValue) Get() interface{} {
	return m.value
}

func (m *mockMutableValue) Set(value interface{}) {
	m.value = value
}

func TestNavController(t *testing.T) {
	// Initialize NavController with empty stack
	backStack := &mockMutableValue{value: []BackStackEntry{}}
	nc := NewNavController(backStack)

	// Test Initial State
	if nc.CurrentEntry() != nil {
		t.Error("Expected nil CurrentEntry for empty stack")
	}

	// Test Navigate
	nc.Navigate("home")
	if nc.CurrentEntry().Route != "home" {
		t.Errorf("Expected route 'home', got '%s'", nc.CurrentEntry().Route)
	}

	nc.Navigate("details")
	if nc.CurrentEntry().Route != "details" {
		t.Errorf("Expected route 'details', got '%s'", nc.CurrentEntry().Route)
	}

	// Test PopBackStack
	popped := nc.PopBackStack()
	if !popped {
		t.Error("Expected PopBackStack to return true")
	}
	if nc.CurrentEntry().Route != "home" {
		t.Errorf("Expected route 'home' after pop, got '%s'", nc.CurrentEntry().Route)
	}

	// Test PopBackStack to empty
	popped = nc.PopBackStack()
	if !popped {
		t.Error("Expected PopBackStack to return true for last item")
	}
	if nc.CurrentEntry() != nil {
		t.Error("Expected nil CurrentEntry after popping last item")
	}

	// Test PopBackStack on empty
	popped = nc.PopBackStack()
	if popped {
		t.Error("Expected PopBackStack to return false for empty stack")
	}
}

// Simple test for builder
func TestNavGraphBuilder(t *testing.T) {
	builder := NewNavGraphBuilder()

	dummy := func(c api.Composer) api.Composer { return c }

	builder.Composable("route1", dummy)

	if _, ok := builder.destinations["route1"]; !ok {
		t.Error("Expected route1 to be present")
	}
}
