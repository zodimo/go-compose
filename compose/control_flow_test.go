package compose_test

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
)

// MockComposable records whether it was called
func MockComposable(called *bool) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		*called = true
		return c
	}
}

func TestConditionalComposables(t *testing.T) {
	// Setup composer
	mockStore := store.NewPersistentState(map[string]state.MutableValue{})
	c := compose.NewComposer(mockStore)

	t.Run("If", func(t *testing.T) {
		t.Run("True condition", func(t *testing.T) {
			trueCalled := false
			falseCalled := false

			// Must execute the returned composable
			c.If(true, MockComposable(&trueCalled), MockComposable(&falseCalled))(c)

			if !trueCalled {
				t.Error("Expected true branch to be called")
			}
			if falseCalled {
				t.Error("Expected false branch not to be called")
			}
		})

		t.Run("False condition", func(t *testing.T) {
			trueCalled := false
			falseCalled := false

			c.If(false, MockComposable(&trueCalled), MockComposable(&falseCalled))(c)

			if trueCalled {
				t.Error("Expected true branch not to be called")
			}
			if !falseCalled {
				t.Error("Expected false branch to be called")
			}
		})
	})

	t.Run("When", func(t *testing.T) {
		t.Run("True condition", func(t *testing.T) {
			called := false
			c.When(true, MockComposable(&called))(c)
			if !called {
				t.Error("Expected branch to be called")
			}
		})

		t.Run("False condition", func(t *testing.T) {
			called := false
			c.When(false, MockComposable(&called))(c)
			if called {
				t.Error("Expected branch not to be called")
			}
		})
	})

	t.Run("Else", func(t *testing.T) {
		t.Run("True condition (Else skipped)", func(t *testing.T) {
			called := false
			// Else(true) means the condition was true, so we do NOT execute the else block
			c.Else(true, MockComposable(&called))(c)
			if called {
				t.Error("Expected branch not to be called")
			}
		})

		t.Run("False condition (Else executed)", func(t *testing.T) {
			called := false
			// Else(false) means the condition was false, so we DO execute the else block
			c.Else(false, MockComposable(&called))(c)
			if !called {
				t.Error("Expected branch to be called")
			}
		})
	})

	t.Run("Range", func(t *testing.T) {
		count := 5
		calledCount := 0

		c.Range(count, func(i int) compose.Composable {
			return func(c compose.Composer) compose.Composer {
				calledCount++
				return c
			}
		})(c)

		if calledCount != count {
			t.Errorf("Expected Range to call %d times, but called %d", count, calledCount)
		}
	})

	t.Run("Key", func(t *testing.T) {
		called := false
		key := "my-key"

		// We execute Key, which should start a block, execute content, and end block
		c.Key(key, func(c compose.Composer) compose.Composer {
			called = true
			return c
		})(c)

		if !called {
			t.Error("Expected content inside Key to be called")
		}
	})
}
