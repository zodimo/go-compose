package text

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/intl"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var _ Annotation = (*SpanStyle)(nil)

// SpanStyle configuration for a text span.
type SpanStyle struct {
	Color                  graphics.Color
	Brush                  graphics.Brush
	Alpha                  float32
	FontSize               unit.TextUnit
	FontWeight             font.FontWeight
	FontStyle              font.FontStyle
	FontSynthesis          font.FontSynthesis
	FontFamily             font.FontFamily
	FontFeatureSettings    string
	LetterSpacing          unit.TextUnit
	BaselineShift          style.BaselineShift
	TextGeometricTransform style.TextGeometricTransform
	LocaleList             intl.LocaleList
	Background             graphics.Color
	TextDecoration         style.TextDecoration
	Shadow                 graphics.Shadow
	PlatformStyle          *PlatformSpanStyle
	DrawStyle              graphics.DrawStyle
}

// NewSpanStyle creates a new SpanStyle with default values.
func NewSpanStyle() SpanStyle {
	return SpanStyle{
		Color:    graphics.ColorUnspecified,
		Brush:    graphics.BrushUnspecified,
		Alpha:    float32(math.NaN()),
		FontSize: unit.TextUnitUnspecified,
		// FontWeight:             font.FontWeightUnspecified,
		// FontStyle:              font.FontStyleUnspecified,
		// FontSynthesis:          font.FontSynthesisUnspecified,
		// FontFamily:             font.FontFamilyUnspecified,
		// FontFeatureSettings: "",
		LetterSpacing: unit.TextUnitUnspecified,
		BaselineShift: style.BaselineShiftUnspecified,
		// TextGeometricTransform: style.TextGeometricTransformUnspecified,
		// LocaleList:             intl.LocaleListUnspecified,
		// Background:             graphics.ColorUnspecified,
		// TextDecoration:         style.TextDecorationUnspecified,
		// Shadow:                 graphics.ShadowUnspecified,
		PlatformStyle: nil,
		// DrawStyle:              graphics.DrawStyleUnspecified,
	}
}

func (s SpanStyle) isAnnotation() {}

func (s SpanStyle) Merge(other SpanStyle) SpanStyle {
	return SpanStyle{
		Color:                  s.Color.TakeOrElse(other.Color),
		Brush:                  takeBrushOrElse(s.Brush, other.Brush),
		Alpha:                  ifNaN(s.Alpha, other.Alpha),
		FontSize:               s.FontSize.TakeOrElse(other.FontSize),
		FontWeight:             takeFontWeightOrElse(s.FontWeight, other.FontWeight),
		FontStyle:              takeFontStyleOrElse(s.FontStyle, other.FontStyle),
		FontSynthesis:          takeFontSynthesisOrElse(s.FontSynthesis, other.FontSynthesis),
		FontFamily:             takeFontFamilyOrElse(s.FontFamily, other.FontFamily),
		FontFeatureSettings:    takeStringOrElse(s.FontFeatureSettings, other.FontFeatureSettings),
		LetterSpacing:          s.LetterSpacing.TakeOrElse(other.LetterSpacing),
		BaselineShift:          takeBaselineShiftOrElse(s.BaselineShift, other.BaselineShift),
		TextGeometricTransform: takeTextGeometricTransformOrElse(s.TextGeometricTransform, other.TextGeometricTransform),
		LocaleList:             takeLocaleListOrElse(s.LocaleList, other.LocaleList),
		Background:             s.Background.TakeOrElse(other.Background),
		TextDecoration:         takeTextDecorationOrElse(s.TextDecoration, other.TextDecoration),
		Shadow:                 takeShadowOrElse(s.Shadow, other.Shadow),
		PlatformStyle:          s.PlatformStyle.Merge(other.PlatformStyle),
		DrawStyle:              takeDrawStyleOrElse(s.DrawStyle, other.DrawStyle),
	}
}

// Equals returns true if the other SpanStyle is equal to this one.
func (s SpanStyle) Equals(other SpanStyle) bool {

	if s.Color != other.Color {
		return false
	}
	// XOR
	if !(s.Brush == graphics.BrushUnspecified && other.Brush == graphics.BrushUnspecified) {

	}
	if !s.Brush.Equal(other.Brush) {
		return false
	}
	if !(floatutils.IsNaN(s.Alpha) && floatutils.IsNaN(other.Alpha) || s.Alpha == other.Alpha) {
		return false
	}
	return s.FontSize.Equals(other.FontSize) &&
		s.FontWeight == other.FontWeight &&
		s.FontStyle == other.FontStyle &&
		s.FontSynthesis == other.FontSynthesis &&
		s.FontFamily == other.FontFamily &&
		s.FontFeatureSettings == other.FontFeatureSettings &&
		s.LetterSpacing.Equals(other.LetterSpacing) &&
		baselineShiftEquals(s.BaselineShift, other.BaselineShift) &&
		s.TextGeometricTransform == other.TextGeometricTransform &&
		s.LocaleList.Equals(other.LocaleList) &&
		s.Background == other.Background &&
		s.TextDecoration == other.TextDecoration &&
		s.Shadow == other.Shadow &&
		s.PlatformStyle.Equals(other.PlatformStyle) &&
		s.DrawStyle == other.DrawStyle
}

