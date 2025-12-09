package modifiers

import (
	"go-compose-dev/compose/ui/graphics"
	"go-compose-dev/internal/modifier"
	"image/color"
)

type BackgroundOptions struct {
	Shape graphics.Shape
}

func DefaultBackgroundOptions() BackgroundOptions {
	return BackgroundOptions{
		Shape: graphics.ShapeRectangle,
	}
}

type BackgroundOption func(options *BackgroundOptions)

func BackgroundColor(color color.Color, options ...BackgroundOption) Modifier {

	opt := DefaultBackgroundOptions()
	for _, option := range options {
		option(&opt)
	}
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&backgroundElement{
				background: Background{
					color: color,
					shape: opt.Shape,
				},
			},
		),
		modifier.NewInspectorInfo(
			"background",
			map[string]any{
				"color":   color,
				"options": opt,
			},
		),
	)
}
