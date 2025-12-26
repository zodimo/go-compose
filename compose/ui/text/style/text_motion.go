package style

import (
	"fmt"

	"github.com/zodimo/go-ternary"
)

// TextMotion defines ways to render and place glyphs to provide readability
// and smooth animations for text.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextMotion.kt

// Linearity defines the possible valid configurations for text linearity.
// Both font hinting and Linear text cannot be enabled at the same time.
type Linearity int

const (
	LinearityUnspecified Linearity = 0
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
	case LinearityUnspecified:
		return "Linearity.Unspecified"
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

func (l Linearity) TakeOrElse(defaultLinearity Linearity) Linearity {
	if l != LinearityUnspecified {
		return l
	}
	return defaultLinearity
}

func (l Linearity) IsSpecified() bool {
	return l != LinearityUnspecified
}

type SubpixelTextPositioning int

const (
	SubpixelTextPositioningUnspecified SubpixelTextPositioning = 0
	SubpixelTextPositioningTrue        SubpixelTextPositioning = 1
	SubpixelTextPositioningFalse       SubpixelTextPositioning = 2
)

func (s SubpixelTextPositioning) String() string {
	switch s {
	case SubpixelTextPositioningUnspecified:
		return "Unspecified"
	case SubpixelTextPositioningTrue:
		return "True"
	case SubpixelTextPositioningFalse:
		return "False"
	default:
		return "Invalid"
	}
}

func (s SubpixelTextPositioning) TakeOrElse(defaultSubpixelTextPositioning SubpixelTextPositioning) SubpixelTextPositioning {
	if s != SubpixelTextPositioningUnspecified {
		return s
	}
	return defaultSubpixelTextPositioning
}

func (s SubpixelTextPositioning) IsSpecified() bool {
	return s != SubpixelTextPositioningUnspecified
}

// TextMotion configuration for text rendering.
type TextMotion struct {
	// Linearity defines the text linearity mode.
	Linearity Linearity
	// SubpixelTextPositioning enables subpixel text positioning for smoother animations.
	SubpixelTextPositioning SubpixelTextPositioning
}

var (
	TextMotionUnspecified = &TextMotion{
		Linearity:               LinearityUnspecified,
		SubpixelTextPositioning: SubpixelTextPositioningUnspecified,
	}
	// TextMotionStatic optimizes glyph shaping, placement, and overall rendering
	// for maximum readability. Intended for text that is not animated.
	// This is the default TextMotion.
	TextMotionStatic = &TextMotion{
		Linearity:               LinearityFontHinting,
		SubpixelTextPositioning: SubpixelTextPositioningFalse,
	}

	// TextMotionAnimated renders text for maximum linearity which provides smooth
	// animations for text. Trade-off is the readability of the text on some low
	// DPI devices. Use this TextMotion if you are planning to scale, translate,
	// or rotate text.
	TextMotionAnimated = &TextMotion{
		Linearity:               LinearityLinear,
		SubpixelTextPositioning: SubpixelTextPositioningTrue,
	}
)

type TextMotionOptions struct {
	Linearity               Linearity
	SubpixelTextPositioning SubpixelTextPositioning
}

type TextMotionOption func(*TextMotionOptions)

func WithLinearity(linearity Linearity) TextMotionOption {
	return func(options *TextMotionOptions) {
		options.Linearity = linearity
	}
}

func WithSubpixelTextPositioning(subpixelTextPositioning bool) TextMotionOption {
	return func(options *TextMotionOptions) {
		options.SubpixelTextPositioning = ternary.Ternary(subpixelTextPositioning, SubpixelTextPositioningTrue, SubpixelTextPositioningFalse)
	}
}

// Copy creates a copy of the TextMotion with optional modifications.
func (t TextMotion) Copy(options ...TextMotionOption) *TextMotion {
	opt := TextMotionOptions{
		Linearity:               LinearityUnspecified,
		SubpixelTextPositioning: SubpixelTextPositioningUnspecified,
	}
	for _, option := range options {
		option(&opt)
	}

	return &TextMotion{
		Linearity:               opt.Linearity.TakeOrElse(t.Linearity),
		SubpixelTextPositioning: opt.SubpixelTextPositioning.TakeOrElse(t.SubpixelTextPositioning),
	}
}

func StringTextMotion(s *TextMotion) string {
	s = CoalesceTextMotion(s, TextMotionUnspecified)

	switch s {
	case TextMotionUnspecified:
		return "TextMotion.Unspecified"
	case TextMotionStatic:
		return "TextMotion.Static"
	case TextMotionAnimated:
		return "TextMotion.Animated"
	default:
		return fmt.Sprintf("TextMotion(%s, subpixel=%s)", s.Linearity, s.SubpixelTextPositioning)
	}

}

func IsSpecifiedTextMotion(s *TextMotion) bool {
	return s != nil && s != TextMotionUnspecified
}
func TakeOrElseTextMotion(s, def *TextMotion) *TextMotion {
	if !IsSpecifiedTextMotion(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameTextMotion(a, b *TextMotion) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextMotionUnspecified
	}
	if b == nil {
		return a == TextMotionUnspecified
	}
	return a == b
}

func EqualTextMotion(a, b *TextMotion) bool {
	if !SameTextMotion(a, b) {
		return SemanticEqualTextMotion(a, b)
	}
	return true
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextMotion(a, b *TextMotion) bool {

	a = CoalesceTextMotion(a, TextMotionUnspecified)
	b = CoalesceTextMotion(b, TextMotionUnspecified)

	return a.Linearity == b.Linearity &&
		a.SubpixelTextPositioning == b.SubpixelTextPositioning
}

func CoalesceTextMotion(ptr, def *TextMotion) *TextMotion {
	if ptr == nil {
		return def
	}
	return ptr
}
