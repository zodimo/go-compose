package divider

import (
	"image"

	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gioUnit "gioui.org/unit"
)

const Material3DivideNodeID = "Material3Divider"

func Divider(options ...DividerOption) Composable {
	return func(c Composer) Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		theme := material3.Theme(c)
		opts.Color = opts.Color.TakeOrElse(theme.ColorScheme().OutlineVariant)

		c.StartBlock(Material3DivideNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(widgetConstructor(opts))

		return c.EndBlock()
	}
}

func widgetConstructor(options DividerOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			thickness := gtx.Dp(gioUnit.Dp(options.Thickness))
			if thickness < 1 {
				thickness = 1
			}

			// Dividers fill the width
			width := gtx.Constraints.Min.X
			if gtx.Constraints.Max.X > width {
				width = gtx.Constraints.Max.X // Or Min/Max strategy? Usually divider fills parent width.
			}

			// Size
			size := image.Pt(width, thickness)

			// Resolve Color
			resolvedColor := graphics.ColorToNRGBA(options.Color)

			// Draw
			shape := clip.Rect{Max: size}.Push(gtx.Ops)
			paint.ColorOp{Color: resolvedColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			shape.Pop()

			return layoutnode.LayoutDimensions{Size: size}
		}
	})

}
