package unit

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// Dp is a value class that represents a density-independent pixel.
//
// The value is stored as a float32.
type Dp float32

// dpUnspecified is the constant for an unspecified Dp. It is deliberately
// package-private: the only public way to obtain it is DpCompanion.Unspecified().
var dpUnspecified = Dp(floatutils.Float32Unspecified)

// DpInfinity represents infinite dp dimension.
var DpInfinity = Dp(floatutils.FloatInfinite)

// DpHairline is the constant for a hairline Dp.
const DpHairline Dp = 0

// NewDp creates a new Dp.
func NewDp(value float32) Dp {
	return Dp(value)
}

// Value returns the float32 value.
func (d Dp) Value() float32 {
	return float32(d)
}

// Add adds two Dps.
func (d Dp) Add(other Dp) Dp {
	return Dp(d + other)
}

// Subtract subtracts other Dp from this Dp.
func (d Dp) Subtract(other Dp) Dp {
	return Dp(d - other)
}

// Negate returns -d.
func (d Dp) Negate() Dp {
	return Dp(-d)
}

// Div returns d / other (scalar).
func (d Dp) Div(other float32) Dp {
	return Dp(float32(d) / other)
}

// DivInt returns d / other (int).
func (d Dp) DivInt(other int) Dp {
	return Dp(float32(d) / float32(other))
}

// DivDp returns d / other (scalar).
func (d Dp) DivDp(other Dp) float32 {
	return float32(d) / float32(other)
}

// Times returns d * other (scalar).
func (d Dp) Times(other float32) Dp {
	return Dp(float32(d) * other)
}

// TimesInt returns d * other (int).
func (d Dp) TimesInt(other int) Dp {
	return Dp(float32(d) * float32(other))
}

// CompareTo compares this Dp to another Dp.
func (d Dp) CompareTo(other Dp) int {
	if d.IsUnspecified() || other.IsUnspecified() {
		if d.IsUnspecified() && other.IsUnspecified() {
			return 0
		}
		if d.IsUnspecified() {
			return 1 // NaN is "larger"
		}
		return -1
	}
	if floatutils.Float32Equals(float32(d), float32(other), floatutils.Float32EqualityThreshold) {
		return 0
	}
	if d < other {
		return -1
	}
	return 1
}

// String returns the string representation.
func (d Dp) String() string {
	if d.IsUnspecified() {
		return "Dp{Unspecified}"
	}
	return fmt.Sprintf("Dp{%.1f}", d)
}

// IsSpecified checks if the Dp is specified (not NaN).
func (d Dp) IsSpecified() bool {
	return isSpecifiedDp(d)
}

// IsUnspecified checks if the Dp is unspecified (NaN).
func (d Dp) IsUnspecified() bool {
	return isUnspecifiedDp(d)
}

// IsFinite returns true when it is finite.
func (d Dp) IsFinite() bool {
	return !math.IsInf(float64(d), 0) && !math.IsNaN(float64(d))
}

// === Sentinel operations ===

func isSpecifiedDp(d Dp) bool {
	return floatutils.IsSpecified(d)
}

func isUnspecifiedDp(d Dp) bool {
	return floatutils.IsUnspecified(d)
}

func takeOrElseDp(d, def Dp) Dp {
	if isSpecifiedDp(d) {
		return d
	}
	return def
}

// mergeDp layers b onto a: a value from b wins only where b is specified.
func mergeDp(a, b Dp) Dp {
	if isUnspecifiedDp(b) {
		return a
	}
	return b
}

func stringDp(d Dp) string {
	if isUnspecifiedDp(d) {
		return "Dp{Unspecified}"
	}
	return fmt.Sprintf("Dp{%.1f}", d)
}

func sameDp(a, b Dp) bool {
	if isUnspecifiedDp(a) && isUnspecifiedDp(b) {
		return true
	}
	return a == b
}

