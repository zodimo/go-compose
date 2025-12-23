package graphics

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/theme"
)

var BrushUnspecified Brush = nil

// Brush is the interface for all brush types used for drawing.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Brush.kt
type Brush interface {
	ApplyTo(size geometry.Size, p *Paint, alpha float32)
	IntrinsicSize() geometry.Size
	Equal(other Brush) bool
	isBrush()
}

// SolidColor is a Brush that represents a single solid color.
type SolidColor struct {
	Value Color
}

func (s SolidColor) isBrush() {}

func (s SolidColor) ApplyTo(size geometry.Size, p *Paint, alpha float32) {
	p.Alpha = 1.0 // DefaultAlpha
	// In Kotlin: if (alpha != DefaultAlpha) value.copy(alpha = value.alpha * alpha) else value
	// For now, we just pass the color. Handling alpha on ColorDescriptor might require AppendUpdate.
	// Since we don't have full Color logic validation yet, we'll set Paint fields.
	p.Color = s.Value
	if alpha != 1.0 {
		// p.Color = p.Color.SetOpacity(alpha)? This depends on ColorDescriptor interface.
		// For now we set Paint.Alpha which might be used by the renderer.
		p.Alpha = alpha
	}
	p.Shader = nil
}

func (s SolidColor) IntrinsicSize() geometry.Size {
	return geometry.SizeUnspecified
}

func (s SolidColor) Equal(other Brush) bool {
	if other, ok := other.(SolidColor); ok {
		// Use ColorDescriptor comparison
		// return s.Value.Compare(other.Value)
		// But s.Value is an interface.
		return s.Value.Compare(other.Value)
	}
	return false
}

// NewSolidColor creates a new SolidColor brush from a Color.
func NewSolidColor(color Color) SolidColor {
	return SolidColor{Value: color}
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
	p.Color = nil // Or unspecified
}

// LinearGradient Brush implementation.
type LinearGradient struct {
	Colors   []Color
	Stops    []float32
	Start    geometry.Offset
	End      geometry.Offset
	TileMode TileMode
}

func (l LinearGradient) isBrush() {}

func (l LinearGradient) ApplyTo(size geometry.Size, p *Paint, alpha float32) {
	applyToShaderBrush(l, size, p, alpha)
}

func (l LinearGradient) IntrinsicSize() geometry.Size {
	width := float32(math.NaN())
	height := float32(math.NaN())
	if l.Start.IsFinite() && l.End.IsFinite() {
		width = float32(math.Abs(float64(l.Start.X() - l.End.X())))
		height = float32(math.Abs(float64(l.Start.Y() - l.End.Y())))
	}
	return geometry.NewSize(width, height)
}

func (l LinearGradient) CreateShader(size geometry.Size) Shader {
	startX := l.Start.X()
	if startX == float32(math.Inf(1)) {
		startX = size.Width()
	}
	startY := l.Start.Y()
	if startY == float32(math.Inf(1)) {
		startY = size.Height()
	}
	endX := l.End.X()
	if endX == float32(math.Inf(1)) {
		endX = size.Width()
	}
	endY := l.End.Y()
	if endY == float32(math.Inf(1)) {
		endY = size.Height()
	}
	return LinearGradientShader{
		Colors:     l.Colors,
		ColorStops: l.Stops,
		From:       geometry.NewOffset(startX, startY),
		To:         geometry.NewOffset(endX, endY),
		TileMode:   l.TileMode,
	}
}

func (l LinearGradient) Equal(other Brush) bool {
	o, ok := other.(LinearGradient)
	if !ok {
		return false
	}
	if len(l.Colors) != len(o.Colors) {
		return false
	}
	for i := range l.Colors {
		if !l.Colors[i].Compare(o.Colors[i]) {
			return false
		}
	}
	// Check stops, start, end, tileMode
	if !float32SliceEqual(l.Stops, o.Stops) {
		return false
	}
	if !l.Start.Equal(o.Start) {
		return false
	}
	if !l.End.Equal(o.End) {
		return false
	}
	if l.TileMode != o.TileMode {
		return false
	}
	return true
}

// RadialGradient Brush implementation.
type RadialGradient struct {
	Colors   []Color
	Stops    []float32
	Center   geometry.Offset
	Radius   float32
	TileMode TileMode
}

func (r RadialGradient) isBrush() {}

func (r RadialGradient) ApplyTo(size geometry.Size, p *Paint, alpha float32) {
	applyToShaderBrush(r, size, p, alpha)
}

func (r RadialGradient) IntrinsicSize() geometry.Size {
	if !math.IsInf(float64(r.Radius), 0) && !math.IsNaN(float64(r.Radius)) {
		return geometry.NewSize(r.Radius*2, r.Radius*2)
	}
	return geometry.SizeUnspecified
}

