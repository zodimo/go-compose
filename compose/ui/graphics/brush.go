package graphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// brushUnspecified is the sentinel for Brush. It is deliberately package-private:
// the only public way to obtain it is BrushCompanion.Unspecified().
var brushUnspecified Brush = &SolidColor{Value: colorUnspecified}

// Brush is the interface for all brush types used for drawing.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Brush.kt
type Brush interface {
	ApplyTo(size geometry.Size, p *Paint, alpha float32)
	IntrinsicSize() geometry.Size

	//closed interface
	isBrush()
}

// ShaderBrush is a Brush implementation that wraps a shader.
type ShaderBrush interface {
	Brush
	CreateShader(size geometry.Size) Shader
}

// shaderBrushBase provides common implementation for ShaderBrush types.
// Go doesn't have abstract classes, so we duplicate ApplyTo or use embedding if struct.
// But ShaderBrush is an interface here to allow LinearGradient etc to implement it.
// We can make LinearGradient implement ApplyTo directly.

func applyToShaderBrush(s ShaderBrush, size geometry.Size, p *Paint, alpha float32) {
	p.Shader = s.CreateShader(size)
	p.Alpha = alpha
	p.Color = colorUnspecified
}

// Constructors (Top level functions instead of companion object)

// Helpers

func float32SliceEqual(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Compatibility constants and helpers

func IsSolidColor(b Brush) bool {
	_, ok := b.(*SolidColor)
	return ok
}

func IsShaderBrush(b Brush) bool {
	_, ok := b.(ShaderBrush)
	return ok
}

func AsSolidColor(b Brush) *SolidColor {
	if sc, ok := b.(*SolidColor); ok {
		return sc
	}
	return nil
}

func AsShaderBrush(b Brush) ShaderBrush {
	if sb, ok := b.(ShaderBrush); ok {
		return sb
	}
	return nil
}

// ShaderBrushForTest is an exported ShaderBrush implementation for cross-package testing.
// Use NewShaderBrushForTest() to create instances in tests.
type ShaderBrushForTest struct{}

func NewShaderBrushForTest() ShaderBrushForTest {
	return ShaderBrushForTest{}
}

func (s ShaderBrushForTest) isBrush()                                            {}
func (s ShaderBrushForTest) CreateShader(size geometry.Size) Shader              { return nil }
func (s ShaderBrushForTest) ApplyTo(size geometry.Size, p *Paint, alpha float32) {}
func (s ShaderBrushForTest) IntrinsicSize() geometry.Size {
	return geometry.SizeCompanion.Unspecified()
}
func (s ShaderBrushForTest) Equal(other Brush) bool {
	_, ok := other.(ShaderBrushForTest)
	return ok
}

// LerpBrush linearly interpolates between two brushes.
func LerpBrush(start, stop Brush, fraction float32) Brush {

	if fraction == 0 {
		return start
	}
	if fraction == 1 {
		return stop
	}

	if start == nil && stop == nil {
		return nil
	}
	if start == nil {
		start = SolidColor{Value: colorUnspecified} // Or Transparent? Kotlin uses Transparent
	}
	if stop == nil {
		stop = SolidColor{Value: colorUnspecified}
	}

	// Case 1: Both SolidColor
	if s1, ok1 := start.(*SolidColor); ok1 {
		if s2, ok2 := stop.(*SolidColor); ok2 {
			return &SolidColor{Value: LerpColor(s1.Value, s2.Value, fraction)}
		}
	}

	// Case 2: SolidColor to Gradient or Gradient to SolidColor
	// In Kotlin, SolidColor is converted to a Gradient of same colors.
	// We need to handle this.
	// But first, let's handle same-type gradients.

	if l1, ok1 := start.(*LinearGradient); ok1 {
		if l2, ok2 := stop.(*LinearGradient); ok2 {
			return &LinearGradient{
				Colors:   LerpColors(l1.Colors, l2.Colors, fraction),
				Stops:    lerp.FloatList32(l1.Stops, l2.Stops, fraction),
				Start:    geometry.LerpOffset(l1.Start, l2.Start, fraction),
				End:      geometry.LerpOffset(l1.End, l2.End, fraction),
				TileMode: l1.TileMode, // or l2.TileMode based on t < 0.5
			}
		}
	}

	if r1, ok1 := start.(*RadialGradient); ok1 {
		if r2, ok2 := stop.(*RadialGradient); ok2 {
			return &RadialGradient{
				Colors:   LerpColors(r1.Colors, r2.Colors, fraction),
				Stops:    lerp.FloatList32(r1.Stops, r2.Stops, fraction),
				Center:   geometry.LerpOffset(r1.Center, r2.Center, fraction),
				Radius:   lerp.Float32(r1.Radius, r2.Radius, fraction),
				TileMode: r1.TileMode,
			}
		}
	}

	// Case 3: Mixed types (SolidColor <-> Gradient)
	// TODO: fully implement mixed type interpolation by promoting SolidColor to Gradient
	// For now, discrete switch
	if fraction < 0.5 {
		return start
	}
	return stop
}

// Helpers for lerping lists

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// === Sentinel operations ===

func isSpecifiedBrush(s Brush) bool {
	return s != nil && s != brushUnspecified
}

func isUnspecifiedBrush(s Brush) bool {
	return !isSpecifiedBrush(s)
}

func takeOrElseBrush(s, def Brush) Brush {
	if !isSpecifiedBrush(s) {
		return def
	}
	return s
}

func coalesceBrush(ptr, def Brush) Brush {
	if ptr == nil {
		return def
	}
	return ptr
}

// sameBrush is identity (2 ns).
func sameBrush(a, b Brush) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == brushUnspecified
	}
	if b == nil {
		return a == brushUnspecified
	}
	// must be pointers
	return a == b
}

