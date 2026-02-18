package snackbar

import (
	"context"
	"sync"
	"time"

	"gioui.org/widget"
)

// SnackbarResult represents the outcome of a snackbar being shown.
// Matches Kotlin's SnackbarResult enum.
type SnackbarResult int

const (
	// SnackbarDismissed indicates the snackbar was dismissed by timeout or by the user.
	SnackbarDismissed SnackbarResult = iota
	// SnackbarActionPerformed indicates the action button was clicked.
	SnackbarActionPerformed
)

// SnackbarDuration controls how long a snackbar will be shown.
// Matches Kotlin's SnackbarDuration enum.
type SnackbarDuration int

const (
	// SnackbarDurationShort shows the snackbar for a short period (4s).
	SnackbarDurationShort SnackbarDuration = iota
	// SnackbarDurationLong shows the snackbar for a longer period (10s).
	SnackbarDurationLong
	// SnackbarDurationIndefinite shows the snackbar until explicitly dismissed or action is clicked.
	SnackbarDurationIndefinite
)

// ToDuration converts a SnackbarDuration to a time.Duration.
// Returns 0 for Indefinite (meaning no auto-dismiss).
func (d SnackbarDuration) ToDuration() time.Duration {
	switch d {
	case SnackbarDurationShort:
		return 4 * time.Second
	case SnackbarDurationLong:
		return 10 * time.Second
	case SnackbarDurationIndefinite:
		return 0
	default:
		return 4 * time.Second
	}
}

// SnackbarVisuals holds the visual representation for a particular Snackbar.
// Matches Kotlin's SnackbarVisuals interface.
type SnackbarVisuals struct {
	Message           string
	ActionLabel       string
	WithDismissAction bool
	Duration          SnackbarDuration
}

// SnackbarData represents the data of one particular Snackbar as managed by SnackbarHostState.
// Matches Kotlin's SnackbarData interface.
type SnackbarData struct {
	Context context.Context

	Visuals  SnackbarVisuals
	onResult func(SnackbarResult) // optional callback invoked when dismissed/actioned
	resultCh chan SnackbarResult  // signals the queue processor to advance
	once     sync.Once

	// Gio clickable state — persists across frames for immediate-mode rendering
	ActionClickable  widget.Clickable
	DismissClickable widget.Clickable
}

// newSnackbarData creates a new SnackbarData with a result channel and optional callback.
func newSnackbarData(ctx context.Context, visuals SnackbarVisuals, onResult func(SnackbarResult)) *SnackbarData {
	return &SnackbarData{
		Context:  ctx,
		Visuals:  visuals,
		onResult: onResult,
		resultCh: make(chan SnackbarResult, 1),
	}
}

// PerformAction notifies that the action on the Snackbar was clicked.
func (d *SnackbarData) PerformAction() {
	d.once.Do(func() {
		d.resultCh <- SnackbarActionPerformed
		if d.onResult != nil {
			d.onResult(SnackbarActionPerformed)
		}
	})
}

// Dismiss notifies that the Snackbar was dismissed (by timeout or user).
func (d *SnackbarData) Dismiss() {
	d.once.Do(func() {
		d.resultCh <- SnackbarDismissed
		if d.onResult != nil {
			d.onResult(SnackbarDismissed)
		}
	})
}
