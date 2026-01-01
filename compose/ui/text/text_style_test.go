package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

func TestTextStyleUnspecified(t *testing.T) {
	ts := TextStyleUnspecified

	if ts == nil {
		t.Fatal("TextStyleUnspecified should not be nil")
	}
	// Verify it uses unspecified span and paragraph styles
	if ts.spanStyle != SpanStyleUnspecified {
		t.Error("TextStyleUnspecified.spanStyle should be SpanStyleUnspecified")
	}
	if ts.paragraphStyle != ParagraphStyleUnspecified {
		t.Error("TextStyleUnspecified.paragraphStyle should be ParagraphStyleUnspecified")
	}
}

func TestTextStyle_Getters(t *testing.T) {
	// Create a TextStyle with specific values via options
	testColor := graphics.ColorRed
	testBg := graphics.ColorBlue
	testFontSize := unit.Sp(24)
	testFontWeight := font.FontWeightBold
	testFontStyle := font.FontStyleItalic
	testLetterSpacing := unit.Sp(2)
	testTextAlign := style.TextAlignMiddle
	testTextDirection := style.TextDirectionRtl
	testLineHeight := unit.Sp(32)
	testLineBreak := style.LineBreakSimple

	ts := TextStyleFromOptions(
		WithColor(testColor),
		WithBackground(testBg),
		WithFontSize(testFontSize),
		WithFontWeight(testFontWeight),
		WithFontStyle(testFontStyle),
		WithLetterSpacing(testLetterSpacing),
		WithTextAlign(testTextAlign),
		WithTextDirection(testTextDirection),
		WithLineHeight(testLineHeight),
		WithLineBreak(testLineBreak),
	)

	// Test span style getters
	if ts.Color() != testColor {
		t.Errorf("Color() = %v, want %v", ts.Color(), testColor)
	}
	if ts.Background() != testBg {
		t.Errorf("Background() = %v, want %v", ts.Background(), testBg)
	}
	if ts.FontSize() != testFontSize {
		t.Errorf("FontSize() = %v, want %v", ts.FontSize(), testFontSize)
	}
	if ts.FontWeight() != testFontWeight {
		t.Errorf("FontWeight() = %v, want %v", ts.FontWeight(), testFontWeight)
	}
	if ts.FontStyle() != testFontStyle {
		t.Errorf("FontStyle() = %v, want %v", ts.FontStyle(), testFontStyle)
	}
	if ts.LetterSpacing() != testLetterSpacing {
		t.Errorf("LetterSpacing() = %v, want %v", ts.LetterSpacing(), testLetterSpacing)
	}

	// Test paragraph style getters
	if ts.TextAlign() != testTextAlign {
		t.Errorf("TextAlign() = %v, want %v", ts.TextAlign(), testTextAlign)
	}
	if ts.TextDirection() != testTextDirection {
		t.Errorf("TextDirection() = %v, want %v", ts.TextDirection(), testTextDirection)
	}
	if ts.LineHeight() != testLineHeight {
		t.Errorf("LineHeight() = %v, want %v", ts.LineHeight(), testLineHeight)
	}
	if ts.LineBreak() != testLineBreak {
		t.Errorf("LineBreak() = %v, want %v", ts.LineBreak(), testLineBreak)
	}
}

