package style

import "fmt"

// TextMotion defines ways to render and place glyphs to provide readability
// and smooth animations for text.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextMotion.kt

// Linearity defines the possible valid configurations for text linearity.
// Both font hinting and Linear text cannot be enabled at the same time.
type Linearity int

const (
	// LinearityLinear is equal to applying LINEAR_TEXT_FLAG and turning hinting off.
	LinearityLinear Linearity = 1
	// LinearityFontHinting is equal to removing LINEAR_TEXT_FLAG and turning hinting on.
	LinearityFontHinting Linearity = 2
	// LinearityNone is equal to removing LINEAR_TEXT_FLAG and turning hinting off.
	LinearityNone Linearity = 3
)

// String returns a string representation of the Linearity.
func (l Linearity) String() string {
	switch l {
	case LinearityLinear:
		return "Linearity.Linear"
	case LinearityFontHinting:
		return "Linearity.FontHinting"
	case LinearityNone:
		return "Linearity.None"
	default:
		return "Invalid"
	}
}

// TextMotion configuration for text rendering.
type TextMotion struct {
	// Linearity defines the text linearity mode.
	Linearity Linearity
	// SubpixelTextPositioning enables subpixel text positioning for smoother animations.
	SubpixelTextPositioning bool
}

var (
	// TextMotionStatic optimizes glyph shaping, placement, and overall rendering
	// for maximum readability. Intended for text that is not animated.
	// This is the default TextMotion.
	TextMotionStatic = TextMotion{
		Linearity:               LinearityFontHinting,
		SubpixelTextPositioning: false,
	}

	// TextMotionAnimated renders text for maximum linearity which provides smooth
	// animations for text. Trade-off is the readability of the text on some low
	// DPI devices. Use this TextMotion if you are planning to scale, translate,
	// or rotate text.
	TextMotionAnimated = TextMotion{
		Linearity:               LinearityLinear,
		SubpixelTextPositioning: true,
	}
)

// Copy creates a copy of the TextMotion with optional modifications.
func (t TextMotion) Copy(linearity *Linearity, subpixelTextPositioning *bool) TextMotion {
	result := t
	if linearity != nil {
		result.Linearity = *linearity
	}
	if subpixelTextPositioning != nil {
		result.SubpixelTextPositioning = *subpixelTextPositioning
	}
	return result
}

// Equals returns true if the other TextMotion is equal to this one.
func (t TextMotion) Equals(other TextMotion) bool {
	return t.Linearity == other.Linearity &&
		t.SubpixelTextPositioning == other.SubpixelTextPositioning
}

// String returns a string representation of the TextMotion.
func (t TextMotion) String() string {
	switch t {
	case TextMotionStatic:
		return "TextMotion.Static"
	case TextMotionAnimated:
		return "TextMotion.Animated"
	default:
		return fmt.Sprintf("TextMotion(%s, subpixel=%v)", t.Linearity, t.SubpixelTextPositioning)
	}
}
