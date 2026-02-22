package state

import (
	"sync"
	"testing"
	"time"
)

func TestMutableValue_AtomicUpdates(t *testing.T) {
	mv := NewMutableValue(10, nil, nil)

	// CompareAndSet
	if !mv.CompareAndSet(10, 20) {
		t.Error("CompareAndSet(10, 20) failed")
	}
	if mv.Get() != 20 {
		t.Errorf("Expected 20, got %v", mv.Get())
	}
	if mv.CompareAndSet(10, 30) {
		t.Error("CompareAndSet(10, 30) succeeded unexpectedly")
	}
	if mv.Get() != 20 {
		t.Errorf("Expected 20, got %v", mv.Get())
	}

	// Update
	mv.Update(func(v any) any {
		return v.(int) + 5
	})
	if mv.Get() != 25 {
		t.Errorf("Expected 25, got %v", mv.Get())
	}

	// UpdateAndGet
	res := mv.UpdateAndGet(func(v any) any {
		return v.(int) * 2
	})
	if res != 50 || mv.Get() != 50 {
		t.Errorf("UpdateAndGet failed. Res: %v, Value: %v", res, mv.Get())
	}

	// GetAndUpdate
	prev := mv.GetAndUpdate(func(v any) any {
		return v.(int) + 1
	})
	if prev != 50 || mv.Get() != 51 {
		t.Errorf("GetAndUpdate failed. Prev: %v, Value: %v", prev, mv.Get())
	}
}

func TestMutableValueTyped_AtomicUpdates(t *testing.T) {
	mv := newMutableState(10)

	// CompareAndSet
	if !mv.CompareAndSet(10, 20) {
		t.Error("CompareAndSet(10, 20) failed")
	}
	if mv.Get() != 20 {
		t.Errorf("Expected 20, got %v", mv.Get())
	}
	if mv.CompareAndSet(10, 30) {
		t.Error("CompareAndSet(10, 30) succeeded unexpectedly")
	}
	if mv.Get() != 20 {
		t.Errorf("Expected 20, got %v", mv.Get())
	}

	// Update
	mv.Update(func(v int) int {
		return v + 5
	})
	if mv.Get() != 25 {
		t.Errorf("Expected 25, got %v", mv.Get())
	}

	// UpdateAndGet
	res := mv.UpdateAndGet(func(v int) int {
		return v * 2
	})
	if res != 50 || mv.Get() != 50 {
		t.Errorf("UpdateAndGet failed. Res: %v, Value: %v", res, mv.Get())
	}

	// GetAndUpdate
	prev := mv.GetAndUpdate(func(v int) int {
		return v + 1
	})
	if prev != 50 || mv.Get() != 51 {
		t.Errorf("GetAndUpdate failed. Prev: %v, Value: %v", prev, mv.Get())
	}
}

func TestMutableValue_Concurrency(t *testing.T) {
	mv := newMutableState(0)
	var wg sync.WaitGroup
	routines := 50
	increments := 100

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				mv.Update(func(c int) int {
					return c + 1
				})
			}
		}()
	}
	wg.Wait()

	expected := routines * increments
	if mv.Get() != expected {
		t.Errorf("Concurrent updates failed. Expected %d, got %d", expected, mv.Get())
	}
}

func TestWrappers(t *testing.T) {
	// Typed wrapper (MutableValueTyped wrapped as MutableValueTypedWrapper which implements MutableValueTyped? No, wrapper is ... implicit?)
	// Let's test the wrappers explicitly.

	// 1. Typed -> Untyped (Unwrap)
	typed := newMutableState(10)
	untyped := typed.Unwrap() // returns MutableValue

	untyped.Update(func(current any) any {
		return current.(int) + 10
	})
	if typed.Get() != 20 {
		t.Errorf("Wrappers: Untyped Update failed to affect underlying typed value")
	}

	// 2. Untyped -> Typed (Wrapper)
	// Create untyped
	baseUntyped := NewMutableValue(100, nil, nil)
	// Wrap as typed
	wrappedTyped, err := MutableValueToTyped[int](baseUntyped)
	if err != nil {
		t.Fatalf("MutableValueToTyped failed: %v", err)
	}

	wrappedTyped.Update(func(current int) int {
		return current + 50
	})
	if baseUntyped.Get() != 150 {
		t.Errorf("Wrappers: Typed wrapper Update failed to affect underlying untyped value")
	}
}

// TestMutableValue_SubscriptionIsHot ensures that subscribers are notified
// immediately when the value changes, without needing to call Get().
func TestMutableValue_SubscriptionIsHot(t *testing.T) {
	// 1. Create a MutableValue
	mv := NewMutableValue(0, nil, nil)

	// 2. Setup a channel to receive notifications
	notifyCh := make(chan struct{}, 10)
	sub := mv.Subscribe(func() {
		notifyCh <- struct{}{}
	})
	defer sub.Unsubscribe()

	// 3. Update the value
	// This should trigger the subscription IMMEDIATELY (Hot)
	// We are NOT calling mv.Get() here, which proves it's not lazy.
	mv.Set(1)

	// 4. Verify notification received
	select {
	case <-notifyCh:
		// Success: Notification received
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for subscription notification. Subscription might be lazy (waiting for Get?) or broken.")
	}

	// Double check value
	if val := mv.Get(); val != 1 {
		t.Errorf("Expected value 1, got %v", val)
	}

	// 5. Update again ensuring repeated hot updates
	mv.Set(2)
	select {
	case <-notifyCh:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for second notification.")
	}
}

// TestMutableValue_Subscription_DistinctUntilChanged ensures we don't get
// notifications if the value hasn't actually changed (based on comparison policy),
// which is a standard optimization, but distinct from "lazy".
func TestMutableValue_Subscription_DistinctUntilChanged(t *testing.T) {
	mv := NewMutableValue(10, nil, nil)
	count := 0
	mv.Subscribe(func() {
		count++
	})

	mv.Set(10) // No change
	if count != 0 {
		t.Errorf("Received notification for unchanged value")
	}

	mv.Set(11) // Change
	if count != 1 {
		t.Errorf("Expected 1 notification, got %d", count)
	}
}
