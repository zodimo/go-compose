package style

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextForegroundStyle.kt

// TextForegroundStyle represents possible ways to draw Text e.g. color, brush.
// This interface aims to unify unspecified versions of complementary drawing styles.
//
// Guarantees:
//   - If Color() is not Unspecified, Brush() is nil.
//   - If Brush() is not nil, Color() is Unspecified.
//   - Both Color() can be Unspecified and Brush() nil, indicating nothing is specified.
//   - SolidColor brushes are stored as regular Colors.
type TextForegroundStyle interface {
	// Color returns the color for this style. Returns Unspecified if using a brush.
	Color() graphics.Color
	// Brush returns the brush for this style. Returns nil if using a color.
	Brush() graphics.Brush
	// Alpha returns the alpha value. Returns NaN if unspecified.
	Alpha() float32
	// Merge merges this style with another, returning the result.
	Merge(other TextForegroundStyle) TextForegroundStyle
	// TakeOrElse returns this style if it's not Unspecified, otherwise calls other().
	TakeOrElse(other func() TextForegroundStyle) TextForegroundStyle

	isTextForegroundStyle()
}

// --- Unspecified singleton ---

type textForegroundStyleUnspecified struct{}

func (u textForegroundStyleUnspecified) Color() graphics.Color {
	return graphics.ColorUnspecified
}

func (u textForegroundStyleUnspecified) Brush() graphics.Brush {
	return nil
}

func (u textForegroundStyleUnspecified) Alpha() float32 {
	return float32(math.NaN())
}

func (u textForegroundStyleUnspecified) Merge(other TextForegroundStyle) TextForegroundStyle {
	return defaultMerge(u, other)
}

func (u textForegroundStyleUnspecified) TakeOrElse(other func() TextForegroundStyle) TextForegroundStyle {
	return other()
}

func (u textForegroundStyleUnspecified) isTextForegroundStyle() {}

// TextForegroundStyleUnspecified is the unspecified TextForegroundStyle.
var TextForegroundStyleUnspecified TextForegroundStyle = textForegroundStyleUnspecified{}

// --- ColorStyle ---

type colorStyle struct {
	value graphics.Color
}

func (c colorStyle) Color() graphics.Color {
	return c.value
}

func (c colorStyle) Brush() graphics.Brush {
	return nil
}

func (c colorStyle) Alpha() float32 {
	// TODO: extract alpha from color when ColorDescriptor supports it
	return 1.0
}

func (c colorStyle) Merge(other TextForegroundStyle) TextForegroundStyle {
	return defaultMerge(c, other)
}

func (c colorStyle) TakeOrElse(other func() TextForegroundStyle) TextForegroundStyle {
	return c
}

func (c colorStyle) isTextForegroundStyle() {}

// --- BrushStyle ---

type brushStyle struct {
	value graphics.ShaderBrush
	alpha float32
}

func (b brushStyle) Color() graphics.Color {
	return graphics.ColorUnspecified
}

func (b brushStyle) Brush() graphics.Brush {
	return b.value
}

func (b brushStyle) Alpha() float32 {
	return b.alpha
}

func (b brushStyle) Merge(other TextForegroundStyle) TextForegroundStyle {
	return defaultMerge(b, other)
}

func (b brushStyle) TakeOrElse(other func() TextForegroundStyle) TextForegroundStyle {
	return b
}

func (b brushStyle) isTextForegroundStyle() {}

// IsBrushStyle returns true if the style is a BrushStyle.
func IsBrushStyle(s TextForegroundStyle) bool {
	_, ok := s.(brushStyle)
	return ok
}

// --- Factory functions ---

// TextForegroundStyleFromColor creates a TextForegroundStyle from a Color.
// Returns Unspecified if the color is unspecified.
func TextForegroundStyleFromColor(color graphics.Color) TextForegroundStyle {
	if isColorSpecified(color) {
		return colorStyle{value: color}
	}
	return TextForegroundStyleUnspecified
}

