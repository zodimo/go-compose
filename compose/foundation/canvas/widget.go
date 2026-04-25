package canvas

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

func widgetConstructor(
	canvas graphics.Canvas,
	density unit.Density,
	layoutDirection unit.LayoutDirection,
	onDraw func(drawscope graphics.DrawScope),
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			w, h := float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)

			drawscope := graphics.NewCanvasDrawScope(
				canvas,
				geometry.NewSize(
					w,
					h,
				),
				density,
				layoutDirection,
			)
			onDraw(drawscope)

			return layoutnode.LayoutDimensions{
				Size: image.Point{
					X: int(w),
					Y: int(h),
				},
			}
		}
	})
}
