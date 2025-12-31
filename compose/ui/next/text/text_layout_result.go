package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// TextLayoutInput holds the parameters used for computing a text layout result.
//
// Note: The Kotlin source includes deprecated resourceLoader field and
// DeprecatedBridgeFontResourceLoader class. These are omitted in the Go port.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/TextLayoutResult.kt
type TextLayoutInput struct {
	// Text is the annotated string used for computing text layout.
	Text AnnotatedString

	// Style is the text style used for computing this text layout.
	Style TextStyle

	// Placeholders is a list of placeholders inserted into text layout that reserve space
	// to embed icons or custom emojis.
	Placeholders []Range[Placeholder]

	// MaxLines is the maximum number of lines for the text layout.
	MaxLines int

	// SoftWrap indicates whether the text should break at soft line breaks.
	SoftWrap bool

	// Overflow indicates how visual overflow should be handled.
	Overflow style.TextOverFlow

	// Density is the density used for computing this text layout.
	Density unit.Density

	// LayoutDirection is the layout direction used for computing this text layout.
	LayoutDirection unit.LayoutDirection

	// FontFamilyResolver is the font resolver used for computing this text layout.
	// This is a placeholder interface for platform-specific implementations.
	FontFamilyResolver interface{}

	// Constraints are the layout constraints for this text layout.
	Constraints unit.Constraints
}

// NewTextLayoutInput creates a new TextLayoutInput with the given parameters.
func NewTextLayoutInput(
	text AnnotatedString,
	textStyle TextStyle,
	placeholders []Range[Placeholder],
	maxLines int,
	softWrap bool,
	overflow style.TextOverFlow,
	density unit.Density,
	layoutDirection unit.LayoutDirection,
	fontFamilyResolver interface{},
	constraints unit.Constraints,
) TextLayoutInput {
	return TextLayoutInput{
		Text:               text,
		Style:              textStyle,
		Placeholders:       placeholders,
		MaxLines:           maxLines,
		SoftWrap:           softWrap,
		Overflow:           overflow,
		Density:            density,
		LayoutDirection:    layoutDirection,
		FontFamilyResolver: fontFamilyResolver,
		Constraints:        constraints,
	}
}

// Equals checks equality with another TextLayoutInput.
func (i TextLayoutInput) Equals(other TextLayoutInput) bool {
	if i.Text.String() != other.Text.String() {
		return false
	}
	// Note: TextStyle comparison would need proper Equals method
	if i.MaxLines != other.MaxLines {
		return false
	}
	if i.SoftWrap != other.SoftWrap {
		return false
	}
	if i.Overflow != other.Overflow {
		return false
	}
	if i.LayoutDirection != other.LayoutDirection {
		return false
	}
	if i.Constraints != other.Constraints {
		return false
	}
	if len(i.Placeholders) != len(other.Placeholders) {
		return false
	}
	return true
}

// String returns a string representation of TextLayoutInput.
func (i TextLayoutInput) String() string {
	return fmt.Sprintf(
		"TextLayoutInput(text=%s, style=%v, placeholders=%v, maxLines=%d, softWrap=%t, overflow=%s, layoutDirection=%s, constraints=%s)",
		i.Text, i.Style, i.Placeholders, i.MaxLines, i.SoftWrap, i.Overflow, i.LayoutDirection, i.Constraints,
	)
}

// MultiParagraph is the interface for multi-paragraph text layout.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/MultiParagraph.kt
type MultiParagraph interface {
	// LineCount returns the number of lines in the text layout.
	LineCount() int

	// IsLineEllipsized returns true if the given line is ellipsized.
	IsLineEllipsized(lineIndex int) bool

	// DidExceedMaxLines returns true if the text exceeded the maximum number of lines.
	DidExceedMaxLines() bool

	// FirstBaseline returns the distance from the top to the alphabetic baseline of the first line.
	FirstBaseline() float32

	// LastBaseline returns the distance from the top to the alphabetic baseline of the last line.
	LastBaseline() float32

	// Height returns the total height of the text layout.
	Height() float32

	// Width returns the total width of the text layout.
	Width() float32

	// PlaceholderRects returns bounding boxes for placeholders.
	PlaceholderRects() []*geometry.Rect

	// GetLineStart returns the start offset of the given line, inclusive.
	GetLineStart(lineIndex int) int

	// GetLineEnd returns the end offset of the given line.
	// If visibleEnd is true, trailing whitespaces are not counted.
	GetLineEnd(lineIndex int, visibleEnd bool) int

	// GetLineTop returns the top y coordinate of the given line.
	GetLineTop(lineIndex int) float32

	// GetLineBaseline returns the baseline y coordinate of the given line.
	GetLineBaseline(lineIndex int) float32

	// GetLineBottom returns the bottom y coordinate of the given line.
	GetLineBottom(lineIndex int) float32

	// GetLineLeft returns the left x coordinate of the given line.
	GetLineLeft(lineIndex int) float32

	// GetLineRight returns the right x coordinate of the given line.
	GetLineRight(lineIndex int) float32

	// GetLineForOffset returns the line number for the specified text offset.
	GetLineForOffset(offset int) int

	// GetLineForVerticalPosition returns the line number closest to the given vertical position.
	GetLineForVerticalPosition(vertical float32) int

	// GetHorizontalPosition returns the horizontal position for the specified text offset.
	GetHorizontalPosition(offset int, usePrimaryDirection bool) float32

	// GetParagraphDirection returns the text direction of the paragraph containing the offset.
	GetParagraphDirection(offset int) style.ResolvedTextDirection

	// GetBidiRunDirection returns the text direction of the BiDi run at the offset.
	GetBidiRunDirection(offset int) style.ResolvedTextDirection

	// GetOffsetForPosition returns the character offset closest to the given graphical position.
	GetOffsetForPosition(position geometry.Offset) int

	// GetBoundingBox returns the bounding box of the character at the given offset.
	GetBoundingBox(offset int) geometry.Rect

	// GetWordBoundary returns the text range of the word at the given offset.
	GetWordBoundary(offset int) TextRange

	// GetCursorRect returns the rectangle of the cursor area at the given offset.
	GetCursorRect(offset int) geometry.Rect

	// GetPathForRange returns a path that encloses the given text range.
	GetPathForRange(start, end int) graphics.Path
}

