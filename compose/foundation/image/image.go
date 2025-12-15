package image

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer
type Modifier = modifier.Modifier

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
		// Correct way to apply modifier in internal composer:
		// We cast modifier.Modifier (internal) to match what C.Modifier expects (internal).
		// Since we imported internal/modifier as modifier, it matches.
		c.Modifier(func(m modifier.Modifier) modifier.Modifier {
			return m.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(imageWidgetConstructor(imageResource, opts))
		return c.EndBlock()
	}
}

func imageWidgetConstructor(resource ImageResource, options ImageOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			srcSize := layoutnode.LayoutDimensions{Size: resource.ImageOp.Size()}
			dstSize := layoutnode.LayoutDimensions{Size: gtx.Constraints.Max}
			_ = options.ContentScale.Scale(srcSize, dstSize)

			// scalefactor to fit and scale

			widget := ImageWidget{
				Src: resource.ImageOp,
				// Scale: options.ContentScale.Scale(srcSize, dstSize),
			}
			return widget.Layout(gtx)
		}
	})
}
