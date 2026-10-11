package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
)

func TestParagraphStyle_String(t *testing.T) {
	s := &ParagraphStyle{
		TextAlign:  style.TextAlignCenter,
		LineHeight: unit.Sp(12),
	}
	str := StringParagraphStyle(s)

	if str == "" {
		t.Error("String() returned empty string")
	}
	// Verify some content
	// TextAlign=Center
	// LineHeight=12.sp
	// TextDirection=Unspecified (default)
}

func TestParagraphStyle_Merge(t *testing.T) {
	s1 := &ParagraphStyle{
		TextAlign:  style.TextAlignLeft,
		LineHeight: unit.Sp(10),
	}
	s2 := &ParagraphStyle{
		TextAlign:  style.TextAlignRight,
		LineHeight: unit.TextUnitUnspecified,
	}

	merged := MergeParagraphStyle(s1, s2)

	if merged.TextAlign != style.TextAlignRight {
		t.Errorf("Expected TextAlign Right, got %v", merged.TextAlign)
	}
	if merged.LineHeight != unit.Sp(10) {
		t.Errorf("Expected LineHeight 10.sp, got %v", merged.LineHeight)
	}

	// Test nil
	if MergeParagraphStyle(s1, nil) != s1 {
		t.Error("Merge with nil should return original")
	}
}

func TestParagraphStyle_Lerp(t *testing.T) {
	s1 := &ParagraphStyle{
		LineHeight: unit.Sp(10),
	}
	s2 := &ParagraphStyle{
		LineHeight: unit.Sp(20),
	}

	lerped := LerpParagraphStyle(s1, s2, 0.5)

	if lerped.LineHeight != unit.Sp(15) {
		t.Errorf("Expected LineHeight 15.sp, got %v", lerped.LineHeight)
	}

	// Discrete lerp test
	s3 := &ParagraphStyle{
		TextAlign: style.TextAlignLeft,
	}
	s4 := &ParagraphStyle{
		TextAlign: style.TextAlignRight,
	}
	lerpedDiscrete := LerpParagraphStyle(s3, s4, 0.4)
	if lerpedDiscrete.TextAlign != style.TextAlignLeft {
		t.Errorf("Expected TextAlign Left at 0.4, got %v", lerpedDiscrete.TextAlign)
	}
	lerpedDiscrete = LerpParagraphStyle(s3, s4, 0.6)
	if lerpedDiscrete.TextAlign != style.TextAlignRight {
		t.Errorf("Expected TextAlign Right at 0.6, got %v", lerpedDiscrete.TextAlign)
	}
}

func TestResolveParagraphStyleDefaults(t *testing.T) {
	s := &ParagraphStyle{} // All defaults/unspecified
	resolved := ResolveParagraphStyleDefaults(s, unit.LayoutDirectionLtr)

	if resolved.TextAlign != style.TextAlignStart {
		t.Error("Expected default TextAlign Start")
	}
	if resolved.TextDirection != style.TextDirectionLtr {
		t.Error("Expected default TextDirection Ltr (from LayoutDirection)")
	}
	if resolved.LineHeight.IsSpecified() {
		// Should be Unspecified as it defaults to Unspecified in implementation?
		// "lineHeight = if (style.lineHeight.isUnspecified) DefaultLineHeight else style.lineHeight"
		// DefaultLineHeight is Unspecified.
		t.Error("Expected default LineHeight Unspecified")
	}
	// Check TextIndent default
	if resolved.TextIndent != style.TextIndentNone { // Pointer comparison might fail if not same instance, but TextIndentNone is global var
		// If TextIndentNone is a var, it might be same instance if assigned directly.
		// In ResolveParagraphStyleDefaults: textIndent = &style.TextIndentNone.
		// Wait, &style.TextIndentNone takes address of global. Correct.
		// But equality might need value check or SameTextIndent.
		if !style.SameTextIndent(resolved.TextIndent, style.TextIndentNone) {
			t.Error("Expected TextIndentNone")
		}
	}
}
