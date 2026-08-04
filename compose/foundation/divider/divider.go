package divider

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gioUnit "gioui.org/unit"
)

const FoundationDivideNodeID = "FoundationDivider"

func Divider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		opts.Color = opts.Color.TakeOrElse(graphics.ColorBlack)
		opts.Axis = opts.Axis.TakeOrElse(layout.AxisHorizontal)

		c.StartBlock(FoundationDivideNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(widgetConstructor(opts, opts.Axis == layout.AxisHorizontal))

		return c.EndBlock()
	}
}

func HorizontalDivider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		opts.Color = opts.Color.TakeOrElse(graphics.ColorBlack)
		opts.Axis = layout.AxisHorizontal

		return Divider(options...)(c)
	}
}

func VerticalDivider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		opts.Color = opts.Color.TakeOrElse(graphics.ColorBlack)
		opts.Axis = layout.AxisVertical

		return Divider(options...)(c)
	}
}

func widgetConstructor(options DividerOptions, isHorizontal bool) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			thickness := gtx.ToGio().Dp(gioUnit.Dp(options.Thickness))
			if thickness < 1 {
				thickness = 1
			}

			// Size
			var size image.Point

			if isHorizontal {
				// Dividers fill the width
				width := gtx.ToGio().Constraints.Min.X
				if gtx.ToGio().Constraints.Max.X > width {
					width = gtx.ToGio().Constraints.Max.X // Or Min/Max strategy? Usually divider fills parent width.
				}
				size = image.Pt(width, thickness)
			} else {
				// Dividers fill the height
				height := gtx.ToGio().Constraints.Min.Y
				if gtx.ToGio().Constraints.Max.Y > height {
					height = gtx.ToGio().Constraints.Max.Y // Usually divider fills parent height.
				}
				size = image.Pt(thickness, height)
			}

			// Resolve Color
			resolvedColor := graphics.ColorToNRGBA(options.Color)

			// Draw
			shape := clip.Rect{Max: size}.Push(gtx.ToGio().Ops)
			paint.ColorOp{Color: resolvedColor}.Add(gtx.ToGio().Ops)
			paint.PaintOp{}.Add(gtx.ToGio().Ops)
			shape.Pop()

			return layoutnode.LayoutDimensions{Size: size}
		}
	})
}
