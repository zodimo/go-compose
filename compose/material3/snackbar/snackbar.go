package snackbar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/material3/tokens"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// Snackbar displays a brief message at the bottom of the screen.
//
// To display a snackbar, use [SnackbarHostState.ShowSnackbar].
func Snackbar(
	data *SnackbarData,
	options ...SnackbarOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultOptions()
		for _, opt := range options {
			if opt != nil {
				opt(&opts)
			}
		}

		theme := material3.Theme(c)
		// colorScheme := theme.ColorScheme()
		typography := theme.Typography()
		snackbarDefaults := SnackbarDefaults(theme)

		// Defaults from Material 3 spec
		containerColor := opts.ContainerColor.TakeOrElse(snackbarDefaults.Color)
		contentColor := opts.ContentColor.TakeOrElse(snackbarDefaults.ContentColor)
		// actionColor := opts.ActionColor.TakeOrElse(snackbarDefaults.ActionColor)
		actionContentColor := opts.ActionContentColor.TakeOrElse(snackbarDefaults.ActionContentColor)
		dismissActionContentColor := opts.DismissActionContentColor.TakeOrElse(snackbarDefaults.DismissActionContentColor)
		shape := shape.TakeOrElseShape(opts.Shape, snackbarDefaults.Shape)

		return box.Box(
			surface.Surface(
				func(c api.Composer) api.Composer {
					return row.Row(
						c.Sequence(
							// Message
							// Weight(1) ensures the text takes up available space, pushing actions to the end.
							column.Column(
								func(c api.Composer) api.Composer {
									text.Text(
										data.Visuals.Message,
										text.WithTextStyle(typography.BodyMedium),
										text.WithColor(contentColor),
									)(c)
									return c
								},
								column.WithModifier(
									weight.Weight(1).
										Then(padding.Vertical(12, 12)). // Vertical padding for single/multi-line
										Then(padding.Horizontal(16, 16)),
								),
								column.WithAlignment(column.Middle), // Center text vertically
							),
							// Actions
							c.When(
								data.Visuals.ActionLabel != "",
								button.Text(
									func() { data.PerformAction() },
									data.Visuals.ActionLabel,
									// button.WithColor(actionColor), // Text button doesn't support container color override via this API yet, and it's usually transparent.
									button.WithContentColor(actionContentColor),
									button.WithModifier(padding.End(8)),
								),
							),
							c.When(
								data.Visuals.WithDismissAction,
								icon.Icon(
									icon.SymbolClose,
									icon.WithColor(dismissActionContentColor),
									icon.WithModifier(padding.End(8).Then(
										clickable.OnClick(func() {
											data.Dismiss()
										}),
									)),
								),
							),
							c.When(
								data.Visuals.ActionLabel == "" && data.Visuals.WithDismissAction == false,
								spacer.Width(8),
							),
						),
						row.WithModifier(
							size.MaxWidth(600),
							// size.FillMaxWidth(), // Don't force full width, let it wrap content up to max.
						),
						row.WithAlignment(row.Middle),
					)(c)
				},
				surface.WithColor(containerColor),
				surface.WithContentColor(contentColor),
				surface.WithShape(shape),
				surface.WithShadowElevation(tokens.ElevationTokens.Level3),
				surface.WithModifier(
					padding.All(8). // Margin around the snackbar
							Then(opts.Modifier),
				),
			),
			box.WithModifier(padding.All(16)),
		)(c)
	}
}

type SnackbarHostState interface {
	ShowSnackbar(message string, options ...SnackbarOption)
	CurrentSnackbarData() *SnackbarData
	isSnackbar()
}

// SnackbarHostState controls the queue and the current Snackbar being shown inside the SnackbarHost.
//
// Guarantees to show at most one snackbar at a time. When ShowSnackbar is called,
// the message is enqueued. A processing goroutine is spun up to display items
// one at a time, waiting for each to be dismissed before showing the next.
// The goroutine exits when the queue is empty.
//
// The currentSnackbar is an observable MutableValueTyped — reading it in a composable
// registers observation, and setting it triggers recomposition.
type snackbarHostState struct {
	// currentSnackbar is observable state matching Kotlin's mutableStateOf<SnackbarData?>(null).
	currentSnackbar state.MutableValueTyped[*SnackbarData]

	coroutineScope effect.CoroutineScope

	mu         sync.Mutex
	queue      []*SnackbarData
	processing bool // true when the queue processing goroutine is running
}

// RememberSnackbarHostState creates a remembered SnackbarHostState tied to the composer lifecycle.
func RememberSnackbarHostState(c compose.Composer) SnackbarHostState {
	key := c.GenerateID()
	path := c.GetPath()

	snackbarDataPath := fmt.Sprintf("%d/%s/snackbarData", key, path)

	currentSnackbar := state.MustRemember(c, snackbarDataPath, func() *SnackbarData {
		return nil
	})
	coroutineScope := effect.RememberCoroutineScope(c)

	snackbarHostStatePath := fmt.Sprintf("%d/%s/snackbarHostState", key, path)
	snackbatHostState := state.MustRemember(c, snackbarHostStatePath, func() *snackbarHostState {
		return &snackbarHostState{
			currentSnackbar: currentSnackbar,
			coroutineScope:  coroutineScope,
		}
	})

	return snackbatHostState.Get()
}

