package internal

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// TextLayoutState manages text layout for TextField.
//
// This is a stub interface that documents the operational semantics.
// A full implementation requires integration with the text layout engine.
//
// Operational Semantics:
//   - Caches TextLayoutResult for efficient re-layout
//   - Manages coordinate translations between decoration box and inner text
//   - Provides offset↔position conversion for gesture handling
//   - Updates when text, style, or constraints change
//
// Coordinate Spaces:
//   - Window coordinates: Absolute screen position
//   - Decoration coordinates: Relative to the decorated text field container
//   - TextLayout coordinates: Relative to the text layout rectangle
//   - Core coordinates: Internal text rendering coordinates
//
// This is a port of androidx.compose.foundation.text.input.internal.TextLayoutState.
type TextLayoutState interface {
	// UpdateNonMeasureInputs updates properties that don't require re-measurement.
	// Called when text field properties change outside the measure phase.
	//
	// Parameters:
	//   - textFieldState: The transformed text field state
	//   - textStyle: The text style to apply
	//   - singleLine: Whether the field is single-line
	//   - softWrap: Whether to enable soft wrapping
	UpdateNonMeasureInputs(
		textFieldState TransformedTextFieldStateInterface,
		textStyle TextStyle,
		singleLine bool,
		softWrap bool,
	)

	// LayoutWithNewMeasureInputs performs layout and returns the result.
	// Called during measure phase with current constraints.
	//
	// Parameters:
	//   - density: The display density
	//   - layoutDirection: LTR or RTL
	//   - constraints: Size constraints for the layout
	//
	// Returns the text layout result which contains:
	//   - Measured size
	//   - Line information
	//   - Bounding boxes for text ranges
	LayoutWithNewMeasureInputs(
		density float64,
		layoutDirection LayoutDirection,
		constraints Constraints,
	) TextLayoutResult

	// GetOffsetForPosition converts a screen position to a text offset.
	// Returns the character offset closest to the position.
	// Returns -1 if layout has not been performed.
	//
	// Parameters:
	//   - position: Position in decoration box coordinates
	//   - coerceInVisibleBounds: If true, coerce position to visible text bounds
	GetOffsetForPosition(position Offset, coerceInVisibleBounds bool) int

	// IsPositionOnText returns true if the position is over rendered text.
	// Returns false if position is in empty space left/right of text.
	IsPositionOnText(position Offset) bool

	// CoercedInVisibleBoundsOfInputText coerces a position to visible text bounds.
	// Used when clicks happen outside the visible inner text field.
	CoercedInVisibleBoundsOfInputText(offset Offset) Offset
}

// Stub types for dependencies (would be defined in other packages)

// TransformedTextFieldStateInterface is the interface for transformed state.
// See TransformedTextFieldState for the full implementation requirements.
type TransformedTextFieldStateInterface interface {
	OutputText() *text.TextRange
	VisualText() string
}

// TextStyle represents text styling configuration.
// This would be imported from compose/ui/text in a full implementation.
type TextStyle struct {
	// Stub - would include font, size, color, etc.
}

// LayoutDirection represents text direction (LTR or RTL).
type LayoutDirection int

const (
	LayoutDirectionLTR LayoutDirection = iota
	LayoutDirectionRTL
)

// Constraints represents size constraints for layout.
type Constraints struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

// TextLayoutResult contains the result of text layout.
type TextLayoutResult struct {
	// Stub - would include size, line info, bounding boxes, etc.
	Size Size
}

// Size represents dimensions.
type Size struct {
	Width  int
	Height int
}

// Offset represents a 2D position.
type Offset struct {
	X float64
	Y float64
}
