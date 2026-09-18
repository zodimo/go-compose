package formengine

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestControl_SetAndLifecycle(t *testing.T) {
	t.Run("Set updates value and marks dirty", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		c.Set("b")
		if c.Value() != "b" {
			t.Fatalf("Value() = %q, want %q", c.Value(), "b")
		}
		if !c.IsDirty() {
			t.Fatal("IsDirty() = false after Set, want true")
		}
		if c.IsPristine() {
			t.Fatal("IsPristine() = true after Set, want false")
		}
	})

	t.Run("MarkAsPristine clears dirty", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		c.Set("b")
		c.MarkAsPristine()
		if c.IsDirty() {
			t.Fatal("IsDirty() = true after MarkAsPristine, want false")
		}
		if !c.IsPristine() {
			t.Fatal("IsPristine() = false after MarkAsPristine, want true")
		}
	})

	t.Run("MarkAsTouched and MarkAsUntouched", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		if c.IsTouched() {
			t.Fatal("fresh control should not be touched")
		}
		c.MarkAsTouched()
		if !c.IsTouched() {
			t.Fatal("IsTouched() = false after MarkAsTouched, want true")
		}
		c.MarkAsUntouched()
		if c.IsTouched() {
			t.Fatal("IsTouched() = true after MarkAsUntouched, want false")
		}
	})

	t.Run("Set on disabled control is a no-op", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		c.SetDisabled(true)
		c.Set("b")
		if c.Value() != "a" {
			t.Fatalf("Value() = %q after Set on disabled, want %q", c.Value(), "a")
		}
		if c.IsDirty() {
			t.Fatal("disabled control became dirty")
		}
	})
}

func TestControl_ValidatorAggregation(t *testing.T) {
	c := NewControl(
		NewPlainValueStore(""),
		"",
		Required(""),
		MinLength(3),
	)

	c.Set("")

	if c.Validate() {
		t.Fatal("Validate() = true with two failing validators, want false")
	}
	if len(c.errors) != 2 {
		t.Fatalf("expected 2 stored errors, got %d", len(c.errors))
	}

	errs := c.Errors()
	if len(errs) != 1 {
		t.Fatalf("Errors() should have one entry (joined), got %d", len(errs))
	}
	msg := errs[c.Path()]
	if !strings.Contains(msg, "value is required") || !strings.Contains(msg, "minimum length not met") {
		t.Fatalf("Errors() joined message = %q, want both validator messages", msg)
	}
}

func TestControl_ReentrancyNoDeadlock(t *testing.T) {
	// A validator that calls back into the control's read methods while
	// validation runs. If Set held the node lock during validation, this would
	// deadlock.
	var c *Control[string]
	c = NewControl(NewPlainValueStore(""), "", func(v string) error {
		_ = c.Value()
		_ = c.Status()
		_ = c.IsDirty()
		_ = c.IsEnabled()
		if v == "" {
			return errors.New("empty")
		}
		return nil
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.Set("hello")
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock: Set with a reentrant validator did not return")
	}
}

func TestControl_ValueListenerAndUnsubscribe(t *testing.T) {
	t.Run("value listener fires on Set", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		var got string
		c.OnValueChange(func(v string) { got = v })
		c.Set("b")
		if got != "b" {
			t.Fatalf("value listener received %q, want %q", got, "b")
		}
	})

	t.Run("unsubscribe stops delivery", func(t *testing.T) {
		c := NewControl(NewPlainValueStore("a"), "a")
		count := 0
		sub := c.OnValueChange(func(v string) { count++ })
		sub.Unsubscribe()
		c.Set("b")
		if count != 0 {
			t.Fatalf("listener fired %d times after Unsubscribe, want 0", count)
		}
	})
}

func TestControl_Reset(t *testing.T) {
	c := NewControl(NewPlainValueStore("initial"), "initial")
	c.Set("changed")
	c.MarkAsTouched()

	c.Reset()

	if c.Value() != "initial" {
		t.Fatalf("Value() = %q after Reset, want %q", c.Value(), "initial")
	}
	if c.IsDirty() {
		t.Fatal("IsDirty() = true after Reset, want false")
	}
	if c.IsTouched() {
		t.Fatal("IsTouched() = true after Reset, want false")
	}
	if !c.IsPristine() {
		t.Fatal("IsPristine() = false after Reset, want true")
	}
	if c.Errors() != nil {
		t.Fatalf("Errors() = %v after Reset, want nil", c.Errors())
	}
}