// TextForegroundStyleFromBrush creates a TextForegroundStyle from a Brush and alpha.
// If brush is nil, returns Unspecified.
// If brush is a SolidColor, returns a ColorStyle with the color modulated by alpha.
// If brush is a ShaderBrush, returns a BrushStyle.
func TextForegroundStyleFromBrush(brush graphics.Brush, alpha float32) TextForegroundStyle {
	if brush == nil {
		return TextForegroundStyleUnspecified
	}

	if sc := graphics.AsSolidColor(brush); sc != nil {
		modulatedColor := ModulateColor(sc.Value, alpha)
		return TextForegroundStyleFromColor(modulatedColor)
	}

	if sb := graphics.AsShaderBrush(brush); sb != nil {
		return brushStyle{value: sb, alpha: alpha}
	}

	return TextForegroundStyleUnspecified
}

// --- Merge logic ---

// defaultMerge implements the merge logic for TextForegroundStyle.
// This prevents Color or Unspecified TextForegroundStyle from overriding an existing Brush.
func defaultMerge(current, other TextForegroundStyle) TextForegroundStyle {
	switch {
	case IsBrushStyle(other) && IsBrushStyle(current):
		// Both are brush styles, merge with alpha fallback
		otherBrush := other.(brushStyle)
		currentBrush := current.(brushStyle)
		newAlpha := takeOrElseFloat(otherBrush.alpha, func() float32 { return currentBrush.alpha })
		return brushStyle{value: otherBrush.value, alpha: newAlpha}

	case IsBrushStyle(other) && !IsBrushStyle(current):
		// Other is brush, current is not - use other
		return other

	case !IsBrushStyle(other) && IsBrushStyle(current):
		// Other is not brush, current is brush - keep current (preserve brush)
		return current

	default:
		// Neither is brush - use takeOrElse
		return other.TakeOrElse(func() TextForegroundStyle { return current })
	}
}

// --- Lerp function ---

// LerpTextForegroundStyle linearly interpolates between two TextForegroundStyles.
// If both styles are not BrushStyles, lerps the color values.
// Otherwise, lerps discretely.
func LerpTextForegroundStyle(start, stop TextForegroundStyle, fraction float32) TextForegroundStyle {
	if !IsBrushStyle(start) && !IsBrushStyle(stop) {
		// Both are colors - lerp colors
		// TODO: implement color lerp when available
		// For now, use discrete lerp
		return lerpDiscrete(start, stop, fraction)
	}

	if IsBrushStyle(start) && IsBrushStyle(stop) {
		// Both are brushes - lerp alpha, discrete brush
		startBrush := start.(brushStyle)
		stopBrush := stop.(brushStyle)
		lerpedBrush := lerpDiscrete(startBrush.value, stopBrush.value, fraction)
		lerpedAlpha := lerpFloat(startBrush.alpha, stopBrush.alpha, fraction)
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
	// Use SetOpacity to apply the alpha modulation
	// The color's alpha is multiplied by the given alpha
	return color.SetOpacity(alphaToOpacity(alpha))
}

// isColorSpecified checks if a color is specified (not the Unspecified sentinel).
func isColorSpecified(color graphics.Color) bool {
	return color.TakeOrElse(func() graphics.Color { return nil }) != nil
}

// takeOrElseFloat returns the value if it's not NaN, otherwise calls block().
func takeOrElseFloat(value float32, block func() float32) float32 {
	if math.IsNaN(float64(value)) {
		return block()
	}
	return value
}

// lerpFloat linearly interpolates between two floats.
func lerpFloat(a, b, fraction float32) float32 {
	return a + (b-a)*fraction
}

// alphaToOpacity converts a 0-1 alpha to an OpacityLevel.
// This is a simple mapping for the SetOpacity interface.
func alphaToOpacity(alpha float32) graphics.OpacityLevel {
	// Clamp alpha to 0-1 range
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	// Map to opacity level (0-100 scale typically)
	return graphics.OpacityLevel(alpha * 100)
}