func (s *snackbarHostState) isSnackbar() {}

// ShowSnackbar enqueues a snackbar message and returns immediately (non-blocking).
//
// The context can be used for cancellation — if ctx is cancelled before the snackbar
// is shown, it will be removed from the queue.
//
// Use WithOnResult to receive a callback when the snackbar is dismissed or its action
// is performed.
func (s *snackbarHostState) ShowSnackbar(message string, options ...SnackbarOption) {
	opts := DefaultOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}
	opts.resolveDefaults()

	visuals := SnackbarVisuals{
		Message:           message,
		ActionLabel:       opts.ActionLabel,
		WithDismissAction: opts.WithDismissAction,
		Duration:          opts.Duration,
	}

	data := newSnackbarData(opts.Context, visuals, opts.OnResult)

	// If context is already cancelled, don't enqueue
	if opts.Context.Err() != nil {
		return
	}

	s.mu.Lock()
	s.queue = append(s.queue, data)
	shouldStart := !s.processing
	if shouldStart {
		s.processing = true
	}
	s.mu.Unlock()

	// Watch for context cancellation to remove from queue before being shown
	if opts.Context != nil && opts.Context != context.Background() {
		go func() {
			<-opts.Context.Done()
			s.removeFromQueue(data)
		}()
	}

	// Start the processing goroutine if not already running
	if shouldStart {
		go s.processQueue()
	}
}

// processQueue processes snackbar items one at a time.
// It sets the current snackbar, waits for it to be dismissed (via resultCh),
// then advances to the next item. Exits when the queue is empty.
func (s *snackbarHostState) processQueue() {
	for {
		select {
		case <-s.coroutineScope.Context().Done():
			s.mu.Lock()
			s.processing = false
			s.currentSnackbar.Set(nil)
			s.mu.Unlock()
			return
		default:
		}
		s.mu.Lock()
		if len(s.queue) == 0 {
			s.processing = false
			s.currentSnackbar.Set(nil)
			s.mu.Unlock()
			return
		}
		// Dequeue next item
		next := s.queue[0]
		s.queue = s.queue[1:]
		s.mu.Unlock()

		// Set as current — triggers recomposition via observable state
		s.currentSnackbar.Set(next)

		// Wait for this snackbar to be dismissed or actioned
		<-next.resultCh

		// Clear current before potentially setting the next one
		s.currentSnackbar.Set(nil)
	}
}

// removeFromQueue removes a specific snackbar from the queue if it hasn't been shown yet.
func (s *snackbarHostState) removeFromQueue(data *SnackbarData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.queue {
		if d == data {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return
		}
	}
}

// CurrentSnackbarData returns the currently displayed snackbar data, or nil if none.
// Reading this in a composable registers observation — recomposition will be triggered
// when the current snackbar changes.
func (s *snackbarHostState) CurrentSnackbarData() *SnackbarData {
	return s.currentSnackbar.Get()
}

// SnackbarHost is a composable that displays snackbars from the given SnackbarHostState.
//
// It shows at most one snackbar at a time. Duration-based auto-dismiss is managed
// via LaunchedEffect keyed on the current snackbar.
func SnackbarHost(hostState SnackbarHostState) compose.Composable {
	return func(c compose.Composer) compose.Composer {

		// Reading CurrentSnackbarData() calls Get() on the MutableValueTyped,
		// which triggers NotifyRead — the compose framework now observes this state.
		current := hostState.CurrentSnackbarData()

		if current != nil {
			// LaunchedEffect keyed on current snackbar — handles auto-dismiss after duration
			effect.LaunchedEffect(func(ctx context.Context) {
				dur := current.Visuals.Duration.ToDuration()
				if dur > 0 && current.Context.Err() == nil {
					select {
					case <-time.After(dur):
						current.Dismiss()
					case <-current.Context.Done():
						current.Dismiss()
					case <-ctx.Done():
						// Effect cancelled (e.g., recomposition with new key)
					}
				}
				// If duration is 0 (Indefinite), we just wait for explicit dismiss or cancelled context

				select {
				case <-current.Context.Done():
					current.Dismiss()
				case <-current.resultCh:
					// noop
				}

			}, current)(c)

			// Wrap in a Box aligned to bottom-center (S), matching Material 3 spec
			return box.Box(
				func(c compose.Composer) compose.Composer {
					return Snackbar(current)(c)
				},
				box.WithAlignment(box.S),
				box.WithModifier(size.FillMax()),
			)(c)
		}

		return c
	}
}

// SnackbarHostWithContent is a composable that displays snackbars with a custom content builder.
// This matches Kotlin's SnackbarHost(hostState, snackbar = { data -> ... }) pattern.
func SnackbarHostWithContent(hostState SnackbarHostState, content func(data *SnackbarData) api.Composable) compose.Composable {
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
					case <-ctx.Done():
					}
				}
			}, current)(c)

			// Render custom content, wrapped in a box aligned to bottom-center
			box.Box(
				func(c compose.Composer) compose.Composer {
					content(current)(c)
					return c
				},
				box.WithAlignment(box.S),
				box.WithModifier(size.FillMax()),
			)(c)
		}

		return c.EndBlock()
	}
}
