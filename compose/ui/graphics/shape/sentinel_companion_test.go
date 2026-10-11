package shape

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for ShapeCompanion (struct/interface) ===

func TestShapeCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, ShapeCompanion, Shape(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := Shape(&CutCornerShape{Radius: 4})
		override := Shape(&CutCornerShape{Radius: 8})
		sentinel.AssertMergeDirection(t, ShapeCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, ShapeCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		def := Shape(&CutCornerShape{Radius: 4})
		sentinel.AssertCoalesceFoldsNil(t, ShapeCompanion, def)
	})
}