// semanticEqualDp treats two unspecified Dps as equal; comparing the NaN
// representation with == would otherwise report them as different.
func semanticEqualDp(a, b Dp) bool {
	if isUnspecifiedDp(a) && isUnspecifiedDp(b) {
		return true
	}
	return a == b
}

func equalDp(a, b Dp) bool {
	if isUnspecifiedDp(a) && isUnspecifiedDp(b) {
		return true
	}
	return floatutils.Float32Equals(float32(a), float32(b), floatutils.Float32EqualityThreshold)
}

func copyDp(d Dp) Dp {
	return d
}

// DpCompanion is the single exported sentinel-family symbol for Dp.
var DpCompanion = companions.NewScalarCompanion(
	func() Dp { return dpUnspecified },
	stringDp,
	isSpecifiedDp,
	isUnspecifiedDp,
	takeOrElseDp,
	mergeDp,
	sameDp,
	semanticEqualDp,
	equalDp,
	copyDp,
)

// MinDp returns the smaller of two Dps.
func MinDp(a, b Dp) Dp {
	if a < b {
		return a
	}
	return b
}

// MaxDp returns the larger of two Dps.
func MaxDp(a, b Dp) Dp {
	if a > b {
		return a
	}
	return b
}

// CoerceIn ensures value is in range [minimumValue, maximumValue].
func (d Dp) CoerceIn(minimumValue, maximumValue Dp) Dp {
	if d < minimumValue {
		return minimumValue
	}
	if d > maximumValue {
		return maximumValue
	}
	return d
}

// CoerceAtLeast ensures value is at least minimumValue.
func (d Dp) CoerceAtLeast(minimumValue Dp) Dp {
	if d < minimumValue {
		return minimumValue
	}
	return d
}

// CoerceAtMost ensures value is at most maximumValue.
func (d Dp) CoerceAtMost(maximumValue Dp) Dp {
	if d > maximumValue {
		return maximumValue
	}
	return d
}

// LerpDp linearly interpolates between two Dps.
func LerpDp(start, stop Dp, fraction float32) Dp {
	return Dp(lerpBetween(float32(start), float32(stop), float64(fraction)))
}

// DpOffset represents a 2D offset using Dp.
type DpOffset struct {
	X Dp
	Y Dp
}

// NewDpOffset constructs a DpOffset.
func NewDpOffset(x, y Dp) DpOffset {
	return DpOffset{X: x, Y: y}
}

var DpOffsetZero = DpOffset{X: 0, Y: 0}

// dpOffsetUnspecified is the sentinel for DpOffset, package-private so the only
// public door is DpOffsetCompanion.Unspecified().
var dpOffsetUnspecified = DpOffset{X: dpUnspecified, Y: dpUnspecified}

func (o DpOffset) Copy(x, y Dp) DpOffset {
	return DpOffset{X: x, Y: y}
}

func (o DpOffset) Subtract(other DpOffset) DpOffset {
	return DpOffset{X: o.X - other.X, Y: o.Y - other.Y}
}

func (o DpOffset) Add(other DpOffset) DpOffset {
	return DpOffset{X: o.X + other.X, Y: o.Y + other.Y}
}

func (o DpOffset) String() string {
	if o.IsSpecified() {
		return fmt.Sprintf("(%v, %v)", o.X, o.Y)
	}
	return "DpOffset.Unspecified"
}

func (o DpOffset) IsSpecified() bool {
	return isSpecifiedDpOffset(o)
}

// === DpOffset sentinel operations ===

func isSpecifiedDpOffset(o DpOffset) bool {
	return isSpecifiedDp(o.X) && isSpecifiedDp(o.Y)
}

// isUnspecifiedDpOffset reports true when either component is unspecified.
func isUnspecifiedDpOffset(o DpOffset) bool {
	return isUnspecifiedDp(o.X) || isUnspecifiedDp(o.Y)
	// Kotlin checks packedValue equality to NaN_NaN, so both must be NaN roughly.
	// However DpOffset.kt says: packedValue == 0x7fc00000_7fc00000L
	// which implies both are NaN.
	// But practically, if one is NaN, it's likely unspecified.
	// Let's stick strictly to "both must be valid" for IsSpecified.
}

