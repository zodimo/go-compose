package state

import (
	"sync"
	"testing"

	"github.com/zodimo/go-compose/state/internal"
)

// Integration tests to verify MutableValue and TypedMutableValue
// work correctly with DerivedState from the state package.

func TestMutableValue_WithDerivedState_BasicCalculation(t *testing.T) {
	// Setup: MutableValue -> DerivedState
	mv := NewMutableValue(10, nil, func(a, b any) bool { return a == b })

	calculatedCalls := 0
	derived := DerivedStateOf(func() int {
		calculatedCalls++
		return mv.Get().(int) * 2
	})

	// Initial get
	if got := derived.Get(); got != 20 {
		t.Errorf("Expected 20, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Second get, should be cached
	if got := derived.Get(); got != 20 {
		t.Errorf("Expected 20, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation (cached), got %d", calculatedCalls)
	}

	// Update MutableValue
	mv.Set(15)

	// Third get, should recalculate (push invalidation triggered by Set)
	if got := derived.Get(); got != 30 {
		t.Errorf("Expected 30, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}
}

func TestMutableValue_WithDerivedState_PushInvalidation(t *testing.T) {
	mv := NewMutableValue(1, nil, func(a, b any) bool { return a == b })

	derived := DerivedStateOf(func() int {
		return mv.Get().(int) * 2
	})

	// Initial get to establish subscription
	if got := derived.Get(); got != 2 {
		t.Errorf("Expected 2, got %v", got)
	}

	// After initial get, should not be invalid
	if derived.IsInvalid() {
		t.Error("Expected derived to be valid after Get()")
	}

	// Change the MutableValue - should push invalidation
	mv.Set(5)

	// Should now be invalid (push happened)
	if !derived.IsInvalid() {
		t.Error("Expected derived to be invalid after MutableValue changed")
	}

	// Get should recalculate
	if got := derived.Get(); got != 10 {
		t.Errorf("Expected 10, got %v", got)
	}

	// Should be valid again
	if derived.IsInvalid() {
		t.Error("Expected derived to be valid after Get()")
	}
}

func TestMutableValue_WithDerivedState_MultipleDependencies(t *testing.T) {
	mv1 := NewMutableValue(10, nil, func(a, b any) bool { return a == b })
	mv2 := NewMutableValue(20, nil, func(a, b any) bool { return a == b })

	calculatedCalls := 0
	derived := DerivedStateOf(func() int {
		calculatedCalls++
		return mv1.Get().(int) + mv2.Get().(int)
	})

	// Initial: 10 + 20 = 30
	if got := derived.Get(); got != 30 {
		t.Errorf("Expected 30, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Update mv1
	mv1.Set(15)
	if got := derived.Get(); got != 35 {
		t.Errorf("Expected 35, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}

	// Update mv2
	mv2.Set(25)
	if got := derived.Get(); got != 40 {
		t.Errorf("Expected 40, got %v", got)
	}
	if calculatedCalls != 3 {
		t.Errorf("Expected 3 calculations, got %d", calculatedCalls)
	}
}

func TestTypedMutableValue_WithDerivedState_BasicCalculation(t *testing.T) {
	// Setup: MutableValueTypedWrapper -> DerivedState
	mv := NewMutableValue(10, nil, func(a, b any) bool { return a == b })
	typed, err := internal.MutableValueToTyped[int](mv)
	if err != nil {
		t.Fatalf("Failed to wrap: %v", err)
	}

	calculatedCalls := 0
	derived := DerivedStateOf(func() int {
		calculatedCalls++
		return typed.Get() * 2
	})

	// Initial get
	if got := derived.Get(); got != 20 {
		t.Errorf("Expected 20, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Update via typed wrapper
	typed.Set(15)

	// Should recalculate
	if got := derived.Get(); got != 30 {
		t.Errorf("Expected 30, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}
}

func TestTypedMutableValue_WithDerivedState_ChainedDerivedStates(t *testing.T) {
	// Test chain: TypedMutableValue -> DerivedA -> DerivedB
	mv := NewMutableValue(1, nil, func(a, b any) bool { return a == b })
	typed, _ := internal.MutableValueToTyped[int](mv)

	derivedA := DerivedStateOf(func() int {
		return typed.Get() + 10
	})

	derivedB := DerivedStateOf(func() int {
		return derivedA.Get() + 100
	})

	// Initial: (1 + 10) + 100 = 111
	if got := derivedB.Get(); got != 111 {
		t.Errorf("Expected 111, got %v", got)
	}

	// Update typed value
	typed.Set(2)

	// Should propagate through chain: (2 + 10) + 100 = 112
	if got := derivedB.Get(); got != 112 {
		t.Errorf("Expected 112, got %v", got)
	}
}

func TestMutableValue_WithDerivedState_ConcurrentAccess(t *testing.T) {
	mv := NewMutableValue(0, nil, func(a, b any) bool { return a == b })

	derived := DerivedStateOf(func() int {
		return mv.Get().(int) * 2
	})

	var wg sync.WaitGroup
	iterations := 100

	// Multiple readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = derived.Get()
			}
		}()
	}

	// Writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			mv.Set(j)
		}
	}()

	wg.Wait()

	// Should complete without race or panic
	// Final value should be (iterations-1) * 2
	expected := (iterations - 1) * 2
	if got := derived.Get(); got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestTypedMutableValue_WithDerivedState_ConcurrentAccess(t *testing.T) {
	mv := NewMutableValue(0, nil, func(a, b any) bool { return a == b })
	typed, _ := internal.MutableValueToTyped[int](mv)

	derived := DerivedStateOf(func() int {
		return typed.Get() * 2
	})

	var wg sync.WaitGroup
	iterations := 100

	// Multiple readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = derived.Get()
			}
		}()
	}

	// Writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			typed.Set(j)
		}
	}()

	wg.Wait()

	// Should complete without race or panic
	expected := (iterations - 1) * 2
	if got := derived.Get(); got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestMutableValue_WithDerivedState_SubscriberNotification(t *testing.T) {
	mv := NewMutableValue(10, nil, func(a, b any) bool { return a == b })

	derived := DerivedStateOf(func() int {
		return mv.Get().(int) * 2
	})

	// Initial get
	_ = derived.Get()

	// Track value-change notifications
	notified := false
	derived.Subscribe(func() {
		notified = true
	})

	// Change MutableValue - derived value will change from 20 to 30
	mv.Set(15)

	// Trigger recalculation by calling Get()
	_ = derived.Get()

	// Subscriber should have been notified because value changed
	if !notified {
		t.Error("Expected subscriber to be notified when derived value changed")
	}
}

func TestMutableValue_WithDerivedState_NoNotificationIfValueSame(t *testing.T) {
	mv := NewMutableValue(10, nil, func(a, b any) bool { return a == b })

	// Derived state that always returns a constant when input > 5
	derived := DerivedStateOf(func() int {
		val := mv.Get().(int)
		if val > 5 {
			return 1 // Returns 1 for any value > 5
		}
		return 0
	})

	// Initial get: 10 > 5 -> returns 1
	if got := derived.Get(); got != 1 {
		t.Errorf("Expected 1, got %v", got)
	}

	// Track value-change notifications
	notified := false
	derived.Subscribe(func() {
		notified = true
	})

	// Change MutableValue to another value > 5 - derived still returns 1
	mv.Set(20)
	_ = derived.Get()

	// Subscriber should NOT be notified because derived value didn't change
	if notified {
		t.Error("Subscriber should NOT be notified when derived value didn't change")
	}
}