// semanticEqualBrush is field-by-field equality (20 ns).
func semanticEqualBrush(a, b Brush) bool {

	a = coalesceBrush(a, brushUnspecified)
	b = coalesceBrush(b, brushUnspecified)

	switch brush := a.(type) {
	case *SolidColor:
		if brushB, ok := b.(*SolidColor); ok {
			return brush.Value == brushB.Value
		}
		return false
	case *LinearGradient:
		if brushB, ok := b.(*LinearGradient); ok {
			return EqualLinearGradient(brush, brushB)
		}
	case *RadialGradient:
		if brush2, ok := b.(*RadialGradient); ok {
			return EqualRadialGradient(brush, brush2)
		}
	case *SweepGradient:
		if brush2, ok := b.(*SweepGradient); ok {
			return EqualSweepGradient(brush, brush2)
		}
	default:
		panic(fmt.Sprintf("unhandled brush type, %T", a))
	}
	return false
}

func equalBrush(a, b Brush) bool {
	if !sameBrush(a, b) {
		return semanticEqualBrush(a, b)
	}
	return true
}

func mergeBrush(a, b Brush) Brush {
	a = coalesceBrush(a, brushUnspecified)
	b = coalesceBrush(b, brushUnspecified)
	if a == brushUnspecified {
		return b
	}
	if b == brushUnspecified {
		return a
	}
	return b
}

func stringBrush(b Brush) string {
	if !isSpecifiedBrush(b) {
		return "Brush{Unspecified}"
	}
	return fmt.Sprintf("Brush{%T}", b)
}

func copyBrush(b Brush, options ...func(Brush)) Brush {
	return b
}

// BrushCompanion is the single exported sentinel-family symbol for Brush.
var BrushCompanion = companions.NewStructCompanion(
	func() Brush { return brushUnspecified },
	stringBrush,
	isSpecifiedBrush,
	isUnspecifiedBrush,
	takeOrElseBrush,
	mergeBrush,
	sameBrush,
	semanticEqualBrush,
	equalBrush,
	coalesceBrush,
	copyBrush,
)

// Short for IsSpecifiedShadow
func IsSpecifiedBrush(s Brush) bool {
	return BrushCompanion.IsSpecified(s)
}
func TakeOrElseBrush(s, def Brush) Brush {
	return takeOrElseBrush(s, def)
}

func CoalesceBrush(ptr, def Brush) Brush {
	return coalesceBrush(ptr, def)
}

// Identity (2 ns)
func SameBrush(a, b Brush) bool {
	return sameBrush(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualBrush(a, b Brush) bool {
	return semanticEqualBrush(a, b)
}

func EqualBrush(a, b Brush) bool {
	return equalBrush(a, b)
}
