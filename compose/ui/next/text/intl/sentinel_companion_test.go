package intl

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for LocaleListCompanion (struct/pointer) ===

func TestLocaleListCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		sentinel.AssertSentinelIsNotZero(t, LocaleListCompanion, (*LocaleList)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		base := &LocaleList{}
		override := &LocaleList{}
		sentinel.AssertMergeDirection(t, LocaleListCompanion, base, override)
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, LocaleListCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, LocaleListCompanion, LocaleListEmpty)
	})
	t.Run("copy_leaves_source_unchanged", func(t *testing.T) {
		source := &LocaleList{}
		sentinel.AssertCopyLeavesSource(t, LocaleListCompanion, source,
			func(l *LocaleList) { l.List = append(l.List, Locale{}) },
			func(l *LocaleList) any {
				if l == nil {
					return -1
				}
				return len(l.List)
			})
	})
}
