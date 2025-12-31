package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

// PlaceholderVerticalAlign specifies how a placeholder is vertically aligned within a text line.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/Placeholder.kt
type PlaceholderVerticalAlign int

const (
	// PlaceholderVerticalAlignAboveBaseline aligns the bottom of the placeholder with the baseline.
	PlaceholderVerticalAlignAboveBaseline PlaceholderVerticalAlign = iota + 1

	// PlaceholderVerticalAlignTop aligns the top of the placeholder with the top of the entire line.
	PlaceholderVerticalAlignTop

	// PlaceholderVerticalAlignBottom aligns the bottom of the placeholder with the bottom of the entire line.
	PlaceholderVerticalAlignBottom

	// PlaceholderVerticalAlignCenter aligns the center of the placeholder with the center of the entire line.
	PlaceholderVerticalAlignCenter

	// PlaceholderVerticalAlignTextTop aligns the top of the placeholder with the top of the proceeding text.
	// Different from Top when there are texts with different font sizes in the same line.
	PlaceholderVerticalAlignTextTop

	// PlaceholderVerticalAlignTextBottom aligns the bottom of the placeholder with the bottom of the proceeding text.
	// Different from Bottom when there are texts with different font sizes in the same line.
	PlaceholderVerticalAlignTextBottom

	// PlaceholderVerticalAlignTextCenter aligns the center of the placeholder with the center of the proceeding text.
	// Different from Center when there are texts with different font sizes in the same line.
	PlaceholderVerticalAlignTextCenter
)

// String returns the string representation of the PlaceholderVerticalAlign.
func (p PlaceholderVerticalAlign) String() string {
	switch p {
	case PlaceholderVerticalAlignAboveBaseline:
		return "AboveBaseline"
	case PlaceholderVerticalAlignTop:
		return "Top"
	case PlaceholderVerticalAlignBottom:
		return "Bottom"
	case PlaceholderVerticalAlignCenter:
		return "Center"
	case PlaceholderVerticalAlignTextTop:
		return "TextTop"
	case PlaceholderVerticalAlignTextBottom:
		return "TextBottom"
	case PlaceholderVerticalAlignTextCenter:
		return "TextCenter"
	default:
		return "Invalid"
	}
}

// Placeholder is a rectangle box inserted into text, which tells the text processor
// to leave an empty space. It is typically used to insert inline images, custom emojis,
// etc. into the text paragraph.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/Placeholder.kt
type Placeholder struct {
	// Width of the placeholder, must be specified in sp or em.
	Width unit.TextUnit

	// Height of the placeholder, must be specified in sp or em.
	Height unit.TextUnit

	// PlaceholderVerticalAlign specifies the vertical alignment of the placeholder within the text line.
	PlaceholderVerticalAlign PlaceholderVerticalAlign
}

// NewPlaceholder creates a new Placeholder with validation.
// Panics if width or height is unspecified.
func NewPlaceholder(width, height unit.TextUnit, placeholderVerticalAlign PlaceholderVerticalAlign) Placeholder {
	if !width.IsSpecified() {
		panic("width cannot be TextUnit.Unspecified")
	}
	if !height.IsSpecified() {
		panic("height cannot be TextUnit.Unspecified")
	}
	return Placeholder{
		Width:                    width,
		Height:                   height,
		PlaceholderVerticalAlign: placeholderVerticalAlign,
	}
}

// Copy creates a copy of the Placeholder with optional field overrides.
func (p Placeholder) Copy(width, height unit.TextUnit, placeholderVerticalAlign PlaceholderVerticalAlign) Placeholder {
	return Placeholder{
		Width:                    width,
		Height:                   height,
		PlaceholderVerticalAlign: placeholderVerticalAlign,
	}
}

// Equals checks equality with another Placeholder.
func (p Placeholder) Equals(other Placeholder) bool {
	return p.Width == other.Width &&
		p.Height == other.Height &&
		p.PlaceholderVerticalAlign == other.PlaceholderVerticalAlign
}

// String returns a string representation of the Placeholder.
func (p Placeholder) String() string {
	return fmt.Sprintf("Placeholder(width=%s, height=%s, placeholderVerticalAlign=%s)",
		p.Width, p.Height, p.PlaceholderVerticalAlign)
}
