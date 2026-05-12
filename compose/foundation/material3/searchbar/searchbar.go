package searchbar

import (
	"fmt"

	"git.sr.ht/~schnwalter/gio-mw/widget/search"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

const Material3SearchBarNodeID = "Material3SearchBar"

type HandlerWrapper struct {
	Func func(string)
}

type OnSubmitWrapper struct {
	Func func()
}

type SearchBarStateTracker struct {
	LastValue string
}

// SearchBar implements a Material Design 3 search bar.
// It wraps gio-mw's widget/search.
func SearchBar(
	value string,
	onValueChange func(string),
	options ...SearchBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultSearchBarOptions()
		for _, opt := range options {
			opt(&opts)
		}

		key := c.GenerateID()
		path := c.GetPath()

		// Handler wrapper
		handlerWrapperState := c.State(fmt.Sprintf("%d/%s/handler_wrapper", key, path), func() any {
			return &HandlerWrapper{Func: onValueChange}
		})
		handlerWrapper := handlerWrapperState.Get().(*HandlerWrapper)
		handlerWrapper.Func = onValueChange

		// OnSubmit wrapper
		var onSubmitWrapper *OnSubmitWrapper
		if opts.OnSubmit != nil {
			onSubmitWrapperState := c.State(fmt.Sprintf("%d/%s/onsubmit_wrapper", key, path), func() any {
				return &OnSubmitWrapper{Func: opts.OnSubmit}
			})
			onSubmitWrapper = onSubmitWrapperState.Get().(*OnSubmitWrapper)
			onSubmitWrapper.Func = opts.OnSubmit
		}

		// Search widget state
		searchVal := c.State(fmt.Sprintf("%d/%s/searchbar", key, path), func() any {
			return search.Bar()
		})
		searchWidget := searchVal.Get().(*search.Search)

		// State tracker
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker", key, path), func() any {
			return &SearchBarStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*SearchBarStateTracker)

		// Apply properties
		searchWidget.SupportingText = opts.SupportingText
		searchWidget.LeadingIcon = opts.LeadingIcon
		searchWidget.TrailingIcon = opts.TrailingIcon
		// searchWidget.Disabled = !opts.Enabled // Not available in Search struct based on godoc

		c.StartBlock(Material3SearchBarNodeID)
		c.Modifier(func(m Modifier) Modifier {
			return m.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(searchBarWidgetConstructor(searchWidget, value, opts, handlerWrapper, onSubmitWrapper, tracker))

		return c.EndBlock()
	}
}

func searchBarWidgetConstructor(
	s *search.Search,
	value string,
	opts SearchBarOptions,
	handler *HandlerWrapper,
	onSubmitHandler *OnSubmitWrapper,
	tracker *SearchBarStateTracker,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Sync External Change
			if value != tracker.LastValue {
				if s.GetText() != value {
					s.SetText(value)
				}
				tracker.LastValue = value
			}

			// 2. Drive Events
			if s.Submitted(gtx) {
				if onSubmitHandler != nil && onSubmitHandler.Func != nil {
					onSubmitHandler.Func()
				}
			}

			// 3. Detect User Change
			newText := s.GetText()
			if newText != value && opts.Enabled {
				if handler.Func != nil {
					handler.Func(newText)
				}
			}

			// 4. Layout
			return s.Layout(gtx)
		}
	})
}
