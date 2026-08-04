package image

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"gioui.org/op/paint"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

type ImageResource = graphics.ImageResource

func Image(imageResource ImageResource, options ...ImageOption) Composable {
	opts := DefaultImageOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	return func(c Composer) Composer {
		c.StartBlock("Image")

		c.Modifier(func(m ui.Modifier) ui.Modifier {
			return m.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(imageWidgetConstructor(imageResource, opts))
		return c.EndBlock()
	}
}

func imageWidgetConstructor(resource ImageResource, options ImageOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			widget := ImageWidget{
				Src:          paint.NewImageOp(resource.Image()),
				ContentScale: options.ContentScale,
				Alignment:    options.Alignment,
				Alpha:        options.Alpha,
			}
			return widget.Layout(gtx)
		}
	})
}
