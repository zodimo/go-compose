package textfield

import (
	"fmt"

	"git.sr.ht/~schnwalter/gio-mw/widget/input"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

const Material3FilledTextFieldNodeID = "Material3FilledTextField"

// Filled implements the Filled Material Design 3 text field.
// It wraps gio-mw's FilledTextInput/FilledTextArea.
func Filled(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {

	opts := DefaultTextFieldOptions()
	for _, opt := range options {
		opt(&opts)
	}

	opts.Colors = ResolveTextFieldColors(opts.Colors)

	return func(c Composer) Composer {

		key := c.GenerateID()
		path := c.GetPath()

		// Handler wrapper to avoid closure capture issues
		handlerWrapperState := c.State(fmt.Sprintf("%d/%s/handler_wrapper", key, path), func() any {
			return &HandlerWrapper{Func: onValueChange}
		})
		handlerWrapper := handlerWrapperState.Get().(*HandlerWrapper)
		handlerWrapper.Func = onValueChange

		// OnSubmit wrapper to avoid closure capture issues
		var onSubmitWrapper *OnSubmitWrapper
		if opts.OnSubmit != nil {
			onSubmitWrapperState := c.State(fmt.Sprintf("%d/%s/onsubmit_wrapper", key, path), func() any {
				return &OnSubmitWrapper{Func: opts.OnSubmit}
			})
			onSubmitWrapper = onSubmitWrapperState.Get().(*OnSubmitWrapper)
			onSubmitWrapper.Func = opts.OnSubmit
		}

		// Input widget state
		// We include SingleLine in key to recreate widget if mode changes
		inputStatePath := fmt.Sprintf("%d/%s/textfield/s%v", key, path, opts.SingleLine)
		inputVal := c.State(inputStatePath, func() any {
			if opts.SingleLine {
				return input.FilledTextInput()
			}
			return input.FilledTextArea()
		})
		inp := inputVal.Get().(*input.Input)

		// State tracker for synchronization
		// We tie the tracker to the same key as input so they reset together
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker/s%v", key, path, opts.SingleLine), func() any {
			// Initialize with empty string assuming fresh Editor starts empty
			return &TextFieldStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*TextFieldStateTracker)

		// Apply properties to the widget
		// These assignments are cheap and happen during composition
		inp.LabelText = opts.Label
		inp.SupportingText = sentinel.TakeOrElseString(opts.SupportingText, "")
		inp.Disabled = !opts.Enabled
		inp.Error = opts.IsError

		c.StartBlock(Material3FilledTextFieldNodeID)
		c.Modifier(func(m Modifier) Modifier {
			return m.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(textFieldWidgetConstructor(inp, value, opts, handlerWrapper, onSubmitWrapper, tracker))

		return c.EndBlock()
	}
}

func textFieldWidgetConstructor(
	inp *input.Input,
	value string,
	opts TextFieldOptions,
	handler *HandlerWrapper,
	onSubmitHandler *OnSubmitWrapper,
	tracker *TextFieldStateTracker,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Sync External Change
			// If the incoming value property differs from what we Last saw,
			// it means there was an external update (or initialization).
			if value != tracker.LastValue {
				// Only call SetText if the content actually differs.
				// This prevents resetting the cursor position if the editor already contains the text
				// (e.g. from user input that just propagated back to us).
				if inp.Editor.GetText() != value {
					inp.Editor.SetText(value)
				}
				tracker.LastValue = value
			}

			// 2. Drive Editor Events
			// Submitted() processes all events for the editor and returns true if Enter was pressed.
			if inp.Editor.Submitted(gtx) {
				if onSubmitHandler != nil && onSubmitHandler.Func != nil {
					onSubmitHandler.Func()
				}
			}

			// 3. Detect User Change
			// We compare the editor's current internal text with the prop value.
			newText := inp.Editor.GetText()
			if newText != value {
				if handler.Func != nil {
					handler.Func(newText)
				}
				// We do NOT update tracker.LastValue here.
				// We wait for the value to come back via the prop 'value'.
				// This avoids "fighting" or oscillation during the frame latency window.
			}

			// 4. Layout
			return inp.Layout(gtx)
		}
	})
}
