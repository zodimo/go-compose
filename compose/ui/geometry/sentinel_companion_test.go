package geometry

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for OffsetCompanion, SizeCompanion, CornerRadiusCompanion ===
//
// These assert the semantics that survive review for the go-sentinel-helper
// contract: sentinel != zero, merge direction (both ways), two unspecified values
// are semantically equal (the packed NaN sentinel breaks naive ==), and copy is
// identity for these value types.

func TestOffsetCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, OffsetCompanion, OffsetZero, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, OffsetCompanion, NewOffset(10, 20), NewOffset(30, 40))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, OffsetCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		o := NewOffset(10, 20)
		if got := OffsetCompanion.Copy(o); got != o {
			t.Errorf("Copy(%v) = %v, want the same value", o, got)
		}
	})
}

func TestSizeCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, SizeCompanion, SizeZero, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, SizeCompanion, NewSize(10, 20), NewSize(30, 40))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, SizeCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		s := NewSize(10, 20)
		if got := SizeCompanion.Copy(s); got != s {
			t.Errorf("Copy(%v) = %v, want the same value", s, got)
		}
	})
}

func TestCornerRadiusCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, CornerRadiusCompanion, CornerRadiusZero, true)
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, CornerRadiusCompanion, NewCornerRadius(4, 4), NewCornerRadius(8, 8))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, CornerRadiusCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		c := NewCornerRadius(4, 4)
		if got := CornerRadiusCompanion.Copy(c); got != c {
			t.Errorf("Copy(%v) = %v, want the same value", c, got)
		}
	})
}