// TextLayoutResult holds the result of text layout computation.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/TextLayoutResult.kt
type TextLayoutResult struct {
	layoutInput    TextLayoutInput
	multiParagraph MultiParagraph
	size           unit.IntSize

	// Cached from multiParagraph
	firstBaseline    float32
	lastBaseline     float32
	placeholderRects []*geometry.Rect
}

// NewTextLayoutResult creates a new TextLayoutResult.
func NewTextLayoutResult(layoutInput TextLayoutInput, multiParagraph MultiParagraph, size unit.IntSize) TextLayoutResult {
	result := TextLayoutResult{
		layoutInput:    layoutInput,
		multiParagraph: multiParagraph,
		size:           size,
	}
	if multiParagraph != nil {
		result.firstBaseline = multiParagraph.FirstBaseline()
		result.lastBaseline = multiParagraph.LastBaseline()
		result.placeholderRects = multiParagraph.PlaceholderRects()
	}
	return result
}

// LayoutInput returns the parameters used for computing this text layout result.
func (r TextLayoutResult) LayoutInput() TextLayoutInput {
	return r.layoutInput
}

// MultiParagraph returns the underlying multi-paragraph object.
func (r TextLayoutResult) MultiParagraph() MultiParagraph {
	return r.multiParagraph
}

// Size returns the size of this text layout in pixels.
func (r TextLayoutResult) Size() unit.IntSize {
	return r.size
}

// FirstBaseline returns the distance from the top to the alphabetic baseline of the first line.
func (r TextLayoutResult) FirstBaseline() float32 {
	return r.firstBaseline
}

// LastBaseline returns the distance from the top to the alphabetic baseline of the last line.
func (r TextLayoutResult) LastBaseline() float32 {
	return r.lastBaseline
}

// PlaceholderRects returns the bounding boxes for placeholders.
func (r TextLayoutResult) PlaceholderRects() []*geometry.Rect {
	return r.placeholderRects
}

// DidOverflowHeight returns true if the text is too tall and couldn't fit with given height.
func (r TextLayoutResult) DidOverflowHeight() bool {
	if r.multiParagraph == nil {
		return false
	}
	return r.multiParagraph.DidExceedMaxLines() || float32(r.size.Height) < r.multiParagraph.Height()
}

// DidOverflowWidth returns true if the text is too wide and couldn't fit with given width.
func (r TextLayoutResult) DidOverflowWidth() bool {
	if r.multiParagraph == nil {
		return false
	}
	return float32(r.size.Width) < r.multiParagraph.Width()
}

// HasVisualOverflow returns true if either vertical or horizontal overflow happens.
func (r TextLayoutResult) HasVisualOverflow() bool {
	return r.DidOverflowWidth() || r.DidOverflowHeight()
}

// LineCount returns the number of lines in this text layout.
func (r TextLayoutResult) LineCount() int {
	if r.multiParagraph == nil {
		return 0
	}
	return r.multiParagraph.LineCount()
}

// GetLineStart returns the start offset of the given line, inclusive.
func (r TextLayoutResult) GetLineStart(lineIndex int) int {
	return r.multiParagraph.GetLineStart(lineIndex)
}

// GetLineEnd returns the end offset of the given line.
func (r TextLayoutResult) GetLineEnd(lineIndex int, visibleEnd bool) int {
	return r.multiParagraph.GetLineEnd(lineIndex, visibleEnd)
}

// IsLineEllipsized returns true if the given line is ellipsized.
func (r TextLayoutResult) IsLineEllipsized(lineIndex int) bool {
	if r.multiParagraph == nil {
		return false
	}
	return r.multiParagraph.IsLineEllipsized(lineIndex)
}

