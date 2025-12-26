package style

import (
	"fmt"
)

var LineHeightStyleUnspecified *LineHeightStyle = &LineHeightStyle{
	Alignment: LineHeightStyleAlignmentUnspecified,
	Trim:      LineHeightStyleTrimUnspecified,
	Mode:      LineHeightStyleModeUnspecified,
}

// LineHeightStyle configuration for line height.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/LineHeightStyle.kt
type LineHeightStyle struct {
	Alignment *LineHeightStyleAlignment
	Trim      LineHeightStyleTrim
	Mode      LineHeightStyleMode
}

// LineHeightStyleAlignment defines how to align the line in the space provided by the line height.
type LineHeightStyleAlignment struct {
	TopRatio float32
}

var (
	LineHeightStyleAlignmentUnspecified = &LineHeightStyleAlignment{TopRatio: -1}

	// LineHeightStyleAlignmentTop aligns the line to the top of the space reserved for that line.
	LineHeightStyleAlignmentTop = &LineHeightStyleAlignment{TopRatio: 0}
	// LineHeightStyleAlignmentCenter aligns the line to the center of the space reserved for the line.
	LineHeightStyleAlignmentCenter = &LineHeightStyleAlignment{TopRatio: 0.5}
	// LineHeightStyleAlignmentProportional aligns the line proportional to the ascent and descent values of the line.
	LineHeightStyleAlignmentProportional = &LineHeightStyleAlignment{TopRatio: -1}
	// LineHeightStyleAlignmentBottom aligns the line to the bottom of the space reserved for that line.
	LineHeightStyleAlignmentBottom = &LineHeightStyleAlignment{TopRatio: 1}
)

func StringLineHeightStyleAlignment(a *LineHeightStyleAlignment) string {
	switch a {
	case LineHeightStyleAlignmentUnspecified:
		return "LineHeightStyle.Alignment.Unspecified"
	case LineHeightStyleAlignmentTop:
		return "LineHeightStyle.Alignment.Top"
	case LineHeightStyleAlignmentCenter:
		return "LineHeightStyle.Alignment.Center"
	case LineHeightStyleAlignmentProportional:
		return "LineHeightStyle.Alignment.Proportional"
	case LineHeightStyleAlignmentBottom:
		return "LineHeightStyle.Alignment.Bottom"
	default:
		return "LineHeightStyle.Alignment(topRatio = " + float32ToString(a.TopRatio) + ")"
	}
}

func IsSpecifiedLineHeightStyleAlignment(a *LineHeightStyleAlignment) bool {
	return a != nil && a != LineHeightStyleAlignmentUnspecified
}
func TakeOrElseLineHeightStyleAlignment(a, b *LineHeightStyleAlignment) *LineHeightStyleAlignment {
	if !IsSpecifiedLineHeightStyleAlignment(a) {
		return b
	}
	return a
}

// Helper for float to string conversion, simple implementation
func float32ToString(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

// LineHeightStyleTrim defines whether to trim the extra space from the top of the first line and the bottom of the last line of text.
type LineHeightStyleTrim int

const (
	flagTrimTop    = 0x00000001
	flagTrimBottom = 0x00000010

	LineHeightStyleTrimFirstLineTop   LineHeightStyleTrim = flagTrimTop
	LineHeightStyleTrimLastLineBottom LineHeightStyleTrim = flagTrimBottom
	LineHeightStyleTrimBoth           LineHeightStyleTrim = flagTrimTop | flagTrimBottom
	LineHeightStyleTrimNone           LineHeightStyleTrim = 0
	LineHeightStyleTrimUnspecified    LineHeightStyleTrim = -1
)

func (t LineHeightStyleTrim) IsTrimFirstLineTop() bool {
	return t&flagTrimTop > 0
}

func (t LineHeightStyleTrim) IsTrimLastLineBottom() bool {
	return t&flagTrimBottom > 0
}
func (t LineHeightStyleTrim) IsTrimUnspecified() bool {
	return t == LineHeightStyleTrimUnspecified
}

func (t LineHeightStyleTrim) TakeOrElse(other LineHeightStyleTrim) LineHeightStyleTrim {
	if t.IsTrimUnspecified() {
		return other
	}
	return t
}

// LineHeightStyleMode defines if the specified line height value should be enforced.
type LineHeightStyleMode int

const (
	LineHeightStyleModeUnspecified LineHeightStyleMode = -1
	// LineHeightStyleModeFixed guarantees that taller glyphs won't be trimmed at the boundaries.
	LineHeightStyleModeFixed LineHeightStyleMode = 0
	// LineHeightStyleModeMinimum prevents the overflow of tall glyphs in middle lines.
	LineHeightStyleModeMinimum LineHeightStyleMode = 1
	// LineHeightStyleModeTight gets rid of the safety rails that are added by Fixed.
	LineHeightStyleModeTight LineHeightStyleMode = 2
)

func (m LineHeightStyleMode) IsSpecified() bool {
	return m != LineHeightStyleModeUnspecified
}
func (m LineHeightStyleMode) TakeOrElse(other LineHeightStyleMode) LineHeightStyleMode {
	if !m.IsSpecified() {
		return other
	}
	return m
}

var DefaultLineHeightStyle = &LineHeightStyle{
	Alignment: LineHeightStyleAlignmentProportional,
	Trim:      LineHeightStyleTrimBoth,
	Mode:      LineHeightStyleModeFixed,
}

func IsSpecifiedLineHeightStyle(s *LineHeightStyle) bool {
	return s != nil && s != LineHeightStyleUnspecified
}

// Identity (2 ns)
func SameLineHeightStyle(a, b *LineHeightStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == LineHeightStyleUnspecified
	}
	if b == nil {
		return a == LineHeightStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualLineHeightStyle(a, b *LineHeightStyle) bool {

	a = CoalesceLineHeightStyle(a, LineHeightStyleUnspecified)
	b = CoalesceLineHeightStyle(b, LineHeightStyleUnspecified)

	return a.Alignment == b.Alignment &&
		a.Trim == b.Trim &&
		a.Mode == b.Mode
}

func EqualLineHeightStyle(a, b *LineHeightStyle) bool {
	if !SameLineHeightStyle(a, b) {
		return SemanticEqualLineHeightStyle(a, b)
	}
	return true
}

func CoalesceLineHeightStyle(ptr, def *LineHeightStyle) *LineHeightStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func MergeLineHeightStyle(a, b *LineHeightStyle) *LineHeightStyle {
	a = CoalesceLineHeightStyle(a, LineHeightStyleUnspecified)
	b = CoalesceLineHeightStyle(b, LineHeightStyleUnspecified)

	if a == LineHeightStyleUnspecified {
		return b
	}
	if b == LineHeightStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &LineHeightStyle{
		Alignment: TakeOrElseLineHeightStyleAlignment(a.Alignment, b.Alignment),
		Trim:      a.Trim.TakeOrElse(b.Trim),
		Mode:      a.Mode.TakeOrElse(b.Mode),
	}
}

func TakeOrElseLineHeightStyle(s, def *LineHeightStyle) *LineHeightStyle {
	if !IsSpecifiedLineHeightStyle(s) {
		return def
	}
	return s
}
