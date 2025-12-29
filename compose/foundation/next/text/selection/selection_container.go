package selection

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/modifier"
)

// Composable represents a composable function type.
type Composable = compose.Composable

type SectionContainerOptions struct {
	Modifier modifier.Modifier
}
type SectionContainerOption func(*SectionContainerOptions)

func WithModifier(modifier modifier.Modifier) SectionContainerOption {
	return func(options *SectionContainerOptions) {
		options.Modifier = modifier
	}
}

// SelectionContainer enables text selection for its direct or indirect children.
//
// Use of a lazy layout, such as LazyRow or LazyColumn, within a SelectionContainer
// has undefined behavior on text items that aren't composed. For example, texts that
// aren't composed will not be included in copy operations and select all will not
// expand the selection to include them.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionContainer.kt
func SelectionContainer(content Composable, options ...SectionContainerOption) Composable {
	opts := &SectionContainerOptions{
		Modifier: modifier.EmptyModifier,
	}
	for _, option := range options {
		option(opts)
	}

	return func(c compose.Composer) compose.Composer {
		// Create a mutable state for selection
		var selection *Selection

		c.StartBlock("SelectionContainer")

		// Call internal selection container with state management
		internalSelectionContainer(
			c,
			opts.Modifier,
			selection,
			func(newSelection *Selection) {
				selection = newSelection
			},
			content,
		)

		return c.EndBlock()
	}
}

// internalSelectionContainer is the internal implementation of SelectionContainer.
//
// The selection composable wraps composables and lets them be selectable.
// It paints the selection area with start and end handles.
func internalSelectionContainer(
	c compose.Composer,
	mod modifier.Modifier,
	selection *Selection,
	onSelectionChange func(*Selection),
	children Composable,
) {
	// Create selection registrar
	// In production, this would use rememberSaveable
	registrar := NewSelectionRegistrarImpl()

	// Create selection manager
	manager := NewSelectionManager(registrar)
	manager.SetOnSelectionChange(onSelectionChange)
	manager.SetSelection(selection)

	// Get clipboard and haptic feedback from composition locals
	// manager.HapticFeedback = LocalHapticFeedback.current
	// manager.OnCopyHandler = ...
	// manager.TextToolbar = LocalTextToolbar.current

	// Apply modifier
	c.Modifier(func(m modifier.Modifier) modifier.Modifier {
		return m.Then(mod).Then(manager.Modifier())
	})

	// Provide the selection registrar to children
	compose.CompositionLocalProvider1(
		LocalSelectionRegistrar,
		SelectionRegistrar(registrar),
		func(c compose.Composer) compose.Composer {
			// Render children
			children(c)

			// Render selection handles if in touch mode with active selection
			if manager.IsInTouchMode() && manager.HasFocus() && !manager.IsTriviallyCollapsedSelection() {
				sel := manager.Selection()
				if sel != nil {
					// Render start handle
					renderSelectionHandle(c, manager, true, sel)
					// Render end handle
					renderSelectionHandle(c, manager, false, sel)
				}
			}

			return c
		},
	)(c)
}

// renderSelectionHandle renders a selection handle.
// This is a placeholder that will be expanded with actual handle rendering.
func renderSelectionHandle(c compose.Composer, manager *SelectionManager, isStartHandle bool, selection *Selection) {
	// TODO: Implement handle rendering using SelectionHandle composable
	// This requires:
	// - Getting handle position from manager
	// - Getting direction from selection
	// - Creating pointer input handler for dragging
}