func takeOrElseDpOffset(o, def DpOffset) DpOffset {
	if isSpecifiedDpOffset(o) {
		return o
	}
	return def
}

// mergeDpOffset layers b onto a: a value from b wins only where b is specified.
func mergeDpOffset(a, b DpOffset) DpOffset {
	return DpOffset{X: mergeDp(a.X, b.X), Y: mergeDp(a.Y, b.Y)}
}

func stringDpOffset(o DpOffset) string {
	if isSpecifiedDpOffset(o) {
		return fmt.Sprintf("(%v, %v)", o.X, o.Y)
	}
	return "DpOffset.Unspecified"
}

func sameDpOffset(a, b DpOffset) bool {
	return sameDp(a.X, b.X) && sameDp(a.Y, b.Y)
}

func semanticEqualDpOffset(a, b DpOffset) bool {
	return semanticEqualDp(a.X, b.X) && semanticEqualDp(a.Y, b.Y)
}

func equalDpOffset(a, b DpOffset) bool {
	return equalDp(a.X, b.X) && equalDp(a.Y, b.Y)
}

func copyDpOffset(o DpOffset) DpOffset {
	return o
}

// DpOffsetCompanion is the single exported sentinel-family symbol for DpOffset.
var DpOffsetCompanion = companions.NewScalarCompanion(
	func() DpOffset { return dpOffsetUnspecified },
	stringDpOffset,
	isSpecifiedDpOffset,
	isUnspecifiedDpOffset,
	takeOrElseDpOffset,
	mergeDpOffset,
	sameDpOffset,
	semanticEqualDpOffset,
	equalDpOffset,
	copyDpOffset,
)

func LerpDpOffset(start, stop DpOffset, fraction float32) DpOffset {
	return DpOffset{
		X: LerpDp(start.X, stop.X, fraction),
		Y: LerpDp(start.Y, stop.Y, fraction),
	}
}

// DpSize represents a size with Dp dimensions.
type DpSize struct {
	Width  Dp
	Height Dp
}

// NewDpSize creates a new DpSize.
func NewDpSize(width, height Dp) DpSize {
	return DpSize{Width: width, Height: height}
}

// dpSizeUnspecified is the sentinel for DpSize, package-private so the only
// public door is DpSizeCompanion.Unspecified().
var dpSizeUnspecified = DpSize{Width: dpUnspecified, Height: dpUnspecified}
var DpSizeZero = DpSize{Width: 0, Height: 0}

func (s DpSize) Copy(width, height Dp) DpSize {
	return DpSize{Width: width, Height: height}
}

func (s DpSize) Subtract(other DpSize) DpSize {
	return DpSize{Width: s.Width - other.Width, Height: s.Height - other.Height}
}

func (s DpSize) Add(other DpSize) DpSize {
	return DpSize{Width: s.Width + other.Width, Height: s.Height + other.Height}
}

func (s DpSize) Times(other float32) DpSize {
	return DpSize{Width: s.Width.Times(other), Height: s.Height.Times(other)}
}

func (s DpSize) TimesInt(other int) DpSize {
	return DpSize{Width: s.Width.TimesInt(other), Height: s.Height.TimesInt(other)}
}

func (s DpSize) Div(other float32) DpSize {
	return DpSize{Width: s.Width.Div(other), Height: s.Height.Div(other)}
}

func (s DpSize) DivInt(other int) DpSize {
	return DpSize{Width: s.Width.DivInt(other), Height: s.Height.DivInt(other)}
}

func (s DpSize) IsSpecified() bool {
	return isSpecifiedDpSize(s)
}

func (s DpSize) String() string {
	if isSpecifiedDpSize(s) {
		return fmt.Sprintf("%v x %v", s.Width, s.Height)
	}
	return "DpSize.Unspecified"
}

