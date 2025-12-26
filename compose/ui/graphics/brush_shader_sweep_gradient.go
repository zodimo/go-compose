package graphics

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

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
	return &SweepGradientShader{
		Center:     geometry.NewOffset(centerX, centerY),
		Colors:     s.Colors,
		ColorStops: s.Stops,
	}
}

func SweepGradientBrush(colors []Color, center geometry.Offset) *SweepGradient {
	return &SweepGradient{
		Colors: colors,
		Center: center,
	}
}

func SemanticEqualSweepGradient(a, b *SweepGradient) bool {
	a = CoalesceBrush(a, BrushUnspecified).(*SweepGradient)
	b = CoalesceBrush(b, BrushUnspecified).(*SweepGradient)

	//center
	if !a.Center.Equal(b.Center) {
		return false
	}
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

	return true
}

func EqualSweepGradient(a, b *SweepGradient) bool {
	if !SameBrush(a, b) {
		return SemanticEqualSweepGradient(a, b)
	}
	return true
}
