package font

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for FontWeightCompanion and FontStyleCompanion (scalar) ===

func TestFontWeightCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, FontWeightCompanion, FontWeight(0), true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, FontWeightCompanion, FontWeightNormal, FontWeightBold)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, FontWeightCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		w := FontWeightBold
		if got := FontWeightCompanion.Copy(w); got != w {
			t.Errorf("Copy(%v) = %v, want the same value", w, got)
		}
	})
}

func TestFontStyleCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		// FontStyleNormal is 0, which is a real value distinct from the -1 sentinel.
		sentinel.AssertSentinelIsNotZero(t, FontStyleCompanion, FontStyleNormal, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, FontStyleCompanion, FontStyleNormal, FontStyleItalic)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, FontStyleCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		s := FontStyleItalic
		if got := FontStyleCompanion.Copy(s); got != s {
			t.Errorf("Copy(%v) = %v, want the same value", s, got)
		}
	})
}
