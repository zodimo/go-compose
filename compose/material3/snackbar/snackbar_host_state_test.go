package snackbar

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose"
)

func TestShowSnackbar_BlocksUntilDismissed(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	done := make(chan struct{})
	var result SnackbarResult
	var err error

	go func() {
		result, err = hostState.ShowSnackbar(context.Background(), "test message")
		close(done)
	}()

	// Give the goroutine time to start and block
	time.Sleep(50 * time.Millisecond)

	// Verify snackbar is showing
	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar data to be non-nil")
	}
	if current.Visuals.Message != "test message" {
		t.Fatalf("expected message 'test message', got %q", current.Visuals.Message)
	}

	// Verify caller is still blocked
	select {
	case <-done:
		t.Fatal("ShowSnackbar should still be blocking")
	default:
	}

	// Dismiss the snackbar
	current.Dismiss()

	// Wait for the caller to unblock
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("ShowSnackbar should have returned after dismiss")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != SnackbarDismissed {
		t.Fatalf("expected SnackbarDismissed, got %v", result)
	}
}

func TestShowSnackbar_ActionPerformed(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	done := make(chan struct{})
	var result SnackbarResult

	go func() {
		result, _ = hostState.ShowSnackbar(context.Background(), "action msg",
			WithActionLabel("undo"))
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar data")
	}
	if current.Visuals.ActionLabel != "undo" {
		t.Fatalf("expected action label 'undo', got %q", current.Visuals.ActionLabel)
	}

	current.PerformAction()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("ShowSnackbar should have returned after action")
	}

	if result != SnackbarActionPerformed {
		t.Fatalf("expected SnackbarActionPerformed, got %v", result)
	}
}

func TestShowSnackbar_FIFOQueue(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	var mu sync.Mutex
	var resultOrder []string

	// Launch 5 snackbars concurrently, staggered to ensure FIFO ordering
	const count = 5
	allDone := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		msg := string(rune('A' + i)) // "A", "B", "C", "D", "E"
		time.Sleep(10 * time.Millisecond)
		go func(msg string) {
			defer wg.Done()
			hostState.ShowSnackbar(context.Background(), msg)
			mu.Lock()
			resultOrder = append(resultOrder, msg)
			mu.Unlock()
		}(msg)
	}

	go func() {
		wg.Wait()
		close(allDone)
	}()

	// Process the queue: dismiss each one in order
	for i := 0; i < count; i++ {
		// Wait for a current snackbar to appear
		var current *SnackbarData
		for attempt := 0; attempt < 100; attempt++ {
			current = hostState.CurrentSnackbarData()
			if current != nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if current == nil {
			t.Fatalf("iteration %d: no current snackbar appeared", i)
		}

		expectedMsg := string(rune('A' + i))
		if current.Visuals.Message != expectedMsg {
			t.Fatalf("iteration %d: expected message %q, got %q", i, expectedMsg, current.Visuals.Message)
		}

		current.Dismiss()
		hostState.advanceQueue()

		// Small delay to let the dismissed goroutine complete
		time.Sleep(20 * time.Millisecond)
	}

	select {
	case <-allDone:
	case <-time.After(5 * time.Second):
		t.Fatal("not all snackbars completed in time")
	}

	// Verify FIFO order
	expected := []string{"A", "B", "C", "D", "E"}
	if len(resultOrder) != len(expected) {
		t.Fatalf("expected %d results, got %d: %v", len(expected), len(resultOrder), resultOrder)
	}
	for i, v := range expected {
		if resultOrder[i] != v {
			t.Fatalf("position %d: expected %q, got %q (full order: %v)", i, v, resultOrder[i], resultOrder)
		}
	}
}

