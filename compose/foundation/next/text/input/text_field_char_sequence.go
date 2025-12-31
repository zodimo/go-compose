package input

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// PlacedAnnotation is an alias for text.Range[text.Annotation].
// This represents an annotation with its position range in the text.
type PlacedAnnotation = text.Range[text.Annotation]

// TextHighlight represents a highlighted range of text.
// This may be used to display handwriting gesture previews from the IME.
type TextHighlight struct {
	Type  TextHighlightType
	Range text.TextRange
}

// TextFieldCharSequence is an immutable snapshot of the contents of a TextFieldState.
//
// This type directly represents the text being edited along with:
//   - selection: The current cursor position or selection range
//   - composition: The range being composed by the IME (if any)
//   - highlight: The highlighted range for handwriting gesture previews (if any)
//   - composingAnnotations: Annotations attached to the composing region
//   - outputAnnotations: Annotations from OutputTransformation
//
// This is a port of androidx.compose.foundation.text.input.TextFieldCharSequence.
type TextFieldCharSequence struct {
	text                 string
	selection            text.TextRange
	composition          *text.TextRange
	highlight            *TextHighlight
	composingAnnotations []PlacedAnnotation
	outputAnnotations    []PlacedAnnotation
}

// NewTextFieldCharSequence creates a new TextFieldCharSequence with the given text and selection.
func NewTextFieldCharSequence(txt string, selection text.TextRange) *TextFieldCharSequence {
	return &TextFieldCharSequence{
		text:      txt,
		selection: selection.CoerceIn(0, len(txt)),
	}
}

// NewTextFieldCharSequenceWithComposition creates a TextFieldCharSequence with composition range.
func NewTextFieldCharSequenceWithComposition(
	txt string,
	selection text.TextRange,
	composition *text.TextRange,
) *TextFieldCharSequence {
	tfcs := NewTextFieldCharSequence(txt, selection)
	if composition != nil {
		coerced := composition.CoerceIn(0, len(txt))
		tfcs.composition = &coerced
	}
	return tfcs
}

// NewTextFieldCharSequenceFull creates a TextFieldCharSequence with all fields.
func NewTextFieldCharSequenceFull(
	txt string,
	selection text.TextRange,
	composition *text.TextRange,
	highlight *TextHighlight,
	composingAnnotations []PlacedAnnotation,
	outputAnnotations []PlacedAnnotation,
) *TextFieldCharSequence {
	tfcs := NewTextFieldCharSequenceWithComposition(txt, selection, composition)
	if highlight != nil {
		coercedRange := highlight.Range.CoerceIn(0, len(txt))
		tfcs.highlight = &TextHighlight{
			Type:  highlight.Type,
			Range: coercedRange,
		}
	}
	tfcs.composingAnnotations = composingAnnotations
	tfcs.outputAnnotations = outputAnnotations
	return tfcs
}

// Text returns the plain text content.
func (t *TextFieldCharSequence) Text() string {
	return t.text
}

// Len returns the length of the text in bytes.
func (t *TextFieldCharSequence) Len() int {
	return len(t.text)
}

// Selection returns the current selection range.
// If the selection is collapsed, it represents the cursor location.
func (t *TextFieldCharSequence) Selection() text.TextRange {
	return t.selection
}

// Composition returns the current composing region dictated by the IME.
// Returns nil if there is no composing region.
func (t *TextFieldCharSequence) Composition() *text.TextRange {
	return t.composition
}

// Highlight returns the current highlight for handwriting gesture previews.
// Returns nil if there is no highlight.
func (t *TextFieldCharSequence) Highlight() *TextHighlight {
	return t.highlight
}

// ComposingAnnotations returns annotations attached to the composing region.
// These are usually styling cues like underline or different background colors.
func (t *TextFieldCharSequence) ComposingAnnotations() []PlacedAnnotation {
	return t.composingAnnotations
}

// OutputAnnotations returns annotations added by OutputTransformation.
func (t *TextFieldCharSequence) OutputAnnotations() []PlacedAnnotation {
	return t.outputAnnotations
}

// CharAt returns the byte at the given index.
func (t *TextFieldCharSequence) CharAt(index int) byte {
	return t.text[index]
}

// SubSequence returns a substring from startIndex (inclusive) to endIndex (exclusive).
func (t *TextFieldCharSequence) SubSequence(startIndex, endIndex int) string {
	return t.text[startIndex:endIndex]
}

// String returns the plain text content (implements fmt.Stringer).
func (t *TextFieldCharSequence) String() string {
	return t.text
}

// ContentEquals returns true if the text content equals the given string.
func (t *TextFieldCharSequence) ContentEquals(other string) bool {
	return t.text == other
}

// ShouldShowSelection returns whether to show the cursor or selection and associated handles.
// When there is a handwriting gesture preview highlight, the cursor or selection should be hidden.
func (t *TextFieldCharSequence) ShouldShowSelection() bool {
	return t.highlight == nil
}

// Equals returns true if other has the same contents, selection, composition, and highlight.
func (t *TextFieldCharSequence) Equals(other *TextFieldCharSequence) bool {
	if t == other {
		return true
	}
	if other == nil {
		return false
	}
	if !t.selection.Equals(other.selection) {
		return false
	}
	if !textRangePtrEquals(t.composition, other.composition) {
		return false
	}
	if !textHighlightPtrEquals(t.highlight, other.highlight) {
		return false
	}
	if !placedAnnotationsEqual(t.composingAnnotations, other.composingAnnotations) {
		return false
	}
	return t.text == other.text
}

// textRangePtrEquals compares two *text.TextRange for equality.
func textRangePtrEquals(a, b *text.TextRange) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(*b)
}

// textHighlightPtrEquals compares two *TextHighlight for equality.
func textHighlightPtrEquals(a, b *TextHighlight) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Type == b.Type && a.Range.Equals(b.Range)
}

// placedAnnotationsEqual compares two slices of PlacedAnnotation for equality.
func placedAnnotationsEqual(a, b []PlacedAnnotation) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Start != b[i].Start || a[i].End != b[i].End || a[i].Tag != b[i].Tag {
			return false
		}
		// Note: We don't deep compare annotation items as they may be complex types
	}
	return true
}

// GetTextBeforeSelection returns the text before the selection.
// maxChars is the maximum number of characters (inclusive) before the selection minimum.
func (t *TextFieldCharSequence) GetTextBeforeSelection(maxChars int) string {
	start := t.selection.Min() - maxChars
	if start < 0 {
		start = 0
	}
	return t.text[start:t.selection.Min()]
}

// GetTextAfterSelection returns the text after the selection.
// maxChars is the maximum number of characters (exclusive) after the selection maximum.
func (t *TextFieldCharSequence) GetTextAfterSelection(maxChars int) string {
	end := t.selection.Max() + maxChars
	if end > len(t.text) {
		end = len(t.text)
	}
	return t.text[t.selection.Max():end]
}

// GetSelectedText returns the currently selected text.
func (t *TextFieldCharSequence) GetSelectedText() string {
	return t.text[t.selection.Min():t.selection.Max()]
}

// ToCharArray copies the contents from [sourceStartIndex, sourceEndIndex) into destination
// starting at destinationOffset.
func (t *TextFieldCharSequence) ToCharArray(
	destination []byte,
	destinationOffset int,
	sourceStartIndex int,
	sourceEndIndex int,
) {
	copy(destination[destinationOffset:], t.text[sourceStartIndex:sourceEndIndex])
}
