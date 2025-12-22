package style

import "fmt"

// LineHeightStyle configuration for line height.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/LineHeightStyle.kt
type LineHeightStyle struct {
	Alignment LineHeightStyleAlignment
	Trim      LineHeightStyleTrim
	Mode      LineHeightStyleMode
}

// LineHeightStyleAlignment defines how to align the line in the space provided by the line height.
type LineHeightStyleAlignment struct {
	TopRatio float32
}

var (
	// LineHeightStyleAlignmentTop aligns the line to the top of the space reserved for that line.
	LineHeightStyleAlignmentTop = LineHeightStyleAlignment{TopRatio: 0}
	// LineHeightStyleAlignmentCenter aligns the line to the center of the space reserved for the line.
	LineHeightStyleAlignmentCenter = LineHeightStyleAlignment{TopRatio: 0.5}
	// LineHeightStyleAlignmentProportional aligns the line proportional to the ascent and descent values of the line.
	LineHeightStyleAlignmentProportional = LineHeightStyleAlignment{TopRatio: -1}
	// LineHeightStyleAlignmentBottom aligns the line to the bottom of the space reserved for that line.
	LineHeightStyleAlignmentBottom = LineHeightStyleAlignment{TopRatio: 1}
)

func (a LineHeightStyleAlignment) String() string {
	switch a {
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
)

func (t LineHeightStyleTrim) IsTrimFirstLineTop() bool {
	return t&flagTrimTop > 0
}

func (t LineHeightStyleTrim) IsTrimLastLineBottom() bool {
	return t&flagTrimBottom > 0
}

// LineHeightStyleMode defines if the specified line height value should be enforced.
type LineHeightStyleMode int

const (
	// LineHeightStyleModeFixed guarantees that taller glyphs won't be trimmed at the boundaries.
	LineHeightStyleModeFixed LineHeightStyleMode = 0
	// LineHeightStyleModeMinimum prevents the overflow of tall glyphs in middle lines.
	LineHeightStyleModeMinimum LineHeightStyleMode = 1
	// LineHeightStyleModeTight gets rid of the safety rails that are added by Fixed.
	LineHeightStyleModeTight LineHeightStyleMode = 2
)

var DefaultLineHeightStyle = LineHeightStyle{
	Alignment: LineHeightStyleAlignmentProportional,
	Trim:      LineHeightStyleTrimBoth,
	Mode:      LineHeightStyleModeFixed,
}

func (l LineHeightStyle) Equals(other LineHeightStyle) bool {
	return l.Alignment == other.Alignment &&
		l.Trim == other.Trim &&
		l.Mode == other.Mode
}
