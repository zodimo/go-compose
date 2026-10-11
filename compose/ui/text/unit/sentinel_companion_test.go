package unit

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for TextUnitCompanion ===
//
// TextUnit is a packed value type whose unspecified sentinel carries a NaN
// payload; `same` compares the packed bits while `semanticEqual` treats two
// unspecified values as equal.

func TestTextUnitCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		// TextUnit's zero value is an unspecified raw type, so the sentinel is not
		// meaningfully distinct from zero — assert only the transposition guard.
		sentinel.AssertSentinelIsNotZero(t, TextUnitCompanion, TextUnit{}, false)
	})
	t.Run("unspecified_is_not_specified", func(t *testing.T) {
		u := TextUnitCompanion.Unspecified()
		if TextUnitCompanion.IsSpecified(u) {
			t.Error("IsSpecified(Unspecified()) = true, want false")
		}
		if !TextUnitCompanion.IsUnspecified(u) {
			t.Error("IsUnspecified(Unspecified()) = false, want true")
		}
	})
	t.Run("merge_direction", func(t *testing.T) {
		sentinel.AssertMergeDirection(t, TextUnitCompanion, Sp(10), Sp(20))
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextUnitCompanion)
	})
	t.Run("copy_is_identity", func(t *testing.T) {
		tu := Sp(12)
		if got := TextUnitCompanion.Copy(tu); got != tu {
			t.Errorf("Copy(%v) = %v, want the same value", tu, got)
		}
	})
	t.Run("string_reads_as_unspecified", func(t *testing.T) {
		if got := TextUnitCompanion.String(TextUnitCompanion.Unspecified()); got != "TextUnit{Unspecified}" {
			t.Errorf("String(Unspecified()) = %q, want TextUnit{Unspecified}", got)
		}
	})
}
