package input

// TextFieldDecorator adds decorations around the text field.
//
// A decorator wraps the inner text field with custom UI elements like:
//   - Leading/trailing icons
//   - Labels and placeholder text
//   - Helper/error messages
//   - Borders and backgrounds
//
// The decorator automatically expands the hit target area of the text field
// to include the decorated region.
//
// This is a port of androidx.compose.foundation.text.input.TextFieldDecorator.
//
// Operational Semantics:
//   - The Decoration function receives an innerTextField callback that renders
//     the actual text editing surface.
//   - The innerTextField MUST be called exactly once within Decoration.
//   - Calling innerTextField zero times or more than once is an error.
//   - The decorator controls the layout and positioning of the inner text field
//     relative to decorations.
//   - Touch events on any part of the decorated area will focus the text field.
//
// Example usage (pseudo-code):
//
//	decorator := func(innerTextField func()) {
//	    // Render leading icon
//	    renderIcon()
//	    // Render the actual text field
//	    innerTextField()
//	    // Render trailing icon
//	    renderClearButton()
//	}
type TextFieldDecorator interface {
	// Decoration renders decorations around the inner text field.
	//
	// Parameters:
	//   - innerTextField: A callback that renders the core text editing surface.
	//     This MUST be called exactly once.
	//
	// The decorator controls the layout of the inner text field relative to
	// any decorations (icons, labels, etc.).
	Decoration(innerTextField Composable) Composable
}

// TextFieldDecoratorFunc is a function type that implements TextFieldDecorator.
type TextFieldDecoratorFunc func(innerTextField Composable) Composable

// Decoration implements TextFieldDecorator.
func (f TextFieldDecoratorFunc) Decoration(innerTextField Composable) Composable {
	return f(innerTextField)
}

// NoDecorationDecorator renders only the inner text field with no decorations.
var NoDecorationDecorator TextFieldDecorator = TextFieldDecoratorFunc(func(innerTextField Composable) Composable {
	return innerTextField
})
