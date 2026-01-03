package animation

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/background"
)

type AnimatedBackgroundElement struct {
	Anim  *animation.VisibilityAnimation
	Color graphics.Color
	Shape background.Shape
}

func (e AnimatedBackgroundElement) Create() Node {
	return NewAnimatedBackgroundNode(e)
}

func (e AnimatedBackgroundElement) Update(node Node) {
	n := node.(*AnimatedBackgroundNode)
	n.element = e
}

func (e AnimatedBackgroundElement) Equals(other Element) bool {
	o, ok := other.(AnimatedBackgroundElement)
	if !ok {
		return false
	}
	return o.Anim == e.Anim && e.Color == o.Color && shape.EqualShape(e.Shape, o.Shape)
}

// AnimatedBackground creates a modifier that paints a background with alpha scaled by animation.
func AnimatedBackground(anim *animation.VisibilityAnimation, color graphics.Color, shape background.Shape) Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			AnimatedBackgroundElement{
				Anim:  anim,
				Color: color,
				Shape: shape,
			},
		),
		modifier.NewInspectorInfo(
			"AnimatedBackground",
			map[string]any{
				"Anim":  anim,
				"Color": color,
			},
		),
	)
}
