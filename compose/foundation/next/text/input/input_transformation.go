package input

import (
	"strings"
	"unicode"
)

// InputTransformation filters/transforms user input after it's received.
//
// Input transformations run AFTER user input is received but BEFORE it's
// committed to the TextFieldState. They can:
//   - Reject changes (by calling buffer.RevertAllChanges())
//   - Modify the input (e.g., convert to uppercase)
//   - Adjust the selection
//
// Transformations are called for all user-initiated input:
//   - Keyboard input (hardware and software)
//   - Paste operations
//   - Drag and drop
//   - Accessibility input
//   - IME input
//
// This is a port of androidx.compose.foundation.text.input.InputTransformation.
type InputTransformation interface {
	// TransformInput modifies the buffer to filter/transform input.
	// Call buffer.RevertAllChanges() to reject all changes.
	TransformInput(buffer *TextFieldBuffer)
}

// InputTransformationFunc is a function type that implements InputTransformation.
type InputTransformationFunc func(buffer *TextFieldBuffer)

// TransformInput implements InputTransformation.
func (f InputTransformationFunc) TransformInput(buffer *TextFieldBuffer) {
	f(buffer)
}

// ChainedInputTransformation combines two transformations.
type ChainedInputTransformation struct {
	First  InputTransformation
	Second InputTransformation
}

// TransformInput implements InputTransformation by running first then second.
func (c *ChainedInputTransformation) TransformInput(buffer *TextFieldBuffer) {
	c.First.TransformInput(buffer)
	c.Second.TransformInput(buffer)
}

// ChainTransformations creates a transformation that runs first then second.
func ChainTransformations(first, second InputTransformation) InputTransformation {
	if first == nil {
		return second
	}
	if second == nil {
		return first
	}
	return &ChainedInputTransformation{First: first, Second: second}
}

// Then creates a chained transformation that runs this transformation first,
// then the other transformation.
func (f InputTransformationFunc) Then(other InputTransformation) InputTransformation {
	return ChainTransformations(f, other)
}

// MaxLengthInputTransformation truncates input that exceeds maxLength.
type MaxLengthInputTransformation struct {
	MaxLength int
}

// TransformInput implements InputTransformation.
func (m *MaxLengthInputTransformation) TransformInput(buffer *TextFieldBuffer) {
	if buffer.Length() > m.MaxLength {
		// Truncate to max length
		buffer.Delete(m.MaxLength, buffer.Length())
	}
}

// MaxLengthTransformation creates an InputTransformation that rejects input
// exceeding the specified length.
func MaxLengthTransformation(maxLength int) InputTransformation {
	return &MaxLengthInputTransformation{MaxLength: maxLength}
}

// AllCapsInputTransformation converts inserted text to uppercase.
type AllCapsInputTransformation struct {
	// Locale is currently unused but reserved for locale-aware case conversion.
	Locale string
}

// TransformInput implements InputTransformation.
func (a *AllCapsInputTransformation) TransformInput(buffer *TextFieldBuffer) {
	// Get the changes and convert inserted text to uppercase
	changes := buffer.Changes()
	if changes.ChangeCount() == 0 {
		return
	}

	// We need to convert each changed region to uppercase
	// Since we're in the middle of an edit, we work with the current text
	currentText := buffer.String()
	upperText := strings.ToUpper(currentText)

	if currentText != upperText {
		// Replace with uppercase version, preserving selection
		selection := buffer.Selection()
		buffer.Replace(0, buffer.Length(), upperText)
		buffer.SetSelection(selection.Start, selection.End)
	}
}

// AllCapsTransformation creates an InputTransformation that converts
// inserted text to uppercase.
func AllCapsTransformation() InputTransformation {
	return &AllCapsInputTransformation{}
}

// AllCapsTransformationWithLocale creates an InputTransformation that converts
// inserted text to uppercase using the specified locale.
func AllCapsTransformationWithLocale(locale string) InputTransformation {
	return &AllCapsInputTransformation{Locale: locale}
}

// FilterInputTransformation rejects characters that don't pass the filter.
type FilterInputTransformation struct {
	Filter func(r rune) bool
}

// TransformInput implements InputTransformation.
func (f *FilterInputTransformation) TransformInput(buffer *TextFieldBuffer) {
	currentText := buffer.String()
	var filtered strings.Builder
	filtered.Grow(len(currentText))

	for _, r := range currentText {
		if f.Filter(r) {
			filtered.WriteRune(r)
		}
	}

	filteredStr := filtered.String()
	if currentText != filteredStr {
		selection := buffer.Selection()
		buffer.Replace(0, buffer.Length(), filteredStr)
		// Adjust selection - it may now be beyond the new text length
		newLen := buffer.Length()
		newStart := selection.Start
		newEnd := selection.End
		if newStart > newLen {
			newStart = newLen
		}
		if newEnd > newLen {
			newEnd = newLen
		}
		buffer.SetSelection(newStart, newEnd)
	}
}

// DigitsOnlyTransformation creates an InputTransformation that only allows digits.
func DigitsOnlyTransformation() InputTransformation {
	return &FilterInputTransformation{
		Filter: unicode.IsDigit,
	}
}

// AlphanumericOnlyTransformation creates an InputTransformation that only allows
// letters and digits.
func AlphanumericOnlyTransformation() InputTransformation {
	return &FilterInputTransformation{
		Filter: func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsDigit(r)
		},
	}
}

// ByValueInputTransformation transforms based on comparing current and proposed values.
type ByValueInputTransformation struct {
	Transform func(current, proposed string) string
}

// TransformInput implements InputTransformation.
func (b *ByValueInputTransformation) TransformInput(buffer *TextFieldBuffer) {
	originalText := buffer.OriginalValue().Text()
	proposedText := buffer.String()

	finalText := b.Transform(originalText, proposedText)
	if finalText != proposedText {
		selection := buffer.Selection()
		buffer.Replace(0, buffer.Length(), finalText)
		// Adjust selection to be within bounds
		newLen := buffer.Length()
		newStart := selection.Start
		newEnd := selection.End
		if newStart > newLen {
			newStart = newLen
		}
		if newEnd > newLen {
			newEnd = newLen
		}
		buffer.SetSelection(newStart, newEnd)
	}
}

// InputTransformationByValue creates a transformation from a comparison function.
// The function receives the current text and the proposed text (after user input),
// and returns the final text that should be used.
func InputTransformationByValue(transform func(current, proposed string) string) InputTransformation {
	return &ByValueInputTransformation{Transform: transform}
}