// IsDefault returns true if this SpanStyle has all default/unspecified values.
func (s SpanStyle) IsDefault() bool {
	defaultStyle := NewSpanStyle()
	return s.Color == defaultStyle.Color &&
		s.Brush == defaultStyle.Brush &&
		(floatutils.IsNaN(s.Alpha) && floatutils.IsNaN(defaultStyle.Alpha)) &&
		s.FontSize.Equals(defaultStyle.FontSize) &&
		s.FontWeight == defaultStyle.FontWeight &&
		s.FontStyle == defaultStyle.FontStyle &&
		s.FontSynthesis == defaultStyle.FontSynthesis &&
		s.FontFamily == defaultStyle.FontFamily &&
		s.FontFeatureSettings == defaultStyle.FontFeatureSettings &&
		s.LetterSpacing.Equals(defaultStyle.LetterSpacing) &&
		baselineShiftEquals(s.BaselineShift, defaultStyle.BaselineShift) &&
		s.TextGeometricTransform == defaultStyle.TextGeometricTransform &&
		s.LocaleList.Equals(defaultStyle.LocaleList) &&
		s.Background == defaultStyle.Background &&
		s.TextDecoration == defaultStyle.TextDecoration &&
		s.Shadow == defaultStyle.Shadow &&
		s.PlatformStyle == defaultStyle.PlatformStyle &&
		s.DrawStyle == defaultStyle.DrawStyle
}

// Helpers

func takeBrushOrElse(a, b graphics.Brush) graphics.Brush {
	if a != nil {
		return a
	}
	return b
}

func ifNaN(a, b float32) float32 {
	if floatutils.IsNaN(a) {
		return b
	}
	return a
}

func baselineShiftEquals(a, b style.BaselineShift) bool {
	return a == b || (!a.IsSpecified() && !b.IsSpecified())
}

func takeFontWeightOrElse(a, b font.FontWeight) font.FontWeight {
	if a.Weight() == 0 { // Assuming 0 is unset/invalid for FontWeight which starts at 100
		return b
	}
	return a
}

func takeFontStyleOrElse(a, b font.FontStyle) font.FontStyle {
	// This is tricky if enum 0 is a valid value (Normal).
	// Usually we need a comprehensive way to represent 'Unspecified' for enums.
	// For now assuming Go default 0 is what we have if not set explicitly, but FontStyleNormal is 0.
	// In Kotlin they use nullables. Here we might need a separate 'Unspecified' constant or pointer.
	// Ideally we should update font.go to have an Unspecified constant.
	// Let's assume for now 0 is Normal and we can't distinguish unset.
	// Or check if we can change FontStyle definition to include unspecified.
	return a
}

func takeFontSynthesisOrElse(a, b font.FontSynthesis) font.FontSynthesis {
	// Similar issue with zero values.
	return a
}

func takeFontFamilyOrElse(a, b font.FontFamily) font.FontFamily {
	if a == nil {
		return b
	}
	return a
}

func takeStringOrElse(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func takeBaselineShiftOrElse(a, b style.BaselineShift) style.BaselineShift {
	return a // Simplify for compilation, fix logic later
}

func takeTextGeometricTransformOrElse(a, b style.TextGeometricTransform) style.TextGeometricTransform {
	return a
}

func takeLocaleListOrElse(a, b intl.LocaleList) intl.LocaleList {
	return a
}

func takeTextDecorationOrElse(a, b style.TextDecoration) style.TextDecoration {
	return a
}

func takeShadowOrElse(a, b graphics.Shadow) graphics.Shadow {
	return a
}

func takeDrawStyleOrElse(a, b graphics.DrawStyle) graphics.DrawStyle {
	if a == nil {
		return b
	}
	return a
}

// LerpSpanStyle interpolates between two SpanStyles.
func LerpSpanStyle(start, stop SpanStyle, fraction float32) SpanStyle {
	return SpanStyle{
		Color:                  graphics.Lerp(start.Color, stop.Color, fraction),
		Brush:                  LerpDiscrete(start.Brush, stop.Brush, fraction),
		Alpha:                  lerp.Between32(start.Alpha, stop.Alpha, fraction),
		FontSize:               unit.LerpTextUnit(start.FontSize, stop.FontSize, fraction),
		FontWeight:             font.LerpFontWeight(start.FontWeight, stop.FontWeight, fraction),
		FontStyle:              LerpDiscrete(start.FontStyle, stop.FontStyle, fraction),
		FontSynthesis:          LerpDiscrete(start.FontSynthesis, stop.FontSynthesis, fraction),
		FontFamily:             LerpDiscrete(start.FontFamily, stop.FontFamily, fraction),
		FontFeatureSettings:    LerpDiscrete(start.FontFeatureSettings, stop.FontFeatureSettings, fraction),
		LetterSpacing:          unit.LerpTextUnit(start.LetterSpacing, stop.LetterSpacing, fraction),
		BaselineShift:          style.LerpBaselineShift(start.BaselineShift, stop.BaselineShift, fraction),
		TextGeometricTransform: style.LerpGeometricTransform(start.TextGeometricTransform, stop.TextGeometricTransform, fraction),
		LocaleList:             LerpDiscrete(start.LocaleList, stop.LocaleList, fraction),
		Background:             graphics.Lerp(start.Background, stop.Background, fraction),
		TextDecoration:         LerpDiscrete(start.TextDecoration, stop.TextDecoration, fraction),
		Shadow:                 graphics.LerpShadow(start.Shadow, stop.Shadow, fraction),
		PlatformStyle:          lerpPlatformSpanStyle(start.PlatformStyle, stop.PlatformStyle, fraction),
		DrawStyle:              LerpDiscrete(start.DrawStyle, stop.DrawStyle, fraction),
	}
}
