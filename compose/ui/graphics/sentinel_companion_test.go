package graphics

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for ColorCompanion (scalar) ===

func TestColorCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, ColorCompanion, Color(0), true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, ColorCompanion, ColorRed, ColorBlue)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, ColorCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		c := ColorRed
		if got := ColorCompanion.Copy(c); got != c {
			t.Errorf("Copy(%v) = %v, want the same value", c, got)
		}
	})
	t.Run("string_reads_as_unspecified", func(t *testing.T) {
		if got := ColorCompanion.String(ColorCompanion.Unspecified()); got != "Color.Unspecified" {
			t.Errorf("String(Unspecified()) = %q, want Color.Unspecified", got)
		}
	})
}

// === Conformance tests for ShadowCompanion (struct/pointer) ===

func TestShadowCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		// The zero *Shadow is nil, which is itself unspecified; only guard the
		// transposition here.
		sentinel.AssertSentinelIsNotZero(t, ShadowCompanion, (*Shadow)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := NewShadow(ColorRed, ZeroOffset, 2)
		override := NewShadow(ColorBlue, ZeroOffset, 4)
		sentinel.AssertMergeDirection(t, ShadowCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, ShadowCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, ShadowCompanion, ShadowNone)
	})
	t.Run("copy_leaves_source_unchanged", func(t *testing.T) {
		source := NewShadow(ColorRed, ZeroOffset, 2)
		sentinel.AssertCopyLeavesSource(t, ShadowCompanion, source,
			func(s *Shadow) { s.BlurRadius = 9 },
			func(s *Shadow) any {
				if s == nil {
					return "nil"
				}
				return s.BlurRadius
			})
	})
}

// === Conformance tests for BrushCompanion (struct/interface) ===

func TestBrushCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, BrushCompanion, Brush(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &SolidColor{Value: ColorRed}
		override := &SolidColor{Value: ColorBlue}
		sentinel.AssertMergeDirection(t, BrushCompanion, Brush(base), Brush(override))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, BrushCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		def := Brush(&SolidColor{Value: ColorRed})
		sentinel.AssertCoalesceFoldsNil(t, BrushCompanion, def)
	})
}
