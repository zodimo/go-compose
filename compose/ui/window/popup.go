package window

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/layout"
	"gioui.org/op"
)

// PopupAlignment determines how the popup is aligned relative to its anchor point.
// For now, we assume the anchor point is the top-left of where this Popup is called.
type PopupAlignment int

const (
	AlignTopLeft PopupAlignment = iota
)

// Popup shows content overlaid on top of other content.
// It uses op.Defer to effectively "break out" of the z-order (though it respects parent clipping).
// The content is laid out with loose constraints (0 to Max).
func Popup(
	content api.Composable,
	options ...PopupOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultPopupOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		c.StartBlock("Popup")
		c.WithComposable(content)
		c.SetWidgetConstructor(popupWidgetConstructor(opts))
		return c.EndBlock()
	}
}

func popupWidgetConstructor(opts PopupOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Record the content layout
			macro := op.Record(gtx.Ops)

			// 2. Prepare context for popup content
			// We give it the same constraints but reset Min to 0 because popup content
			// typically sizes itself.
			// TODO: Consider if we should pass unbounded Max constraints?
			// Usually Popups are constrained by the window size, but here we are constrained
			// by the parent *unless* we were truly global.
			// Since we act as an overlay in-place, we inherit parent max constraints.
			pGtx := gtx
			pGtx.Constraints.Min = image.Point{}

			// Apply offset if needed
			// Converting DpOffset to pixels requires metric, but unit.DpOffset doesn't have a conversion method directly in older gio?
			// Let's check unit.DpOffset usage. It's usually X, Y unit.Dp.
			xPx := pGtx.Dp(unit.DpToGioUnit(opts.OffsetX))
			yPx := pGtx.Dp(unit.DpToGioUnit(opts.OffsetY))

			op.Offset(image.Pt(xPx, yPx)).Add(pGtx.Ops)

			// Layout children
			// We expect usually one child for Popup, but we handle all.
			// We don't use stack/flex here, we just layout them on top of each other (z-stack)
			// effectively.
			for _, child := range node.Children() {
				childLayoutNode := child.(layoutnode.NodeCoordinator)
				_ = childLayoutNode.Layout(pGtx)
			}

			// 3. Stop recording
			call := macro.Stop()

			// 4. Defer the execution
			op.Defer(gtx.Ops, call)

			// 5. Return empty dimensions because the popup doesn't effectively take space in the flow
			return layout.Dimensions{}
		}
	})
}
