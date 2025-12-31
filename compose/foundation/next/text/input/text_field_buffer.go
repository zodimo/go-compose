package input

import (
	"strings"
	"unicode/utf8"

	"github.com/zodimo/go-compose/compose/foundation/next/text/input/internal"
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// TextFieldBuffer is a mutable buffer for text field editing operations.
//
// It provides methods for replacing, inserting, deleting, and appending text,
// as well as managing the selection/cursor state. All changes are tracked
// for undo/redo and input transformation purposes.
//
// TextFieldBuffer is used within TextFieldState.Edit() blocks to perform
// atomic editing operations. It should not be retained across edit sessions.
//
// This is a port of androidx.compose.foundation.text.input.TextFieldBuffer.
type TextFieldBuffer struct {
	buffer                  *internal.PartialGapBuffer
	changeTracker           *internal.ChangeTracker
	offsetMappingCalculator *internal.OffsetMappingCalculator
	originalValue           *TextFieldCharSequence

	selection            text.TextRange
	composition          *text.TextRange
	highlight            *TextHighlight
	composingAnnotations []PlacedAnnotation
}

// NewTextFieldBuffer creates a new TextFieldBuffer initialized with the given value.
func NewTextFieldBuffer(initialValue *TextFieldCharSequence) *TextFieldBuffer {
	return NewTextFieldBufferWithTracker(initialValue, nil, nil)
}

// NewTextFieldBufferWithTracker creates a TextFieldBuffer with an optional change tracker
// and offset mapping calculator.
func NewTextFieldBufferWithTracker(
	initialValue *TextFieldCharSequence,
	initialChanges *internal.ChangeTracker,
	offsetMappingCalculator *internal.OffsetMappingCalculator,
) *TextFieldBuffer {
	return &TextFieldBuffer{
		buffer:                  internal.NewPartialGapBuffer(initialValue.Text()),
		changeTracker:           internal.NewChangeTracker(initialChanges),
		offsetMappingCalculator: offsetMappingCalculator,
		originalValue:           initialValue,
		selection:               initialValue.Selection(),
		composition:             initialValue.Composition(),
		highlight:               initialValue.Highlight(),
		composingAnnotations:    initialValue.ComposingAnnotations(),
	}
}

// Length returns the number of characters (runes) in the buffer.
func (b *TextFieldBuffer) Length() int {
	return b.buffer.Length()
}

// CharAt returns the rune at the given index.
func (b *TextFieldBuffer) CharAt(index int) rune {
	return b.buffer.Get(index)
}

// String returns the current text content.
func (b *TextFieldBuffer) String() string {
	return b.buffer.String()
}

// Selection returns the current selection range.
func (b *TextFieldBuffer) Selection() text.TextRange {
	return b.selection
}

// SetSelection sets the selection range. Both start and end must be valid.
func (b *TextFieldBuffer) SetSelection(start, end int) {
	b.setSelectionCoerced(start, end)
}

// setSelectionCoerced sets the selection, coercing values to valid range.
func (b *TextFieldBuffer) setSelectionCoerced(start, end int) {
	length := b.Length()
	coercedStart := coerceIn(start, 0, length)
	coercedEnd := coerceIn(end, 0, length)
	b.selection = text.NewTextRange(coercedStart, coercedEnd)
}

// Composition returns the current composition range (IME input in progress).
func (b *TextFieldBuffer) Composition() *text.TextRange {
	return b.composition
}

// SetComposition sets or clears the composition range.
func (b *TextFieldBuffer) SetComposition(composition *text.TextRange) {
	if composition == nil {
		b.composition = nil
		return
	}
	coerced := composition.CoerceIn(0, b.Length())
	b.composition = &coerced
}

// Highlight returns the current highlight for handwriting gesture previews.
func (b *TextFieldBuffer) Highlight() *TextHighlight {
	return b.highlight
}

// SetHighlight sets or clears the highlight.
func (b *TextFieldBuffer) SetHighlight(highlight *TextHighlight) {
	if highlight == nil {
		b.highlight = nil
		return
	}
	coercedRange := highlight.Range.CoerceIn(0, b.Length())
	b.highlight = &TextHighlight{
		Type:  highlight.Type,
		Range: coercedRange,
	}
}

// ComposingAnnotations returns annotations attached to the composing region.
func (b *TextFieldBuffer) ComposingAnnotations() []PlacedAnnotation {
	return b.composingAnnotations
}

// SetComposingAnnotations sets the annotations for the composing region.
func (b *TextFieldBuffer) SetComposingAnnotations(annotations []PlacedAnnotation) {
	b.composingAnnotations = annotations
}

// OriginalValue returns the original value before any modifications.
func (b *TextFieldBuffer) OriginalValue() *TextFieldCharSequence {
	return b.originalValue
}

// OriginalSelection returns the selection from the original value.
func (b *TextFieldBuffer) OriginalSelection() text.TextRange {
	return b.originalValue.Selection()
}

// HasSelection returns true if text is selected (selection is not collapsed).
func (b *TextFieldBuffer) HasSelection() bool {
	return !b.selection.Collapsed()
}

// Changes returns the list of changes made to this buffer.
func (b *TextFieldBuffer) Changes() internal.ChangeList {
	return b.changeTracker
}

// Replace replaces text from start (inclusive) to end (exclusive) with the given text.
// Selection is adjusted to account for the text change.
func (b *TextFieldBuffer) Replace(start, end int, replacement string) {
	b.replaceInternal(start, end, replacement)
}

// replaceInternal performs the actual replacement with change tracking.
func (b *TextFieldBuffer) replaceInternal(start, end int, replacement string) {
	// Calculate replacement metrics
	preLength := end - start
	postLength := utf8.RuneCountInString(replacement)

	// Record the change
	b.changeTracker.TrackChange(start, end, postLength)

	if b.offsetMappingCalculator != nil {
		b.offsetMappingCalculator.RecordEditOperation(start, end, postLength)
	}

	// Perform the replacement
	b.buffer.Replace(start, end, replacement)

	// Adjust selection
	b.selection = adjustTextRange(b.selection, start, preLength, postLength)

	// Clear composition if text changes
	if preLength > 0 || postLength > 0 {
		b.composition = nil
		b.highlight = nil
	}
}

// Append adds text at the end of the buffer.
func (b *TextFieldBuffer) Append(text string) *TextFieldBuffer {
	length := b.Length()
	b.Replace(length, length, text)
	return b
}

// Insert adds text at the given index.
func (b *TextFieldBuffer) Insert(index int, text string) {
	b.Replace(index, index, text)
}

// Delete removes text from start to end.
func (b *TextFieldBuffer) Delete(start, end int) {
	b.Replace(start, end, "")
}

// DeleteSelectedText deletes the currently selected text.
// Returns true if text was deleted.
func (b *TextFieldBuffer) DeleteSelectedText() bool {
	if b.selection.Collapsed() {
		return false
	}
	b.Delete(b.selection.Min(), b.selection.Max())
	return true
}

// PlaceCursorBeforeCharAt places the cursor before the character at the given index.
func (b *TextFieldBuffer) PlaceCursorBeforeCharAt(index int) {
	b.setSelectionCoerced(index, index)
}

// PlaceCursorAfterCharAt places the cursor after the character at the given index.
func (b *TextFieldBuffer) PlaceCursorAfterCharAt(index int) {
	b.setSelectionCoerced(index+1, index+1)
}

// PlaceCursorAtEnd places the cursor at the end of the text.
func (b *TextFieldBuffer) PlaceCursorAtEnd() {
	length := b.Length()
	b.setSelectionCoerced(length, length)
}

// SelectCharsIn selects the characters in the given range.
func (b *TextFieldBuffer) SelectCharsIn(rangeValue text.TextRange) {
	b.setSelectionCoerced(rangeValue.Start, rangeValue.End)
}

// SelectAll selects all text in the buffer.
func (b *TextFieldBuffer) SelectAll() {
	b.setSelectionCoerced(0, b.Length())
}

// RevertAllChanges restores the buffer to its original state.
func (b *TextFieldBuffer) RevertAllChanges() {
	b.Replace(0, b.Length(), b.originalValue.Text())
	b.selection = b.originalValue.Selection()
	b.composition = b.originalValue.Composition()
	b.highlight = b.originalValue.Highlight()
}

// ToTextFieldCharSequence creates an immutable snapshot of the current state.
func (b *TextFieldBuffer) ToTextFieldCharSequence() *TextFieldCharSequence {
	return NewTextFieldCharSequenceFull(
		b.String(),
		b.selection,
		b.composition,
		b.highlight,
		b.composingAnnotations,
		nil, // outputAnnotations are set by OutputTransformation
	)
}

// ContentEquals returns true if the content equals the given string.
func (b *TextFieldBuffer) ContentEquals(other string) bool {
	return b.buffer.ContentEquals(other)
}

// SubSequence returns a substring from start to end.
func (b *TextFieldBuffer) SubSequence(start, end int) string {
	return b.buffer.SubSequence(start, end)
}

// adjustTextRange adjusts a text range after a replacement operation.
func adjustTextRange(rangeValue text.TextRange, replaceStart, preLength, postLength int) text.TextRange {
	replaceEnd := replaceStart + preLength
	delta := postLength - preLength

	newStart := adjustOffset(rangeValue.Start, replaceStart, replaceEnd, delta)
	newEnd := adjustOffset(rangeValue.End, replaceStart, replaceEnd, delta)

	return text.NewTextRange(newStart, newEnd)
}

// adjustOffset adjusts a single offset after a replacement operation.
func adjustOffset(offset, replaceStart, replaceEnd, delta int) int {
	if offset < replaceStart {
		return offset
	}
	if offset > replaceEnd {
		return offset + delta
	}
	// Offset is within the replaced range - clamp to end of replacement
	return replaceStart + max(0, delta+offset-replaceEnd+preLength(replaceStart, replaceEnd))
}

// preLength helper for adjustOffset calculation
func preLength(start, end int) int {
	return end - start
}

// coerceIn clamps value to [minVal, maxVal].
func coerceIn(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetTextBeforeSelection returns the text before the selection.
func (b *TextFieldBuffer) GetTextBeforeSelection(maxChars int) string {
	start := b.selection.Min() - maxChars
	if start < 0 {
		start = 0
	}
	return b.SubSequence(start, b.selection.Min())
}

// GetTextAfterSelection returns the text after the selection.
func (b *TextFieldBuffer) GetTextAfterSelection(maxChars int) string {
	end := b.selection.Max() + maxChars
	length := b.Length()
	if end > length {
		end = length
	}
	return b.SubSequence(b.selection.Max(), end)
}

// GetSelectedText returns the currently selected text.
func (b *TextFieldBuffer) GetSelectedText() string {
	return b.SubSequence(b.selection.Min(), b.selection.Max())
}

// Implements Appendable interface (like Java's Appendable / Kotlin's Appendable)

// WriteString implements io.StringWriter, appending the string.
func (b *TextFieldBuffer) WriteString(s string) (int, error) {
	b.Append(s)
	return len(s), nil
}

// WriteRune appends a single rune.
func (b *TextFieldBuffer) WriteRune(r rune) (int, error) {
	b.Append(string(r))
	return utf8.RuneLen(r), nil
}

// Write implements io.Writer, appending bytes as UTF-8 string.
func (b *TextFieldBuffer) Write(p []byte) (int, error) {
	b.Append(string(p))
	return len(p), nil
}

// AddStyle adds a SpanStyle annotation to the specified range.
// This allows styling a portion of the text buffer.
func (b *TextFieldBuffer) AddStyle(style text.SpanStyle, start, end int) {
	// For now, we store as composing annotations. In the full implementation,
	// this would integrate with the output annotation system.
	// This is a stub that allows the API to be used.
	_ = style
	_ = start
	_ = end
}

// Builder helpers for fluent text building

// AppendLine appends text followed by a newline.
func (b *TextFieldBuffer) AppendLine(text string) *TextFieldBuffer {
	b.Append(text)
	b.Append("\n")
	return b
}

// Clear removes all text and resets selection to start.
func (b *TextFieldBuffer) Clear() {
	b.Replace(0, b.Length(), "")
}

// ReplaceAll replaces all occurrences of old with new.
func (b *TextFieldBuffer) ReplaceAll(old, new string) int {
	content := b.String()
	count := strings.Count(content, old)
	if count == 0 {
		return 0
	}
	newContent := strings.ReplaceAll(content, old, new)
	b.Replace(0, b.Length(), newContent)
	return count
}
