package effect

import (
	"sync"
	"testing"
)

func TestKeysEqual_BothNil(t *testing.T) {
	if !keysEqual(nil, nil) {
		t.Error("expected nil slices to be equal")
	}
}

func TestKeysEqual_DifferentLengths(t *testing.T) {
	if keysEqual([]any{1}, []any{1, 2}) {
		t.Error("expected different length slices to not be equal")
	}
}

func TestKeysEqual_EmptySlices(t *testing.T) {
	if !keysEqual([]any{}, []any{}) {
		t.Error("expected empty slices to be equal")
	}
}

func TestKeysEqual_SameStringValues(t *testing.T) {
	if !keysEqual([]any{"hello", "world"}, []any{"hello", "world"}) {
		t.Error("expected same string values to be equal")
	}
}

func TestKeysEqual_DifferentStringValues(t *testing.T) {
	if keysEqual([]any{"hello"}, []any{"world"}) {
		t.Error("expected different string values to not be equal")
	}
}

func TestKeysEqual_SameIntValues(t *testing.T) {
	if !keysEqual([]any{1, 2, 3}, []any{1, 2, 3}) {
		t.Error("expected same int values to be equal")
	}
}

func TestKeysEqual_DifferentIntValues(t *testing.T) {
	if keysEqual([]any{1}, []any{2}) {
		t.Error("expected different int values to not be equal")
	}
}

func TestKeyEqual_BothNil(t *testing.T) {
	if !keyEqual(nil, nil) {
		t.Error("expected nil keys to be equal")
	}
}

func TestKeyEqual_OneNil(t *testing.T) {
	if keyEqual(nil, "hello") {
		t.Error("expected nil and non-nil to not be equal")
	}
	if keyEqual("hello", nil) {
		t.Error("expected non-nil and nil to not be equal")
	}
}

func TestKeyEqual_DifferentTypes(t *testing.T) {
	if keyEqual(1, "1") {
		t.Error("expected different types to not be equal")
	}
}

// TestKeyEqual_SamePointer verifies that pointers to the same object
// are considered equal (identity comparison, not deep comparison).
func TestKeyEqual_SamePointer(t *testing.T) {
	type myStruct struct {
		mu    sync.Mutex
		value int
	}
	obj := &myStruct{value: 42}
	if !keyEqual(obj, obj) {
		t.Error("expected same pointer to be equal")
	}
}

// TestKeyEqual_DifferentPointers verifies that pointers to different objects
// are considered not equal, even if the objects have the same value.
func TestKeyEqual_DifferentPointers(t *testing.T) {
	type myStruct struct {
		value int
	}
	a := &myStruct{value: 42}
	b := &myStruct{value: 42}
	if keyEqual(a, b) {
		t.Error("expected different pointers to not be equal, even with same values")
	}
}

// TestKeyEqual_PointerWithMutex verifies the core fix: a struct containing
// a sync.Mutex should use pointer identity, not deep comparison.
// With reflect.DeepEqual, the mutex's internal state would cause false negatives.
func TestKeyEqual_PointerWithMutex(t *testing.T) {
	type stateFlow struct {
		mu    sync.Mutex
		value int
	}
	flow := &stateFlow{value: 1}

	// Lock and unlock to change mutex internal state
	flow.mu.Lock()
	flow.mu.Unlock()

	// The same pointer should still be considered equal despite
	// internal mutex state changes
	if !keyEqual(flow, flow) {
		t.Error("expected same pointer to be equal despite mutex state changes")
	}
}

func TestKeyEqual_SameChannel(t *testing.T) {
	ch := make(chan int)
	if !keyEqual(ch, ch) {
		t.Error("expected same channel to be equal")
	}
}

func TestKeyEqual_DifferentChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	if keyEqual(ch1, ch2) {
		t.Error("expected different channels to not be equal")
	}
}

func TestKeyEqual_SameFunc(t *testing.T) {
	fn := func() {}
	if !keyEqual(fn, fn) {
		t.Error("expected same func to be equal")
	}
}

// TestKeysEqual_MixedTypes verifies comparison of mixed value and pointer types.
func TestKeysEqual_MixedTypes(t *testing.T) {
	obj := &struct{ x int }{x: 1}
	keys1 := []any{"route_key", 42, obj}
	keys2 := []any{"route_key", 42, obj}
	if !keysEqual(keys1, keys2) {
		t.Error("expected mixed keys with same values and pointers to be equal")
	}
}

// TestKeysEqual_MixedTypes_DifferentPointer verifies that different pointers
// in a mixed key slice cause inequality.
func TestKeysEqual_MixedTypes_DifferentPointer(t *testing.T) {
	obj1 := &struct{ x int }{x: 1}
	obj2 := &struct{ x int }{x: 1}
	keys1 := []any{"route_key", 42, obj1}
	keys2 := []any{"route_key", 42, obj2}
	if keysEqual(keys1, keys2) {
		t.Error("expected mixed keys with different pointers to not be equal")
	}
}
