package border

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme"
)

type BorderData struct {
	Width Dp
	Shape Shape
	Color theme.ColorDescriptor
}

type BorderElement struct {
	borderData BorderData
}

func (e *BorderElement) Create() Node {
	return NewBorderNode(*e)
}

func (e *BorderElement) Update(node Node) {
	n := node.(*BorderNode)
	n.borderData = e.borderData
}

func (e *BorderElement) Equals(other Element) bool {
	if otherEle, ok := other.(*BorderElement); ok {
		return e.borderData.Width == otherEle.borderData.Width &&
			e.borderData.Shape == otherEle.borderData.Shape &&
			e.borderData.Color.Compare(otherEle.borderData.Color)
	}
	return false
}

func Border(width Dp, colorDesc theme.ColorDescriptor, shape Shape) Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&BorderElement{
				borderData: BorderData{
					Width: width,
					Shape: shape,
					Color: colorDesc,
				},
			},
		),
		modifier.NewInspectorInfo(
			"border",
			map[string]any{
				"width":     width,
				"shape":     shape,
				"colorDesc": colorDesc,
			},
		),
	)
}

// Border with defaults usually needs width and color at least?
// For Shape, default to Rectangle if nil? Or caller handles it.
// Default to Rectangle if nil.
func Simple(width Dp, colorDesc theme.ColorDescriptor) Modifier {
	return Border(width, colorDesc, shape.ShapeRectangle)
}