// GetLineTop returns the top y coordinate of the given line.
func (r TextLayoutResult) GetLineTop(lineIndex int) float32 {
	return r.multiParagraph.GetLineTop(lineIndex)
}

// GetLineBaseline returns the baseline y coordinate of the given line.
func (r TextLayoutResult) GetLineBaseline(lineIndex int) float32 {
	return r.multiParagraph.GetLineBaseline(lineIndex)
}

// GetLineBottom returns the bottom y coordinate of the given line.
func (r TextLayoutResult) GetLineBottom(lineIndex int) float32 {
	return r.multiParagraph.GetLineBottom(lineIndex)
}

// GetLineLeft returns the left x coordinate of the given line.
func (r TextLayoutResult) GetLineLeft(lineIndex int) float32 {
	return r.multiParagraph.GetLineLeft(lineIndex)
}

// GetLineRight returns the right x coordinate of the given line.
func (r TextLayoutResult) GetLineRight(lineIndex int) float32 {
	return r.multiParagraph.GetLineRight(lineIndex)
}

// GetLineForOffset returns the line number for the specified text offset.
func (r TextLayoutResult) GetLineForOffset(offset int) int {
	return r.multiParagraph.GetLineForOffset(offset)
}

// GetLineForVerticalPosition returns the line number closest to the given vertical position.
func (r TextLayoutResult) GetLineForVerticalPosition(vertical float32) int {
	return r.multiParagraph.GetLineForVerticalPosition(vertical)
}

// GetHorizontalPosition returns the horizontal position for the specified text offset.
func (r TextLayoutResult) GetHorizontalPosition(offset int, usePrimaryDirection bool) float32 {
	return r.multiParagraph.GetHorizontalPosition(offset, usePrimaryDirection)
}

// GetParagraphDirection returns the text direction of the paragraph containing the offset.
func (r TextLayoutResult) GetParagraphDirection(offset int) style.ResolvedTextDirection {
	return r.multiParagraph.GetParagraphDirection(offset)
}

// GetBidiRunDirection returns the text direction of the BiDi run at the offset.
func (r TextLayoutResult) GetBidiRunDirection(offset int) style.ResolvedTextDirection {
	return r.multiParagraph.GetBidiRunDirection(offset)
}

// GetOffsetForPosition returns the character offset closest to the given graphical position.
func (r TextLayoutResult) GetOffsetForPosition(position geometry.Offset) int {
	return r.multiParagraph.GetOffsetForPosition(position)
}

// GetBoundingBox returns the bounding box of the character at the given offset.
func (r TextLayoutResult) GetBoundingBox(offset int) geometry.Rect {
	return r.multiParagraph.GetBoundingBox(offset)
}

// GetWordBoundary returns the text range of the word at the given offset.
func (r TextLayoutResult) GetWordBoundary(offset int) TextRange {
	return r.multiParagraph.GetWordBoundary(offset)
}

// GetCursorRect returns the rectangle of the cursor area at the given offset.
func (r TextLayoutResult) GetCursorRect(offset int) geometry.Rect {
	return r.multiParagraph.GetCursorRect(offset)
}

// GetPathForRange returns a path that encloses the given text range.
func (r TextLayoutResult) GetPathForRange(start, end int) graphics.Path {
	return r.multiParagraph.GetPathForRange(start, end)
}

// Copy creates a copy of the TextLayoutResult with optional overrides.
func (r TextLayoutResult) Copy(layoutInput *TextLayoutInput, size *unit.IntSize) TextLayoutResult {
	newLayoutInput := r.layoutInput
	if layoutInput != nil {
		newLayoutInput = *layoutInput
	}
	newSize := r.size
	if size != nil {
		newSize = *size
	}
	return TextLayoutResult{
		layoutInput:      newLayoutInput,
		multiParagraph:   r.multiParagraph,
		size:             newSize,
		firstBaseline:    r.firstBaseline,
		lastBaseline:     r.lastBaseline,
		placeholderRects: r.placeholderRects,
	}
}

// Equals checks equality with another TextLayoutResult.
func (r TextLayoutResult) Equals(other TextLayoutResult) bool {
	if !r.layoutInput.Equals(other.layoutInput) {
		return false
	}
	if r.size != other.size {
		return false
	}
	if r.firstBaseline != other.firstBaseline {
		return false
	}
	if r.lastBaseline != other.lastBaseline {
		return false
	}
	if len(r.placeholderRects) != len(other.placeholderRects) {
		return false
	}
	return true
}

// String returns a string representation of TextLayoutResult.
func (r TextLayoutResult) String() string {
	return fmt.Sprintf(
		"TextLayoutResult(layoutInput=%s, size=%v, firstBaseline=%f, lastBaseline=%f, placeholderRects=%v)",
		r.layoutInput, r.size, r.firstBaseline, r.lastBaseline, r.placeholderRects,
	)
}
