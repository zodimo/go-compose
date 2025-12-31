package style

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// TextForegroundStyleUnspecified is the unspecified TextForegroundStyle.
var TextForegroundStyleUnspecified *TextForegroundStyle = &TextForegroundStyle{
	Color: graphics.ColorUnspecified,
	Brush: nil,
	Alpha: floatutils.Float32Unspecified,
}

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/TextForegroundStyle.kt

// TextForegroundStyle represents possible ways to draw Text e.g. color, brush.
// This interface aims to unify unspecified versions of complementary drawing styles.
//
// Guarantees:
//   - If Color() is not Unspecified, Brush() is nil.
//   - If Brush() is not nil, Color() is Unspecified.
//   - Both Color() can be Unspecified and Brush() nil, indicating nothing is specified.
//   - SolidColor brushes are stored as regular Colors.

type TextForegroundStyle struct {
	Color graphics.Color
	Brush graphics.Brush
	Alpha float32
}

func (s TextForegroundStyle) isBrushStyle() bool {
	return s.Brush != nil
}

// --- Factory functions ---

// TextForegroundStyleFromColor creates a TextForegroundStyle from a Color.
// Returns Unspecified if the color is unspecified.
func TextForegroundStyleFromColor(color graphics.Color) *TextForegroundStyle {
	if color.IsSpecified() {
		return &TextForegroundStyle{
			Color: color,
			Brush: nil,
			Alpha: floatutils.Float32Unspecified,
		}
	}
	return TextForegroundStyleUnspecified
}

// TextForegroundStyleFromBrush creates a TextForegroundStyle from a Brush and alpha.
// If brush is nil, returns Unspecified.
// If brush is a SolidColor, returns a ColorStyle with the color modulated by alpha.
// If brush is a ShaderBrush, returns a BrushStyle.
func TextForegroundStyleFromBrush(brush graphics.Brush, alpha float32) *TextForegroundStyle {
	if brush == nil {
		return TextForegroundStyleUnspecified
	}

	if sc := graphics.AsSolidColor(brush); sc != nil {
		modulatedColor := ModulateColor(sc.Value, alpha)
		return TextForegroundStyleFromColor(modulatedColor)
	}

	if sb := graphics.AsShaderBrush(brush); sb != nil {
		return &TextForegroundStyle{
			Color: graphics.ColorUnspecified,
			Brush: sb,
			Alpha: alpha,
		}
	}

	return TextForegroundStyleUnspecified
}

// --- Lerp function ---

// LerpTextForegroundStyle linearly interpolates between two TextForegroundStyles.
// If both styles are not BrushStyles, lerps the color values.
// Otherwise, lerps discretely.
func LerpTextForegroundStyle(start, stop *TextForegroundStyle, fraction float32) *TextForegroundStyle {

	start = CoalesceTextForegroundStyle(start, TextForegroundStyleUnspecified)
	stop = CoalesceTextForegroundStyle(stop, TextForegroundStyleUnspecified)

	if !start.isBrushStyle() && !stop.isBrushStyle() {
		// Both are colors - lerp colors
		c1 := start.Color
		c2 := stop.Color
		return TextForegroundStyleFromColor(graphics.LerpColor(c1, c2, fraction))
	}

	if start.isBrushStyle() && stop.isBrushStyle() {
		// Both are brushes - lerp alpha, discrete brush

		lerpedBrush := lerpDiscrete(start.Brush, stop.Brush, fraction)
		lerpedAlpha := lerp.Between32(start.Alpha, stop.Alpha, fraction)
		return TextForegroundStyleFromBrush(lerpedBrush, lerpedAlpha)
	}

	// Mixed types - discrete lerp
	return lerpDiscrete(start, stop, fraction)
}

// --- Helper functions ---

// ModulateColor modulates a color's alpha by the given alpha value.
// If alpha is NaN or >= 1.0, returns the color unchanged.
// Otherwise, multiplies the color's alpha by the given alpha.
func ModulateColor(color graphics.Color, alpha float32) graphics.Color {
	if math.IsNaN(float64(alpha)) || alpha >= 1.0 {
		return color
	}
	return color.Copy(color.Alpha()*alpha, color.Red(), color.Green(), color.Blue())
}

// isColorSpecified checks if a color is specified (not the Unspecified sentinel).
func isColorSpecified(color graphics.Color) bool {
	return color != graphics.ColorUnspecified
}

// takeOrElseFloat returns the value if it's not NaN, otherwise calls block().
func takeOrElseFloat(value float32, block func() float32) float32 {
	if math.IsNaN(float64(value)) {
		return block()
	}
	return value
}

// MergeTextForegroundStyle implements the merge logic for TextForegroundStyle.
// This prevents Color or Unspecified TextForegroundStyle from overriding an existing Brush.
func MergeTextForegroundStyle(current, other *TextForegroundStyle) *TextForegroundStyle {
	current = CoalesceTextForegroundStyle(current, TextForegroundStyleUnspecified)
	other = CoalesceTextForegroundStyle(other, TextForegroundStyleUnspecified)

	switch {
	case other.isBrushStyle() && current.isBrushStyle():
		// Both are brush styles, merge with alpha fallback

		newAlpha := floatutils.TakeOrElse(other.Alpha, current.Alpha)
		return TextForegroundStyleFromBrush(other.Brush, newAlpha)

	case other.isBrushStyle() && !current.isBrushStyle():
		// Other is brush, current is not - use other
		return other

	case !other.isBrushStyle() && current.isBrushStyle():
		// Other is not brush, current is brush - keep current (preserve brush)
		return current

	default:
		// Neither is brush - use takeOrElse
		return TakeOrElseTextForegroundStyle(other, current)
	}
}

func CoalesceTextForegroundStyle(ptr, def *TextForegroundStyle) *TextForegroundStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func TakeOrElseTextForegroundStyle(style, defaultStyle *TextForegroundStyle) *TextForegroundStyle {
	if style == nil || style == TextForegroundStyleUnspecified {
		return defaultStyle
	}
	return style
}

// Identity (2 ns)
func SameTextForegroundStyle(a, b *TextForegroundStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextForegroundStyleUnspecified
	}
	if b == nil {
		return a == TextForegroundStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextForegroundStyle(a, b *TextForegroundStyle) bool {

	a = CoalesceTextForegroundStyle(a, TextForegroundStyleUnspecified)
	b = CoalesceTextForegroundStyle(b, TextForegroundStyleUnspecified)

	return a.Color == b.Color &&
		graphics.EqualBrush(a.Brush, b.Brush) &&
		floatutils.Float32Equals(a.Alpha, b.Alpha, floatutils.Float32EqualityThreshold)
}

func EqualTextForegroundStyle(a, b *TextForegroundStyle) bool {
	if !SameTextForegroundStyle(a, b) {
		return SemanticEqualTextForegroundStyle(a, b)
	}
	return true
}
