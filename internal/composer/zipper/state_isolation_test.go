package zipper

import (
	"testing"

	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/store"
)

// TestIfStateBranchIsolation verifies that state created within different branches
// of c.If are isolated from each other (different keys), preventing type mismatches
// when switching between branches.
func TestIfStateBranchIsolation(t *testing.T) {
	// Create persistent state store
	ps := store.NewPersistentState()

	// Type A state (simulating LoginState)
	type TypeA struct {
		Name string
	}

	// Type B state (simulating ControlPanelState)
	type TypeB struct {
		Value int
	}

	// First composition: condition = false, use TypeA in ifFalse branch
	func() {
		c := NewComposer(api.ComposerWithStore(ps))
		c.StartBlock("root")

		condition := false
		c.If(condition,
			// ifTrue branch (TypeB) - NOT executed
			func(c Composer) Composer {
				stateValue := c.State("inner_state", func() any {
					return &TypeB{Value: 42}
				})
				_ = stateValue.Get().(*TypeB)
				return c
			},
			// ifFalse branch (TypeA) - executed
			func(c Composer) Composer {
				stateValue := c.State("inner_state", func() any {
					return &TypeA{Name: "test"}
				})
				result := stateValue.Get().(*TypeA)
				if result.Name != "test" {
					t.Errorf("Expected Name='test', got %q", result.Name)
				}
				return c
			},
		)(c)

		c.EndBlock()
	}()

	// Second composition: condition = true, use TypeB in ifTrue branch
	// This should NOT panic because the prefixes are different
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic during second composition (condition=true): %v", r)
			}
		}()

		c := NewComposer(api.ComposerWithStore(ps))
		c.StartBlock("root")

		condition := true
		c.If(condition,
			// ifTrue branch (TypeB) - executed
			func(c Composer) Composer {
				stateValue := c.State("inner_state", func() any {
					return &TypeB{Value: 42}
				})
				result := stateValue.Get().(*TypeB)
				if result.Value != 42 {
					t.Errorf("Expected Value=42, got %d", result.Value)
				}
				return c
			},
			// ifFalse branch (TypeA) - NOT executed
			func(c Composer) Composer {
				stateValue := c.State("inner_state", func() any {
					return &TypeA{Name: "test"}
				})
				_ = stateValue.Get().(*TypeA)
				return c
			},
		)(c)

		c.EndBlock()
	}()
}

// TestKeyPrefixStackScoping verifies that nested c.Key calls create properly
// scoped prefixes for all GenerateID calls within.
func TestKeyPrefixStackScoping(t *testing.T) {
	ps := store.NewPersistentState()

	c := NewComposer(api.ComposerWithStore(ps))
	c.StartBlock("root")

	// Collect IDs at different scoping levels
	var outerID, innerID1, innerID2 Identifier

	// Outer ID (no prefix)
	outerID = c.GenerateID()

	// Nested Key scopes
	c.Key("scope_a", func(c Composer) Composer {
		innerID1 = c.GenerateID()

		c.Key("scope_b", func(c Composer) Composer {
			innerID2 = c.GenerateID()
			return c
		})(c)

		return c
	})(c)

	c.EndBlock()

	// Verify IDs are different (scoped)
	if outerID.String() == innerID1.String() {
		t.Errorf("outerID should differ from innerID1 (scoped)")
	}
	if innerID1.String() == innerID2.String() {
		t.Errorf("innerID1 should differ from innerID2 (nested scope)")
	}

	t.Logf("outerID: %s", outerID.String())
	t.Logf("innerID1 (scope_a): %s", innerID1.String())
	t.Logf("innerID2 (scope_a/scope_b): %s", innerID2.String())
}

// TestIfGeneratesDistinctStateKeys verifies that the same state key name
// used in both If branches results in different actual keys.
func TestIfGeneratesDistinctStateKeys(t *testing.T) {
	ps := store.NewPersistentState()

	// Track which states are created
	stateKeys := make(map[string]bool)

	// Wrapper to capture state keys
	originalGetState := ps.GetState
	_ = originalGetState // Note: we can't easily intercept, so we'll test via behavior

	// First: execute ifFalse branch
	func() {
		c := NewComposer(api.ComposerWithStore(ps))
		c.StartBlock("root")

		c.If(false,
			func(c Composer) Composer {
				c.State("shared_key", func() any { return "type_true" })
				return c
			},
			func(c Composer) Composer {
				val := c.State("shared_key", func() any { return "type_false" })
				if val.Get().(string) != "type_false" {
					t.Errorf("Expected 'type_false', got %v", val.Get())
				}
				return c
			},
		)(c)

		c.EndBlock()
	}()

	// Second: execute ifTrue branch - should get its own state
	func() {
		c := NewComposer(api.ComposerWithStore(ps))
		c.StartBlock("root")

		c.If(true,
			func(c Composer) Composer {
				val := c.State("shared_key", func() any { return "type_true" })
				// This should be "type_true" (new state), NOT "type_false" (from other branch)
				if val.Get().(string) != "type_true" {
					t.Errorf("Expected 'type_true' (isolated state), got %v", val.Get())
				}
				return c
			},
			func(c Composer) Composer {
				c.State("shared_key", func() any { return "type_false" })
				return c
			},
		)(c)

		c.EndBlock()
	}()

	_ = stateKeys
}
