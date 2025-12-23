package style

import (
	"math"
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// --- Tests for TextForegroundStyleUnspecified ---

func TestTextForegroundStyleUnspecified(t *testing.T) {
	u := TextForegroundStyleUnspecified

	// Color should be unspecified
	if u.Color() != graphics.ColorUnspecified {
		t.Errorf("Expected unspecified color")
	}

	// Brush should be nil
	if u.Brush() != nil {
		t.Errorf("Expected nil brush, got %v", u.Brush())
	}

	// Alpha should be NaN
	if !math.IsNaN(float64(u.Alpha())) {
		t.Errorf("Expected NaN alpha, got %v", u.Alpha())
	}
}

func TestTextForegroundStyleUnspecified_TakeOrElse(t *testing.T) {
	u := TextForegroundStyleUnspecified

	// A specified color
	specifiedColor := graphics.ColorRed
	fallback := TextForegroundStyleFromColor(specifiedColor)

	result := u.TakeOrElse(func() TextForegroundStyle { return fallback })

	// Should return the fallback since u is unspecified
	if !isColorSpecified(result.Color()) {
		t.Errorf("Expected specified color from fallback")
	}
}

// --- Tests for ColorStyle ---

func TestTextForegroundStyleFromColor_Specified(t *testing.T) {
	specifiedColor := graphics.ColorRed
	style := TextForegroundStyleFromColor(specifiedColor)

	// Color should be the specified color
	if !isColorSpecified(style.Color()) {
		t.Errorf("Expected specified color")
	}

	// Brush should be nil
	if style.Brush() != nil {
		t.Errorf("Expected nil brush")
	}

	// Should not be a BrushStyle
	if IsBrushStyle(style) {
		t.Errorf("Expected style to not be a BrushStyle")
	}
}

func TestTextForegroundStyleFromColor_Unspecified(t *testing.T) {
	style := TextForegroundStyleFromColor(graphics.ColorUnspecified)

	// Should return Unspecified
	if style != TextForegroundStyleUnspecified {
		t.Errorf("Expected TextForegroundStyleUnspecified")
	}
}

func TestColorStyle_TakeOrElse(t *testing.T) {
	specifiedColor := graphics.ColorRed
	style := TextForegroundStyleFromColor(specifiedColor)

	// TakeOrElse should return the style itself
	result := style.TakeOrElse(func() TextForegroundStyle { return TextForegroundStyleUnspecified })

	if !isColorSpecified(result.Color()) {
		t.Errorf("Expected specified color to be returned")
	}
}

// --- Tests for BrushStyle ---

func TestTextForegroundStyleFromBrush_Nil(t *testing.T) {
	style := TextForegroundStyleFromBrush(nil, 1.0)

	// Should return Unspecified
	if style != TextForegroundStyleUnspecified {
		t.Errorf("Expected TextForegroundStyleUnspecified for nil brush")
	}
}

func TestTextForegroundStyleFromBrush_SolidColor(t *testing.T) {
	color := graphics.ColorRed
	brush := graphics.NewSolidColor(color)

	style := TextForegroundStyleFromBrush(brush, 1.0)

	// Should NOT be a BrushStyle (SolidColor converts to ColorStyle)
	if IsBrushStyle(style) {
		t.Errorf("SolidColor brush should result in ColorStyle, not BrushStyle")
	}

	// Color should be specified
	if !isColorSpecified(style.Color()) {
		t.Errorf("Expected specified color")
	}
}

func TestTextForegroundStyleFromBrush_ShaderBrush(t *testing.T) {
	brush := graphics.NewShaderBrushForTest()

	style := TextForegroundStyleFromBrush(brush, 0.5)

	// Should be a BrushStyle
	if !IsBrushStyle(style) {
		t.Errorf("ShaderBrush should result in BrushStyle")
	}

	// Color should be unspecified
	if style.Color() != graphics.ColorUnspecified {
		t.Errorf("Expected unspecified color for BrushStyle")
	}

	// Brush should not be nil
	if style.Brush() == nil {
		t.Errorf("Expected non-nil brush")
	}

	// Alpha should be 0.5
	if style.Alpha() != 0.5 {
		t.Errorf("Expected alpha 0.5, got %v", style.Alpha())
	}
}

func TestBrushStyle_TakeOrElse(t *testing.T) {
	brush := graphics.NewShaderBrushForTest()
	style := TextForegroundStyleFromBrush(brush, 0.8)

	// TakeOrElse should return the style itself
	result := style.TakeOrElse(func() TextForegroundStyle { return TextForegroundStyleUnspecified })

	if !IsBrushStyle(result) {
		t.Errorf("Expected BrushStyle to be returned")
	}
}

// --- Tests for Merge ---

func TestMerge_BothBrushStyles(t *testing.T) {
	brush1 := graphics.NewShaderBrushForTest()
	brush2 := graphics.NewShaderBrushForTest()

	style1 := TextForegroundStyleFromBrush(brush1, 0.5)
	style2 := TextForegroundStyleFromBrush(brush2, 0.8)

	merged := style1.Merge(style2)

	// Should be a BrushStyle with style2's brush and alpha
	if !IsBrushStyle(merged) {
		t.Errorf("Merged should be BrushStyle")
	}
	if merged.Alpha() != 0.8 {
		t.Errorf("Expected alpha 0.8, got %v", merged.Alpha())
	}
}

func TestMerge_BrushStyleAndColorStyle(t *testing.T) {
	brush := graphics.NewShaderBrushForTest()
	brushStyle := TextForegroundStyleFromBrush(brush, 0.5)

	color := graphics.ColorRed
	colorStyle := TextForegroundStyleFromColor(color)

	// Merging color over brush should preserve brush
	merged := brushStyle.Merge(colorStyle)
	if !IsBrushStyle(merged) {
		t.Errorf("Brush should be preserved when merging color over it")
	}

	// Merging brush over color should use brush
	merged2 := colorStyle.Merge(brushStyle)
	if !IsBrushStyle(merged2) {
		t.Errorf("Brush should override color")
	}
}

func TestMerge_ColorStyleAndUnspecified(t *testing.T) {
	color := graphics.ColorRed
	colorStyle := TextForegroundStyleFromColor(color)

	merged := colorStyle.Merge(TextForegroundStyleUnspecified)

	// Unspecified should fallback to colorStyle
	if !isColorSpecified(merged.Color()) {
		t.Errorf("Expected specified color after merge with unspecified")
	}
}

// --- Tests for Lerp ---

func TestLerpTextForegroundStyle_BothUnspecified(t *testing.T) {
	result := LerpTextForegroundStyle(TextForegroundStyleUnspecified, TextForegroundStyleUnspecified, 0.5)

	if result != TextForegroundStyleUnspecified {
		t.Errorf("Lerping two unspecified should return unspecified")
	}
}

func TestLerpTextForegroundStyle_BothBrushStyles(t *testing.T) {
	brush1 := graphics.NewShaderBrushForTest()
	brush2 := graphics.NewShaderBrushForTest()

	style1 := TextForegroundStyleFromBrush(brush1, 0.2)
	style2 := TextForegroundStyleFromBrush(brush2, 0.8)

	result := LerpTextForegroundStyle(style1, style2, 0.5)

	// Result should be a BrushStyle with lerped alpha
	if !IsBrushStyle(result) {
		t.Errorf("Lerped result should be BrushStyle")
	}

	expectedAlpha := lerpFloat(0.2, 0.8, 0.5) // 0.5
	if result.Alpha() != expectedAlpha {
		t.Errorf("Expected alpha %v, got %v", expectedAlpha, result.Alpha())
	}
}

func TestLerpTextForegroundStyle_Discrete(t *testing.T) {
	brush := graphics.NewShaderBrushForTest()
	brushStyle := TextForegroundStyleFromBrush(brush, 0.5)

	color := graphics.ColorRed
	colorStyle := TextForegroundStyleFromColor(color)

	// Lerp at 0.25 should return brushStyle (discrete snap < 0.5)
	result1 := LerpTextForegroundStyle(brushStyle, colorStyle, 0.25)
	if !IsBrushStyle(result1) {
		t.Errorf("Expected brushStyle at fraction 0.25")
	}

	// Lerp at 0.75 should return colorStyle (discrete snap >= 0.5)
	result2 := LerpTextForegroundStyle(brushStyle, colorStyle, 0.75)
	if IsBrushStyle(result2) {
		t.Errorf("Expected colorStyle at fraction 0.75")
	}
}

// --- Tests for ModulateColor ---

func TestModulateColor_NaN(t *testing.T) {
	color := graphics.ColorRed
	result := ModulateColor(color, float32(math.NaN()))

	// Should return the color unchanged
	if color != result {
		t.Errorf("Expected unchanged color for NaN alpha")
	}
}

func TestModulateColor_FullAlpha(t *testing.T) {
	color := graphics.ColorRed
	result := ModulateColor(color, 1.0)

	// Should return the color unchanged
	if color != result {
		t.Errorf("Expected unchanged color for alpha 1.0")
	}
}

func TestModulateColor_HalfAlpha(t *testing.T) {
	color := graphics.ColorRed
	result := ModulateColor(color, 0.5)

	// Should return a modified color
	if result == color {
		t.Errorf("Expected modified color")
	}
	// Roughly check alpha
	if result.Alpha() == color.Alpha() {
		t.Errorf("Expected modified alpha")
	}
}
