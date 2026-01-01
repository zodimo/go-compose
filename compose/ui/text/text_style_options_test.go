package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// === Span Style Options ===

func TestWithColor(t *testing.T) {
	ts := TextStyleFromOptions(WithColor(graphics.ColorRed))

	if ts.Color() != graphics.ColorRed {
		t.Errorf("WithColor() Color() = %v, want %v", ts.Color(), graphics.ColorRed)
	}
}

func TestWithFontSize(t *testing.T) {
	testSize := unit.Sp(24)
	ts := TextStyleFromOptions(WithFontSize(testSize))

	if ts.FontSize() != testSize {
		t.Errorf("WithFontSize() FontSize() = %v, want %v", ts.FontSize(), testSize)
	}
}

func TestWithFontWeight(t *testing.T) {
	ts := TextStyleFromOptions(WithFontWeight(font.FontWeightBold))

	if ts.FontWeight() != font.FontWeightBold {
		t.Errorf("WithFontWeight() FontWeight() = %v, want %v", ts.FontWeight(), font.FontWeightBold)
	}
}

func TestWithFontStyle(t *testing.T) {
	ts := TextStyleFromOptions(WithFontStyle(font.FontStyleItalic))

	if ts.FontStyle() != font.FontStyleItalic {
		t.Errorf("WithFontStyle() FontStyle() = %v, want %v", ts.FontStyle(), font.FontStyleItalic)
	}
}

func TestWithFontFamily(t *testing.T) {
	customFamily := font.FontFamilyDefault
	ts := TextStyleFromOptions(WithFontFamily(customFamily))

	// Direct comparison since we're using a known constant
	if ts.FontFamily() != customFamily {
		t.Errorf("WithFontFamily() FontFamily() = %v, want %v", ts.FontFamily(), customFamily)
	}
}

func TestWithLetterSpacing(t *testing.T) {
	testSpacing := unit.Sp(2.5)
	ts := TextStyleFromOptions(WithLetterSpacing(testSpacing))

	if ts.LetterSpacing() != testSpacing {
		t.Errorf("WithLetterSpacing() LetterSpacing() = %v, want %v", ts.LetterSpacing(), testSpacing)
	}
}

func TestWithBackground(t *testing.T) {
	ts := TextStyleFromOptions(WithBackground(graphics.ColorBlue))

	if ts.Background() != graphics.ColorBlue {
		t.Errorf("WithBackground() Background() = %v, want %v", ts.Background(), graphics.ColorBlue)
	}
}

func TestWithTextDecoration(t *testing.T) {
	decoration := style.TextDecorationUnderline
	ts := TextStyleFromOptions(WithTextDecoration(decoration))

	if ts.TextDecoration() != decoration {
		t.Errorf("WithTextDecoration() TextDecoration() = %v, want %v", ts.TextDecoration(), decoration)
	}
}

func TestWithShadow(t *testing.T) {
	// Create a test shadow
	shadow := graphics.NewShadow(
		graphics.ColorBlack,
		geometry.NewOffset(5, 5),
		10,
	)
	ts := TextStyleFromOptions(WithShadow(shadow))

	if ts.spanStyle.Shadow() != shadow {
		t.Errorf("WithShadow() Shadow() = %v, want %v", ts.spanStyle.Shadow(), shadow)
	}
}

// === Paragraph Style Options ===

func TestWithTextAlign(t *testing.T) {
	testCases := []style.TextAlign{
		style.TextAlignStart,
		style.TextAlignEnd,
		style.TextAlignMiddle,
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(t *testing.T) {
			ts := TextStyleFromOptions(WithTextAlign(tc))
			if ts.TextAlign() != tc {
				t.Errorf("WithTextAlign(%v) TextAlign() = %v, want %v", tc, ts.TextAlign(), tc)
			}
		})
	}
}

func TestWithTextDirection(t *testing.T) {
	testCases := []struct {
		name      string
		direction style.TextDirection
	}{
		{"Ltr", style.TextDirectionLtr},
		{"Rtl", style.TextDirectionRtl},
		{"Content", style.TextDirectionContent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := TextStyleFromOptions(WithTextDirection(tc.direction))
			if ts.TextDirection() != tc.direction {
				t.Errorf("WithTextDirection(%v) TextDirection() = %v, want %v", tc.direction, ts.TextDirection(), tc.direction)
			}
		})
	}
}

func TestWithLineHeight(t *testing.T) {
	testHeight := unit.Sp(28)
	ts := TextStyleFromOptions(WithLineHeight(testHeight))

	if ts.LineHeight() != testHeight {
		t.Errorf("WithLineHeight() LineHeight() = %v, want %v", ts.LineHeight(), testHeight)
	}
}

