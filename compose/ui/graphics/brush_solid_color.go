package graphics

import "github.com/zodimo/go-compose/compose/ui/geometry"

// SolidColor is a Brush that represents a single solid color.
type SolidColor struct {
	Value Color
}

func (s SolidColor) isBrush() {}

func (s SolidColor) ApplyTo(size geometry.Size, p *Paint, alpha float32) {
	p.Alpha = 1.0 // DefaultAlpha
	// In Kotlin: if (alpha != DefaultAlpha) value.copy(alpha = value.alpha * alpha) else value
	// For now, we just pass the color. Handling alpha on Color can use the Copy receiver or graphics.SetOpacity api
	// Since we don't have full Color logic validation yet, we'll set Paint fields.
	p.Color = s.Value
	if alpha != 1.0 {
		// p.Color = p.Color.SetOpacity(alpha)? This depends on graphics.Color interface.
		// For now we set Paint.Alpha which might be used by the renderer.
		p.Alpha = alpha
	}
	p.Shader = nil
}

func (s SolidColor) IntrinsicSize() geometry.Size {
	return geometry.SizeCompanion.Unspecified()
}

// NewSolidColor creates a new SolidColor brush from a Color.
func NewSolidColor(color Color) *SolidColor {
	return &SolidColor{Value: color}
}

func SemanticEqualSolidColor(a, b *SolidColor) bool {
	a = CoalesceBrush(a, brushUnspecified).(*SolidColor)
	b = CoalesceBrush(b, brushUnspecified).(*SolidColor)

	return a.Value == b.Value
}

func EqualSolidColor(a, b *SolidColor) bool {
	if !SameBrush(a, b) {
		return SemanticEqualSolidColor(a, b)
	}
	return true
}
