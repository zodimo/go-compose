package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text/font"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for the companions in next/text ===

func TestTextStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, TextStyleCompanion, (*TextStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &TextStyle{}
		override := &TextStyle{}
		sentinel.AssertMergeDirection(t, TextStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, TextStyleCompanion, &TextStyle{})
	})
}

func TestSpanStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, SpanStyleCompanion, (*SpanStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		// SpanStyle merge combines field-by-field into a new SpanStyle
		// (b wins where specified), so the result is never pointer-identical
		// to the override. Assert the merged field instead.
		base := &SpanStyle{FontWeight: font.FontWeightNormal}
		override := &SpanStyle{FontWeight: font.FontWeightBold}

		if got := SpanStyleCompanion.Merge(base, SpanStyleCompanion.Unspecified()); got == nil || got.FontWeight != font.FontWeightNormal {
			t.Errorf("Merge(base, Unspecified()).FontWeight = %v, want %v", got, font.FontWeightNormal)
		}
		if got := SpanStyleCompanion.Merge(SpanStyleCompanion.Unspecified(), override); got == nil || got.FontWeight != font.FontWeightBold {
			t.Errorf("Merge(Unspecified(), override).FontWeight = %v, want %v", got, font.FontWeightBold)
		}
		if got := SpanStyleCompanion.Merge(base, override); got == nil || got.FontWeight != font.FontWeightBold {
			t.Errorf("Merge(base, override).FontWeight = %v, want the override's %v", got, font.FontWeightBold)
		}
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, SpanStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, SpanStyleCompanion, &SpanStyle{})
	})
}

func TestParagraphStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, ParagraphStyleCompanion, (*ParagraphStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &ParagraphStyle{TextDirection: style.TextDirectionLtr}
		override := &ParagraphStyle{TextDirection: style.TextDirectionRtl}
		sentinel.AssertMergeDirection(t, ParagraphStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, ParagraphStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, ParagraphStyleCompanion, &ParagraphStyle{})
	})
}

// === Conformance tests for the platform companions (struct/pointer) ===

func TestPlatformTextStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, PlatformTextStyleCompanion, (*PlatformTextStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &PlatformTextStyle{SpanStyle: &PlatformSpanStyle{}}
		override := &PlatformTextStyle{ParagraphStyle: &PlatformParagraphStyle{}}
		sentinel.AssertMergeDirection(t, PlatformTextStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, PlatformTextStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, PlatformTextStyleCompanion, &PlatformTextStyle{})
	})
}

func TestPlatformSpanStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, PlatformSpanStyleCompanion, (*PlatformSpanStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &PlatformSpanStyle{}
		override := &PlatformSpanStyle{}
		sentinel.AssertMergeDirection(t, PlatformSpanStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, PlatformSpanStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, PlatformSpanStyleCompanion, &PlatformSpanStyle{})
	})
}

func TestPlatformParagraphStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, PlatformParagraphStyleCompanion, (*PlatformParagraphStyle)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &PlatformParagraphStyle{}
		override := &PlatformParagraphStyle{}
		sentinel.AssertMergeDirection(t, PlatformParagraphStyleCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, PlatformParagraphStyleCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, PlatformParagraphStyleCompanion, &PlatformParagraphStyle{})
	})
}