func (r RadialGradient) CreateShader(size geometry.Size) Shader {
	centerX := r.Center.X()
	centerY := r.Center.Y()
	if r.Center.IsUnspecified() {
		center := size.Center()
		centerX = center.X()
		centerY = center.Y()
	} else {
		if centerX == float32(math.Inf(1)) {
			centerX = size.Width()
		}
		if centerY == float32(math.Inf(1)) {
			centerY = size.Height()
		}
	}
	radius := r.Radius
	if radius == float32(math.Inf(1)) {
		radius = size.MinDimension() / 2
	}

	return RadialGradientShader{
		Colors:     r.Colors,
		ColorStops: r.Stops,
		Center:     geometry.NewOffset(centerX, centerY),
		Radius:     radius,
		TileMode:   r.TileMode,
	}
}

func (r RadialGradient) Equal(other Brush) bool {
	o, ok := other.(RadialGradient)
	if !ok {
		return false
	}
	if len(r.Colors) != len(o.Colors) {
		return false
	}
	for i := range r.Colors {
		if !r.Colors[i].Compare(o.Colors[i]) {
			return false
		}
	}
	if !float32SliceEqual(r.Stops, o.Stops) {
		return false
	}
	if !r.Center.Equal(o.Center) {
		return false
	}
	if r.Radius != o.Radius {
		return false
	}
	if r.TileMode != o.TileMode {
		return false
	}
	return true
}

// SweepGradient Brush implementation.
type SweepGradient struct {
	Center geometry.Offset
	Colors []Color
	Stops  []float32
}

func (s SweepGradient) isBrush() {}

func (s SweepGradient) ApplyTo(size geometry.Size, p *Paint, alpha float32) {
	applyToShaderBrush(s, size, p, alpha)
}

func (s SweepGradient) IntrinsicSize() geometry.Size {
	return geometry.SizeUnspecified
}

func (s SweepGradient) CreateShader(size geometry.Size) Shader {
	centerX := s.Center.X()
	centerY := s.Center.Y()
	if s.Center.IsUnspecified() {
		center := size.Center()
		centerX = center.X()
		centerY = center.Y()
	} else {
		if centerX == float32(math.Inf(1)) {
			centerX = size.Width()
		}
		if centerY == float32(math.Inf(1)) {
			centerY = size.Height()
		}
	}
	return SweepGradientShader{
		Center:     geometry.NewOffset(centerX, centerY),
		Colors:     s.Colors,
		ColorStops: s.Stops,
	}
}

func (s SweepGradient) Equal(other Brush) bool {
	o, ok := other.(SweepGradient)
	if !ok {
		return false
	}
	if !s.Center.Equal(o.Center) {
		return false
	}
	if len(s.Colors) != len(o.Colors) {
		return false
	}
	for i := range s.Colors {
		if !s.Colors[i].Compare(o.Colors[i]) {
			return false
		}
	}
	if !float32SliceEqual(s.Stops, o.Stops) {
		return false
	}
	return true
}

// Constructors (Top level functions instead of companion object)

func LinearGradientBrush(colors []Color, start, end geometry.Offset, tileMode TileMode) LinearGradient {
	// Defaults are handled by caller or explicitly passed.
	// In Kotlin: start=Zero, end=Infinite, tileMode=Clamp
	return LinearGradient{
		Colors:   colors,
		Start:    start,
		End:      end,
		TileMode: tileMode,
	}
}

func LinearGradientBrushWithStops(colorStops []struct {
	Stop  float32
	Color Color
}, start, end geometry.Offset, tileMode TileMode) LinearGradient {
	colors := make([]Color, len(colorStops))
	stops := make([]float32, len(colorStops))
	for i, cs := range colorStops {
		colors[i] = cs.Color
		stops[i] = cs.Stop
	}
	return LinearGradient{
		Colors:   colors,
		Stops:    stops,
		Start:    start,
		End:      end,
		TileMode: tileMode,
	}
}

func HorizontalGradient(colors []Color, startX, endX float32, tileMode TileMode) LinearGradient {
	return LinearGradientBrush(
		colors,
		geometry.NewOffset(startX, 0.0),
		geometry.NewOffset(endX, 0.0),
		tileMode,
	)
}

func VerticalGradient(colors []Color, startY, endY float32, tileMode TileMode) LinearGradient {
	return LinearGradientBrush(
		colors,
		geometry.NewOffset(0.0, startY),
		geometry.NewOffset(0.0, endY),
		tileMode,
	)
}

func RadialGradientBrush(colors []Color, center geometry.Offset, radius float32, tileMode TileMode) RadialGradient {
	// Defaults: center=Unspecified, radius=Infinite, tileMode=Clamp
	return RadialGradient{
		Colors:   colors,
		Center:   center,
		Radius:   radius,
		TileMode: tileMode,
	}
}

func SweepGradientBrush(colors []Color, center geometry.Offset) SweepGradient {
	return SweepGradient{
		Colors: colors,
		Center: center,
	}
}

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
	_, ok := b.(SolidColor)
	return ok
}

func IsShaderBrush(b Brush) bool {
	_, ok := b.(ShaderBrush)
	return ok
}