func TestWithLineBreak(t *testing.T) {
	testCases := []style.LineBreak{
		style.LineBreakSimple,
		style.LineBreakHeading,
		style.LineBreakParagraph,
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(t *testing.T) {
			ts := TextStyleFromOptions(WithLineBreak(tc))
			if ts.LineBreak() != tc {
				t.Errorf("WithLineBreak(%v) LineBreak() = %v, want %v", tc, ts.LineBreak(), tc)
			}
		})
	}
}

// === Chaining Options ===

func TestTextStyleOptions_Chaining(t *testing.T) {
	ts := TextStyleFromOptions(
		// Span style options
		WithColor(graphics.ColorRed),
		WithFontSize(unit.Sp(20)),
		WithFontWeight(font.FontWeightBold),
		WithFontStyle(font.FontStyleItalic),
		WithLetterSpacing(unit.Sp(1)),
		WithBackground(graphics.ColorYellow),
		// Paragraph style options
		WithTextAlign(style.TextAlignMiddle),
		WithTextDirection(style.TextDirectionRtl),
		WithLineHeight(unit.Sp(30)),
		WithLineBreak(style.LineBreakSimple),
	)

	// Verify all span style values
	if ts.Color() != graphics.ColorRed {
		t.Errorf("Chained Color() = %v, want %v", ts.Color(), graphics.ColorRed)
	}
	if ts.FontSize() != unit.Sp(20) {
		t.Errorf("Chained FontSize() = %v, want %v", ts.FontSize(), unit.Sp(20))
	}
	if ts.FontWeight() != font.FontWeightBold {
		t.Errorf("Chained FontWeight() = %v, want %v", ts.FontWeight(), font.FontWeightBold)
	}
	if ts.FontStyle() != font.FontStyleItalic {
		t.Errorf("Chained FontStyle() = %v, want %v", ts.FontStyle(), font.FontStyleItalic)
	}
	if ts.LetterSpacing() != unit.Sp(1) {
		t.Errorf("Chained LetterSpacing() = %v, want %v", ts.LetterSpacing(), unit.Sp(1))
	}
	if ts.Background() != graphics.ColorYellow {
		t.Errorf("Chained Background() = %v, want %v", ts.Background(), graphics.ColorYellow)
	}

	// Verify all paragraph style values
	if ts.TextAlign() != style.TextAlignMiddle {
		t.Errorf("Chained TextAlign() = %v, want %v", ts.TextAlign(), style.TextAlignMiddle)
	}
	if ts.TextDirection() != style.TextDirectionRtl {
		t.Errorf("Chained TextDirection() = %v, want %v", ts.TextDirection(), style.TextDirectionRtl)
	}
	if ts.LineHeight() != unit.Sp(30) {
		t.Errorf("Chained LineHeight() = %v, want %v", ts.LineHeight(), unit.Sp(30))
	}
	if ts.LineBreak() != style.LineBreakSimple {
		t.Errorf("Chained LineBreak() = %v, want %v", ts.LineBreak(), style.LineBreakSimple)
	}
}

// === Options Applied to Copy ===

func TestTextStyleOptions_AppliedToCopy(t *testing.T) {
	original := TextStyleFromOptions(
		WithColor(graphics.ColorRed),
		WithFontSize(unit.Sp(16)),
	)

	// Apply new options via CopyTextStyle
	copied := CopyTextStyle(original,
		WithColor(graphics.ColorBlue),
		WithFontWeight(font.FontWeightBold),
	)

	// Original should be unchanged
	if original.Color() != graphics.ColorRed {
		t.Error("Original should be unchanged after copy")
	}
	if original.FontWeight() != font.FontWeightUnspecified {
		t.Error("Original FontWeight should be unchanged")
	}

	// Copy should have new values
	if copied.Color() != graphics.ColorBlue {
		t.Errorf("Copied Color() = %v, want %v", copied.Color(), graphics.ColorBlue)
	}
	if copied.FontWeight() != font.FontWeightBold {
		t.Errorf("Copied FontWeight() = %v, want %v", copied.FontWeight(), font.FontWeightBold)
	}
	// Copy should preserve unmodified values from original
	if copied.FontSize() != unit.Sp(16) {
		t.Errorf("Copied FontSize() = %v, want %v", copied.FontSize(), unit.Sp(16))
	}
}

// === Options with Nil/Unspecified Base ===

func TestTextStyleOptions_FromNil(t *testing.T) {
	// Options should handle nil base gracefully via CoalesceTextStyle
	ts := TextStyleFromOptions(WithColor(graphics.ColorGreen))

	if ts.Color() != graphics.ColorGreen {
		t.Errorf("TextStyleFromOptions with nil base Color() = %v, want %v", ts.Color(), graphics.ColorGreen)
	}
}
