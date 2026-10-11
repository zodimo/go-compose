package style

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for the scalar companions in next/text/style ===

func TestTextAlignCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextAlignCompanion, TextAlignCenter, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, TextAlignCompanion, TextAlignLeft, TextAlignRight)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextAlignCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		v := TextAlignCenter
		if got := TextAlignCompanion.Copy(v); got != v {
			t.Errorf("Copy(%v) = %v, want the same value", v, got)
		}
	})
}

func TestTextDirectionCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextDirectionCompanion, TextDirectionLtr, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, TextDirectionCompanion, TextDirectionLtr, TextDirectionRtl)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextDirectionCompanion)
	})
}

func TestLineBreakCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, LineBreakCompanion, LineBreakSimple, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, LineBreakCompanion, LineBreakSimple, LineBreakParagraph)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, LineBreakCompanion)
	})
}

func TestHyphensCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, HyphensCompanion, HyphensNone, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, HyphensCompanion, HyphensNone, HyphensAuto)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, HyphensCompanion)
	})
}

func TestBaselineShiftCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		// BaselineShiftNone is 0, a real value distinct from the NaN sentinel.
		sentinel.AssertSentinelIsNotZero(t, BaselineShiftCompanion, BaselineShiftNone, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, BaselineShiftCompanion, BaselineShiftSubscript, BaselineShiftSuperscript)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, BaselineShiftCompanion)
	})
}

// === Conformance tests for the struct companions in next/text/style ===

func TestTextIndentCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextIndentCompanion, (*TextIndent)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := NewTextIndent(unit.Sp(4), unit.Sp(4))
		override := NewTextIndent(unit.Sp(8), unit.Sp(8))
		sentinel.AssertMergeDirection(t, TextIndentCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextIndentCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, TextIndentCompanion, TextIndentNone)
	})
}

func TestLineHeightStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, LineHeightStyleCompanion, (*LineHeightStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &LineHeightStyle{Trim: LineHeightStyleTrimFirstLineTop, Mode: LineHeightStyleModeFixed}
		override := &LineHeightStyle{Trim: LineHeightStyleTrimBoth, Mode: LineHeightStyleModeMinimum}
		sentinel.AssertMergeDirection(t, LineHeightStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, LineHeightStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, LineHeightStyleCompanion, DefaultLineHeightStyle)
	})
}

func TestTextGeometricTransformCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextGeometricTransformCompanion, (*TextGeometricTransform)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &TextGeometricTransform{ScaleX: 1, SkewX: 0}
		override := &TextGeometricTransform{ScaleX: 2, SkewX: 1}
		sentinel.AssertMergeDirection(t, TextGeometricTransformCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextGeometricTransformCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, TextGeometricTransformCompanion, TextGeometricTransformNone)
	})
}

func TestTextForegroundStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextForegroundStyleCompanion, (*TextForegroundStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := TextForegroundStyleFromColor(graphics.ColorRed)
		override := TextForegroundStyleFromColor(graphics.ColorBlue)
		sentinel.AssertMergeDirection(t, TextForegroundStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextForegroundStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, TextForegroundStyleCompanion, TextForegroundStyleFromColor(graphics.ColorRed))
	})
}
