package style

import (
	"testing"

	"github.com/zodimo/go-compose/internal/testutils/sentinel"
)

// === Conformance tests for TextDecorationCompanion (struct/pointer) ===

func TestTextDecorationCompanionConformance(t *testing.T) {
	t.Run("sentinel_is_not_zero", func(t *testing.T) {
		// TextDecorationUnspecified is a distinct pointer, not the nil zero value.
		sentinel.AssertSentinelIsNotZero(t, TextDecorationCompanion, (*TextDecoration)(nil), false)
	})
	t.Run("merge_direction", func(t *testing.T) {
		// TextDecoration merge is additive: it ORs the masks. The sentinel
		// contributes nothing; two specified decorations combine.
		base := TextDecorationUnderline
		override := TextDecorationLineThrough
		if got := TextDecorationCompanion.Merge(base, TextDecorationCompanion.Unspecified()); !TextDecorationCompanion.SemanticEqual(got, base) {
			t.Errorf("Merge(base, Unspecified()) != base — the sentinel must contribute nothing")
		}
		if got := TextDecorationCompanion.Merge(TextDecorationCompanion.Unspecified(), override); !TextDecorationCompanion.SemanticEqual(got, override) {
			t.Errorf("Merge(Unspecified(), override) != override — a specified value must win over the sentinel")
		}
		combined := TextDecorationCompanion.Merge(base, override)
		if !(combined.Contains(base) && combined.Contains(override)) {
			t.Errorf("Merge(underline, lineThrough) = %v, want a decoration containing both", combined)
		}
	})
	t.Run("unspecified_semantic_equal", func(t *testing.T) {
		sentinel.AssertUnspecifiedSemanticEqual(t, TextDecorationCompanion)
	})
	t.Run("coalesce_folds_nil_onto_default", func(t *testing.T) {
		sentinel.AssertCoalesceFoldsNil(t, TextDecorationCompanion, TextDecorationNone)
	})
	t.Run("copy_leaves_source_unchanged", func(t *testing.T) {
		source := &TextDecoration{mask: TextDecorationMaskUnderline}
		sentinel.AssertCopyLeavesSource(t, TextDecorationCompanion, source,
			func(d *TextDecoration) { d.mask |= TextDecorationMaskLineThrough },
			func(d *TextDecoration) any {
				if d == nil {
					return TextDecorationMaskUnspecified
				}
				return d.mask
			})
	})
}
