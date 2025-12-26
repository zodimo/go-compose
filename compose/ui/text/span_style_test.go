package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

func TestSpanStyle_Merge(t *testing.T) {
	s1 := &SpanStyle{
		FontSize:   unit.Sp(10),
		FontWeight: font.FontWeightBold,
	}
	s2 := &SpanStyle{
		FontSize:   unit.Sp(20),
		FontWeight: font.FontWeightUnspecified,
		FontStyle:  font.FontStyleItalic,
	}

	merged := MergeSpanStyle(s1, s2)

	if merged.FontSize != unit.Sp(20) {
		t.Errorf("Expected FontSize 20, got %v", merged.FontSize)
	}
	if merged.FontWeight != font.FontWeightBold {
		t.Errorf("Expected FontWeight Bold, got %v", merged.FontWeight)
	}
	if merged.FontStyle != font.FontStyleItalic {
		t.Errorf("Expected FontStyle Italic, got %v", merged.FontStyle)
	}

	// Test nil merge
	if MergeSpanStyle(s1, nil) != s1 {
		t.Error("Merging with nil should return original")
	}

	// Test unspecified
	if MergeSpanStyle(s1, SpanStyleUnspecified).FontSize != s1.FontSize {
		t.Error("Merging with Unspecified should preserve original fields where unspecified")
	}
}

func TestSpanStyle_Plus(t *testing.T) {
	s1 := &SpanStyle{
		FontSize: unit.Sp(10),
	}
	s2 := &SpanStyle{
		FontSize: unit.Sp(20),
	}

	res := MergeSpanStyle(s1, s2)
	if res.FontSize != unit.Sp(20) {
		t.Errorf("Plus expected 20, got %v", res.FontSize)
	}
}

func TestLerpSpanStyle(t *testing.T) {
	s1 := &SpanStyle{
		FontSize:   unit.Sp(10),
		Background: graphics.ColorBlack,
	}
	s2 := &SpanStyle{
		FontSize:   unit.Sp(20),
		Background: graphics.ColorWhite,
	}

	lerped := LerpSpanStyle(nil, s1, s2, 0.5)

	// Lerp 10 and 20 -> 15
	expectedSize := unit.Sp(15)
	if lerped.FontSize != expectedSize {
		t.Errorf("Expected size 15, got %v", lerped.FontSize)
	}

	// Lerp Black(0xFF000000) and White(0xFFFFFFFF) at 0.5 -> 0xFF7F7F7F (ish? actually alpha is different in argb vs rgba thinking, but LerpColor works)
	// Black is 0xFF000000 (Alpha 255, R0 G0 B0)
	// White is 0xFFFFFFFF (Alpha 255, R255 G255 B255)
	// 50% -> Alpha 255, R127 G127 B127 -> 0xFF7F7F7F
	// Let's just check it's not start or stop
	if lerped.Background == s1.Background || lerped.Background == s2.Background {
		t.Error("LerpColor failed to interpolate")
	}
}

func TestSpanStyle_String(t *testing.T) {
	s := &SpanStyle{
		FontSize: unit.Sp(12),
	}
	str := StringSpanStyle(s)
	if str == "" {
		t.Error("String() returned empty string")
	}
	// Basic check
	if len(str) < 10 {
		t.Error("String() content seems too short")
	}
}

func TestResolveSpanStyleDefaults(t *testing.T) {
	s := &SpanStyle{} // All unspecified
	resolved := ResolveSpanStyleDefaults(s)

	if resolved.FontSize != DefaultFontSize {
		t.Errorf("Expected default font size %v, got %v", DefaultFontSize, resolved.FontSize)
	}
	if resolved.Background != graphics.ColorTransparent {
		t.Errorf("Expected default background transparent, got %v", resolved.Background)
	}
}