func TestIsSpecifiedTextStyle(t *testing.T) {
	tests := []struct {
		name     string
		style    *TextStyle
		expected bool
	}{
		{"nil", nil, false},
		{"Unspecified", TextStyleUnspecified, false},
		{"CustomStyle", TextStyleFromOptions(WithColor(graphics.ColorRed)), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSpecifiedTextStyle(tc.style)
			if result != tc.expected {
				t.Errorf("IsSpecifiedTextStyle(%v) = %v, want %v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestTakeOrElseTextStyle(t *testing.T) {
	defaultStyle := TextStyleFromOptions(WithColor(graphics.ColorBlue))
	customStyle := TextStyleFromOptions(WithColor(graphics.ColorRed))

	tests := []struct {
		name     string
		style    *TextStyle
		def      *TextStyle
		expected *TextStyle
	}{
		{"nil_returns_default", nil, defaultStyle, defaultStyle},
		{"unspecified_returns_default", TextStyleUnspecified, defaultStyle, defaultStyle},
		{"custom_returns_custom", customStyle, defaultStyle, customStyle},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TakeOrElseTextStyle(tc.style, tc.def)
			if result != tc.expected {
				t.Errorf("TakeOrElseTextStyle() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestCopyTextStyle(t *testing.T) {
	original := TextStyleFromOptions(
		WithColor(graphics.ColorRed),
		WithFontSize(unit.Sp(20)),
	)

	// Copy with modifications
	copied := CopyTextStyle(original, WithColor(graphics.ColorBlue))

	// Original should be unchanged
	if original.Color() != graphics.ColorRed {
		t.Error("Original TextStyle should be unchanged after copy")
	}

	// Copied should have new color
	if copied.Color() != graphics.ColorBlue {
		t.Errorf("Copied TextStyle Color() = %v, want %v", copied.Color(), graphics.ColorBlue)
	}

	// Copied should inherit other properties
	if copied.FontSize() != unit.Sp(20) {
		t.Errorf("Copied TextStyle FontSize() = %v, want %v", copied.FontSize(), unit.Sp(20))
	}
}

func TestMergeTextStyle(t *testing.T) {
	t.Run("both_nil", func(t *testing.T) {
		result := MergeTextStyle(nil, nil)
		// Both nil coalesced to Unspecified, returns Unspecified
		if result != TextStyleUnspecified {
			t.Error("MergeTextStyle(nil, nil) should return TextStyleUnspecified")
		}
	})

	t.Run("a_unspecified", func(t *testing.T) {
		b := TextStyleFromOptions(WithColor(graphics.ColorRed))
		result := MergeTextStyle(TextStyleUnspecified, b)
		if result != b {
			t.Error("MergeTextStyle with a=Unspecified should return b")
		}
	})

	t.Run("b_unspecified", func(t *testing.T) {
		a := TextStyleFromOptions(WithColor(graphics.ColorRed))
		result := MergeTextStyle(a, TextStyleUnspecified)
		if result != a {
			t.Error("MergeTextStyle with b=Unspecified should return a")
		}
	})

	t.Run("both_specified", func(t *testing.T) {
		a := TextStyleFromOptions(
			WithColor(graphics.ColorRed),
			WithFontSize(unit.Sp(16)),
		)
		b := TextStyleFromOptions(
			WithColor(graphics.ColorBlue),
			WithFontWeight(font.FontWeightBold),
		)

		result := MergeTextStyle(a, b)

		// b's color should override a's
		if result.Color() != graphics.ColorBlue {
			t.Errorf("Merged Color() = %v, want %v", result.Color(), graphics.ColorBlue)
		}
		// a's fontSize should be preserved
		if result.FontSize() != unit.Sp(16) {
			t.Errorf("Merged FontSize() = %v, want %v", result.FontSize(), unit.Sp(16))
		}
		// b's fontWeight should be included
		if result.FontWeight() != font.FontWeightBold {
			t.Errorf("Merged FontWeight() = %v, want %v", result.FontWeight(), font.FontWeightBold)
		}
	})
}

func TestCoalesceTextStyle(t *testing.T) {
	defaultStyle := TextStyleFromOptions(WithColor(graphics.ColorBlue))
	customStyle := TextStyleFromOptions(WithColor(graphics.ColorRed))

	tests := []struct {
		name     string
		ptr      *TextStyle
		def      *TextStyle
		expected *TextStyle
	}{
		{"nil_returns_default", nil, defaultStyle, defaultStyle},
		{"non_nil_returns_ptr", customStyle, defaultStyle, customStyle},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CoalesceTextStyle(tc.ptr, tc.def)
			if result != tc.expected {
				t.Errorf("CoalesceTextStyle() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTextStyleFromOptions(t *testing.T) {
	t.Run("no_options", func(t *testing.T) {
		ts := TextStyleFromOptions()
		// Should have unspecified values
		if ts.Color() != graphics.ColorUnspecified {
			t.Errorf("TextStyleFromOptions() with no options Color() = %v, want Unspecified", ts.Color())
		}
	})

	t.Run("with_options", func(t *testing.T) {
		ts := TextStyleFromOptions(
			WithColor(graphics.ColorRed),
			WithFontSize(unit.Sp(18)),
			WithTextAlign(style.TextAlignMiddle),
		)

		if ts.Color() != graphics.ColorRed {
			t.Errorf("Color() = %v, want %v", ts.Color(), graphics.ColorRed)
		}
		if ts.FontSize() != unit.Sp(18) {
			t.Errorf("FontSize() = %v, want %v", ts.FontSize(), unit.Sp(18))
		}
		if ts.TextAlign() != style.TextAlignMiddle {
			t.Errorf("TextAlign() = %v, want %v", ts.TextAlign(), style.TextAlignMiddle)
		}
	})
}

func TestTextStyleResolveDefaults(t *testing.T) {
	unspecified := TextStyleUnspecified
	resolved := TextStyleResolveDefaults(unspecified, unit.LayoutDirectionLtr)

	// Should resolve to default values
	if resolved.FontSize() != DefaultFontSize {
		t.Errorf("Resolved FontSize() = %v, want %v", resolved.FontSize(), DefaultFontSize)
	}
	if resolved.FontWeight() != DefaultFontWeight {
		t.Errorf("Resolved FontWeight() = %v, want %v", resolved.FontWeight(), DefaultFontWeight)
	}
	if resolved.FontStyle() != DefaultFontStyle {
		t.Errorf("Resolved FontStyle() = %v, want %v", resolved.FontStyle(), DefaultFontStyle)
	}
	if resolved.TextAlign() != DefaultTextAlign {
		t.Errorf("Resolved TextAlign() = %v, want %v", resolved.TextAlign(), DefaultTextAlign)
	}
}

func TestStringTextStyle(t *testing.T) {
	ts := TextStyleFromOptions(WithColor(graphics.ColorRed))
	str := StringTextStyle(ts)

	// Should contain TextStyle prefix
	if len(str) == 0 {
		t.Error("StringTextStyle() should return non-empty string")
	}
	// Basic structure check
	if str[:10] != "TextStyle{" {
		t.Errorf("StringTextStyle() should start with 'TextStyle{', got: %s", str[:10])
	}
}

func TestTextStyle_ImplementsInterface(t *testing.T) {
	// Verify TextStyle implements TextStyleInterface
	var _ TextStyleInterface = (*TextStyle)(nil)
}
