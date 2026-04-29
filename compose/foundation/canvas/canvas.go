package canvas

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/pkg/api"
)

type CanvasOptions struct {
	Modifier ui.Modifier
}

type CanvasOption func(opts *CanvasOptions)

func WithModifier(m ui.Modifier) CanvasOption {
	return func(opts *CanvasOptions) {
		opts.Modifier = m
	}
}

const FoundationCanvasNodeID = "FoundationCanvas"

// composable
func Canvas(
	onDraw func(scope graphics.DrawScope),
	options ...CanvasOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := CanvasOptions{}
		for _, opt := range options {
			if opt != nil {
				opt(&opts)
			}
		}
		density := platform.LocalDensity.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		c.StartBlock(FoundationCanvasNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(widgetConstructor(
			density,
			layoutDirection,
			onDraw,
		))

		return c.EndBlock()
	}
}