func AsSolidColor(b Brush) *SolidColor {
	if sc, ok := b.(SolidColor); ok {
		return &sc
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
func (s ShaderBrushForTest) IntrinsicSize() geometry.Size                        { return geometry.SizeUnspecified }
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
		start = SolidColor{Value: ColorUnspecified} // Or Transparent? Kotlin uses Transparent
	}
	if stop == nil {
		stop = SolidColor{Value: ColorUnspecified}
	}

	// Case 1: Both SolidColor
	if s1, ok1 := start.(SolidColor); ok1 {
		if s2, ok2 := stop.(SolidColor); ok2 {
			// Resolve colors. For now assume Color type has methods or we use theme helper.
			// But SolidColor.Value is theme.ColorDescriptor.
			// theme.Lerp exists.
			// However theme.Lerp returns LerpColorUpdate.
			// We need a helper to actually lerp two ColorDescriptors if possible,
			// or just wrap them in a SolidColor that resolves at draw time?
			// The theme package seems to rely on updates.
			// But for now, let's look at how SolidColor.Value is used.
			// If we can't easily lerp descriptors directly without context,
			// maybe we just return one or the other if types mismatch or use a specialized Lerp implementation.
			// Actually, Kotlin implementation:
			// return SolidColor(lerp(value, other.value, t))
			// Color.kt lerp:
			// fun lerp(start: Color, stop: Color, fraction: Float): Color
			//
			// In our Go theme package, we have Lerp(stop, fraction).
			// But that's an Update.
			// If we want to return a NEW Brush that represents the interpolated state,
			// we might need a "LerpColor" brush or "LerpColorDescriptor".
			// OR, if the underlying ColorDescriptor supports Lerp (it doesn't seem to verifyably return a new descriptor directly in the interface).
			//
			// HOWEVER, `theme.ColorDescriptor` interface doesn't have `Lerp` method.
			// But `theme` package might have `LerpColors(c1, c2, t)`.
			// Let's assume for now we cannot fully implement Color lerp without more theme helpers
			// so we will stub it or use a discrete switch for 0.5 if we can't lerp.
			//
			// Wait, `theme/color.go` had `Lerp(stop, fraction)` which returned `LerpColorUpdate`.
			// This suggests the "Color" is dynamic.
			// So we can return a SolidColor whose Value is `start.Value.AppendUpdate(theme.Lerp(stop.Value, fraction))`.
			// Yes! That seems to be the design.
			return NewSolidColor(s1.Value.AppendUpdate(theme.Lerp(s2.Value, fraction)))
		}
	}

	// Case 2: SolidColor to Gradient or Gradient to SolidColor
	// In Kotlin, SolidColor is converted to a Gradient of same colors.
	// We need to handle this.
	// But first, let's handle same-type gradients.

	if l1, ok1 := start.(LinearGradient); ok1 {
		if l2, ok2 := stop.(LinearGradient); ok2 {
			return LinearGradient{
				Colors:   lerpColors(l1.Colors, l2.Colors, fraction),
				Stops:    lerpFloats(l1.Stops, l2.Stops, fraction),
				Start:    geometry.LerpOffset(l1.Start, l2.Start, fraction),
				End:      geometry.LerpOffset(l1.End, l2.End, fraction),
				TileMode: l1.TileMode, // or l2.TileMode based on t < 0.5
			}
		}
	}

	if r1, ok1 := start.(RadialGradient); ok1 {
		if r2, ok2 := stop.(RadialGradient); ok2 {
			return RadialGradient{
				Colors:   lerpColors(r1.Colors, r2.Colors, fraction),
				Stops:    lerpFloats(r1.Stops, r2.Stops, fraction),
				Center:   geometry.LerpOffset(r1.Center, r2.Center, fraction),
				Radius:   lerpFloat(r1.Radius, r2.Radius, fraction),
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

func lerpColors(a, b []Color, t float32) []Color {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	res := make([]Color, n)
	for i := 0; i < n; i++ {
		c1 := a[min(i, len(a)-1)]
		c2 := b[min(i, len(b)-1)]
		// Use the AppendUpdate pattern for lerping colors
		res[i] = c1.AppendUpdate(theme.Lerp(c2, t))
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func lerpFloats(a, b []float32, t float32) []float32 {
	if a == nil && b == nil {
		return nil
	}
	if a == nil {
		a = b // Treat as if 'a' has same stops as 'b' but maybe we should ensure size?
		// Kotlin says: "if other == null return null".
		// But here we might be lerping [0, 1] to [0, 0.5, 1].
		// Let's following Kotlin logic simplified:
		// If either is null/empty, we might just return the other or interpolate to default?
		// Kotlin: lerpNullableFloatList
		// if (right == null || left == null) return null
		// This implies if one doesn't have stops (evenly distributed), result doesn't either?
		// Wait, linearGradient without stops implies even distribution.
		if b == nil {
			return nil
		}
		// So if either is nil, return nil (meaning evenly distributed result?)
		// If we lerp from explicit stops to implicit stops, the result should probably be explicit?
		// But Kotlin returns null if either is null.
		return nil
	}
	if b == nil {
		return nil
	}

	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	res := make([]float32, n)
	for i := 0; i < n; i++ {
		f1 := a[min(i, len(a)-1)]
		f2 := b[min(i, len(b)-1)]
		res[i] = (1-t)*f1 + t*f2
	}
	return res
}

func lerpFloat(a, b, t float32) float32 {
	return (1-t)*a + t*b
}
