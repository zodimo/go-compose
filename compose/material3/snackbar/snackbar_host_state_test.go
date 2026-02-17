package snackbar

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose"
)

// newTestHostState creates a SnackbarHostState for testing without a composer.
func newTestHostState() SnackbarHostState {
	return RememberSnackbarHostState(compose.NewComposer())
}

func TestShowSnackbar_NonBlocking(t *testing.T) {
	hostState := newTestHostState()

	// ShowSnackbar should return immediately (non-blocking)
	done := make(chan struct{})
	go func() {
		hostState.ShowSnackbar("test message")
		close(done)
	}()

	select {
	case <-done:
		// Good — returned immediately
	case <-time.After(1 * time.Second):
		t.Fatal("ShowSnackbar should not block")
	}

	// Give the queue goroutine time to set current
	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar data to be non-nil")
	}
	if current.Visuals.Message != "test message" {
		t.Fatalf("expected message 'test message', got %q", current.Visuals.Message)
	}

	// Clean up: dismiss
	current.Dismiss()
	time.Sleep(50 * time.Millisecond)

	if hostState.CurrentSnackbarData() != nil {
		t.Fatal("expected nil after dismiss")
	}
}

func TestShowSnackbar_OnResultCallback(t *testing.T) {
	hostState := newTestHostState()

	var gotResult SnackbarResult
	resultReceived := make(chan struct{})

	hostState.ShowSnackbar("callback msg",
		WithOnResult(func(result SnackbarResult) {
			gotResult = result
			close(resultReceived)
		}),
	)

	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar")
	}

	current.PerformAction()

	select {
	case <-resultReceived:
	case <-time.After(1 * time.Second):
		t.Fatal("callback should have been invoked")
	}

	if gotResult != SnackbarActionPerformed {
		t.Fatalf("expected SnackbarActionPerformed, got %v", gotResult)
	}
}

func TestShowSnackbar_DismissCallsOnResult(t *testing.T) {
	hostState := newTestHostState()

	var gotResult SnackbarResult
	resultReceived := make(chan struct{})

	hostState.ShowSnackbar("dismiss msg",
		WithOnResult(func(result SnackbarResult) {
			gotResult = result
			close(resultReceived)
		}),
	)

	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	current.Dismiss()

	select {
	case <-resultReceived:
	case <-time.After(1 * time.Second):
		t.Fatal("callback should have been invoked on dismiss")
	}

	if gotResult != SnackbarDismissed {
		t.Fatalf("expected SnackbarDismissed, got %v", gotResult)
	}
}

func TestShowSnackbar_FIFOQueue(t *testing.T) {
	hostState := newTestHostState()

	var mu sync.Mutex
	var dismissOrder []string

	// Enqueue 5 snackbars
	for i := 0; i < 5; i++ {
		msg := string(rune('A' + i))
		hostState.ShowSnackbar(msg,
			WithOnResult(func(result SnackbarResult) {
				mu.Lock()
				dismissOrder = append(dismissOrder, msg)
				mu.Unlock()
			}),
		)
	}

	// Process the queue: dismiss each one in order
	for i := 0; i < 5; i++ {
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
		// Let the queue goroutine advance
		time.Sleep(50 * time.Millisecond)
	}

	// Verify FIFO order
	expected := []string{"A", "B", "C", "D", "E"}
	mu.Lock()
	defer mu.Unlock()
	if len(dismissOrder) != len(expected) {
		t.Fatalf("expected %d results, got %d: %v", len(expected), len(dismissOrder), dismissOrder)
	}
	for i, v := range expected {
		if dismissOrder[i] != v {
			t.Fatalf("position %d: expected %q, got %q (full order: %v)", i, v, dismissOrder[i], dismissOrder)
		}
	}
}

func TestShowSnackbar_OnlyOneActive(t *testing.T) {
	hostState := newTestHostState()

	hostState.ShowSnackbar("first")
	hostState.ShowSnackbar("second")

	time.Sleep(50 * time.Millisecond)

	// Only first should be current
	current := hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected current snackbar")
	}
	if current.Visuals.Message != "first" {
		t.Fatalf("expected 'first', got %q", current.Visuals.Message)
	}

	// Dismiss first — queue goroutine should advance to second
	current.Dismiss()
	time.Sleep(50 * time.Millisecond)

	current = hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected second snackbar to be current")
	}
	if current.Visuals.Message != "second" {
		t.Fatalf("expected 'second', got %q", current.Visuals.Message)
	}

	// Dismiss second — should be nil, goroutine should exit
	current.Dismiss()
	time.Sleep(50 * time.Millisecond)

	if hostState.CurrentSnackbarData() != nil {
		t.Fatal("expected no current snackbar after all dismissed")
	}

	if snackbar, ok := hostState.(*snackbarHostState); ok {
		snackbar.mu.Lock()
		processing := snackbar.processing
		snackbar.mu.Unlock()
		if processing {
			t.Fatal("expected processing to be false after queue is empty")
		}
	} else {
		t.Fatal("expected *snackbarHostState")
	}

}

func TestShowSnackbar_ContextCancellation(t *testing.T) {
	hostState := newTestHostState()

	// Show a snackbar that occupies the slot
	hostState.ShowSnackbar("blocker")
	time.Sleep(50 * time.Millisecond)

	// Queue a second snackbar with a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	hostState.ShowSnackbar("cancelled", WithContext(ctx))
	time.Sleep(50 * time.Millisecond)

	// Current should still be "blocker"
	current := hostState.CurrentSnackbarData()
	if current.Visuals.Message != "blocker" {
		t.Fatalf("expected 'blocker', got %q", current.Visuals.Message)
	}

	// Cancel the queued snackbar
	cancel()
	time.Sleep(50 * time.Millisecond)

	// Verify the queue is empty (cancelled item removed)

	if snackbar, ok := hostState.(*snackbarHostState); ok {
		snackbar.mu.Lock()
		queueLen := len(snackbar.queue)
		snackbar.mu.Unlock()
		if queueLen != 0 {
			t.Fatalf("expected empty queue after cancellation, got %d items", queueLen)
		}
	} else {
		t.Fatal("expected *snackbarHostState")
	}

	// Clean up: dismiss the blocker
	current.Dismiss()
	time.Sleep(50 * time.Millisecond)
}

func TestShowSnackbar_QueueGoroutineRestartsAfterEmpty(t *testing.T) {
	hostState := newTestHostState()

	// First round
	hostState.ShowSnackbar("round1")
	time.Sleep(50 * time.Millisecond)

	current := hostState.CurrentSnackbarData()
	current.Dismiss()
	time.Sleep(50 * time.Millisecond)

	// Goroutine should have stopped
	if snackbar, ok := hostState.(*snackbarHostState); ok {
		snackbar.mu.Lock()
		processing := snackbar.processing
		snackbar.mu.Unlock()
		if processing {
			t.Fatal("expected processing to be false")
		}
	} else {
		t.Fatal("expected *snackbarHostState")
	}

	// Second round — goroutine should restart
	hostState.ShowSnackbar("round2")
	time.Sleep(50 * time.Millisecond)

	current = hostState.CurrentSnackbarData()
	if current == nil {
		t.Fatal("expected snackbar after restart")
	}
	if current.Visuals.Message != "round2" {
		t.Fatalf("expected 'round2', got %q", current.Visuals.Message)
	}

	current.Dismiss()
	time.Sleep(50 * time.Millisecond)
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
	data := newSnackbarData(SnackbarVisuals{Message: "test"}, nil)

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
