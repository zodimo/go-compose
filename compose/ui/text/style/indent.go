package style

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

//https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextIndent.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=32

// TextIndent specifies the indentation of a paragraph.
// FirstLine is the amount of indentation applied to the first line.
// RestLine is the amount of indentation applied to every line except the first line.
type TextIndent struct {
	FirstLine unit.TextUnit
	RestLine  unit.TextUnit
}

// TextIndentNone is a constant for no text indent.
var TextIndentNone = TextIndent{FirstLine: unit.Sp(0), RestLine: unit.Sp(0)}

// NewTextIndent creates a new TextIndent with the given first line and rest line indentation.
func NewTextIndent(firstLine, restLine unit.TextUnit) TextIndent {
	return TextIndent{
		FirstLine: firstLine,
		RestLine:  restLine,
	}
}

// Copy returns a new TextIndent with optionally updated fields.
// Pass nil to use the current value for that field.
func (ti TextIndent) Copy(firstLine, restLine *unit.TextUnit) TextIndent {
	fl := ti.FirstLine
	if firstLine != nil {
		fl = *firstLine
	}
	rl := ti.RestLine
	if restLine != nil {
		rl = *restLine
	}
	return TextIndent{
		FirstLine: fl,
		RestLine:  rl,
	}
}

// Equals checks if two TextIndent instances are equal.
func (ti TextIndent) Equals(other TextIndent) bool {
	return ti.FirstLine.Equals(other.FirstLine) && ti.RestLine.Equals(other.RestLine)
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

// String returns a string representation of the TextIndent.
func (ti TextIndent) String() string {
	return fmt.Sprintf("TextIndent(firstLine=%s, restLine=%s)", ti.FirstLine, ti.RestLine)
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
	if a.IsUnspecified() || b.IsUnspecified() {
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
func LerpTextIndent(start, stop TextIndent, fraction float32) TextIndent {
	return TextIndent{
		FirstLine: LerpTextUnitInheritable(start.FirstLine, stop.FirstLine, fraction),
		RestLine:  LerpTextUnitInheritable(start.RestLine, stop.RestLine, fraction),
	}
}