func (s DpSize) Center() DpOffset {
	return DpOffset{
		X: s.Width.Div(2),
		Y: s.Height.Div(2),
	}
}

// === DpSize sentinel operations ===

func isSpecifiedDpSize(s DpSize) bool {
	// Kotlin: packedValue != NaN_NaN. So if ANY part is not NaN, it is NOT Unspecified?
	// Wait, IsUnspecified is packedValue == NaN_NaN.
	// So IsSpecified is != NaN_NaN.
	// If one is NaN and other is 0, it is NOT IsUnspecified, so it IS IsSpecified?
	// No, DpSize is value class.
	// Let's assume strict checks: valid if both valid.
	return isSpecifiedDp(s.Width) && isSpecifiedDp(s.Height)
}

func isUnspecifiedDpSize(s DpSize) bool {
	return isUnspecifiedDp(s.Width) && isUnspecifiedDp(s.Height)
}

func takeOrElseDpSize(s, def DpSize) DpSize {
	if isSpecifiedDpSize(s) {
		return s
	}
	return def
}

// mergeDpSize layers b onto a: a value from b wins only where b is specified.
func mergeDpSize(a, b DpSize) DpSize {
	return DpSize{Width: mergeDp(a.Width, b.Width), Height: mergeDp(a.Height, b.Height)}
}

func stringDpSize(s DpSize) string {
	if isSpecifiedDpSize(s) {
		return fmt.Sprintf("%v x %v", s.Width, s.Height)
	}
	return "DpSize.Unspecified"
}

func sameDpSize(a, b DpSize) bool {
	return sameDp(a.Width, b.Width) && sameDp(a.Height, b.Height)
}

func semanticEqualDpSize(a, b DpSize) bool {
	return semanticEqualDp(a.Width, b.Width) && semanticEqualDp(a.Height, b.Height)
}

func equalDpSize(a, b DpSize) bool {
	return equalDp(a.Width, b.Width) && equalDp(a.Height, b.Height)
}

func copyDpSize(s DpSize) DpSize {
	return s
}

// DpSizeCompanion is the single exported sentinel-family symbol for DpSize.
var DpSizeCompanion = companions.NewScalarCompanion(
	func() DpSize { return dpSizeUnspecified },
	stringDpSize,
	isSpecifiedDpSize,
	isUnspecifiedDpSize,
	takeOrElseDpSize,
	mergeDpSize,
	sameDpSize,
	semanticEqualDpSize,
	equalDpSize,
	copyDpSize,
)

func LerpDpSize(start, stop DpSize, fraction float32) DpSize {
	return DpSize{
		Width:  LerpDp(start.Width, stop.Width, fraction),
		Height: LerpDp(start.Height, stop.Height, fraction),
	}
}

// DpRect represents a rectangle with Dp coordinates.
type DpRect struct {
	Left   Dp
	Top    Dp
	Right  Dp
	Bottom Dp
}

// NewDpRect creates a new DpRect.
func NewDpRect(left, top, right, bottom Dp) DpRect {
	return DpRect{Left: left, Top: top, Right: right, Bottom: bottom}
}

// NewDpRectFromOriginSize creates a DpRect from origin and size.
func NewDpRectFromOriginSize(origin DpOffset, size DpSize) DpRect {
	return DpRect{
		Left:   origin.X,
		Top:    origin.Y,
		Right:  origin.X + size.Width,
		Bottom: origin.Y + size.Height,
	}
}

func (r DpRect) Width() Dp {
	return r.Right - r.Left
}

func (r DpRect) Height() Dp {
	return r.Bottom - r.Top
}

func (r DpRect) Size() DpSize {
	return DpSize{Width: r.Width(), Height: r.Height()}
}

// String returns the string representation of DpRect.
func (r DpRect) String() string {
	return fmt.Sprintf("DpRect(left=%v, top=%v, right=%v, bottom=%v)", r.Left, r.Top, r.Right, r.Bottom)
}
