package shadow

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/modifier"
)

type ShadowData struct {
	Elevation    Dp
	Shape        Shape
	AmbientColor graphics.Color
	SpotColor    graphics.Color
}

type ShadowElement struct {
	shadowData ShadowData
}

func (e *ShadowElement) Create() Node {
	return NewShadowNode(*e)
}

func (e *ShadowElement) Update(node Node) {
	n := node.(*ShadowNode)
	n.shadowData = e.shadowData
}

func (e *ShadowElement) Equals(other Element) bool {
	if otherEle, ok := other.(*ShadowElement); ok {
		return e.shadowData.Elevation == otherEle.shadowData.Elevation &&
			e.shadowData.Shape == otherEle.shadowData.Shape // Pointer comparison for shape might be enough if not changed
	}
	return false
}

func Shadow(elevation Dp, shape Shape, ambientColor, spotColor graphics.Color) Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&ShadowElement{
				shadowData: ShadowData{
					Elevation:    elevation,
					Shape:        shape,
					AmbientColor: ambientColor,
					SpotColor:    spotColor,
				},
			},
		),
		modifier.NewInspectorInfo(
			"shadow",
			map[string]any{
				"elevation":    elevation,
				"shape":        shape,
				"ambientColor": ambientColor,
				"spotColor":    spotColor,
			},
		),
	)
}

// Simple Shadow with defaults using theme Shadow role
func Simple(elevation Dp, shape Shape) Modifier {
	// unspecified allow for theme defaults
	unspecifiedColor := graphics.ColorUnspecified
	return Shadow(elevation, shape, unspecifiedColor, unspecifiedColor)
}
