package style

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// textIndentUnspecified is the sentinel for TextIndent. It is deliberately
// package-private: the only public way to obtain it is
// TextIndentCompanion.Unspecified().
var textIndentUnspecified = &TextIndent{FirstLine: unit.TextUnitCompanion.Unspecified(), RestLine: unit.TextUnitCompanion.Unspecified()}

//https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/TextIndent.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=32

// TextIndent specifies the indentation of a paragraph.
// FirstLine is the amount of indentation applied to the first line.
// RestLine is the amount of indentation applied to every line except the first line.
type TextIndent struct {
	FirstLine unit.TextUnit
	RestLine  unit.TextUnit
}

// TextIndentNone is a constant for no text indent.
var TextIndentNone = &TextIndent{FirstLine: unit.Sp(0), RestLine: unit.Sp(0)}

// NewTextIndent creates a new TextIndent with the given first line and rest line indentation.
func NewTextIndent(firstLine, restLine unit.TextUnit) *TextIndent {
	return &TextIndent{
		FirstLine: firstLine,
		RestLine:  restLine,
	}
}

// Copy returns a new TextIndent with optionally updated fields.
// Pass nil to use the current value for that field.
func (ti TextIndent) Copy(firstLine, restLine *unit.TextUnit) *TextIndent {
	fl := ti.FirstLine
	if firstLine != nil {
		fl = *firstLine
	}
	rl := ti.RestLine
	if restLine != nil {
		rl = *restLine
	}
	return &TextIndent{
		FirstLine: fl,
		RestLine:  rl,
	}
}

// HashCode returns a hash code for the TextIndent.
func (ti TextIndent) HashCode() int {
	result := hashTextUnit(ti.FirstLine)
	result = 31*result + hashTextUnit(ti.RestLine)
	return result
}

// hashTextUnit returns a hash code for a TextUnit.
func hashTextUnit(tu unit.TextUnit) int {
	// Simple hash combining type and value
	return int(tu.Type())*31 + int(tu.Value()*1000)
}

// lerpDiscrete returns a if fraction < 0.5, otherwise b.
//
// This is used for values that cannot be transitioned smoothly (e.g., enums,
// nullable types, or "unspecified" sentinel values). Instead of interpolating,
// the function performs a discrete "snap" at the halfway point.
//
// Example behavior:
//   - fraction 0.0 to 0.49 → returns a
//   - fraction 0.5 to 1.0  → returns b
func lerpDiscrete[T any](a, b T, fraction float32) T {
	if fraction < 0.5 {
		return a
	}
	return b
}

// LerpTextUnitInheritable linearly interpolates between two TextUnits,
// with special handling for Unspecified values.
//
// This function exists to support style merging and animation scenarios where
// TextUnit values may be Unspecified (meaning "inherit from parent style").
// Since there is no meaningful way to linearly interpolate between a concrete
// value (e.g., 16.sp) and an abstract "unspecified" sentinel, this function
// falls back to discrete interpolation when either value is Unspecified.
//
// Behavior:
//   - If BOTH values are specified: performs standard linear interpolation
//   - If EITHER value is Unspecified: performs discrete interpolation (snaps at 0.5)
//
// Discrete interpolation example (animating from Unspecified to 16.sp):
//   - fraction 0.0 to 0.49 → returns Unspecified (inherits from parent)
//   - fraction 0.5 to 1.0  → returns 16.sp
//
// This matches Kotlin's lerpTextUnitInheritable behavior in Jetpack Compose.
func LerpTextUnitInheritable(a, b unit.TextUnit, fraction float32) unit.TextUnit {
	if unit.TextUnitCompanion.IsUnspecified(a) || unit.TextUnitCompanion.IsUnspecified(b) {
		return lerpDiscrete(a, b, fraction)
	}
	return unit.LerpTextUnit(a, b, fraction)
}

