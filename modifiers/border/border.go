package border

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
)

type BorderData struct {
	Width Dp
	Shape Shape
	Color graphics.Color
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
			e.borderData.Color == otherEle.borderData.Color
	}
	return false
}

func Border(width Dp, col graphics.Color, shape Shape) Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&BorderElement{
				borderData: BorderData{
					Width: width,
					Shape: shape,
					Color: col,
				},
			},
		),
		modifier.NewInspectorInfo(
			"border",
			map[string]any{
				"width": width,
				"shape": shape,
				"color": col,
			},
		),
	)
}

// Border with defaults usually needs width and color at least?
// For Shape, default to Rectangle if nil? Or caller handles it.
// Default to Rectangle if nil.
func Simple(width Dp, col graphics.Color) Modifier {
	return Border(width, col, shape.ShapeRectangle)
}
