package textfield

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/textinput"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"gioui.org/widget/material"
)

const Material3OutlinedTextFieldNodeID = "Material3OutlinedTextField"

// Outlined implements the Outlined Material Design 3 text field.
// It uses a custom widget implementation adapted from gio-x.
func Outlined(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {

	opts := DefaultTextFieldOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {

		opts.Colors = ResolveTextFieldColors(c, opts.Colors)
		opts.SupportingText = sentinel.TakeOrElseString(opts.SupportingText, "")

		textStyle := material3.LocalTextStyle.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)
		textStyle = text.TextStyleResolveDefaults(textStyle, layoutDirection)
		_ = textStyle // used for future text styling

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

		// Custom Outlined Widget State
		widgetStatePath := fmt.Sprintf("%d/%s/outlined_widget/s%v", key, path, opts.SingleLine)
		editorVal := c.State(widgetStatePath, func() any {
			return textinput.NewEditor(opts.SingleLine, opts.OnSubmit != nil, 0)
		})
		outEditor := editorVal.Get().(*textinput.Editor)

		// State tracker for synchronization
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker/s%v", key, path, opts.SingleLine), func() any {
			return &TextFieldStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*TextFieldStateTracker)

		args := TextEditorConstructorArgs{
			Editor:          outEditor,
			Value:           value,
			Opts:            opts,
			Handler:         handlerWrapper,
			OnSubmitHandler: onSubmitWrapper,
			Tracker:         tracker,
		}

		c.StartBlock(Material3OutlinedTextFieldNodeID)
		c.Modifier(func(m ui.Modifier) ui.Modifier {
			return m.Then(opts.Modifier)
		})

		// Constructor
		c.SetWidgetConstructor(outlinedTextFieldWidgetConstructor(args))

		return c.EndBlock()
	}
}

type TextEditorConstructorArgs struct {
	Editor          *textinput.Editor
	Value           string
	Opts            TextFieldOptions
	Handler         *HandlerWrapper
	OnSubmitHandler *OnSubmitWrapper
	Tracker         *TextFieldStateTracker
}

func outlinedTextFieldWidgetConstructor(args TextEditorConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {

	w := &TextFieldWidget{
		Editor: args.Editor,
	}

	// Update static properties
	w.Editor.SetSingleLine(args.Opts.SingleLine)
	w.Editor.SetSubmit(args.Opts.OnSubmit != nil)
	w.Helper = args.Opts.SupportingText
	w.SetError(args.Opts.IsError, args.Opts.SupportingText)

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Sync External Change
			if args.Value != args.Tracker.LastValue {
				if w.Editor.Text() != args.Value {
					w.Editor.SetText(args.Value)
				}
				args.Tracker.LastValue = args.Value
			}

			// 2. Events & Layout
			th := material.NewTheme() // Fallback TODO: real theme
			thWrapper := textinput.NewTheme(th)

			// Check for submit events
			hasSubmit := w.Editor.ProcessEvents(gtx)
			if hasSubmit {
				if args.OnSubmitHandler != nil && args.OnSubmitHandler.Func != nil {
					args.OnSubmitHandler.Func()
				}
			}

			// Check for text changes
			currentText := w.Editor.Text()
			if currentText != args.Value {
				if args.Handler.Func != nil {
					args.Handler.Func(currentText)
				}
			}

			w.Colors = args.Opts.Colors

			return w.Layout(gtx, thWrapper, args.Opts.Label)
		}
	})
}
