package style

import (
	"fmt"

	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// lineHeightStyleUnspecified is the sentinel for LineHeightStyle. It is
// deliberately package-private: the only public way to obtain it is
// LineHeightStyleCompanion.Unspecified().
var lineHeightStyleUnspecified = &LineHeightStyle{
	Alignment: lineHeightStyleAlignmentUnspecified,
	Trim:      LineHeightStyleTrimUnspecified,
	Mode:      LineHeightStyleModeUnspecified,
}

// LineHeightStyle configuration for line height.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/LineHeightStyle.kt
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
	lineHeightStyleAlignmentUnspecified = &LineHeightStyleAlignment{TopRatio: -1}

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
	case lineHeightStyleAlignmentUnspecified:
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
	return a != nil && a != lineHeightStyleAlignmentUnspecified
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
	return isSpecifiedLineHeightStyle(s)
}

// Identity (2 ns)
func SameLineHeightStyle(a, b *LineHeightStyle) bool {
	return sameLineHeightStyle(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualLineHeightStyle(a, b *LineHeightStyle) bool {
	return semanticEqualLineHeightStyle(a, b)
}

func EqualLineHeightStyle(a, b *LineHeightStyle) bool {
	return equalLineHeightStyle(a, b)
}

func CoalesceLineHeightStyle(ptr, def *LineHeightStyle) *LineHeightStyle {
	return coalesceLineHeightStyle(ptr, def)
}

func MergeLineHeightStyle(a, b *LineHeightStyle) *LineHeightStyle {
	return mergeLineHeightStyle(a, b)
}

func TakeOrElseLineHeightStyle(s, def *LineHeightStyle) *LineHeightStyle {
	return takeOrElseLineHeightStyle(s, def)
}

// === Sentinel operations ===

func isSpecifiedLineHeightStyle(s *LineHeightStyle) bool {
	return s != nil && s != lineHeightStyleUnspecified
}

func isUnspecifiedLineHeightStyle(s *LineHeightStyle) bool {
	return !isSpecifiedLineHeightStyle(s)
}

func takeOrElseLineHeightStyle(s, def *LineHeightStyle) *LineHeightStyle {
	if !isSpecifiedLineHeightStyle(s) {
		return def
	}
	return s
}

func coalesceLineHeightStyle(ptr, def *LineHeightStyle) *LineHeightStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeLineHeightStyle(a, b *LineHeightStyle) *LineHeightStyle {
	a = coalesceLineHeightStyle(a, lineHeightStyleUnspecified)
	b = coalesceLineHeightStyle(b, lineHeightStyleUnspecified)

	if a == lineHeightStyleUnspecified {
		return b
	}
	if b == lineHeightStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b wins where specified.
	return &LineHeightStyle{
		Alignment: TakeOrElseLineHeightStyleAlignment(b.Alignment, a.Alignment),
		Trim:      b.Trim.TakeOrElse(a.Trim),
		Mode:      b.Mode.TakeOrElse(a.Mode),
	}
}

func stringLineHeightStyle(s *LineHeightStyle) string {
	return StringLineHeightStyle(s)
}

func sameLineHeightStyle(a, b *LineHeightStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == lineHeightStyleUnspecified
	}
	if b == nil {
		return a == lineHeightStyleUnspecified
	}
	return a == b
}

func semanticEqualLineHeightStyle(a, b *LineHeightStyle) bool {

	a = coalesceLineHeightStyle(a, lineHeightStyleUnspecified)
	b = coalesceLineHeightStyle(b, lineHeightStyleUnspecified)

	if a == lineHeightStyleUnspecified && b == lineHeightStyleUnspecified {
		return true
	}

	return a.Alignment == b.Alignment &&
		a.Trim == b.Trim &&
		a.Mode == b.Mode
}

func equalLineHeightStyle(a, b *LineHeightStyle) bool {
	if !sameLineHeightStyle(a, b) {
		return semanticEqualLineHeightStyle(a, b)
	}
	return true
}

// copyLineHeightStyle returns a copy of s; the source is never mutated.
func copyLineHeightStyle(s *LineHeightStyle, options ...func(*LineHeightStyle)) *LineHeightStyle {
	cp := &LineHeightStyle{}
	if isSpecifiedLineHeightStyle(s) {
		cp.Alignment = s.Alignment
		cp.Trim = s.Trim
		cp.Mode = s.Mode
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// LineHeightStyleCompanion is the single exported sentinel-family symbol for LineHeightStyle.
var LineHeightStyleCompanion = companions.NewStructCompanion(
	func() *LineHeightStyle { return lineHeightStyleUnspecified },
	stringLineHeightStyle,
	isSpecifiedLineHeightStyle,
	isUnspecifiedLineHeightStyle,
	takeOrElseLineHeightStyle,
	mergeLineHeightStyle,
	sameLineHeightStyle,
	semanticEqualLineHeightStyle,
	equalLineHeightStyle,
	coalesceLineHeightStyle,
	copyLineHeightStyle,
)

func StringLineHeightStyle(s *LineHeightStyle) string {
	if !IsSpecifiedLineHeightStyle(s) {
		return "LineHeightStyle.Unspecified"
	}
	return "LineHeightStyle(" +
		"Alignment=" + StringLineHeightStyleAlignment(s.Alignment) + ", " +
		"Trim=" + s.Trim.String() + ", " +
		"Mode=" + s.Mode.String() +
		")"
}

func (t LineHeightStyleTrim) String() string {
	switch t {
	case LineHeightStyleTrimFirstLineTop:
		return "LineHeightStyle.Trim.FirstLineTop"
	case LineHeightStyleTrimLastLineBottom:
		return "LineHeightStyle.Trim.LastLineBottom"
	case LineHeightStyleTrimBoth:
		return "LineHeightStyle.Trim.Both"
	case LineHeightStyleTrimNone:
		return "LineHeightStyle.Trim.None"
	case LineHeightStyleTrimUnspecified:
		return "LineHeightStyle.Trim.Unspecified"
	default:
		return "Invalid"
	}
}

func (m LineHeightStyleMode) String() string {
	switch m {
	case LineHeightStyleModeFixed:
		return "LineHeightStyle.Mode.Fixed"
	case LineHeightStyleModeMinimum:
		return "LineHeightStyle.Mode.Minimum"
	case LineHeightStyleModeTight:
		return "LineHeightStyle.Mode.Tight"
	case LineHeightStyleModeUnspecified:
		return "LineHeightStyle.Mode.Unspecified"
	default:
		return "Invalid"
	}
}
