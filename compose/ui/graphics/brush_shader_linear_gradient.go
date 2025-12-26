package graphics

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

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
	return &LinearGradientShader{
		Colors:     l.Colors,
		ColorStops: l.Stops,
		From:       geometry.NewOffset(startX, startY),
		To:         geometry.NewOffset(endX, endY),
		TileMode:   l.TileMode,
	}
}

func LinearGradientBrush(colors []Color, start, end geometry.Offset, tileMode TileMode) *LinearGradient {
	// Defaults are handled by caller or explicitly passed.
	// In Kotlin: start=Zero, end=Infinite, tileMode=Clamp
	return &LinearGradient{
		Colors:   colors,
		Start:    start,
		End:      end,
		TileMode: tileMode,
	}
}

func LinearGradientBrushWithStops(colorStops []struct {
	Stop  float32
	Color Color
}, start, end geometry.Offset, tileMode TileMode) *LinearGradient {
	colors := make([]Color, len(colorStops))
	stops := make([]float32, len(colorStops))
	for i, cs := range colorStops {
		colors[i] = cs.Color
		stops[i] = cs.Stop
	}
	return &LinearGradient{
		Colors:   colors,
		Stops:    stops,
		Start:    start,
		End:      end,
		TileMode: tileMode,
	}
}

func HorizontalGradient(colors []Color, startX, endX float32, tileMode TileMode) *LinearGradient {
	return LinearGradientBrush(
		colors,
		geometry.NewOffset(startX, 0.0),
		geometry.NewOffset(endX, 0.0),
		tileMode,
	)
}

func VerticalGradient(colors []Color, startY, endY float32, tileMode TileMode) *LinearGradient {
	return LinearGradientBrush(
		colors,
		geometry.NewOffset(0.0, startY),
		geometry.NewOffset(0.0, endY),
		tileMode,
	)
}

func SemanticEqualLinearGradient(a, b *LinearGradient) bool {
	a = CoalesceBrush(a, BrushUnspecified).(*LinearGradient)
	b = CoalesceBrush(b, BrushUnspecified).(*LinearGradient)

	// colors
	if len(a.Colors) != len(b.Colors) {
		return false
	}
	for i := range a.Colors {
		if a.Colors[i] != b.Colors[i] {
			return false
		}
	}
	//stops
	if !float32SliceEqual(a.Stops, b.Stops) {
		return false
	}
	//start
	if !a.Start.Equal(b.Start) {
		return false
	}
	//end
	if !a.End.Equal(b.End) {
		return false
	}
	//tileMode
	if a.TileMode != b.TileMode {
		return false
	}

	return true
}

func EqualLinearGradient(a, b *LinearGradient) bool {
	if !SameBrush(a, b) {
		return SemanticEqualLinearGradient(a, b)
	}
	return true
}
