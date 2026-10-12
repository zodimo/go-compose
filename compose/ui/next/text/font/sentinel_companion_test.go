package font

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for the companions in next/text/font ===

func TestFontWeightCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, FontWeightCompanion, FontWeightNormal, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, FontWeightCompanion, FontWeightNormal, FontWeightBold)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, FontWeightCompanion)
	})
}

func TestFontStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, FontStyleCompanion, FontStyleNormal, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, FontStyleCompanion, FontStyleNormal, FontStyleItalic)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, FontStyleCompanion)
	})
}

// === Conformance tests for FontSynthesisCompanion (struct/pointer) ===

func TestFontSynthesisCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, FontSynthesisCompanion, (*FontSynthesis)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := FontSynthesisWeight
		override := FontSynthesisStyle
		sentinel.AssertMergeDirection(t, FontSynthesisCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, FontSynthesisCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, FontSynthesisCompanion, FontSynthesisAll)
	})
	t.Run("copy_leaves_source_unchanged", func(t *testing.T) {
		source := &FontSynthesis{value: synthesisWeightFlag}
		sentinel.AssertCopyLeavesSource(t, FontSynthesisCompanion, source,
			func(f *FontSynthesis) { f.value = synthesisStyleFlag },
			func(f *FontSynthesis) any {
				if f == nil {
					return -1
				}
				return f.value
			})
	})
}