func TestShowSnackbar_OnlyOneActive(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	// Show first snackbar
	go hostState.ShowSnackbar(context.Background(), "first")
	time.Sleep(50 * time.Millisecond)

	// Show second snackbar (should be queued)
	go hostState.ShowSnackbar(context.Background(), "second")
	time.Sleep(50 * time.Millisecond)

	// Only first should be current
	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar")
	}
	if current.Visuals.Message != "first" {
		t.Fatalf("expected 'first', got %q", current.Visuals.Message)
	}

	// Dismiss first — advance queue
	current.Dismiss()
	hostState.advanceQueue()
	time.Sleep(50 * time.Millisecond)

	// Second should now be current
	current = hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected second snackbar to be current")
	}
	if current.Visuals.Message != "second" {
		t.Fatalf("expected 'second', got %q", current.Visuals.Message)
	}

	// Dismiss second — should be nil
	current.Dismiss()
	hostState.advanceQueue()
	time.Sleep(50 * time.Millisecond)

	if hostState.CurrentSnackbarData() != nil {
		t.Fatal("expected no current snackbar after all dismissed")
	}
}

func TestShowSnackbar_ContextCancellation(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	// Show a snackbar that occupies the slot
	go hostState.ShowSnackbar(context.Background(), "blocker")
	time.Sleep(50 * time.Millisecond)

	// Queue a second snackbar with a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	var err error

	go func() {
		_, err = hostState.ShowSnackbar(ctx, "cancelled")
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	// Verify it is queued (current is still "blocker")
	current := hostState.CurrentSnackbarData()
	if current.Visuals.Message != "blocker" {
		t.Fatalf("expected 'blocker', got %q", current.Visuals.Message)
	}

	// Cancel the queued snackbar
	cancel()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("cancelled ShowSnackbar should have returned")
	}

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	// Verify the queue is empty (cancelled item removed)
	hostState.mu.Lock()
	queueLen := len(hostState.queue)
	hostState.mu.Unlock()

	if queueLen != 0 {
		t.Fatalf("expected empty queue after cancellation, got %d items", queueLen)
	}

	// Clean up: dismiss the blocker
	current.Dismiss()
	hostState.advanceQueue()
}

func TestShowSnackbarAsync_DoesNotBlock(t *testing.T) {
	hostState := RemeberSnackbarHostState(compose.NewComposer())

	// ShowSnackbarAsync should return immediately
	done := make(chan struct{})
	go func() {
		hostState.ShowSnackbarAsync("async msg")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("ShowSnackbarAsync should not block")
	}

	// Give the goroutine inside ShowSnackbarAsync time to enqueue
	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected snackbar to be showing")
	}
	if current.Visuals.Message != "async msg" {
		t.Fatalf("expected 'async msg', got %q", current.Visuals.Message)
	}

	// Clean up
	current.Dismiss()
	hostState.advanceQueue()
}

func TestSnackbarDuration_DefaultWithAction(t *testing.T) {
	opts := DefaultOptions()
	WithActionLabel("undo")(&opts)
	opts.resolveDefaults()

	if opts.Duration != SnackbarDurationIndefinite {
		t.Fatalf("expected Indefinite duration when action label present, got %v", opts.Duration)
	}
}

func TestSnackbarDuration_DefaultWithoutAction(t *testing.T) {
	opts := DefaultOptions()
	opts.resolveDefaults()

	if opts.Duration != SnackbarDurationShort {
		t.Fatalf("expected Short duration when no action label, got %v", opts.Duration)
	}
}

func TestSnackbarDuration_ToDuration(t *testing.T) {
	tests := []struct {
		name     string
		dur      SnackbarDuration
		expected time.Duration
	}{
		{"Short", SnackbarDurationShort, 4 * time.Second},
		{"Long", SnackbarDurationLong, 10 * time.Second},
		{"Indefinite", SnackbarDurationIndefinite, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dur.ToDuration(); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestSnackbarData_DismissOnlyOnce(t *testing.T) {
	data := newSnackbarData(SnackbarVisuals{Message: "test"})

	// First dismiss should send result
	data.Dismiss()

	result := <-data.resultCh
	if result != SnackbarDismissed {
		t.Fatalf("expected SnackbarDismissed, got %v", result)
	}

	// Second dismiss should not panic (sync.Once protects)
	data.Dismiss()
	data.PerformAction()
}
