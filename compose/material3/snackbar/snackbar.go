package snackbar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	"gioui.org/layout"
	gioSnackbar "git.sr.ht/~schnwalter/gio-mw/widget/snackbar"
)

// SnackbarHostState controls the queue and the current Snackbar being shown inside the SnackbarHost.
//
// Matches Kotlin's SnackbarHostState: guarantees to show at most one snackbar at a time.
// If ShowSnackbar is called while another snackbar is visible, it will be queued (FIFO)
// and the caller will block until the snackbar is shown and subsequently dismissed.
//
// The currentSnackbar field is a MutableValueTyped, making it observable by the compose framework.
// When it changes, the SnackbarHost composable is automatically recomposed.
type SnackbarHostState struct {
	// currentSnackbar is an observable state — reading it in a composable registers
	// observation, and setting it triggers recomposition. Matches Kotlin's
	// mutableStateOf<SnackbarData?>(null).
	currentSnackbar state.MutableValueTyped[*SnackbarData]

	mu    sync.Mutex
	queue []*SnackbarData
}

// NewSnackbarHostState creates a new SnackbarHostState.
func RemeberSnackbarHostState(c compose.Composer) *SnackbarHostState {

	key := c.GenerateID()
	path := c.GetPath()

	snackbarHostStatePath := fmt.Sprintf("%d/%s/snackbarHostState", key, path)

	currentSnackbar := state.MustRemember(c, snackbarHostStatePath, func() *SnackbarData {
		return nil
	})

	return &SnackbarHostState{
		currentSnackbar: currentSnackbar,
	}
}

// ShowSnackbar shows or queues a snackbar and blocks until it is dismissed.
//
// This matches Kotlin's suspend fun showSnackbar: the caller is suspended (blocked)
// until the snackbar disappears. If another snackbar is already visible, this one
// is queued and will be shown after the current one is dismissed.
//
// If the context is cancelled, the snackbar is removed from the queue (if queued)
// or dismissed (if currently shown), and ctx.Err() is returned.
//
// Must be called from a goroutine (not the main/UI goroutine).
func (s *SnackbarHostState) ShowSnackbar(ctx context.Context, message string, options ...SnackbarOption) (SnackbarResult, error) {
	opts := DefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}
	opts.resolveDefaults()

	visuals := SnackbarVisuals{
		Message:           message,
		ActionLabel:       opts.ActionLabel,
		WithDismissAction: opts.WithDismissAction,
		Duration:          opts.Duration,
	}

	data := newSnackbarData(visuals)

	s.mu.Lock()
	if s.currentSnackbar.Get() == nil {
		// No snackbar currently showing — display immediately
		s.currentSnackbar.Set(data)
	} else {
		// Queue behind the current one
		s.queue = append(s.queue, data)
	}
	s.mu.Unlock()

	// Block until result or context cancellation
	select {
	case result := <-data.resultCh:
		return result, nil
	case <-ctx.Done():
		// Cancel: remove from queue if queued, or dismiss if current
		s.mu.Lock()
		if s.currentSnackbar.Get() == data {
			s.mu.Unlock()
			data.Dismiss()
			// Drain the result channel (Dismiss sends to it)
			<-data.resultCh
			s.advanceQueue()
		} else {
			// Remove from queue
			for i, d := range s.queue {
				if d == data {
					s.queue = append(s.queue[:i], s.queue[i+1:]...)
					break
				}
			}
			s.mu.Unlock()
		}
		return SnackbarDismissed, ctx.Err()
	}
}

// ShowSnackbarAsync is a fire-and-forget convenience that shows a snackbar without blocking.
// This is useful for click handlers that cannot block.
func (s *SnackbarHostState) ShowSnackbarAsync(message string, options ...SnackbarOption) {
	go func() {
		_, _ = s.ShowSnackbar(context.Background(), message, options...)
	}()
}

// CurrentSnackbarData returns the currently displayed snackbar data, or nil if none.
// Reading this in a composable registers observation — recomposition will be triggered
// when the current snackbar changes.
func (s *SnackbarHostState) CurrentSnackbarData() *SnackbarData {
	return s.currentSnackbar.Get()
}

// advanceQueue moves to the next queued snackbar, if any.
// Sets the observable state which triggers recomposition.
func (s *SnackbarHostState) advanceQueue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) > 0 {
		s.currentSnackbar.Set(s.queue[0])
		s.queue = s.queue[1:]
	} else {
		s.currentSnackbar.Set(nil)
	}
}

// SnackbarHost is a composable that displays snackbars from the given SnackbarHostState.
//
// It shows at most one snackbar at a time, using a foundation overlay for rendering.
// Duration-based auto-dismiss is managed via LaunchedEffect.
func SnackbarHost(hostState *SnackbarHostState) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		c.StartBlock("SnackbarHost")

		// Reading CurrentSnackbarData() calls Get() on the MutableValueTyped,
		// which triggers NotifyRead — the compose framework now observes this state.
		current := hostState.CurrentSnackbarData()

		constructor := func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
			return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
				if current == nil {
					return layout.Dimensions{}
				}

				// Render the snackbar styled widget at bottom-center
				snackStyle := gioSnackbar.Plain(current.Visuals.Message)
				return snackStyle.Layout(gtx)
			}
		}

		if current != nil {
			// LaunchedEffect keyed on current snackbar — handles auto-dismiss after duration
			effect.LaunchedEffect(func(ctx context.Context) {
				dur := current.Visuals.Duration.ToDuration()
				if dur > 0 {
					select {
					case <-time.After(dur):
						current.Dismiss()
						hostState.advanceQueue()
					case <-ctx.Done():
						// Effect cancelled (e.g., recomposition with new key)
					}
				}
				// If duration is 0 (Indefinite), we just wait for explicit dismiss
			}, current)(c)
		}

		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(constructor))
		return c.EndBlock()
	}
}

// SnackbarHostWithContent is a composable that displays snackbars with a custom content builder.
// This matches Kotlin's SnackbarHost(hostState, snackbar = { data -> ... }) pattern.
func SnackbarHostWithContent(hostState *SnackbarHostState, content func(data *SnackbarData) api.Composable) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		c.StartBlock("SnackbarHost")

		current := hostState.CurrentSnackbarData()

		if current != nil {
			// LaunchedEffect for auto-dismiss
			effect.LaunchedEffect(func(ctx context.Context) {
				dur := current.Visuals.Duration.ToDuration()
				if dur > 0 {
					select {
					case <-time.After(dur):
						current.Dismiss()
						hostState.advanceQueue()
					case <-ctx.Done():
					}
				}
			}, current)(c)

			// Render custom content, wrapped in a box to fill max size
			box.Box(
				func(c compose.Composer) compose.Composer {
					content(current)(c)
					return c
				},
				box.WithModifier(size.FillMax()),
			)(c)
		}

		return c.EndBlock()
	}
}
