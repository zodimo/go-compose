package unit

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for DpCompanion, DpOffsetCompanion, DpSizeCompanion ===
//
// Dp's sentinel is NaN, so `sentinel != zero` and `unspecified == unspecified`
// are the two semantics that break under a naive `==` implementation.

func TestDpCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, DpCompanion, Dp(0), true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, DpCompanion, Dp(10), Dp(20))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, DpCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		d := Dp(10)
		if got := DpCompanion.Copy(d); got != d {
			t.Errorf("Copy(%v) = %v, want the same value", d, got)
		}
	})
	t.Run("unspecified_is_reported_by_string", func(t *testing.T) {
		if got := DpCompanion.String(DpCompanion.Unspecified()); got != "Dp{Unspecified}" {
			t.Errorf("String(Unspecified()) = %q, want Dp{Unspecified}", got)
		}
	})
}

func TestDpOffsetCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, DpOffsetCompanion, DpOffsetZero, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, DpOffsetCompanion,
			DpOffset{X: 10, Y: 20}, DpOffset{X: 30, Y: 40})
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, DpOffsetCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		o := DpOffset{X: 10, Y: 20}
		if got := DpOffsetCompanion.Copy(o); got != o {
			t.Errorf("Copy(%v) = %v, want the same value", o, got)
		}
	})
}

func TestDpSizeCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, DpSizeCompanion, DpSizeZero, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, DpSizeCompanion,
			DpSize{Width: 10, Height: 20}, DpSize{Width: 30, Height: 40})
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, DpSizeCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		s := DpSize{Width: 10, Height: 20}
		if got := DpSizeCompanion.Copy(s); got != s {
			t.Errorf("Copy(%v) = %v, want the same value", s, got)
		}
	})
}
