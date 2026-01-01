package textfield

import (
	"github.com/zodimo/go-compose/compose/foundation/text/input"
	"github.com/zodimo/go-compose/compose/ui/text"
)

// TextFieldOptions configures a BasicTextField.
type TextFieldOptions struct {
	// Modifier to apply to the text field.
	Modifier Modifier

	// Enabled controls the enabled state. When false, the field is not editable
	// or focusable, and input is not selectable.
	Enabled bool

	// ReadOnly controls the editable state. When true, the field cannot be
	// modified but can be focused and text can be copied.
	ReadOnly bool

	// TextStyle configures typography and graphics for the text content.
	TextStyle *text.TextStyle

	// InputTransformation filters/transforms user input after it's received
	// but before it's committed to the TextFieldState.
	InputTransformation input.InputTransformation

	// OutputTransformation transforms text for visual presentation only,
	// without affecting the underlying TextFieldState.
	OutputTransformation input.OutputTransformation

	// LineLimits specifies text wrapping and height behavior.
	LineLimits input.TextFieldLineLimits

	// Decorator allows adding decorations around the text field.
	Decorator input.TextFieldDecorator
}

// TextFieldOption is a functional option for configuring BasicTextField.
type TextFieldOption func(*TextFieldOptions)

// WithModifier sets the modifier for the text field.
func WithModifier(m Modifier) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Modifier = m
	}
}

// WithEnabled sets the enabled state.
func WithEnabled(enabled bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Enabled = enabled
	}
}

// WithReadOnly sets the read-only state.
func WithReadOnly(readOnly bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.ReadOnly = readOnly
	}
}

// WithTextStyle sets the text style.
func WithTextStyle(ts *text.TextStyle) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.TextStyle = ts
	}
}

// WithInputTransformation sets the input transformation.
func WithInputTransformation(t input.InputTransformation) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.InputTransformation = t
	}
}

// WithOutputTransformation sets the output transformation.
func WithOutputTransformation(t input.OutputTransformation) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.OutputTransformation = t
	}
}

// WithLineLimits sets the line limits for wrapping and height.
func WithLineLimits(limits input.TextFieldLineLimits) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.LineLimits = limits
	}
}

// WithDecorator sets the text field decorator.
func WithDecorator(decorator input.TextFieldDecorator) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Decorator = decorator
	}
}
