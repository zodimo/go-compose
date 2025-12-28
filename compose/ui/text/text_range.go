package text

import "fmt"

// TextRange is an immutable text range, representing a range from Start (inclusive) to End (exclusive).
// End can be smaller than Start; use Min() and Max() to get ordered values.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/TextRange.kt
type TextRange struct {
	Start int
	End   int
}

// TextRangeZero is a TextRange with both start and end at 0.
var TextRangeZero = TextRange{Start: 0, End: 0}

// NewTextRange creates a new TextRange with the given start and end.
// Panics if start or end is negative.
func NewTextRange(start, end int) TextRange {
	if start < 0 || end < 0 {
		panic(fmt.Sprintf("start and end cannot be negative. [start: %d, end: %d]", start, end))
	}
	return TextRange{Start: start, End: end}
}

// NewTextRangeCollapsed creates a TextRange where start equals end.
func NewTextRangeCollapsed(index int) TextRange {
	return NewTextRange(index, index)
}

// Min returns the minimum offset of the range.
func (r TextRange) Min() int {
	if r.Start < r.End {
		return r.Start
	}
	return r.End
}

// Max returns the maximum offset of the range.
func (r TextRange) Max() int {
	if r.Start > r.End {
		return r.Start
	}
	return r.End
}

// Collapsed returns true if the range is collapsed (start == end).
func (r TextRange) Collapsed() bool {
	return r.Start == r.End
}

// Reversed returns true if the start offset is larger than the end offset.
func (r TextRange) Reversed() bool {
	return r.Start > r.End
}

// Length returns the length of the range.
func (r TextRange) Length() int {
	return r.Max() - r.Min()
}

// Intersects returns true if the given range has intersection with this range.
func (r TextRange) Intersects(other TextRange) bool {
	return r.Min() < other.Max() && other.Min() < r.Max()
}

// Contains returns true if this range covers including equals with the given range.
func (r TextRange) Contains(other TextRange) bool {
	return r.Min() <= other.Min() && other.Max() <= r.Max()
}

// ContainsOffset returns true if the given offset is a part of this range.
func (r TextRange) ContainsOffset(offset int) bool {
	return offset >= r.Min() && offset < r.Max()
}

// String returns a string representation of the TextRange.
func (r TextRange) String() string {
	return fmt.Sprintf("TextRange(%d, %d)", r.Start, r.End)
}

// Equals checks equality with another TextRange.
func (r TextRange) Equals(other TextRange) bool {
	return r.Start == other.Start && r.End == other.End
}

// CoerceIn ensures that Start and End values lie in the specified range.
// For each value:
// - if value is smaller than minimumValue, value is replaced by minimumValue
// - if value is greater than maximumValue, value is replaced by maximumValue
func (r TextRange) CoerceIn(minimumValue, maximumValue int) TextRange {
	newStart := coerceIn(r.Start, minimumValue, maximumValue)
	newEnd := coerceIn(r.End, minimumValue, maximumValue)
	if newStart != r.Start || newEnd != r.End {
		return TextRange{Start: newStart, End: newEnd}
	}
	return r
}

// coerceIn clamps value to [min, max].
func coerceIn(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Substring returns the substring of s corresponding to this TextRange.
func (r TextRange) Substring(s string) string {

	strLen := len(s)
	clamp := func(value, min, max int) int {
		if value < min {
			return min
		}
		if value > max {
			return max
		}
		return value
	}
	start := clamp(r.Min(), 0, strLen)
	end := clamp(r.Max(), 0, strLen)

	return s[start:end]
}