// LerpTextIndent linearly interpolates between two TextIndents.
//
// The fraction argument represents position on the timeline, with 0.0 meaning that the
// interpolation has not started, returning start (or something equivalent to start), 1.0
// meaning that the interpolation has finished, returning stop (or something equivalent to
// stop), and values in between meaning that the interpolation is at the relevant point on the
// timeline between start and stop. The interpolation can be extrapolated beyond 0.0 and 1.0, so
// negative values and values greater than 1.0 are valid.
func LerpTextIndent(start, stop *TextIndent, fraction float32) *TextIndent {
	start = CoalesceTextIndent(start, textIndentUnspecified)
	stop = CoalesceTextIndent(stop, textIndentUnspecified)

	return &TextIndent{
		FirstLine: LerpTextUnitInheritable(start.FirstLine, stop.FirstLine, fraction),
		RestLine:  LerpTextUnitInheritable(start.RestLine, stop.RestLine, fraction),
	}
}

// String returns a string representation of the TextIndent.
func StringTextIndent(s *TextIndent) string {
	if !isSpecifiedTextIndent(s) {
		return "TextIndentUnspecified"
	}
	return fmt.Sprintf("TextIndent(firstLine=%s, restLine=%s)",
		s.FirstLine, s.RestLine)
}

func SameTextIndent(a, b *TextIndent) bool {
	return sameTextIndent(a, b)
}

func EqualTextIndent(a, b *TextIndent) bool {
	return equalTextIndent(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextIndent(a, b *TextIndent) bool {
	return semanticEqualTextIndent(a, b)
}

func CoalesceTextIndent(ptr, def *TextIndent) *TextIndent {
	return coalesceTextIndent(ptr, def)
}

func IsSpecifiedTextIndent(indent *TextIndent) bool {
	return TextIndentCompanion.IsSpecified(indent)
}

func TakeOrElseTextIndent(a, b *TextIndent) *TextIndent {
	return takeOrElseTextIndent(a, b)
}

// === Sentinel operations ===

func isSpecifiedTextIndent(indent *TextIndent) bool {
	return indent != nil && indent != textIndentUnspecified
}

func isUnspecifiedTextIndent(indent *TextIndent) bool {
	return !isSpecifiedTextIndent(indent)
}

func takeOrElseTextIndent(a, b *TextIndent) *TextIndent {
	if isSpecifiedTextIndent(a) {
		return a
	}
	return b
}

func coalesceTextIndent(ptr, def *TextIndent) *TextIndent {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeTextIndent(a, b *TextIndent) *TextIndent {
	if isUnspecifiedTextIndent(b) {
		return coalesceTextIndent(a, textIndentUnspecified)
	}
	return b
}

func stringTextIndent(s *TextIndent) string {
	return StringTextIndent(s)
}

func sameTextIndent(a, b *TextIndent) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textIndentUnspecified
	}
	if b == nil {
		return a == textIndentUnspecified
	}
	return a == b
}

func semanticEqualTextIndent(a, b *TextIndent) bool {

	a = coalesceTextIndent(a, textIndentUnspecified)
	b = coalesceTextIndent(b, textIndentUnspecified)

	if a == textIndentUnspecified && b == textIndentUnspecified {
		return true
	}

	return a.FirstLine.Equals(b.FirstLine) &&
		a.RestLine.Equals(b.RestLine)
}

func equalTextIndent(a, b *TextIndent) bool {
	if !sameTextIndent(a, b) {
		return semanticEqualTextIndent(a, b)
	}
	return true
}

// copyTextIndent returns a copy of ti with the given fields set on top; the
// source is never mutated.
func copyTextIndent(ti *TextIndent, options ...func(*TextIndent)) *TextIndent {
	cp := &TextIndent{}
	if isSpecifiedTextIndent(ti) {
		cp.FirstLine = ti.FirstLine
		cp.RestLine = ti.RestLine
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// TextIndentCompanion is the single exported sentinel-family symbol for TextIndent.
var TextIndentCompanion = companions.NewStructCompanion(
	func() *TextIndent { return textIndentUnspecified },
	stringTextIndent,
	isSpecifiedTextIndent,
	isUnspecifiedTextIndent,
	takeOrElseTextIndent,
	mergeTextIndent,
	sameTextIndent,
	semanticEqualTextIndent,
	equalTextIndent,
	coalesceTextIndent,
	copyTextIndent,
)
