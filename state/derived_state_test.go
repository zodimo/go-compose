package state

import (
	"sync"
	"testing"

	"github.com/zodimo/go-compose/state/core"
	"github.com/zodimo/go-compose/state/internal"
)

func newMockMutableValue[T any](value T) *internal.MutableValue {
	mv := NewMutableValue(value, nil, nil)
	return mv.(*internal.MutableValue)
}

func TestDerivedState_BasicCalculation(t *testing.T) {
	mv1 := newMockMutableValue(10)
	mv2 := newMockMutableValue(20)

	calculatedCalls := 0
	derived := DerivedStateOf(func() int {
		calculatedCalls++
		return mv1.Get().(int) + mv2.Get().(int)
	})

	// Initial get
	if got := derived.Get(); got != 30 {
		t.Errorf("Expected 30, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Second get, should be cached
	if got := derived.Get(); got != 30 {
		t.Errorf("Expected 30, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Update dependency
	mv1.Set(15)

	// Third get, should recalculate (push invalidation triggered by Set)
	if got := derived.Get(); got != 35 {
		t.Errorf("Expected 35, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}
}

func TestDerivedState_DependencySwitching(t *testing.T) {
	// Setup: A -> B or C
	toggle := newMockMutableValue(true)
	valB := newMockMutableValue("B")
	valC := newMockMutableValue("C")

	calculatedCalls := 0
	derived := DerivedStateOf(func() string {
		calculatedCalls++
		if toggle.Get().(bool) {
			return valB.Get().(string)
		}
		return valC.Get().(string)
	})

	// 1. Initial: true -> depends on toggle, valB
	if got := derived.Get(); got != "B" {
		t.Errorf("Expected B, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// 2. Modify valC (should NOT trigger recalc on next Get, as it's not a dependency)
	valC.Set("C2")
	if got := derived.Get(); got != "B" {
		t.Errorf("Expected B, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// 3. Toggle to false -> depends on toggle, valC
	toggle.Set(false)
	if got := derived.Get(); got != "C2" {
		t.Errorf("Expected C2, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}

	// 4. Modify valB (should NOT trigger recalc on next Get, as it's no longer a dependency)
	valB.Set("B2")
	if got := derived.Get(); got != "C2" {
		t.Errorf("Expected C2, got %v", got)
	}
	if calculatedCalls != 2 {
		t.Errorf("Expected 2 calculations, got %d", calculatedCalls)
	}
}

func TestDerivedState_DiamondDependency(t *testing.T) {
	// Root -> A, Root -> B
	// Derived -> A + B
	root := newMockMutableValue(1)

	derivedA := DerivedStateOf(func() int {
		return root.Get().(int) * 2
	})

	derivedB := DerivedStateOf(func() int {
		return root.Get().(int) * 3
	})

	calculatedCalls := 0
	derivedFinal := DerivedStateOf(func() int {
		calculatedCalls++
		return derivedA.Get() + derivedB.Get()
	})

	// Initial: (1*2) + (1*3) = 5
	if got := derivedFinal.Get(); got != 5 {
		t.Errorf("Expected 5, got %v", got)
	}
	if calculatedCalls != 1 {
		t.Errorf("Expected 1 calculation, got %d", calculatedCalls)
	}

	// Update root
	root.Set(2)
	// DerivedA -> 4, DerivedB -> 6
	// DerivedFinal -> 10

	if got := derivedFinal.Get(); got != 10 {
		t.Errorf("Expected 10, got %v", got)
	}
}

func TestDerivedState_NestedCaching(t *testing.T) {
	// A -> Derived B -> Derived C
	// When A changes but B's computed value stays the same:
	// - B must recalculate to determine its new value
	// - C must also recalculate (since it was invalidated by B)
	// - BUT C's subscribers should NOT be notified if C's value didn't change

	valA := newMockMutableValue(10)

	calcB := 0
	derivedB := DerivedStateOf(func() int {
		calcB++
		a := valA.Get().(int)
		if a > 5 {
			return 1 // Always 1 if A > 5
		}
		return 0
	})

	calcC := 0
	derivedC := DerivedStateOf(func() int {
		calcC++
		return derivedB.Get() + 100
	})

	// Track if C's subscribers are notified
	cNotified := false
	derivedC.Subscribe(func() {
		cNotified = true
	})

	// Initial: A=10 -> B=1 -> C=101
	if got := derivedC.Get(); got != 101 {
		t.Errorf("Got %d", got)
	}
	if calcB != 1 {
		t.Errorf("B calc %d", calcB)
	}
	if calcC != 1 {
		t.Errorf("C calc %d", calcC)
	}

	// Update A to 20. B will recalculate but still return 1.
	// C will also recalculate (due to invalidation chain), but its value stays 101.
	// C's subscribers should NOT be notified since C's value didn't change.
	cNotified = false
	valA.Set(20)

	if got := derivedC.Get(); got != 101 {
		t.Errorf("Got %d", got)
	}

	if calcB != 2 {
		t.Errorf("B should have recalculated, got %d", calcB)
	}
	// With push-based invalidation, C also recalculates to check if its value changed
	if calcC != 2 {
		t.Errorf("C should have recalculated to verify its value, got %d", calcC)
	}
	// The key optimization: C's subscribers were NOT notified because C's value didn't change
	if cNotified {
		t.Error("C's subscribers should NOT be notified when C's value didn't change")
	}
}

// New tests for subscription pattern

func TestDerivedState_PushInvalidation(t *testing.T) {
	mv := newMockMutableValue(1)

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

	// Change the dependency - should push invalidation
	mv.Set(5)

	// Should now be invalid (push happened)
	if !derived.IsInvalid() {
		t.Error("Expected derived to be invalid after dependency changed")
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

func TestDerivedState_SubscriptionCleanup(t *testing.T) {
	mvA := newMockMutableValue("A")
	mvB := newMockMutableValue("B")
	toggle := newMockMutableValue(true)

	derived := DerivedStateOf(func() string {
		if toggle.Get().(bool) {
			return mvA.Get().(string)
		}
		return mvB.Get().(string)
	})

	// Initial: depends on toggle and mvA
	_ = derived.Get()

	// Check subscription count on mvA
	if internal.MutableValueSubscribers(mvA).Count() != 1 {
		t.Errorf("Expected 1 subscriber on mvA, got %d", internal.MutableValueSubscribers(mvA).Count())
	}
	if internal.MutableValueSubscribers(mvB).Count() != 0 {
		t.Errorf("Expected 0 subscribers on mvB, got %d", internal.MutableValueSubscribers(mvB).Count())
	}

	// Toggle to false: should unsubscribe from mvA, subscribe to mvB
	toggle.Set(false)
	_ = derived.Get()

	if internal.MutableValueSubscribers(mvA).Count() != 0 {
		t.Errorf("Expected 0 subscribers on mvA after toggle, got %d", internal.MutableValueSubscribers(mvA).Count())
	}
	if internal.MutableValueSubscribers(mvB).Count() != 1 {
		t.Errorf("Expected 1 subscriber on mvB after toggle, got %d", internal.MutableValueSubscribers(mvB).Count())
	}
}

func TestDerivedState_ChainedInvalidation(t *testing.T) {
	// Test that invalidation propagates through a chain: A -> B -> C
	root := newMockMutableValue(1)

	derivedB := DerivedStateOf(func() int {
		return root.Get().(int) + 10
	})

	derivedC := DerivedStateOf(func() int {
		return derivedB.Get() + 100
	})

	// Initial calculation
	if got := derivedC.Get(); got != 111 {
		t.Errorf("Expected 111, got %v", got)
	}

	// Track if C gets invalidated (use SubscribeForInvalidation for invalidation chain)
	invalidated := false
	derivedC.SubscribeForInvalidation(func() {
		invalidated = true
	})

	// Change root - should trigger chain invalidation
	root.Set(2)

	if !invalidated {
		t.Error("Expected derivedC to be invalidated via chain")
	}

	if got := derivedC.Get(); got != 112 {
		t.Errorf("Expected 112, got %v", got)
	}
}

func TestSubscriptionManager_MultipleSubscribers(t *testing.T) {
	sm := core.NewSubscriptionManager()

	called1 := false
	called2 := false

	sub1 := sm.Subscribe(func() { called1 = true })
	sub2 := sm.Subscribe(func() { called2 = true })

	sm.NotifyAll()

	if !called1 || !called2 {
		t.Error("Expected both callbacks to be called")
	}

	// Unsubscribe one
	called1 = false
	called2 = false
	sub1.Unsubscribe()

	sm.NotifyAll()

	if called1 {
		t.Error("Expected callback1 to NOT be called after unsubscribe")
	}
	if !called2 {
		t.Error("Expected callback2 to still be called")
	}

	// Unsubscribe again (should be idempotent)
	sub1.Unsubscribe()
	sub2.Unsubscribe()

	if sm.Count() != 0 {
		t.Errorf("Expected 0 subscribers, got %d", sm.Count())
	}
}

func TestDerivedState_ConcurrentAccess(t *testing.T) {
	mv := newMockMutableValue(0)

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
func TestDerivedState_EagerSubscription(t *testing.T) {
	// 1. Setup dependencies using internal mock to avoid import cycle
	root := newMockMutableValue(10)

	calcCount := 0
	derived := DerivedStateOf(func() int {
		calcCount++
		return root.Get().(int) * 2
	})

	// 2. Initial Get to register dependencies
	if got := derived.Get(); got != 20 {
		t.Errorf("Initial Get: expected 20, got %d", got)
	}
	if calcCount != 1 {
		t.Errorf("Initial Calc: expected 1, got %d", calcCount)
	}

	// 3. Subscribe
	notifyCount := 0
	sub := derived.Subscribe(func() {
		notifyCount++
	})
	defer sub.Unsubscribe()

	// 4. Update dependency
	// This should trigger invalidation -> check subscribers -> eager recalculate -> notify subscribers
	root.Set(20)

	// 5. Verify notification happened without calling Get()
	if notifyCount != 1 {
		t.Errorf("Expected 1 notification, got %d", notifyCount)
	}

	// Verify calculation happened
	if calcCount != 2 {
		t.Errorf("Expected 2 calculations, got %d", calcCount)
	}

	// Verify value is correct
	if got := derived.Get(); got != 40 {
		t.Errorf("Expected 40, got %d", got)
	}
	// Calling Get should NOT trigger another calc since it was eagerly calculated
	if calcCount != 2 {
		t.Errorf("Expected 2 calculations (cached), got %d", calcCount)
	}
}

func TestDerivedState_SubscribeActivates(t *testing.T) {
	mv := newMockMutableValue(10)

	calcCount := 0
	derived := DerivedStateOf(func() int {
		calcCount++
		return mv.Get().(int) * 2
	})

	// Subscribe WITHOUT calling Get() first
	// This should ideally trigger calculation to register dependencies
	notifyCount := 0
	derived.Subscribe(func() {
		notifyCount++
	})

	// Before the fix, this would be 0. After fix, it should be 1 because Subscribe forces an initial calculation/validation
	// to ensure dependencies are registered.
	if calcCount == 0 {
		t.Log("Subscribe did not trigger calculation - this means we aren't listening to deps!")
	}

	// Update dependency
	mv.Set(20)

	// Should be notified TWICE:
	// 1. Initial activation (Subscribe -> Get -> Recalculate -> Notify)
	// 2. Dependency change (Set -> Invalidate -> Recalculate -> Notify)
	if notifyCount != 2 {
		t.Errorf("Expected 2 notifications, got %d", notifyCount)
	}

	// 3rd update
	mv.Set(30)
	if notifyCount != 3 {
		t.Errorf("Expected 3 notifications, got %d", notifyCount)
	}
}

func TestDerivedState_ChainedActivation(t *testing.T) {
	mv := newMockMutableValue(10)

	// derived1 depends on mv
	derived1 := DerivedStateOf(func() int {
		return mv.Get().(int) * 2
	})

	// derived2 depends on derived1
	derived2 := DerivedStateOf(func() int {
		return derived1.Get() + 5
	})

	notifyCount := 0
	// Subscribe to the outer derived state
	// This should trigger activation of derived2, which calls derived1.Get(),
	// which activates derived1, which calls mv.Get().
	derived2.Subscribe(func() {
		notifyCount++
	})

	// Should receive initial activation notification
	if notifyCount != 1 {
		t.Errorf("Expected initial activation notification, got %d", notifyCount)
	}

	// Update the bottom-most dependency
	mv.Set(20)

	// derived1 should invalidate -> derived2 should invalidate
	// derived2 should recalculate -> derived1 should recalculate
	// derived2 value changes (10*2+5=25 -> 20*2+5=45) -> notify subscriber
	if notifyCount != 2 {
		t.Errorf("Expected 2 notifications after update, got %d", notifyCount)
	}

	// Verify values
	if got := derived2.Get(); got != 45 {
		t.Errorf("Expected 45, got %v", got)
	}
}
